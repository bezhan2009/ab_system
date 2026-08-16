package validation

import (
	"ab_system/internal/http/dto"
	"ab_system/pkg/errs"
	"encoding/json"
	"fmt"
)

func ValidateEventTypeCreate(et dto.EventType) []errs.FieldError {
	var errors []errs.FieldError

	if et.Title == "" {
		errors = append(errors, errs.NewFieldError("title", "обязательное поле"))
	} else if len(et.Title) > 255 {
		errors = append(errors, errs.NewFieldError("title", "не более 255 символов"))
	}

	if len(et.Description) > 1000 {
		errors = append(errors, errs.NewFieldError("description", "не более 1000 символов"))
	}

	if et.Schema == nil {
		errors = append(errors, errs.NewFieldError("schema", "обязательное поле"))
	} else {
		schemaBytes, err := json.Marshal(et.Schema)
		if err != nil {
			errors = append(errors, errs.NewFieldError("schema", "неверный формат JSON"))
		} else {
			var schema map[string]interface{}
			if err = json.Unmarshal(schemaBytes, &schema); err != nil {
				errors = append(errors, errs.NewFieldError("schema", "должен быть JSON объект"))
			} else {
				for field, fieldType := range schema {
					typeStr, ok := fieldType.(string)
					if !ok {
						errors = append(errors, errs.NewFieldError(
							fmt.Sprintf("schema.%s", field),
							"тип поля должен быть строкой",
						))
						continue
					}

					switch typeStr {
					case "string", "number", "boolean":
					default:
						errors = append(errors, errs.NewFieldError(
							fmt.Sprintf("schema.%s", field),
							fmt.Sprintf("недопустимый тип '%s' (допустимы: string, number, boolean)", typeStr),
						))
					}
				}
			}
		}
	}

	return errors
}

func ValidateEventTypeUpdate(et dto.EventType) []errs.FieldError {
	var errors []errs.FieldError

	if et.Title != "" && len(et.Title) > 255 {
		errors = append(errors, errs.NewFieldError("title", "не более 255 символов"))
	}

	if len(et.Description) > 1000 {
		errors = append(errors, errs.NewFieldError("description", "не более 1000 символов"))
	}

	if et.Schema != nil {
		schemaBytes, err := json.Marshal(et.Schema)
		if err != nil {
			errors = append(errors, errs.NewFieldError("schema", "неверный формат JSON"))
		} else {
			var schema map[string]interface{}
			if err := json.Unmarshal(schemaBytes, &schema); err != nil {
				errors = append(errors, errs.NewFieldError("schema", "должен быть JSON объект"))
			} else {
				for field, fieldType := range schema {
					typeStr, ok := fieldType.(string)
					if !ok {
						errors = append(errors, errs.NewFieldError(
							fmt.Sprintf("schema.%s", field),
							"тип поля должен быть строкой",
						))
						continue
					}

					switch typeStr {
					case "string", "number", "boolean":
					default:
						errors = append(errors, errs.NewFieldError(
							fmt.Sprintf("schema.%s", field),
							fmt.Sprintf("недопустимый тип '%s' (допустимы: string, number, boolean)", typeStr),
						))
					}
				}
			}
		}
	}

	return errors
}
