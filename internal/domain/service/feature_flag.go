package service

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/domain/repository"
	"ab_system/pkg/errs"
	"context"
	"errors"
)

type FeatureFlagService struct {
	featureFlagReader  repository.FeatureFlagReader
	featureFlagWriter  repository.FeatureFlagWriter
	featureFlagDeleter repository.FeatureFlagDeleter
	userReader         repository.UserReader
}

func NewFeatureFlagService(featureFlagReader repository.FeatureFlagReader,
	featureFlagWriter repository.FeatureFlagWriter,
	featureFlagDeleter repository.FeatureFlagDeleter,
	userReader repository.UserReader,
) *FeatureFlagService {
	return &FeatureFlagService{featureFlagReader, featureFlagWriter, featureFlagDeleter, userReader}
}

func (s *FeatureFlagService) GetAllFeatureFlags(ctx context.Context) (featureFlags []models.FeatureFlag, err error) {
	featureFlags, err = s.featureFlagReader.GetAllFeatureFlags(ctx)
	if err != nil {
		return nil, err
	}

	return featureFlags, nil
}

func (s *FeatureFlagService) GetFeatureFlagById(ctx context.Context, featureFlagById string) (featureFlag models.FeatureFlag, err error) {
	featureFlag, err = s.featureFlagReader.GetFeatureFlagById(ctx, featureFlagById)
	if err != nil {
		return featureFlag, err
	}

	return featureFlag, nil
}

func (s *FeatureFlagService) GetFeatureFlagsByKey(ctx context.Context, key string) (featureFlag models.FeatureFlag, err error) {
	featureFlag, err = s.featureFlagReader.GetFeatureFlagByKey(ctx, key)
	if err != nil {
		return featureFlag, err
	}

	return featureFlag, nil
}

func (s *FeatureFlagService) GetFeatureFlagsByOwner(ctx context.Context, owner string) (featureFlags []models.FeatureFlag, err error) {
	featureFlags, err = s.featureFlagReader.GetFeatureFlagsByOwner(ctx, owner)
	if err != nil {
		return featureFlags, err
	}

	return featureFlags, nil
}

func (s *FeatureFlagService) CreateFeatureFlag(ctx context.Context, featureFlag *models.FeatureFlag) (err error) {
	_, err = s.featureFlagReader.GetFeatureFlagByKey(ctx, featureFlag.Key)
	if err == nil {
		return errs.ErrFeatureFlagAlreadyExists
	}
	if !errors.Is(err, errs.ErrRecordNotFound) {
		return err
	}

	err = s.featureFlagWriter.CreateFeatureFlag(ctx, featureFlag)
	if err != nil {
		return err
	}

	return nil
}

func (s *FeatureFlagService) UpdateFeatureFlag(ctx context.Context, featureFlag *models.FeatureFlag) (updatedFeatureFlag *models.FeatureFlag, err error) {
	//if featureFlag.UserID.String() != "" {
	//	_, err = s.userReader.GetUserByID(ctx, featureFlag.UserID.String())
	//	if err != nil {
	//		return err
	//	}
	//}

	if featureFlag.Key != "" {
		f, err := s.featureFlagReader.GetFeatureFlagByKey(ctx, featureFlag.Key)
		if err == nil && f.ID != featureFlag.ID {
			return nil, errs.ErrFeatureFlagAlreadyExists
		}
	}

	updatedFeatureFlag, err = s.featureFlagWriter.UpdateFeatureFlag(ctx, featureFlag)
	if err != nil {
		return nil, err
	}

	return updatedFeatureFlag, nil
}

func (s *FeatureFlagService) DeleteFeatureFlagById(ctx context.Context, featureFlagId string) (err error) {
	_, err = s.featureFlagReader.GetFeatureFlagById(ctx, featureFlagId)
	if err != nil {
		return err
	}

	err = s.featureFlagDeleter.DeleteFeatureFlag(ctx, featureFlagId)
	if err != nil {
		return err
	}

	return nil
}
