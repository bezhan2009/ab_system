package postgres

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/http/middlewares/observability"
	"ab_system/pkg/logger"
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EventTypeRepository struct {
	db *gorm.DB
}

func NewEventTypeRepository(db *gorm.DB) *EventTypeRepository {
	return &EventTypeRepository{
		db: db,
	}
}

func (r *EventTypeRepository) GetAllEventTypes(ctx context.Context) (eventTypes *[]models.EventType, err error) {
	const op = "repository.postgres.GetAllEventTypes"

	if err = r.db.WithContext(ctx).Find(&eventTypes).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while getting all event types: %s",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return eventTypes, nil
}

func (r *EventTypeRepository) GetEventTypeById(ctx context.Context, eventTypeId string) (eventType *models.EventType, err error) {
	const op = "repository.postgres.GetEventTypeById"

	if err = r.db.WithContext(ctx).First(&eventType, "id = ?", eventTypeId).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while getting event type by id: %s",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return eventType, nil
}

func (r *EventTypeRepository) GetEventTypeByName(ctx context.Context, title string) (eventType *models.EventType, err error) {
	const op = "repository.postgres.GetEventTypeByName"

	if err = r.db.WithContext(ctx).First(&eventType, "title = ?", title).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while getting event type by title: %s",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return eventType, nil
}

func (r *EventTypeRepository) CreateEventType(ctx context.Context, eventType *models.EventType) (err error) {
	const op = "repository.postgres.CreateEventType"

	if err = r.db.WithContext(ctx).Create(&eventType).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while creating event type: %s",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}

	return nil
}

func (r *EventTypeRepository) UpdateEventType(ctx context.Context, eventType *models.EventType) (updatedEventType *models.EventType, err error) {
	const op = "repository.postgres.UpdateEventType"

	if err = r.db.WithContext(ctx).Clauses(clause.Returning{}).Updates(eventType).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while updating event type: %s",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return eventType, nil
}

func (r *EventTypeRepository) DeleteEventType(ctx context.Context, eventTypeId string) (err error) {
	const op = "repository.postgres.DeleteEventType"

	if err = r.db.WithContext(ctx).Where("id = ?", eventTypeId).Delete(&models.EventType{}).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while deleting event type by id: %s",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}

	return nil
}
