package dsl

import (
	"strings"
)

func normalizeExpression(expr string) string {
	expr = strings.TrimSpace(expr)
	expr = normalizeKeywords(expr)
	expr = normalizeOperators(expr)
	expr = normalizeSpaces(expr)
	return removeUnnecessaryParentheses(expr)
}

func removeUnnecessaryParentheses(expr string) string {
	expr = removeOuterParentheses(expr)
	return simplifyExpression(expr)
}

func removeOuterParentheses(expr string) string {
	expr = strings.TrimSpace(expr)

	if len(expr) < 2 || expr[0] != '(' || expr[len(expr)-1] != ')' {
		return expr
	}

	balance := 0
	for i, ch := range expr {
		if ch == '(' {
			balance++
		} else if ch == ')' {
			balance--
			if balance == 0 && i != len(expr)-1 {
				return expr
			}
		}
	}

	return strings.TrimSpace(expr[1 : len(expr)-1])
}

func normalizeKeywords(expr string) string {
	var result strings.Builder
	i := 0

	for i < len(expr) {
		ch := expr[i]

		if isLetter(ch) {
			start := i
			for i < len(expr) && (isLetter(expr[i]) || expr[i] == '_' || expr[i] == '.') {
				i++
			}
			word := expr[start:i]
			upperWord := strings.ToUpper(word)

			switch upperWord {
			case "AND", "OR", "NOT":
				result.WriteString(upperWord)
			default:
				result.WriteString(word)
			}
		} else {
			result.WriteByte(ch)
			i++
		}
	}

	return result.String()
}

func normalizeOperators(expr string) string {
	var result strings.Builder
	i := 0
	inString := false

	for i < len(expr) {
		ch := expr[i]

		if ch == '\'' && !inString {
			inString = true
			result.WriteByte(ch)
			i++
			continue
		}

		if ch == '\'' && inString {
			if i == 0 || expr[i-1] != '\\' {
				inString = false
			}
			result.WriteByte(ch)
			i++
			continue
		}

		if inString {
			result.WriteByte(ch)
			i++
			continue
		}

		if isComparisonOperatorStart(ch) {
			op := string(ch)
			if i+1 < len(expr) {
				twoChar := expr[i : i+2]
				if twoChar == ">=" || twoChar == "<=" || twoChar == "!=" || twoChar == "==" {
					op = twoChar
				}
			}

			result.WriteString(" " + op + " ")
			i += len(op)
		} else {
			result.WriteByte(ch)
			i++
		}
	}

	return result.String()
}

func normalizeSpaces(expr string) string {
	fields := strings.Fields(expr)
	return strings.Join(fields, " ")
}

func simplifyExpression(expr string) string {
	if isSimpleComparison(expr) {
		return expr
	}

	if canSplitByOperator(expr, " OR ") {
		parts := splitTopLevelByOperator(expr, " OR ")
		if len(parts) > 1 {
			simplified := make([]string, len(parts))
			for i, part := range parts {
				simplified[i] = simplifyExpression(removeOuterParentheses(part))
			}
			result := strings.Join(simplified, " OR ")
			return removeParenthesesAroundAndInOr(result)
		}
	}

	if canSplitByOperator(expr, " AND ") {
		parts := splitTopLevelByOperator(expr, " AND ")
		if len(parts) > 1 {
			simplified := make([]string, len(parts))
			for i, part := range parts {
				simplified[i] = simplifyExpression(removeOuterParentheses(part))
			}
			return strings.Join(simplified, " AND ")
		}
	}

	if strings.HasPrefix(expr, "NOT ") {
		inner := strings.TrimSpace(expr[4:])
		inner = removeOuterParentheses(inner)
		if containsSpaces(inner) && !strings.HasPrefix(inner, "(") {
			inner = "(" + inner + ")"
		}
		return "NOT " + simplifyExpression(inner)
	}

	if strings.HasPrefix(expr, "(") && strings.HasSuffix(expr, ")") {
		inner := expr[1 : len(expr)-1]
		simplified := simplifyExpression(inner)
		if needsParentheses(simplified) {
			return "(" + simplified + ")"
		}
		return simplified
	}

	return expr
}

func removeParenthesesAroundAndInOr(expr string) string {
	orParts := splitTopLevelByOperator(expr, " OR ")
	if len(orParts) <= 1 {
		return expr
	}

	simplified := make([]string, len(orParts))
	for i, part := range orParts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "(") && strings.HasSuffix(part, ")") {
			inner := part[1 : len(part)-1]
			innerOrParts := splitTopLevelByOperator(inner, " OR ")
			if len(innerOrParts) == 1 {
				simplified[i] = simplifyExpression(inner)
			} else {
				simplified[i] = "(" + simplifyExpression(inner) + ")"
			}
		} else {
			simplified[i] = simplifyExpression(part)
		}
	}
	return strings.Join(simplified, " OR ")
}

func splitTopLevelByOperator(expr, operator string) []string {
	expr = strings.TrimSpace(expr)
	var parts []string
	var current strings.Builder
	balance := 0

	for i := 0; i < len(expr); i++ {
		ch := expr[i]

		if ch == '(' {
			balance++
		} else if ch == ')' {
			balance--
		}

		if balance == 0 && i+len(operator) <= len(expr) &&
			strings.EqualFold(expr[i:i+len(operator)], operator) {
			parts = append(parts, strings.TrimSpace(current.String()))
			current.Reset()
			i += len(operator) - 1
			continue
		}

		current.WriteByte(ch)
	}

	if current.Len() > 0 {
		parts = append(parts, strings.TrimSpace(current.String()))
	}

	if len(parts) == 0 {
		return []string{expr}
	}

	return parts
}
