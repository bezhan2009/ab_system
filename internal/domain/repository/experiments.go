package repository

import (
	"ab_system/internal/domain/models"
	"context"
)

type ExperimentReader interface {
	GetAllExperiments(ctx context.Context) (experiments *[]models.Experiment, err error)
	GetExperimentsByStatus(ctx context.Context, status string) (experiments *[]models.Experiment, err error)
	GetExperimentByID(ctx context.Context, experimentId string) (experiment *models.Experiment, err error)
	GetExperimentByIdArchive(ctx context.Context, experimentID string) (experiment *models.Experiment, err error)
	GetExperimentByTitle(ctx context.Context, title string) (experiment *models.Experiment, err error)
	GetExperimentByTitleLike(ctx context.Context, title string) (experiment *[]models.Experiment, err error)
	GetExperimentByFlag(ctx context.Context, flag string) (experiments *[]models.Experiment, err error)
	GetExperimentByFlagAndStatus(ctx context.Context, flag string, status string) (experiment *[]models.Experiment, err error)
	GetExperimentsWithRampEnabled(ctx context.Context) (experiments *[]models.Experiment, err error)
}

type ExperimentWriter interface {
	CreateExperiment(ctx context.Context, experiment *models.Experiment) (err error)
	UpdateExperiment(ctx context.Context, experiment *models.Experiment) (updatedExperiment *models.Experiment, err error)
}

type ExperimentDeleter interface {
	DeleteExperiment(ctx context.Context, experimentId string) error
}

type ExperimentVersionReader interface {
	GetAllExperimentVersions(ctx context.Context, experimentID string) (experimentVersions *[]models.ExperimentVersion, err error)
	GetExperimentVersionByID(ctx context.Context, experimentVersionId string) (experimentVersion *models.ExperimentVersion, err error)
}

type ExperimentVersionWriter interface {
	CreateExperimentVersion(ctx context.Context, experimentVersion *models.ExperimentVersion) (err error)
}
