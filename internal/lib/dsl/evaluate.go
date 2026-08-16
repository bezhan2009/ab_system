package dsl

import (
	"ab_system/internal/http/dto"
	"fmt"
	"strings"
)

type EvalResult struct {
	Matched       bool
	Description   string
	MissingFields []string
}

func EvaluateDSL(expression string, payload dto.Payload, schema dto.Schema) EvalResult {
	if expression == "" {
		return EvalResult{Matched: false, Description: "пустое выражение"}
	}

	expr := normalizeExpression(expression)

	missing := &[]string{}
	matched, desc := evalOR(expr, payload, schema, missing)

	return EvalResult{
		Matched:       matched,
		Description:   desc,
		MissingFields: *missing,
	}
}

func evalOR(expr string, payload dto.Payload, schema dto.Schema, missing *[]string) (bool, string) {
	parts, err := splitByOperatorOutsideParentheses(expr, " OR ")
	if err != nil || len(parts) <= 1 {
		return evalAND(expr, payload, schema, missing)
	}

	for i, part := range parts {
		ok, desc := evalAND(strings.TrimSpace(part), payload, schema, missing)
		if ok {
			return true, fmt.Sprintf("OR[%d] истина: %s", i+1, desc)
		}
	}

	return false, "все OR-части ложны"
}

func evalAND(expr string, payload dto.Payload, schema dto.Schema, missing *[]string) (bool, string) {
	parts, err := splitByOperatorOutsideParentheses(expr, " AND ")
	if err != nil || len(parts) <= 1 {
		return evalFactor(expr, payload, schema, missing)
	}

	descs := make([]string, 0, len(parts))
	for i, part := range parts {
		ok, desc := evalFactor(strings.TrimSpace(part), payload, schema, missing)
		if !ok {
			return false, fmt.Sprintf("AND[%d] ложь: %s", i+1, desc)
		}
		descs = append(descs, desc)
	}

	return true, strings.Join(descs, " AND ")
}

func evalFactor(expr string, payload dto.Payload, schema dto.Schema, missing *[]string) (bool, string) {
	expr = strings.TrimSpace(expr)

	if strings.HasPrefix(strings.ToUpper(expr), "NOT ") {
		inner := strings.TrimSpace(expr[4:])
		ok, desc := evalFactor(inner, payload, schema, missing)
		return !ok, fmt.Sprintf("NOT(%s)", desc)
	}

	if strings.HasPrefix(expr, "(") && strings.HasSuffix(expr, ")") {
		inner, ok := extractParenInner(expr)
		if !ok {
			return false, "несбалансированные скобки"
		}
		return evalOR(strings.TrimSpace(inner), payload, schema, missing)
	}

	return evalComparison(expr, payload, schema, missing)
}

func evalComparison(expr string, payload dto.Payload, schema dto.Schema, missing *[]string) (bool, string) {
	cmp, err := parseComparison(expr)
	if err != nil {
		return false, fmt.Sprintf("ошибка разбора: %v", err)
	}

	rawValue, exists := payload.Get(cmp.Field)
	if !exists {
		*missing = append(*missing, cmp.Field)
		return false, fmt.Sprintf("поле '%s' отсутствует в payload: false", cmp.Field)
	}

	fieldSchema, schemaKnown := schema.GetSchemeByName(cmp.Field)
	if !schemaKnown {
		return compareAsString(rawValue, cmp), fmt.Sprintf("(нету схемы) %s %s", cmp.Field, cmp.Operator)
	}

	return compareByType(rawValue, cmp, fieldSchema)
}
