package dsl

import (
	"strings"
)

func findErrorPosition(originalExpr, subExpr string, basePos int) int {
	if subExpr == "" {
		return basePos
	}

	cleanOriginal := strings.ReplaceAll(originalExpr, " ", "")
	cleanSubExpr := strings.ReplaceAll(subExpr, " ", "")

	idx := strings.Index(strings.ToUpper(cleanOriginal), strings.ToUpper(cleanSubExpr))
	if idx >= 0 {
		return idx
	}

	return basePos
}

func getNearContext(expr string, position int) string {
	start := maxInt(0, position-10)
	end := minInt(len(expr), position+10)
	return expr[start:end]
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func canSplitByOperator(expr, operator string) bool {
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
			return true
		}
	}
	return false
}

func containsSpaces(s string) bool {
	return strings.Contains(s, " ")
}

func needsParentheses(expr string) bool {
	if strings.HasPrefix(expr, "(") && strings.HasSuffix(expr, ")") {
		return false
	}
	if canSplitByOperator(expr, " OR ") && canSplitByOperator(expr, " AND ") {
		return true
	}
	return false
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isComparisonOperatorStart(ch byte) bool {
	return ch == '>' || ch == '<' || ch == '=' || ch == '!'
}

func containsLogicalOperator(s string) bool {
	sUpper := strings.ToUpper(s)
	return strings.Contains(sUpper, " AND ") ||
		strings.Contains(sUpper, " OR ") ||
		strings.Contains(sUpper, " NOT ")
}

func isSimpleComparison(expr string) bool {
	expr = strings.TrimSpace(expr)

	if strings.HasPrefix(strings.ToUpper(expr), "NOT ") {
		expr = strings.TrimSpace(expr[4:])
	}

	expr = removeOuterParentheses(expr)

	exprUp := strings.ToUpper(expr)
	if strings.Contains(exprUp, " NOT IN ") || strings.Contains(exprUp, " IN ") {
		return !containsLogicalOperator(expr)
	}

	operators := []string{">=", "<=", "!=", ">", "<", "="}
	for _, op := range operators {
		if idx := strings.Index(expr, op); idx != -1 {
			before := strings.TrimSpace(expr[:idx])
			after := strings.TrimSpace(expr[idx+len(op):])
			if !containsLogicalOperator(before) && !containsLogicalOperator(after) {
				return true
			}
		}
	}

	return false
}
