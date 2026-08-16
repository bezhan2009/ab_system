package repository

import (
	"ab_system/internal/domain/models"
	"context"
)

type EventTypeReader interface {
	GetAllEventTypes(ctx context.Context) (eventTypes *[]models.EventType, err error)
	GetEventTypeById(ctx context.Context, eventTypeId string) (eventType *models.EventType, err error)
	GetEventTypeByName(ctx context.Context, name string) (eventType *models.EventType, err error)
}

type EventTypeWriter interface {
	CreateEventType(ctx context.Context, eventType *models.EventType) (err error)
	UpdateEventType(ctx context.Context, eventType *models.EventType) (updatedEventType *models.EventType, err error)
}

type EventTypeDeleter interface {
	DeleteEventType(ctx context.Context, eventTypeId string) (err error)
}
