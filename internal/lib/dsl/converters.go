package dsl

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func toFloat(v any) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case int32:
		return float64(val), nil
	case uint:
		return float64(val), nil
	case uint64:
		return float64(val), nil
	case string:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, fmt.Errorf("не удалось привести к числу: %q", val)
		}
		return f, nil
	default:
		return 0, fmt.Errorf("неподдерживаемый тип для числового сравнения: %T", v)
	}
}

func toBool(v any) (bool, error) {
	switch val := v.(type) {
	case bool:
		return val, nil
	case int, int64, int32, uint, uint64:
		f, _ := toFloat(val)
		return f != 0, nil
	case string:
		lower := strings.ToLower(val)
		if lower == "true" || lower == "1" {
			return true, nil
		}
		if lower == "false" || lower == "0" {
			return false, nil
		}
		return false, fmt.Errorf("не удалось привести к bool: %q", val)
	default:
		return false, fmt.Errorf("неподдерживаемый тип для bool: %T", v)
	}
}

func floatEq(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}
