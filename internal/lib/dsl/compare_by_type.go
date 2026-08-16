package dsl

import (
	"ab_system/internal/http/dto"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func compareByType(raw any, cmp ComparisonExpr, fs dto.FieldSchema) (bool, string) {
	switch fs.Type {
	case dto.FieldTypeNumber:
		return compareNumber(raw, cmp)
	case dto.FieldTypeBool:
		return compareBool(raw, cmp)
	case dto.FieldTypeDate:
		return compareDate(raw, cmp)
	case dto.FieldTypeVersion:
		return compareVersion(raw, cmp)
	default:
		return compareString(raw, cmp)
	}
}

func compareNumber(raw any, cmp ComparisonExpr) (bool, string) {
	actualF, err := toFloat(raw)
	if err != nil {
		return false, fmt.Sprintf("поле '%s': %v", cmp.Field, err)
	}

	if isListOperator(cmp.Operator) {
		vals := make([]float64, 0, len(cmp.Values))
		for _, v := range cmp.Values {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return false, fmt.Sprintf("неверное число в списке: %s", v)
			}
			vals = append(vals, f)
		}

		inList := false

		for _, v := range vals {
			if floatEq(actualF, v) {
				inList = true
				break
			}
		}

		if strings.ToUpper(cmp.Operator) == "IN" {
			return inList, fmt.Sprintf("%s IN %v", cmp.Field, vals)
		}

		return !inList, fmt.Sprintf("%s NOT IN %v", cmp.Field, vals)
	}

	threshold, err := strconv.ParseFloat(cmp.Value, 64)
	if err != nil {
		return false, fmt.Sprintf("неверный порог: %s", cmp.Value)
	}

	var ok bool
	switch cmp.Operator {
	case "=", "==":
		ok = floatEq(actualF, threshold)
	case "!=":
		ok = !floatEq(actualF, threshold)
	case ">":
		ok = actualF > threshold
	case ">=":
		ok = actualF >= threshold
	case "<":
		ok = actualF < threshold
	case "<=":
		ok = actualF <= threshold
	default:
		return false, fmt.Sprintf("неподдерживаемый оператор для числа: %s", cmp.Operator)
	}

	return ok, fmt.Sprintf("%s(%v) %s %v", cmp.Field, actualF, cmp.Operator, threshold)
}

func compareBool(raw any, cmp ComparisonExpr) (bool, string) {
	actualB, err := toBool(raw)
	if err != nil {
		return false, fmt.Sprintf("поле '%s': %v", cmp.Field, err)
	}

	expectedB := strings.ToLower(cmp.Value) == "true" || cmp.Value == "1"

	var ok bool
	switch cmp.Operator {
	case "=", "==":
		ok = actualB == expectedB
	case "!=":
		ok = actualB != expectedB
	default:
		return false, fmt.Sprintf("неподдерживаемый оператор для bool: %s", cmp.Operator)
	}

	return ok, fmt.Sprintf("%s(%v) %s %v", cmp.Field, actualB, cmp.Operator, expectedB)
}

func compareDate(raw any, cmp ComparisonExpr) (bool, string) {
	actualStr := fmt.Sprintf("%v", raw)

	actualD, err := time.Parse(dateLayout, actualStr)
	if err != nil {
		return false, fmt.Sprintf("поле '%s': неверный формат даты '%s'", cmp.Field, actualStr)
	}

	if isListOperator(cmp.Operator) {
		dates := make([]time.Time, 0, len(cmp.Values))
		for _, v := range cmp.Values {
			d, err := time.Parse(dateLayout, v)
			if err != nil {
				return false, fmt.Sprintf("неверная дата в списке: %s", v)
			}
			dates = append(dates, d)
		}

		inList := false

		for _, d := range dates {
			if actualD.Equal(d) {
				inList = true
				break
			}
		}

		if strings.ToUpper(cmp.Operator) == "IN" {
			return inList, fmt.Sprintf("%s IN [dates]", cmp.Field)
		}

		return !inList, fmt.Sprintf("%s NOT IN [dates]", cmp.Field)
	}

	threshold, err := time.Parse(dateLayout, cmp.Value)
	if err != nil {
		return false, fmt.Sprintf("неверная дата порога: %s", cmp.Value)
	}

	var ok bool
	switch cmp.Operator {
	case "=", "==":
		ok = actualD.Equal(threshold)
	case "!=":
		ok = !actualD.Equal(threshold)
	case ">":
		ok = actualD.After(threshold)
	case ">=":
		ok = !actualD.Before(threshold)
	case "<":
		ok = actualD.Before(threshold)
	case "<=":
		ok = !actualD.After(threshold)
	default:
		return false, fmt.Sprintf("неподдерживаемый оператор для даты: %s", cmp.Operator)
	}

	return ok, fmt.Sprintf("%s(%s) %s %s", cmp.Field, actualStr, cmp.Operator, cmp.Value)
}

func compareString(raw any, cmp ComparisonExpr) (bool, string) {
	actual := fmt.Sprintf("%v", raw)

	if isListOperator(cmp.Operator) {
		inList := false
		for _, v := range cmp.Values {
			if actual == v {
				inList = true
				break
			}
		}

		if strings.ToUpper(cmp.Operator) == "IN" {
			return inList, fmt.Sprintf("%s IN %v", cmp.Field, cmp.Values)
		}

		return !inList, fmt.Sprintf("%s NOT IN %v", cmp.Field, cmp.Values)
	}

	var ok bool
	switch cmp.Operator {
	case "=", "==":
		ok = actual == cmp.Value
	case "!=":
		ok = actual != cmp.Value
	default:
		return false, fmt.Sprintf("неподдерживаемый оператор для строки: %s", cmp.Operator)
	}

	return ok, fmt.Sprintf("%s(%q) %s %q", cmp.Field, actual, cmp.Operator, cmp.Value)
}

func compareVersion(raw any, cmp ComparisonExpr) (bool, string) {
	actual := fmt.Sprintf("%v", raw)

	if isListOperator(cmp.Operator) {
		inList := false
		for _, v := range cmp.Values {
			if compareVersionStrings(actual, v) == 0 {
				inList = true
				break
			}
		}

		if strings.ToUpper(cmp.Operator) == "IN" {
			return inList, fmt.Sprintf("%s IN %v", cmp.Field, cmp.Values)
		}

		return !inList, fmt.Sprintf("%s NOT IN %v", cmp.Field, cmp.Values)
	}

	cmpResult := compareVersionStrings(actual, cmp.Value)
	var ok bool
	switch cmp.Operator {
	case "=", "==":
		ok = cmpResult == 0
	case "!=":
		ok = cmpResult != 0
	case ">":
		ok = cmpResult > 0
	case ">=":
		ok = cmpResult >= 0
	case "<":
		ok = cmpResult < 0
	case "<=":
		ok = cmpResult <= 0
	default:
		return false, fmt.Sprintf("неподдерживаемый оператор для версии: %s", cmp.Operator)
	}

	return ok, fmt.Sprintf("%s(%s) %s %s", cmp.Field, actual, cmp.Operator, cmp.Value)
}

func compareVersionStrings(a, b string) int {
	partsA := splitVersion(a)
	partsB := splitVersion(b)

	for len(partsA) < len(partsB) {
		partsA = append(partsA, 0)
	}

	for len(partsB) < len(partsA) {
		partsB = append(partsB, 0)
	}

	for i := range partsA {
		if partsA[i] < partsB[i] {
			return -1
		}
		if partsA[i] > partsB[i] {
			return 1
		}
	}

	return 0
}

func splitVersion(v string) []int {
	parts := strings.Split(v, ".")
	result := make([]int, 0, len(parts))
	for _, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			n = 0
		}

		result = append(result, n)
	}

	return result
}

func compareAsString(raw any, cmp ComparisonExpr) bool {
	actual := fmt.Sprintf("%v", raw)
	if isListOperator(cmp.Operator) {
		for _, v := range cmp.Values {
			if actual == v {
				return strings.ToUpper(cmp.Operator) == "IN"
			}
		}

		return strings.ToUpper(cmp.Operator) == "NOT IN"
	}

	switch cmp.Operator {
	case "=", "==":
		return actual == cmp.Value
	case "!=":
		return actual != cmp.Value
	}

	return false
}
