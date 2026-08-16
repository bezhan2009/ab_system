package dto

import (
	"fmt"
	"strings"
)

type FieldType string

const (
	FieldTypeString  FieldType = "string"
	FieldTypeNumber  FieldType = "number"
	FieldTypeBool    FieldType = "bool"
	FieldTypeDate    FieldType = "date"
	FieldTypeVersion FieldType = "version"
)

type FieldSchema struct {
	Name             string
	Type             FieldType
	AllowedOperators []string
}

type Schema map[string]FieldSchema

type Payload map[string]any

type SchemaValidationError struct {
	Field   string
	Message string
}

func (e SchemaValidationError) Error() string {
	return fmt.Sprintf("ошибка схемы для поля '%s': %s", e.Field, e.Message)
}

func (p Payload) Get(key string) (any, bool) {
	v, ok := p[strings.ToLower(key)]
	if ok {
		return v, true
	}

	lk := strings.ToLower(key)
	for k, val := range p {
		if strings.ToLower(k) == lk {
			return val, true
		}
	}

	return nil, false
}

func DefaultOperators(t FieldType) []string {
	switch t {
	case FieldTypeNumber, FieldTypeDate, FieldTypeVersion:
		return []string{"=", "!=", ">", ">=", "<", "<="}
	case FieldTypeBool:
		return []string{"=", "!="}
	default: // это для типа string
		return []string{"=", "!=", "IN", "NOT IN"}
	}
}

func (s FieldSchema) OperatorAllowed(op string) bool {
	ops := s.AllowedOperators
	if len(ops) == 0 {
		ops = DefaultOperators(s.Type)
	}
	opUp := strings.ToUpper(op)
	for _, allowed := range ops {
		if strings.ToUpper(allowed) == opUp {
			return true
		}
	}
	return false
}

func (sc Schema) GetSchemeByName(name string) (FieldSchema, bool) {
	f, ok := sc[strings.ToLower(name)]
	return f, ok
}

// Функция HardcodedSchema сгенерирована с помощью Исскуственного Интеллекта(рутинная задача)
func HardcodedSchema() map[string]FieldSchema {
	return map[string]FieldSchema{
		// Стандартные поля пользователя
		"user.country": {
			Name:             "user.country",
			Type:             FieldTypeString,
			AllowedOperators: []string{"==", "!=", "IN", "NOT IN"},
		},
		"user.city": {
			Name:             "user.city",
			Type:             FieldTypeString,
			AllowedOperators: []string{"==", "!=", "IN", "NOT IN"},
		},
		"user.language": {
			Name:             "user.language",
			Type:             FieldTypeString,
			AllowedOperators: []string{"==", "!="},
		},
		"user.gender": {
			Name:             "user.gender",
			Type:             FieldTypeString,
			AllowedOperators: []string{"==", "!="},
		},
		"user.age": {
			Name:             "user.age",
			Type:             FieldTypeNumber,
			AllowedOperators: []string{"==", "!=", ">", ">=", "<", "<="},
		},
		"user.is_premium": {
			Name:             "user.is_premium",
			Type:             FieldTypeBool,
			AllowedOperators: []string{"==", "!="},
		},
		"user.registration_date": {
			Name:             "user.registration_date",
			Type:             FieldTypeDate,
			AllowedOperators: []string{"==", "!=", ">", ">=", "<", "<="},
		},
		"user.last_activity_date": {
			Name:             "user.last_activity_date",
			Type:             FieldTypeDate,
			AllowedOperators: []string{"==", "!=", ">", ">=", "<", "<="},
		},
		"user.session_count": {
			Name:             "user.session_count",
			Type:             FieldTypeNumber,
			AllowedOperators: []string{"==", "!=", ">", ">=", "<", "<="},
		},
		"user.total_purchases": {
			Name:             "user.total_purchases",
			Type:             FieldTypeNumber,
			AllowedOperators: []string{"==", "!=", ">", ">=", "<", "<="},
		},

		// Поля устройства и приложения
		"device.platform": {
			Name:             "device.platform",
			Type:             FieldTypeString,
			AllowedOperators: []string{"==", "!="},
		},
		"device.os_version": {
			Name:             "device.os_version",
			Type:             FieldTypeVersion,
			AllowedOperators: []string{"==", "!=", ">", ">=", "<", "<="},
		},
		"device.model": {
			Name:             "device.model",
			Type:             FieldTypeString,
			AllowedOperators: []string{"==", "!=", "IN", "NOT IN"},
		},
		"device.brand": {
			Name:             "device.brand",
			Type:             FieldTypeString,
			AllowedOperators: []string{"==", "!=", "IN", "NOT IN"},
		},
		"app.version": {
			Name:             "app.version",
			Type:             FieldTypeVersion,
			AllowedOperators: []string{"==", "!=", ">", ">=", "<", "<="},
		},
		"app.build": {
			Name:             "app.build",
			Type:             FieldTypeNumber,
			AllowedOperators: []string{"==", "!=", ">", ">=", "<", "<="},
		},
		"app.channel": {
			Name:             "app.channel",
			Type:             FieldTypeString,
			AllowedOperators: []string{"==", "!="},
		},

		// Временные и контекстные
		"time.hour": {
			Name:             "time.hour",
			Type:             FieldTypeNumber,
			AllowedOperators: []string{"==", "!=", ">", ">=", "<", "<="},
		},
		"time.day_of_week": {
			Name:             "time.day_of_week",
			Type:             FieldTypeNumber,
			AllowedOperators: []string{"==", "!=", "IN", "NOT IN"},
		},
		"time.is_weekend": {
			Name:             "time.is_weekend",
			Type:             FieldTypeBool,
			AllowedOperators: []string{"==", "!="},
		},
	}
}
