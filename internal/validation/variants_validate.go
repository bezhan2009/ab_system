package validation

import (
	"ab_system/internal/domain/models"
	"ab_system/pkg/errs"
)

func ValidateVariants(experiment *models.Experiment) (err error) {
	if len(experiment.Variants) == 0 {
		return errs.ErrNoVariants
	}

	totalWeight := 0
	controlCount := 0

	for _, v := range experiment.Variants {
		if v.Weight <= 0 {
			return errs.ErrInvalidVariantWeight
		}

		totalWeight += v.Weight

		if v.IsControl {
			controlCount++
		}

		if v.Title == "" {
			return errs.ErrVariantNameRequired
		}
	}

	if controlCount != 1 {
		return errs.ErrInvalidControlVariantCount
	}

	if totalWeight != experiment.TrafficPercent {
		return errs.ErrVariantWeightSumMismatch
	}

	return nil
}

func VariantsEqual(v1, v2 []models.Variant) bool {
	if len(v1) != len(v2) {
		return false
	}

	map1 := make(map[string]models.Variant)
	for _, v := range v1 {
		map1[v.ID.String()] = v
	}

	for _, v := range v2 {
		if v1v, ok := map1[v.ID.String()]; !ok {
			return false
		} else {
			if v1v.Title != v.Title || v1v.Value != v.Value ||
				v1v.Weight != v.Weight || v1v.IsControl != v.IsControl {
				return false
			}
		}
	}

	return true
}
