package service

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/domain/repository"
	"ab_system/internal/http/dto"
	"ab_system/internal/lib/dsl"
	"ab_system/pkg/errs"
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"sort"
	"time"

	"github.com/google/uuid"
)

type DecisionService struct {
	flagReader     repository.FeatureFlagReader
	expReader      repository.ExperimentReader
	decisionWriter repository.DecisionWriter
}

func NewDecisionService(
	flagReader repository.FeatureFlagReader,
	expReader repository.ExperimentReader,
	decisionWriter repository.DecisionWriter,
) *DecisionService {
	return &DecisionService{
		flagReader:     flagReader,
		expReader:      expReader,
		decisionWriter: decisionWriter,
	}
}

func (s *DecisionService) Decide(ctx context.Context, req *dto.DecideRequest) (*dto.DecideResponse, error) {
	decisionID := uuid.New().String()
	var flagValues []*dto.FlagValue

	for _, flagKey := range req.Flags {
		flag, err := s.flagReader.GetFeatureFlagByKey(ctx, flagKey)
		if err != nil {
			return nil, err
		}

		exps, err := s.expReader.GetExperimentByFlagAndStatus(ctx, flagKey, string(models.StatusRunning))
		if err != nil && !errors.Is(err, errs.ErrRecordNotFound) {
			return nil, err
		}
		var experiment *models.Experiment
		if err == nil && exps != nil && len(*exps) > 0 {
			exp := (*exps)[0]
			experiment = &exp
		}

		var value string
		var expID *uuid.UUID
		var varID *uuid.UUID

		if experiment != nil {
			passes, err := s.checkTargeting(experiment, req.Payload, req.UserID)
			if err != nil {
				return nil, err
			}

			if passes {
				h := fnv.New32a()
				h.Write([]byte(req.UserID + experiment.ID.String()))
				bucket := int(h.Sum32() % 100)

				fmt.Println(bucket)

				if bucket < experiment.TrafficPercent {
					if experiment.GuardrailTriggered {
						for _, v := range experiment.Variants {
							if v.IsControl {
								value = v.Value
								expID = &experiment.ID
								varID = &v.ID
								break
							}
						}
					} else {
						variant := s.selectVariantByBucket(experiment, bucket)
						if variant != nil {
							value = variant.Value
							expID = &experiment.ID
							varID = &variant.ID
						}
					}
				}
			}
		}

		if value == "" {
			value = flag.DefaultValue
		}

		decision := &models.Decision{
			ID:           uuid.MustParse(decisionID),
			UserID:       req.UserID,
			ExperimentID: expID,
			VariantID:    varID,
			FlagKey:      flagKey,
			Value:        value,
			CreatedAt:    time.Now(),
			ExpiresAt:    time.Now().AddDate(0, 0, 30),
		}
		if err = s.decisionWriter.CreateDecision(ctx, decision); err != nil {
			return nil, err
		}

		expIdStr := ""
		varIdStr := ""
		if expID != nil {
			expIdStr = expID.String()
		}
		if varID != nil {
			varIdStr = varID.String()
		}

		flagValues = append(flagValues, &dto.FlagValue{
			FlagKey:      flagKey,
			Value:        value,
			ExperimentID: expIdStr,
			VariantID:    varIdStr,
		})
	}

	return &dto.DecideResponse{
		DecisionID: decisionID,
		Flags:      flagValues,
	}, nil
}

func (s *DecisionService) checkTargeting(exp *models.Experiment, attributes map[string]any, userID string) (bool, error) {
	if exp.TargetingDsl == "" {
		return true, nil
	}

	payload := dto.Payload{}
	for k, v := range attributes {
		payload[k] = v
	}
	payload["user_id"] = userID

	schema := dto.HardcodedSchema()

	result := dsl.EvaluateDSL(exp.TargetingDsl, payload, schema)
	return result.Matched, nil
}

func (s *DecisionService) selectVariantByBucket(exp *models.Experiment, bucket int) *models.Variant {
	if len(exp.Variants) == 0 {
		return nil
	}

	sorted := make([]models.Variant, len(exp.Variants))
	copy(sorted, exp.Variants)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ID.String() < sorted[j].ID.String()
	})

	cumulative := 0
	for _, v := range sorted {
		cumulative += v.Weight
		if bucket < cumulative {
			return &v
		}
	}

	return &sorted[len(sorted)-1]
}
