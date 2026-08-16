package dsl

import (
	"ab_system/internal/http/dto"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "2006-01-02"

func validateFieldWithSchema(cmp ComparisonExpr, schema dto.Schema) error {
	fieldSchema, ok := schema.GetSchemeByName(cmp.Field)
	if !ok {
		return fmt.Errorf("неизвестное поле: %s", cmp.Field)
	}

	if !fieldSchema.OperatorAllowed(cmp.Operator) {
		allowed := fieldSchema.AllowedOperators
		if len(allowed) == 0 {
			allowed = dto.DefaultOperators(fieldSchema.Type)
		}
		return fmt.Errorf(
			"оператор '%s' не допустим для поля '%s' (тип %s); допустимые: %s",
			cmp.Operator, cmp.Field, fieldSchema.Type,
			strings.Join(allowed, ", "),
		)
	}

	if isListOperator(cmp.Operator) {
		for _, v := range cmp.Values {
			if err := validateValueForType(v, fieldSchema.Type, cmp.Field); err != nil {
				return err
			}
		}
		return nil
	}

	return validateValueForType(cmp.Value, fieldSchema.Type, cmp.Field)
}

func validateValueForType(value string, fieldType dto.FieldType, fieldName string) error {
	switch fieldType {
	case dto.FieldTypeNumber:
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return fmt.Errorf("неверное числовое значение '%s' для поля '%s'", value, fieldName)
		}
	case dto.FieldTypeBool:
		lower := strings.ToLower(value)
		if lower != "true" && lower != "false" && lower != "1" && lower != "0" {
			return fmt.Errorf("неверное bool-значение '%s' для поля '%s' (ожидается true/false)", value, fieldName)
		}
	case dto.FieldTypeDate:
		if _, err := time.Parse(dateLayout, value); err != nil {
			return fmt.Errorf("неверный формат даты '%s' для поля '%s' (ожидается YYYY-MM-DD)", value, fieldName)
		}
	case dto.FieldTypeVersion:
		if !isValidVersionString(value) {
			return fmt.Errorf("неверный формат версии '%s' для поля '%s' (ожидается N.N.N)", value, fieldName)
		}
	case dto.FieldTypeString:
		if value == "" {
			return fmt.Errorf("пустое строковое значение для поля '%s'", fieldName)
		}
	}
	return nil
}

func isValidVersionString(v string) bool {
	if v == "" {
		return false
	}
	for _, ch := range v {
		if ch != '.' && (ch < '0' || ch > '9') {
			return false
		}
	}
	return true
}

func isListOperator(op string) bool {
	opUp := strings.ToUpper(op)
	return opUp == "IN" || opUp == "NOT IN"
}
