package service

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/domain/repository"
	"ab_system/pkg/errs"
	"context"
	"errors"
)

type NotificationSettingsService struct {
	reader    repository.NotificationSettingsReader
	writer    repository.NotificationSettingsWriter
	deleter   repository.NotificationSettingsDeleter
	expReader repository.ExperimentReader
}

func NewNotificationSettingsService(
	reader repository.NotificationSettingsReader,
	writer repository.NotificationSettingsWriter,
	deleter repository.NotificationSettingsDeleter,
	expReader repository.ExperimentReader,
) *NotificationSettingsService {
	return &NotificationSettingsService{
		reader:    reader,
		writer:    writer,
		deleter:   deleter,
		expReader: expReader,
	}
}

func (s *NotificationSettingsService) GetNotificationSettingsByExperimentID(ctx context.Context, experimentID string) (*models.NotificationSettings, error) {
	_, err := s.expReader.GetExperimentByID(ctx, experimentID)
	if err != nil {
		return nil, err
	}

	return s.reader.GetNotificationSettingsByExperimentID(ctx, experimentID)
}

func (s *NotificationSettingsService) CreateNotificationSettings(ctx context.Context, settings *models.NotificationSettings) (*models.NotificationSettings, error) {
	_, err := s.expReader.GetExperimentByID(ctx, settings.ExperimentID.String())
	if err != nil {
		return nil, err
	}

	_, err = s.reader.GetNotificationSettingsByExperimentID(ctx, settings.ExperimentID.String())
	if err != nil && !errors.Is(err, errs.ErrRecordNotFound) {
		return nil, err
	}
	if err == nil {
		return nil, errs.ErrAlreadyExists
	}

	if err = s.writer.CreateNotificationSettings(ctx, settings); err != nil {
		return nil, err
	}

	return settings, nil
}

func (s *NotificationSettingsService) UpdateNotificationSettings(ctx context.Context, settings *models.NotificationSettings) (*models.NotificationSettings, error) {
	_, err := s.expReader.GetExperimentByID(ctx, settings.ExperimentID.String())
	if err != nil {
		return nil, err
	}

	_, err = s.reader.GetNotificationSettingsByExperimentID(ctx, settings.ExperimentID.String())
	if err != nil {
		return nil, err
	}

	updated, err := s.writer.UpdateNotificationSettings(ctx, settings)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *NotificationSettingsService) DeleteNotificationSettingsByExperimentID(ctx context.Context, experimentId string) error {
	settings, err := s.GetNotificationSettingsByExperimentID(ctx, experimentId)
	if err != nil {
		return err
	}

	return s.deleter.DeleteNotificationSettings(ctx, settings.ID.String())
}
