package postgres

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/http/middlewares/observability"
	"ab_system/pkg/logger"
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ExperimentsRepository struct {
	db                  *gorm.DB
	dbExperimentsSearch *gorm.DB
}

type ExperimentVersionsRepository struct {
	db *gorm.DB
}

func NewExperimentsRepository(db *gorm.DB) *ExperimentsRepository {
	dbExperimentsSearch := db.Model(&models.Experiment{}).Where("status != 'archived'").
		Preload("RampUp").
		Preload("NotificationSettings").
		Preload("Variants")
	return &ExperimentsRepository{
		db:                  db,
		dbExperimentsSearch: dbExperimentsSearch,
	}
}

func NewExperimentVersionsRepository(db *gorm.DB) *ExperimentVersionsRepository {
	return &ExperimentVersionsRepository{
		db: db,
	}
}

func (r *ExperimentsRepository) GetAllExperiments(ctx context.Context) (experiments *[]models.Experiment, err error) {
	const op = "repository.postgres.GetAllExperiments"

	if err = r.dbExperimentsSearch.WithContext(ctx).Find(&experiments).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while fetching all experiments: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return experiments, nil
}

func (r *ExperimentsRepository) GetExperimentsByStatus(ctx context.Context, status string) (experiments *[]models.Experiment, err error) {
	const op = "repository.postgres.GetExperimentsByStatus"

	if err = r.db.WithContext(ctx).
		Preload("RampUp").
		Preload("NotificationSettings").
		Preload("Variants").
		Where("status = ?", status).
		Find(&experiments).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while fetching experiments: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return experiments, nil
}

func (r *ExperimentsRepository) GetExperimentByIdArchive(ctx context.Context, experimentID string) (experiment *models.Experiment, err error) {
	const op = "repository.postgres.GetExperimentByIdArchive"

	if err = r.db.WithContext(ctx).First(&experiment, "id = ?", experimentID).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while fetching experiment: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return experiment, nil
}

func (r *ExperimentsRepository) GetExperimentByID(ctx context.Context, experimentId string) (experiment *models.Experiment, err error) {
	const op = "repository.postgres.GetExperimentByID"

	if err = r.dbExperimentsSearch.WithContext(ctx).First(&experiment, "id = ?", experimentId).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while fetching experiment: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return experiment, nil
}

func (r *ExperimentsRepository) GetExperimentByTitle(ctx context.Context, title string) (experiment *models.Experiment, err error) {
	const op = "repository.postgres.GetExperimentByTitle"

	if err = r.dbExperimentsSearch.WithContext(ctx).First(&experiment, "title = ?", title).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while fetching experiment: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return experiment, nil
}

func (r *ExperimentsRepository) GetExperimentByTitleLike(ctx context.Context, title string) (experiment *[]models.Experiment, err error) {
	const op = "repository.postgres.GetExperimentByTitleLike"

	if err = r.dbExperimentsSearch.WithContext(ctx).Find(&experiment, "title like ?", "%"+title+"%").Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while fetching experiment: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return experiment, nil
}

func (r *ExperimentsRepository) GetExperimentByFlag(ctx context.Context, flag string) (experiments *[]models.Experiment, err error) {
	const op = "repository.postgres.GetExperimentsByFlagAndStatus"

	if err = r.dbExperimentsSearch.WithContext(ctx).Find(&experiments, "flag_key = ?", flag).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while fetching experiment: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return experiments, nil
}

func (r *ExperimentsRepository) GetExperimentByFlagAndStatus(ctx context.Context, flag string, status string) (experiments *[]models.Experiment, err error) {
	const op = "repository.postgres.GetExperimentByFlagAndStatus"

	err = r.db.WithContext(ctx).
		Preload("RampUp").
		Preload("NotificationSettings").
		Preload("Variants").
		Where("status = ?", status).
		Where("flag_key = ?", flag).
		Find(&experiments).Error

	if err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return experiments, nil
}

func (r *ExperimentsRepository) GetExperimentsWithRampEnabled(ctx context.Context) (experiments *[]models.Experiment, err error) {
	const op = "repository.postgres.GetExperimentsWithRampEnabled"

	err = r.dbExperimentsSearch.WithContext(ctx).
		Joins("JOIN experiment_ramp_ups ON experiment_ramp_ups.experiment_id = experiments.id").
		Where("experiment_ramp_ups.ramp_enabled = ?", true).
		Find(&experiments).Error

	if err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return experiments, nil
}

func (r *ExperimentsRepository) CreateExperiment(ctx context.Context, experiment *models.Experiment) (err error) {
	const op = "repository.postgres.CreateExperiment"

	if err = r.db.WithContext(ctx).Create(&experiment).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while creating experiment: %v",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}

	return nil
}

func (r *ExperimentsRepository) UpdateExperiment(ctx context.Context, experiment *models.Experiment) (updatedExperiment *models.Experiment, err error) {
	const op = "repository.postgres.UpdateExperiment"

	err = r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).
		Clauses(clause.Returning{}).
		Save(experiment).Error

	if err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while updating experiment: %v",
			op, observability.GetTraceID(ctx), err)
		return nil, TranslateGormError(err)
	}

	if err = r.db.WithContext(ctx).
		Preload("Variants").
		Preload("RampUp").
		Preload("NotificationSettings").
		Where("id = ?", experiment.ID).
		First(experiment).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error reloading experiment: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return experiment, nil
}

func (r *ExperimentsRepository) DeleteExperiment(ctx context.Context, experimentId string) (err error) {
	const op = "repository.postgres.DeleteExperiment"

	if err = r.db.WithContext(ctx).Where("id = ?", experimentId).Delete(&models.Experiment{}).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while deleting experiment: %v",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}

	return nil
}

func (r *ExperimentVersionsRepository) GetAllExperimentVersions(ctx context.Context, experimentID string) (experimentVersions *[]models.ExperimentVersion, err error) {
	const op = "repository.postgres.GetAllExperimentVersions"

	if err = r.db.WithContext(ctx).Find(&experimentVersions, "experiment_id = ?", experimentID).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while fetching experiment versions: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return experimentVersions, nil
}

func (r *ExperimentVersionsRepository) GetExperimentVersionByID(ctx context.Context, experimentVersionId string) (experimentVersion *models.ExperimentVersion, err error) {
	const op = "repository.postgres.GetExperimentVersionByID"

	if err = r.db.WithContext(ctx).First(&experimentVersion, "experiment_version_id = ?", experimentVersionId).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while fetching experiment version: %v",
			op, observability.GetTraceID(ctx), err)

		return nil, TranslateGormError(err)
	}

	return experimentVersion, nil
}

func (r *ExperimentVersionsRepository) CreateExperimentVersion(ctx context.Context, experimentVersion *models.ExperimentVersion) (err error) {
	const op = "repository.postgres.CreateExperimentVersion"

	if err = r.db.WithContext(ctx).Create(&experimentVersion).Error; err != nil {
		logger.Error.Printf("[%s] TraceId=%s Error while creating experiment version: %v",
			op, observability.GetTraceID(ctx), err)

		return TranslateGormError(err)
	}

	return nil
}
