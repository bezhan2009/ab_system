package service

import (
	"ab_system/internal/clients/telegram_notifications"
	"ab_system/internal/domain/models"
	"ab_system/internal/domain/repository"
	"ab_system/internal/http/dto"
	"ab_system/internal/notifications/slack"
	"ab_system/internal/validation"
	"ab_system/pkg/errs"
	"ab_system/pkg/logger"
	pkgDto "ab_system/pkg/notifications/dto"
	"ab_system/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ExperimentService struct {
	experimentReader  repository.ExperimentReader
	experimentWriter  repository.ExperimentWriter
	experimentDeleter repository.ExperimentDeleter

	experimentVersionReader repository.ExperimentVersionReader
	experimentVersionWriter repository.ExperimentVersionWriter

	notificationSettingsReader repository.NotificationSettingsReader
	notificationSettingsWriter repository.NotificationSettingsWriter

	userReader repository.UserReader

	telegramClient *telegram_notifications.NotifyClient
	slackNotifier  *slack.Notifier
}

func NewExperimentService(
	experimentReader repository.ExperimentReader,
	experimentWriter repository.ExperimentWriter,
	experimentDeleter repository.ExperimentDeleter,
	experimentVersionReader repository.ExperimentVersionReader,
	experimentVersionWriter repository.ExperimentVersionWriter,
	notificationSettingsReader repository.NotificationSettingsReader,
	notificationSettingsWriter repository.NotificationSettingsWriter,
	userReader repository.UserReader,
	telegramClient *telegram_notifications.NotifyClient,
	slackNotifier *slack.Notifier,
) *ExperimentService {
	return &ExperimentService{
		experimentReader:           experimentReader,
		experimentWriter:           experimentWriter,
		experimentDeleter:          experimentDeleter,
		experimentVersionReader:    experimentVersionReader,
		experimentVersionWriter:    experimentVersionWriter,
		notificationSettingsReader: notificationSettingsReader,
		notificationSettingsWriter: notificationSettingsWriter,
		userReader:                 userReader,
		telegramClient:             telegramClient,
		slackNotifier:              slackNotifier,
	}
}

func (s *ExperimentService) GetAllExperiments(ctx context.Context) (experiments *[]models.Experiment, err error) {
	experiments, err = s.experimentReader.GetAllExperiments(ctx)
	if err != nil {
		return nil, err
	}

	return experiments, nil
}

func (s *ExperimentService) GetExperimentsByStatus(ctx context.Context, status string) (experiments *[]models.Experiment, err error) {
	experiments, err = s.experimentReader.GetExperimentsByStatus(ctx, status)
	if err != nil {
		return nil, err
	}

	return experiments, nil
}

func (s *ExperimentService) GetExperimentByID(ctx context.Context, experimentID string) (experiment *models.Experiment, err error) {
	experiment, err = s.experimentReader.GetExperimentByID(ctx, experimentID)
	if err != nil {
		return nil, err
	}

	return experiment, nil
}

func (s *ExperimentService) GetExperimentByTitle(ctx context.Context, title string) (experiment *models.Experiment, err error) {
	experiment, err = s.experimentReader.GetExperimentByTitle(ctx, title)
	if err != nil {
		return nil, err
	}

	return experiment, nil
}

func (s *ExperimentService) GetExperimentByTitleLike(ctx context.Context, title string) (experiment *[]models.Experiment, err error) {
	experiment, err = s.experimentReader.GetExperimentByTitleLike(ctx, title)
	if err != nil {
		return nil, err
	}

	return experiment, nil
}

func (s *ExperimentService) GetExperimentsByFlag(ctx context.Context, flag string) (experiments *[]models.Experiment, err error) {
	experiments, err = s.experimentReader.GetExperimentByFlag(ctx, flag)
	if err != nil {
		return nil, err
	}

	return experiments, nil
}

func (s *ExperimentService) GetExperimentsByFlagAndStatus(ctx context.Context, flag, status string) (experiment *[]models.Experiment, err error) {
	experiment, err = s.experimentReader.GetExperimentByFlagAndStatus(ctx, flag, status)
	if err != nil {
		return nil, err
	}

	return experiment, nil
}

func (s *ExperimentService) GetExperimentVersionsByID(ctx context.Context, experimentID string) (experimentVersions *[]models.ExperimentVersion, err error) {
	experimentVersions, err = s.experimentVersionReader.GetAllExperimentVersions(ctx, experimentID)
	if err != nil {
		return nil, err
	}

	return experimentVersions, nil
}

func (s *ExperimentService) CreateExperiment(ctx context.Context, experiment *models.Experiment) (err error) {
	existing, err := s.experimentReader.GetExperimentByTitle(ctx, experiment.Title)
	if err == nil && existing != nil {
		return errs.ErrAlreadyExists
	}
	if !errors.Is(err, errs.ErrRecordNotFound) {
		return err
	}

	if err = validation.ValidateVariants(experiment); err != nil {
		return err
	}

	err = s.experimentWriter.CreateExperiment(ctx, experiment)
	if err != nil {
		return err
	}

	return nil
}

func (s *ExperimentService) SendToReview(ctx context.Context, experimentId string) (err error) {
	experiment, err := s.experimentReader.GetExperimentByID(ctx, experimentId)
	if err != nil {
		return err
	}

	if experiment.Status == models.StatusRunning || experiment.Status == models.StatusPaused {
		return errs.ErrExperimentAlreadyHasBeenOnReview
	}

	experiment.Status = models.StatusInReview
	_, err = s.UpdateExperiment(ctx, experiment)
	if err != nil {
		return err
	}

	return nil
}

func (s *ExperimentService) RunExperiment(ctx context.Context, experimentId string) (err error) {
	experiment, err := s.experimentReader.GetExperimentByID(ctx, experimentId)
	if err != nil {
		return err
	}

	experiment.Status = models.StatusRunning

	_, err = s.UpdateExperiment(ctx, experiment)
	if err != nil {
		return err
	}

	return nil
}

func (s *ExperimentService) CompleteExperiment(ctx context.Context, completeExperimentRequest dto.CompleteExperimentRequest, experimentId string) (err error) {
	experiment, err := s.experimentReader.GetExperimentByID(ctx, experimentId)
	if err != nil {
		return err
	}

	if completeExperimentRequest.Conclusion == "rollout" {
		experiment.WinnerVariantID = ""

		for _, v := range experiment.Variants {
			if v.ID.String() == completeExperimentRequest.WinnerVariantID {
				experiment.WinnerVariantID = completeExperimentRequest.WinnerVariantID
			}
		}

		if experiment.WinnerVariantID == "" {
			return errs.ErrInvalidWinnerVariant
		}
	}

	experiment.Status = models.StatusCompleted
	experiment.Comment = completeExperimentRequest.Comment
	experiment.Conclusion = completeExperimentRequest.Conclusion

	_, err = s.UpdateExperiment(ctx, experiment)
	if err != nil {
		return err
	}

	return nil
}

func (s *ExperimentService) checkNoActiveExperiment(ctx context.Context, flagKey, excludeID string) (err error) {
	activeExps, err := s.experimentReader.GetExperimentByFlagAndStatus(ctx, flagKey, string(models.StatusRunning))
	if err != nil && !errors.Is(err, errs.ErrRecordNotFound) {
		return err
	}
	if activeExps == nil || len(*activeExps) == 0 {
		return nil
	}

	exp := (*activeExps)[0]

	if exp.ID.String() != excludeID {
		return errs.ErrActiveExperimentExists
	}

	return nil
}

func (s *ExperimentService) UpdateExperiment(ctx context.Context, experiment *models.Experiment) (*models.Experiment, error) {
	existing, err := s.experimentReader.GetExperimentByID(ctx, experiment.ID.String())
	if err != nil {
		return nil, err
	}

	if existing.Status == models.StatusRunning || existing.Status == models.StatusPaused {
		if experiment.Title != existing.Title ||
			experiment.FlagKey != existing.FlagKey ||
			experiment.TrafficPercent != existing.TrafficPercent ||
			!validation.VariantsEqual(experiment.Variants, existing.Variants) {

			return nil, errs.ErrCannotEditRunningExperiment
		}
	}

	if existing.Status != experiment.Status {
		if err := validation.ValidateStatusTransition(existing.Status, experiment.Status); err != nil {
			return nil, err
		}

		switch experiment.Status {
		case models.StatusRunning:
			if existing.Status != models.StatusApproved && existing.Status != models.StatusPaused {
				return nil, errs.ErrExperimentNotApproved
			}
			if err := s.checkNoActiveExperiment(ctx, experiment.FlagKey, experiment.ID.String()); err != nil {
				return nil, err
			}
		case models.StatusCompleted:
			if experiment.Conclusion == "" {
				return nil, errs.ErrConclusionRequired
			}
		}
	}

	if experiment.Title != "" && experiment.Title != existing.Title {
		dup, err := s.experimentReader.GetExperimentByTitle(ctx, experiment.Title)
		if err == nil && dup != nil && dup.ID != experiment.ID {
			return nil, errs.ErrAlreadyExists
		}
		if !errors.Is(err, errs.ErrRecordNotFound) {
			return nil, err
		}
	}

	if len(experiment.Variants) > 0 {
		existingMap := make(map[uuid.UUID]models.Variant)
		for _, v := range existing.Variants {
			existingMap[v.ID] = v
		}

		incomingIDs := make(map[uuid.UUID]bool)
		for i, v := range experiment.Variants {
			if v.ID != uuid.Nil {
				if _, ok := existingMap[v.ID]; !ok {
					return nil, errs.ErrVariantNotFound
				}
				incomingIDs[v.ID] = true
			} else {
				experiment.Variants[i].ID = uuid.Nil
			}
		}

		for _, v := range existing.Variants {
			if !incomingIDs[v.ID] {
				experiment.Variants = append(experiment.Variants, v)
			}
		}

		if err = validation.ValidateVariants(experiment); err != nil {
			return nil, err
		}
	} else {
		experiment.Variants = existing.Variants
	}

	now := time.Now()
	if experiment.Status == models.StatusRunning && existing.Status != models.StatusRunning {
		experiment.StartedAt = &now
	}
	if experiment.Status == models.StatusCompleted && existing.Status != models.StatusCompleted {
		experiment.CompletedAt = &now
	}

	notifSettings := experiment.NotificationSettings
	experiment.NotificationSettings = models.NotificationSettings{}

	configChanged := validation.IsConfigChanged(existing, experiment)
	if configChanged {
		experiment.Version = existing.Version + 1
	} else {
		experiment.Version = existing.Version
	}

	updated, err := s.experimentWriter.UpdateExperiment(ctx, experiment)
	if err != nil {
		return nil, err
	}

	experiment.NotificationSettings = notifSettings

	if len(notifSettings.ChatIds) > 0 || notifSettings.ID != uuid.Nil {
		existingSettings, err := s.notificationSettingsReader.GetNotificationSettingsByExperimentID(ctx, updated.ID.String())
		if err != nil && !errors.Is(err, errs.ErrRecordNotFound) {
			return nil, err
		}

		if existingSettings == nil {
			newSettings := &models.NotificationSettings{
				ExperimentID: &updated.ID,
				ChatIds:      notifSettings.ChatIds,
			}

			if err = s.notificationSettingsWriter.CreateNotificationSettings(ctx, newSettings); err != nil {
				return nil, err
			}
		} else {
			existingSettings.ChatIds = notifSettings.ChatIds
			if _, err = s.notificationSettingsWriter.UpdateNotificationSettings(ctx, existingSettings); err != nil {
				return nil, err
			}
		}
	}

	if configChanged {
		expDto := dto.Experiment{}
		snapshotBytes, err := json.Marshal(expDto.ToDTO(updated))
		if err != nil {
			return nil, err
		}

		version := &models.ExperimentVersion{
			ExperimentID: updated.ID,
			Version:      updated.Version,
			Snapshot:     datatypes.JSON(snapshotBytes),
			ChangedBy:    updated.OwnerID,
			CreatedAt:    time.Now(),
		}

		if err = s.experimentVersionWriter.CreateExperimentVersion(ctx, version); err != nil {
			return nil, err
		}
	}

	updated, err = s.experimentReader.GetExperimentByID(ctx, updated.ID.String())
	if err != nil {
		return nil, err
	}

	ownerEmail := ""
	owner, err := s.userReader.GetUserByID(ctx, updated.OwnerID)
	if err != nil {
		ownerEmail = updated.OwnerID
	} else {
		ownerEmail = fmt.Sprintf("%s (%s)", owner.Email, updated.OwnerID)
	}

	go func() {
		bgCtx := context.Background()

		_ = SendTelegramNotificationByExperimentStatus(bgCtx, s.telegramClient.Notify, updated, ownerEmail)
	}()

	go func() {
		bgCtx := context.Background()

		webhooks, err := s.notificationSettingsReader.GetSlackWebhooksForExpNotification(bgCtx, updated.ID.String())
		if err != nil {
			logger.Error.Printf("Failed to get slack webhooks for experiment %s: %v", updated.ID, err)
			return
		}

		if len(webhooks) > 0 {
			slackReq := pkgDto.NotifyRequest{
				EventType:    utils.GetEventTypeByExperimentStatus(string(updated.Status)),
				ExperimentID: updated.ID.String(),
				Experiment:   updated.Title,
				FlagKey:      updated.FlagKey,
				UserID:       ownerEmail,
				Status:       string(updated.Status),
			}

			if err = s.slackNotifier.Send(bgCtx, slackReq, webhooks); err != nil {
				logger.Error.Printf("Failed to send slack notification for experiment %s: %v", updated.ID, err)
			}
		}
	}()

	return updated, nil
}

func (s *ExperimentService) ArchiveExperiment(ctx context.Context, experimentID string) (err error) {
	exp, err := s.experimentReader.GetExperimentByIdArchive(ctx, experimentID)
	if err != nil {
		return err
	}

	if exp.Status != models.StatusCompleted {
		return errs.ErrCannotArchiveNotCompleted
	}

	exp.Status = models.StatusArchived

	_, err = s.experimentWriter.UpdateExperiment(ctx, exp)
	if err != nil {
		return err
	}

	return nil
}
