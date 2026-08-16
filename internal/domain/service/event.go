package service

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/domain/repository"
	"ab_system/internal/http/dto"
	"ab_system/pkg/errs"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type EventService struct {
	eventTypeReader repository.EventTypeReader
	eventReader     repository.EventReader
	eventWriter     repository.EventWriter
	decisionReader  repository.DecisionReader
}

func NewEventService(
	eventTypeReader repository.EventTypeReader,
	eventReader repository.EventReader,
	eventWriter repository.EventWriter,
	decisionReader repository.DecisionReader,
) *EventService {
	return &EventService{
		eventTypeReader: eventTypeReader,
		eventReader:     eventReader,
		eventWriter:     eventWriter,
		decisionReader:  decisionReader,
	}
}

func (s *EventService) ProcessEvents(ctx context.Context, req dto.EventBatchRequest) (*dto.EventResponse, error) {
	resp := &dto.EventResponse{
		Accepted:  0,
		Duplicate: 0,
		Rejected:  0,
		Errors:    []errs.FieldError{},
	}

	for i, er := range req {
		err := s.ProcessEvent(ctx, &er)
		if err != nil {
			var multiErr *errs.MultiValidationError
			if errors.As(err, &multiErr) {
				for _, fe := range multiErr.FieldErrors {
					resp.Errors = append(resp.Errors, errs.FieldError{
						Field: fmt.Sprintf("[%d].%s", i, fe.Field),
						Issue: fe.Issue,
					})
				}
				resp.Rejected++
				continue
			}

			var valErr *errs.ValidationError
			if errors.As(err, &valErr) {
				resp.Errors = append(resp.Errors, errs.FieldError{
					Field: fmt.Sprintf("[%d].%s", i, valErr.Field),
					Issue: valErr.Issue,
				})
				resp.Rejected++
				continue
			}

			if errors.Is(err, errs.ErrDuplicateEvent) {
				resp.Duplicate++
				continue
			}

			resp.Errors = append(resp.Errors, errs.FieldError{
				Field: fmt.Sprintf("[%d]", i),
				Issue: err.Error(),
			})
			resp.Rejected++
		} else {
			resp.Accepted++
		}
	}

	return resp, nil
}

func (s *EventService) ProcessEvent(ctx context.Context, er *dto.EventRequest) error {
	var fieldErrors []errs.FieldError

	eventType, err := s.eventTypeReader.GetEventTypeByName(ctx, er.EventType)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			fieldErrors = append(fieldErrors, errs.FieldError{
				Field: "event_type",
				Issue: fmt.Sprintf("неизвестный тип события: %s", er.EventType),
			})
		} else {
			fieldErrors = append(fieldErrors, errs.FieldError{
				Field: "event_type",
				Issue: err.Error(),
			})
		}
		return errs.NewMultiValidationError(fieldErrors)
	}

	if eventType.RequiresDecisionID && er.DecisionID == "" {
		fieldErrors = append(fieldErrors, errs.FieldError{
			Field: "decision_id",
			Issue: fmt.Sprintf("для типа события %s требуется идентификатор решения (decision_id)", er.EventType),
		})
	}

	if eventType.RequiresUserID && er.UserID == "" {
		fieldErrors = append(fieldErrors, errs.FieldError{
			Field: "user_id",
			Issue: fmt.Sprintf("для типа события %s требуется идентификатор пользователя (user_id)", er.EventType),
		})
	}

	if er.EventID == "" {
		fieldErrors = append(fieldErrors, errs.FieldError{
			Field: "event_id",
			Issue: "требуется event_id",
		})
	}

	if len(eventType.Schema) > 0 {
		var schema map[string]interface{}
		if err := json.Unmarshal(eventType.Schema, &schema); err != nil {
			fieldErrors = append(fieldErrors, errs.FieldError{
				Field: "payload",
				Issue: fmt.Sprintf("недопустимая схема для типа события %s", er.EventType),
			})
		} else {
			s.validatePayload(er.Payload, schema, &fieldErrors)
		}
	}

	if len(fieldErrors) > 0 {
		return errs.NewMultiValidationError(fieldErrors)
	}

	_, err = s.eventReader.GetEventById(ctx, er.EventID)
	if err != nil {
		if !errors.Is(err, errs.ErrRecordNotFound) {
			return err
		}
	} else {
		return errs.ErrDuplicateEvent
	}

	event := &models.Event{
		EventID:     er.EventID,
		EventTypeID: eventType.ID.String(),
		DecisionID:  er.DecisionID,
		UserID:      er.UserID,
		ReceivedAt:  time.Now(),
		ClientTime:  er.ClientTime,
	}

	if er.Payload != nil {
		payloadBytes, err := json.Marshal(er.Payload)
		if err != nil {
			return errs.NewValidationError("payload", fmt.Sprintf("не удалось преобразовать payload: %v", err))
		}

		event.Payload = payloadBytes
	}

	if err = s.eventWriter.CreateEvent(ctx, event); err != nil {
		if errors.Is(err, errs.ErrDuplicateEntry) {
			return errs.ErrDuplicateEvent
		}
		return err
	}

	return nil
}

func (s *EventService) validatePayload(payload map[string]interface{}, schema map[string]interface{}, fieldErrors *[]errs.FieldError) {
	for field, expectedType := range schema {
		value, exists := payload[field]
		if !exists {
			*fieldErrors = append(*fieldErrors, errs.FieldError{
				Field: "payload." + field,
				Issue: fmt.Sprintf("отсутствует обязательное поле: %s", field),
			})
			continue
		}

		switch expectedType {
		case "string":
			if _, ok := value.(string); !ok {
				*fieldErrors = append(*fieldErrors, errs.FieldError{
					Field: "payload." + field,
					Issue: fmt.Sprintf("поле %s должно быть строкой", field),
				})
			}
		case "number":
			if _, ok := value.(float64); !ok {
				*fieldErrors = append(*fieldErrors, errs.FieldError{
					Field: "payload." + field,
					Issue: fmt.Sprintf("поле %s должно быть числом", field),
				})
			}
		case "boolean":
			if _, ok := value.(bool); !ok {
				*fieldErrors = append(*fieldErrors, errs.FieldError{
					Field: "payload." + field,
					Issue: fmt.Sprintf("поле %s должно быть булевым", field),
				})
			}
		default:
			*fieldErrors = append(*fieldErrors, errs.FieldError{
				Field: "payload." + field,
				Issue: fmt.Sprintf("не поддерживается тип %v в схеме", expectedType),
			})
		}
	}
}
