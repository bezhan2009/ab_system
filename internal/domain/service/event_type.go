package service

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/domain/repository"
	"ab_system/pkg/errs"
	"context"
	"errors"
)

type EventTypeService struct {
	eventTypeReader  repository.EventTypeReader
	eventTypeWriter  repository.EventTypeWriter
	eventTypeDeleter repository.EventTypeDeleter
}

func NewEventTypeService(
	eventTypeReader repository.EventTypeReader,
	eventTypeWriter repository.EventTypeWriter,
	eventTypeDeleter repository.EventTypeDeleter,
) *EventTypeService {
	return &EventTypeService{
		eventTypeReader:  eventTypeReader,
		eventTypeWriter:  eventTypeWriter,
		eventTypeDeleter: eventTypeDeleter,
	}
}

func (s *EventTypeService) GetAllEventTypes(ctx context.Context) ([]models.EventType, error) {
	types, err := s.eventTypeReader.GetAllEventTypes(ctx)
	if err != nil {
		return nil, err
	}
	if types == nil {
		return []models.EventType{}, nil
	}
	return *types, nil
}

func (s *EventTypeService) GetEventTypeByID(ctx context.Context, id string) (*models.EventType, error) {
	return s.eventTypeReader.GetEventTypeById(ctx, id)
}

func (s *EventTypeService) GetEventTypeByName(ctx context.Context, name string) (*models.EventType, error) {
	return s.eventTypeReader.GetEventTypeByName(ctx, name)
}

func (s *EventTypeService) CreateEventType(ctx context.Context, et *models.EventType) error {
	existing, err := s.eventTypeReader.GetEventTypeByName(ctx, et.Title)
	if err == nil && existing != nil {
		return errs.ErrAlreadyExists
	}

	if !errors.Is(err, errs.ErrRecordNotFound) {
		return err
	}

	return s.eventTypeWriter.CreateEventType(ctx, et)
}

func (s *EventTypeService) UpdateEventType(ctx context.Context, et *models.EventType) (*models.EventType, error) {
	existing, err := s.eventTypeReader.GetEventTypeById(ctx, et.ID.String())
	if err != nil {
		return nil, err
	}

	if et.Title != "" && et.Title != existing.Title {
		dup, err := s.eventTypeReader.GetEventTypeByName(ctx, et.Title)
		if err == nil && dup != nil && dup.ID != et.ID {
			return nil, errs.ErrAlreadyExists
		}

		if !errors.Is(err, errs.ErrRecordNotFound) {
			return nil, err
		}
	}

	return s.eventTypeWriter.UpdateEventType(ctx, et)
}

func (s *EventTypeService) DeleteEventType(ctx context.Context, id string) error {
	return s.eventTypeDeleter.DeleteEventType(ctx, id)
}
