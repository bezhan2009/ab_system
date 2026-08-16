package validation

import (
	"ab_system/internal/http/dto"
	"ab_system/pkg/errs"
	"fmt"
)

func ValidateEventCreate(event dto.EventRequest) []errs.FieldError {
	var errors []errs.FieldError

	if event.EventID == "" {
		errors = append(errors, errs.NewFieldError("event_id", "обязательное поле"))
	} else if len(event.EventID) > 255 {
		errors = append(errors, errs.NewFieldError("event_id", "не более 255 символов"))
	}

	if event.EventType == "" {
		errors = append(errors, errs.NewFieldError("event_type", "обязательное поле"))
	} else if len(event.EventType) > 255 {
		errors = append(errors, errs.NewFieldError("event_type", "не более 255 символов"))
	}

	if event.DecisionID == "" {
		errors = append(errors, errs.NewFieldError("decision_id", "обязательное поле"))
	} else if len(event.DecisionID) > 255 {
		errors = append(errors, errs.NewFieldError("decision_id", "не более 255 символов"))
	}

	if event.UserID == "" {
		errors = append(errors, errs.NewFieldError("user_id", "обязательное поле"))
	} else if len(event.UserID) > 255 {
		errors = append(errors, errs.NewFieldError("user_id", "не более 255 символов"))
	}

	if event.Payload == nil {
		errors = append(errors, errs.NewFieldError("payload", "обязательное поле"))
	} else if len(event.Payload) == 0 {
		errors = append(errors, errs.NewFieldError("payload", "не может быть пустым"))
	}

	return errors
}

func ValidateEventBatch(events dto.EventBatchRequest) []errs.FieldError {
	var errors []errs.FieldError

	if len(events) == 0 {
		errors = append(errors, errs.NewFieldError("events", "должен быть хотя бы один элемент"))
		return errors
	}

	for i, event := range events {
		eventErrors := ValidateEventCreate(event)
		for _, err := range eventErrors {
			errors = append(errors, errs.FieldError{
				Field: fmt.Sprintf("[%d].%s", i, err.Field),
				Issue: err.Issue,
			})
		}
	}

	return errors
}
