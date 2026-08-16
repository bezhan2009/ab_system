package postgres

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/http/middlewares/observability"
	"ab_system/pkg/logger"
	"context"
	"time"

	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (r *EventRepository) GetAllEventsByExperimentId(
	ctx context.Context,
	experimentId string,
) (events *[]models.Event, err error) {

	const op = "repository.postgres.GetAllEventsByExperimentId"

	err = r.db.WithContext(ctx).
		Preload("EventType").
		Joins("JOIN decisions ON events.decision_id = decisions.id").
		Where("decisions.experiment_id = ?", experimentId).
		Order("events.received_at DESC").
		Find(&events).Error

	if err != nil {
		logger.Error.Printf("[%s] TraceId=%s ExperimentId=%s Error: %v",
			op, observability.GetTraceID(ctx), experimentId, err)

		return nil, TranslateGormError(err)
	}

	return events, nil
}

func (r *EventRepository) GetEventById(ctx context.Context, eventId string) (event *models.Event, err error) {
	const op = "repository.postgres.GetEventById"

	if err = r.db.WithContext(ctx).First(&event, "id = ?", eventId).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return event, nil
}

func (r *EventRepository) CountEventsByTypeAndDecisions(
	ctx context.Context,
	eventType string,
	decisionIDs []string,
	from, to time.Time,
	useClientTime bool,
) (count int64, err error) {
	const op = "repository.postgres.CountEventsByTypeAndDecisions"

	var eventTypeModel *models.EventType
	if err = r.db.WithContext(ctx).Where("title = ?", eventType).First(&eventTypeModel).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s EventType=%s Error: %v",
			op, observability.GetTraceID(ctx), eventType, err)
		
		return 0, TranslateGormError(err)
	}

	if len(decisionIDs) == 0 {
		return 0, nil
	}

	query := r.db.WithContext(ctx).
		Model(&models.Event{}).
		Where("event_type_id = ?", eventTypeModel.ID).
		Where("decision_id IN ?", decisionIDs)

	if useClientTime {
		query = query.Where("COALESCE(client_time, received_at) >= ? AND COALESCE(client_time, received_at) < ?", from, to)
	} else {
		query = query.Where("received_at >= ? AND received_at < ?", from, to)
	}

	err = query.Count(&count).Error
	if err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v", op, observability.GetTraceID(ctx), err)
		return 0, TranslateGormError(err)
	}

	return count, nil
}

func (r *EventRepository) AvgFieldByTypeAndDecisions(
	ctx context.Context,
	eventType, field string,
	decisionIDs []string,
	from, to time.Time,
	useClientTime bool,
) (avg float64, err error) {
	const op = "repository.postgres.AvgFieldByTypeAndDecisions"

	if len(decisionIDs) == 0 {
		return 0, nil
	}

	var result float64
	var query string
	if useClientTime {
		query = `
            SELECT COALESCE(AVG((payload->>?)::numeric), 0)
            FROM events
            WHERE event_type = ?
              AND decision_id IN ?
              AND COALESCE(client_time, received_at) >= ? 
              AND COALESCE(client_time, received_at) < ?
        `
	} else {
		query = `
            SELECT COALESCE(AVG((payload->>?)::numeric), 0)
            FROM events
            WHERE event_type = ?
              AND decision_id IN ?
              AND received_at >= ? 
              AND received_at < ?
        `
	}

	err = r.db.WithContext(ctx).Raw(query, field, eventType, decisionIDs, from, to).Scan(&result).Error
	if err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return 0, TranslateGormError(err)
	}

	return result, nil
}

func (r *EventRepository) CreateEvent(ctx context.Context, event *models.Event) (err error) {
	const op = "repository.postgres.CreateEvent"

	if err = r.db.WithContext(ctx).Create(event).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while creating event: %s",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}

	return nil
}

func (r *EventRepository) DeleteEvent(ctx context.Context, eventTypeId string) (err error) {
	const op = "repository.postgres.DeleteEvent"

	if err = r.db.WithContext(ctx).Where("id = ?", eventTypeId).Delete(&models.Event{}).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while deleting event: %s",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}

	return nil
}
