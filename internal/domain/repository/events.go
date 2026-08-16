package repository

import (
	"ab_system/internal/domain/models"
	"context"
	"time"
)

type EventReader interface {
	GetAllEventsByExperimentId(ctx context.Context, experimentId string) (events *[]models.Event, err error)
	GetEventById(ctx context.Context, id string) (event *models.Event, err error)
	CountEventsByTypeAndDecisions(ctx context.Context, eventType string, decisionIDs []string, from, to time.Time, useClientTime bool) (int64, error)
	AvgFieldByTypeAndDecisions(ctx context.Context, eventType, field string, decisionIDs []string, from, to time.Time, useClientTime bool) (float64, error)
}

type EventWriter interface {
	CreateEvent(ctx context.Context, event *models.Event) (err error)
}

type EventDeleter interface {
	DeleteEvent(ctx context.Context, eventId string) (err error)
}
