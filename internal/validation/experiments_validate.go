package validation

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/http/dto"
	"ab_system/internal/lib/dsl"
	"ab_system/pkg/errs"
	"fmt"

	"github.com/google/uuid"
)

func ValidateExperimentCreate(exp *dto.Experiment) []errs.FieldError {
	var errors []errs.FieldError

	if exp.Title == "" {
		errors = append(errors, errs.NewFieldError("title", "обязательное поле"))
	} else if len(exp.Title) > 255 {
		errors = append(errors, errs.NewFieldError("title", "не более 255 символов"))
	}

	if exp.FlagKey == "" {
		errors = append(errors, errs.NewFieldError("flag_key", "обязательное поле"))
	} else if len(exp.FlagKey) > 255 {
		errors = append(errors, errs.NewFieldError("flag_key", "не более 255 символов"))
	}

	if exp.TrafficPercent < 0 || exp.TrafficPercent > 100 {
		errors = append(errors, errs.NewFieldError("traffic_percent", "должен быть от 0 до 100"))
	}

	if len(exp.Variants) == 0 {
		errors = append(errors, errs.NewFieldError("variants", "должен быть хотя бы один вариант"))
	} else {
		totalWeight := 0
		controlCount := 0
		for i, v := range exp.Variants {
			if v.Title == "" {
				errors = append(errors, errs.NewFieldError(fmt.Sprintf("variants[%d].name", i), "обязательное поле"))
			}

			if v.Value == "" {
				errors = append(errors, errs.NewFieldError(fmt.Sprintf("variants[%d].value", i), "обязательное поле"))
			}

			if v.Weight <= 0 {
				errors = append(errors, errs.NewFieldError(fmt.Sprintf("variants[%d].weight", i), "должен быть положительным"))
			}

			if v.IsControl {
				controlCount++
			}

			totalWeight += v.Weight
		}

		if controlCount != 1 {
			errors = append(errors, errs.NewFieldError("variants", "должен быть ровно один контрольный вариант"))
		}
		if totalWeight != exp.TrafficPercent {
			errors = append(errors, errs.NewFieldError("variants", fmt.Sprintf("сумма весов (%d) должна равняться traffic_percent (%d)", totalWeight, exp.TrafficPercent)))
		}
	}

	resDslValidation := dsl.ValidateDSL(exp.TargetingDsl, dto.HardcodedSchema())

	if !resDslValidation.IsValid {
		for _, err := range resDslValidation.Errors {
			errors = append(errors, errs.NewFieldError("targeting_dsl", err.Error()))
		}
	} else if *resDslValidation.NormalizedExpression != "" {
		exp.TargetingDsl = *resDslValidation.NormalizedExpression
	}

	return errors
}

func ValidateExperimentUpdate(exp *dto.Experiment) []errs.FieldError {
	var errors []errs.FieldError

	if exp.Title != "" && len(exp.Title) > 255 {
		errors = append(errors, errs.NewFieldError("title", "не более 255 символов"))
	}

	if exp.FlagKey != "" && len(exp.FlagKey) > 255 {
		errors = append(errors, errs.NewFieldError("flag_key", "не более 255 символов"))
	}

	if exp.TrafficPercent < 0 || exp.TrafficPercent > 100 {
		errors = append(errors, errs.NewFieldError("traffic_percent", "должен быть от 0 до 100"))
	}

	if len(exp.Variants) > 0 {
		totalWeight := 0
		controlCount := 0
		for i, v := range exp.Variants {
			if v.Title == "" {
				errors = append(errors, errs.NewFieldError(fmt.Sprintf("variants[%d].title", i), "обязательное поле"))
			}

			if v.Value == "" {
				errors = append(errors, errs.NewFieldError(fmt.Sprintf("variants[%d].value", i), "обязательное поле"))
			}

			if v.Weight <= 0 {
				errors = append(errors, errs.NewFieldError(fmt.Sprintf("variants[%d].weight", i), "должен быть положительным"))
			}

			if v.IsControl {
				controlCount++
			}

			totalWeight += v.Weight
		}

		if exp.TrafficPercent != 0 {
			if controlCount != 1 {
				errors = append(errors, errs.NewFieldError("variants", "должен быть ровно один контрольный вариант"))
			}
			if totalWeight != exp.TrafficPercent {
				errors = append(errors, errs.NewFieldError("variants", fmt.Sprintf("сумма весов (%d) должна равняться traffic_percent (%d)", totalWeight, exp.TrafficPercent)))
			}
		}
	}

	if exp.TargetingDsl != "" {
		resDslValidation := dsl.ValidateDSL(exp.TargetingDsl, dto.HardcodedSchema())

		if !resDslValidation.IsValid {
			for _, err := range resDslValidation.Errors {
				errors = append(errors, errs.NewFieldError("targeting_dsl", err.Error()))
			}
		} else if *resDslValidation.NormalizedExpression != "" {
			exp.TargetingDsl = *resDslValidation.NormalizedExpression
		}
	}

	if exp.Status != "" {
		validStatuses := map[string]bool{
			"draft": true, "in_review": true, "approved": true,
			"running": true, "paused": true, "completed": true,
			"archived": true, "rejected": true,
		}
		if !validStatuses[exp.Status] {
			errors = append(errors, errs.NewFieldError("status", "недопустимый статус"))
		}
	}

	if exp.Status == "completed" {
		errors = append(errors, ValidateCompleteExperiment(&dto.CompleteExperimentRequest{
			Conclusion:      exp.Conclusion,
			Comment:         exp.Comment,
			WinnerVariantID: exp.WinnerVariantID,
		})...)
	}

	return errors
}

func ValidateCompleteExperiment(req *dto.CompleteExperimentRequest) []errs.FieldError {
	var errors []errs.FieldError

	if req.Conclusion == "" {
		errors = append(errors, errs.NewFieldError("conclusion", "обязательное поле"))
	} else {
		validConclusions := map[string]bool{
			"rollout":   true,
			"rollback":  true,
			"no_effect": true,
		}
		if !validConclusions[req.Conclusion] {
			errors = append(errors, errs.NewFieldError("conclusion", "допустимые значения: rollout, rollback, no_effect"))
		}
	}

	if req.Conclusion == "rollout" {
		if req.WinnerVariantID == "" {
			errors = append(errors, errs.NewFieldError("winner_variant_id", "обязательное поле при выборе победителя"))
		} else {
			if _, err := uuid.Parse(req.WinnerVariantID); err != nil {
				errors = append(errors, errs.NewFieldError("winner_variant_id", "неверный формат UUID"))
			}
		}
	} else {
		if req.WinnerVariantID != "" {
			errors = append(errors, errs.NewFieldError("winner_variant_id", "поле не должно быть заполнено при данном исходе"))
		}
	}

	if len(req.Comment) > 1000 {
		errors = append(errors, errs.NewFieldError("comment", "не более 1000 символов"))
	}

	return errors
}

func IsConfigChanged(existing, updated *models.Experiment) bool {
	if existing.Title != updated.Title ||
		existing.FlagKey != updated.FlagKey ||
		existing.TrafficPercent != updated.TrafficPercent ||
		existing.TargetingDsl != updated.TargetingDsl {
		return true
	}

	if len(existing.Variants) != len(updated.Variants) {
		return true
	}

	existingMap := make(map[string]models.Variant)
	for _, v := range existing.Variants {
		existingMap[v.ID.String()] = v
	}

	for _, v := range updated.Variants {
		if ev, ok := existingMap[v.ID.String()]; !ok {
			return true
		} else {
			if ev.Title != v.Title || ev.Value != v.Value || ev.Weight != v.Weight || ev.IsControl != v.IsControl {
				return true
			}
		}
	}

	return false
}

func ValidateStatusTransition(oldStatus, newStatus models.ExperimentStatus) error {
	allowedTransitions := map[models.ExperimentStatus][]models.ExperimentStatus{
		models.StatusDraft: {
			models.StatusInReview,
			models.StatusArchived,
		},
		models.StatusInReview: {
			//models.StatusApproved,
			//models.StatusRejected,
			models.StatusDraft,
		},
		models.StatusRejected: {
			models.StatusDraft,
		},
		models.StatusApproved: {
			models.StatusRunning,
		},
		models.StatusRunning: {
			models.StatusPaused,
			models.StatusCompleted,
		},
		models.StatusPaused: {
			models.StatusRunning,
			models.StatusCompleted,
		},
		models.StatusCompleted: {
			models.StatusArchived,
		},
		models.StatusArchived: {},
	}

	allowed, ok := allowedTransitions[oldStatus]
	if !ok {
		return errs.NewValidationError("status", "неизвестный статус")
	}

	for _, status := range allowed {
		if status == newStatus {
			return nil
		}
	}

	return errs.ErrInvalidStatusTransition
}
