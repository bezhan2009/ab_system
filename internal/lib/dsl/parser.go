package dsl

import (
	"fmt"
	"strings"
)

type ComparisonExpr struct {
	Field    string
	Operator string
	Value    string
	Values   []string
}

func parseComparison(expr string) (ComparisonExpr, error) {
	expr = strings.TrimSpace(expr)

	if cmp, ok, err := tryParseInOperator(expr, "NOT IN"); ok || err != nil {
		return cmp, err
	}

	if cmp, ok, err := tryParseInOperator(expr, "IN"); ok || err != nil {
		return cmp, err
	}

	twoCharOps := []string{">=", "<=", "!=", "=="}
	for _, op := range twoCharOps {
		if idx := strings.Index(expr, op); idx != -1 {
			return parseScalarComparison(expr, idx, op)
		}
	}

	oneCharOps := []string{">", "<", "="}
	for _, op := range oneCharOps {
		if idx := strings.Index(expr, op); idx != -1 {
			return parseScalarComparison(expr, idx, op)
		}
	}

	return ComparisonExpr{}, fmt.Errorf("не найден оператор в выражении: %s", expr)
}

func tryParseInOperator(expr, op string) (ComparisonExpr, bool, error) {
	exprUp := strings.ToUpper(expr)
	opUp := strings.ToUpper(op)

	idx := strings.Index(exprUp, " "+opUp+" ")
	if idx == -1 {
		return ComparisonExpr{}, false, nil
	}

	field := strings.TrimSpace(expr[:idx])
	rest := strings.TrimSpace(expr[idx+len(op)+2:])

	if !strings.HasPrefix(rest, "(") {
		return ComparisonExpr{}, false, fmt.Errorf("после %s ожидается список в скобках: %s", op, rest)
	}

	closeIdx := strings.LastIndex(rest, ")")
	if closeIdx == -1 {
		return ComparisonExpr{}, false, fmt.Errorf("не закрыта скобка в списке %s: %s", op, rest)
	}

	inner := rest[1:closeIdx]
	rawItems := strings.Split(inner, ",")
	values := make([]string, 0, len(rawItems))
	for _, item := range rawItems {
		item = strings.TrimSpace(item)
		item = unquoteString(item)
		if item == "" {
			return ComparisonExpr{}, false, fmt.Errorf("пустой элемент в списке %s", op)
		}
		values = append(values, item)
	}

	if len(values) == 0 {
		return ComparisonExpr{}, false, fmt.Errorf("пустой список значений для %s", op)
	}

	operatorStr := op
	return ComparisonExpr{
		Field:    field,
		Operator: operatorStr,
		Values:   values,
	}, true, nil
}

func parseScalarComparison(expr string, idx int, op string) (ComparisonExpr, error) {
	field := strings.TrimSpace(expr[:idx])
	rawValue := strings.TrimSpace(expr[idx+len(op):])

	if field == "" {
		return ComparisonExpr{}, fmt.Errorf("пустое имя поля")
	}

	value, err := parseValue(rawValue)
	if err != nil {
		return ComparisonExpr{}, err
	}

	return ComparisonExpr{
		Field:    field,
		Operator: op,
		Value:    value,
	}, nil
}

func parseValue(raw string) (string, error) {
	raw = strings.TrimSpace(raw)

	if raw == "" {
		return "", fmt.Errorf("пустое значение")
	}

	if strings.HasPrefix(raw, "'") {
		return parseQuotedString(raw)
	}

	if strings.ContainsAny(raw, " \t\n") {
		return "", fmt.Errorf("значение содержит пробелы без кавычек: %s", raw)
	}

	return raw, nil
}

func parseQuotedString(raw string) (string, error) {
	if !strings.HasPrefix(raw, "'") {
		return "", fmt.Errorf("ожидается открывающая кавычка")
	}

	var sb strings.Builder
	i := 1
	for i < len(raw) {
		ch := raw[i]
		if ch == '\\' && i+1 < len(raw) {
			next := raw[i+1]
			switch next {
			case '\'':
				sb.WriteByte('\'')
			case '\\':
				sb.WriteByte('\\')
			default:
				sb.WriteByte('\\')
				sb.WriteByte(next)
			}
			i += 2
			continue
		}
		if ch == '\'' {
			afterQuote := strings.TrimSpace(raw[i+1:])
			if afterQuote != "" {
				return "", fmt.Errorf("лишние символы после строки: %s", raw[i+1:])
			}
			return sb.String(), nil
		}
		sb.WriteByte(ch)
		i++
	}

	return "", fmt.Errorf("не найдена закрывающая кавычка: %s", raw)
}

func unquoteString(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'") && len(s) >= 2 {
		return s[1 : len(s)-1]
	}
	return s
}

func splitByOperatorOutsideParentheses(expr, operator string) ([]string, error) {
	parts := []string{}
	current := strings.Builder{}
	balance := 0
	inString := false

	for i := 0; i < len(expr); i++ {
		ch := expr[i]

		if ch == '\'' && !inString {
			inString = true
			current.WriteByte(ch)
			continue
		}
		if ch == '\'' && inString {
			if i > 0 && expr[i-1] != '\\' {
				inString = false
			}
			current.WriteByte(ch)
			continue
		}
		if inString {
			current.WriteByte(ch)
			continue
		}

		if ch == '(' {
			balance++
			current.WriteByte(ch)
		} else if ch == ')' {
			if balance == 0 {
				return nil, fmt.Errorf("лишняя закрывающая скобка в позиции %d", i)
			}
			balance--
			current.WriteByte(ch)
		} else if balance == 0 && i+len(operator) <= len(expr) &&
			strings.EqualFold(expr[i:i+len(operator)], operator) {
			parts = append(parts, strings.TrimSpace(current.String()))
			current.Reset()
			i += len(operator) - 1
		} else {
			current.WriteByte(ch)
		}
	}

	if inString {
		return nil, fmt.Errorf("незакрытая строковая кавычка")
	}

	if current.Len() > 0 || len(parts) > 0 {
		parts = append(parts, strings.TrimSpace(current.String()))
	}

	if balance != 0 {
		return nil, fmt.Errorf("несбалансированные скобки")
	}

	return parts, nil
}
