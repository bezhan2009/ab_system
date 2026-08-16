package dsl

import (
	"ab_system/internal/http/dto"
	"strings"
)

type ValidationError struct {
	Code     string
	Message  string
	Position *int
	Near     *string
}

func (e ValidationError) Error() string {
	return e.Message
}

type ValidationResult struct {
	IsValid              bool
	NormalizedExpression *string
	Errors               []ValidationError
}

// Error codes
const (
	ErrCodeParseError      = "DSL_PARSE_ERROR"
	ErrCodeInvalidField    = "DSL_INVALID_FIELD"
	ErrCodeInvalidOperator = "DSL_INVALID_OPERATOR"
	ErrCodeInvalidValue    = "DSL_INVALID_VALUE"
)

// Warning codes
const (
	ErrCodeDslIsEmpty = "DSL_IS_EMPTY"
)

func ValidateDSL(expression string, schema dto.Schema) ValidationResult {
	if expression == "" {
		pos := 0
		near := ""
		normalizedExpression := ""

		return ValidationResult{
			IsValid:              true,
			NormalizedExpression: &normalizedExpression,
			Errors: []ValidationError{{
				Code:     ErrCodeDslIsEmpty,
				Message:  "пустое выражение",
				Position: &pos,
				Near:     &near,
			}},
		}
	}

	normalized := normalizeExpression(expression)

	result := validateExprFull(normalized, expression, 0, schema)

	if result.IsValid {
		result.NormalizedExpression = &normalized
	}

	return result
}

func validateExprFull(expr, orig string, basePos int, schema dto.Schema) ValidationResult {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		pos := basePos

		return ValidationResult{
			IsValid: false,
			Errors: []ValidationError{{
				Code:     ErrCodeParseError,
				Message:  "пустое выражение после нормализации",
				Position: &pos,
				Near:     &orig,
			}},
		}
	}

	return validateOR(expr, orig, basePos, schema)
}

func validateOR(expr, orig string, basePos int, schema dto.Schema) ValidationResult {
	parts, err := splitByOperatorOutsideParentheses(expr, " OR ")
	if err != nil {
		return singleError(ErrCodeParseError, err.Error(), orig, expr, basePos)
	}

	if len(parts) == 1 {
		return validateAND(parts[0], orig, findErrorPosition(orig, parts[0], basePos), schema)
	}

	return collectResults(parts, orig, expr, basePos, schema, validateAND)
}

func validateAND(expr, orig string, basePos int, schema dto.Schema) ValidationResult {
	parts, err := splitByOperatorOutsideParentheses(expr, " AND ")
	if err != nil {
		return singleError(ErrCodeParseError, err.Error(), orig, expr, basePos)
	}

	if len(parts) == 1 {
		return validateFactor(strings.TrimSpace(parts[0]), orig, findErrorPosition(orig, parts[0], basePos), schema)
	}

	return collectResults(parts, orig, expr, basePos, schema, validateFactor)
}

func validateFactor(expr, orig string, basePos int, schema dto.Schema) ValidationResult {
	expr = strings.TrimSpace(expr)

	if strings.HasPrefix(strings.ToUpper(expr), "NOT ") {
		inner := strings.TrimSpace(expr[4:])

		return validateFactor(inner, orig, findErrorPosition(orig, inner, basePos), schema)
	}

	if strings.HasPrefix(expr, "(") && strings.HasSuffix(expr, ")") {
		inner, ok := extractParenInner(expr)
		if !ok {
			return singleError(ErrCodeParseError, "несбалансированные скобки", orig, expr, basePos)
		}

		return validateOR(strings.TrimSpace(inner), orig, findErrorPosition(orig, inner, basePos), schema)
	}

	return validateComparison(expr, orig, basePos, schema)
}

func validateComparison(expr, orig string, basePos int, schema dto.Schema) ValidationResult {
	cmp, err := parseComparison(expr)
	if err != nil {
		return singleError(ErrCodeParseError, err.Error(), orig, expr, basePos)
	}

	if err = validateFieldWithSchema(cmp, schema); err != nil {
		code := ErrCodeParseError
		switch {
		case strings.HasPrefix(err.Error(), "неизвестное поле"):
			code = ErrCodeInvalidField
		case strings.HasPrefix(err.Error(), "оператор"):
			code = ErrCodeInvalidOperator
		default:
			code = ErrCodeInvalidValue
		}

		return singleError(code, err.Error(), orig, expr, basePos)
	}

	return ValidationResult{IsValid: true, Errors: []ValidationError{}}
}

type validateFunc func(expr, orig string, basePos int, schema dto.Schema) ValidationResult

func collectResults(parts []string, orig, parentExpr string, basePos int, schema dto.Schema, fn validateFunc) ValidationResult {
	var allErrors []ValidationError
	allValid := true

	for _, part := range parts {
		part = strings.TrimSpace(part)
		pos := findErrorPosition(orig, part, basePos)
		result := fn(part, orig, pos, schema)
		if !result.IsValid {
			allValid = false
			allErrors = append(allErrors, result.Errors...)
		}
	}

	return ValidationResult{IsValid: allValid, Errors: allErrors}
}

func singleError(code, message, orig, subExpr string, basePos int) ValidationResult {
	pos := findErrorPosition(orig, subExpr, basePos)
	near := getNearContext(orig, pos)

	return ValidationResult{
		IsValid: false,
		Errors: []ValidationError{{
			Code:     code,
			Message:  message,
			Position: &pos,
			Near:     &near,
		}},
	}
}

func extractParenInner(expr string) (string, bool) {
	if len(expr) < 2 || expr[0] != '(' || expr[len(expr)-1] != ')' {
		return "", false
	}

	balance := 0

	for i, ch := range expr {
		if ch == '(' {
			balance++
		} else if ch == ')' {
			balance--
			if balance == 0 && i != len(expr)-1 {
				return "", false
			}
		}
	}

	if balance != 0 {
		return "", false
	}

	return expr[1 : len(expr)-1], true
}
