package repository

import (
	"ab_system/internal/domain/models"
	"context"
	"time"
)

type DecisionReader interface {
	GetAllExperimentDecisions(ctx context.Context, experimentId string) (decisions *models.Decision, err error)
	GetDecisionByUserId(ctx context.Context, userId string) (decision *models.Decision, err error)
	GetDecisionsByExperimentAndTime(ctx context.Context, experimentID string, from, to time.Time) ([]models.Decision, error)
}

type DecisionWriter interface {
	CreateDecision(ctx context.Context, decision *models.Decision) (err error)
	UpdateDecision(ctx context.Context, decision *models.Decision) (updatedDecision *models.Decision, err error)
}

type DecisionDeleter interface {
	DeleteDecision(ctx context.Context, decisionId string) (err error)
}
