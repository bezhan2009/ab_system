package validation

import (
	"ab_system/internal/http/dto"
	"ab_system/pkg/errs"
)

func ValidateFeatureFlagCreate(featureFlag dto.FeatureFlag) []errs.FieldError {
	var errors []errs.FieldError

	if featureFlag.Key == "" {
		errors = append(errors, errs.NewFieldError("key", errs.ErrRequiredField.Error()))
	}

	if featureFlag.DefaultValue == "" {
		errors = append(errors, errs.NewFieldError("default_value", errs.ErrRequiredField.Error()))
	}

	if featureFlag.UserID == "" {
		errors = append(errors, errs.NewFieldError("usere_id", errs.ErrRequiredField.Error()))
	}

	if featureFlag.Type == "" {
		errors = append(errors, errs.NewFieldError("type", errs.ErrRequiredField.Error()))
	}

	typeStatuses := map[string]bool{
		"string":  true,
		"number":  true,
		"boolean": true,
	}

	if !typeStatuses[featureFlag.Type] {
		errors = append(errors, errs.NewFieldError("type", errs.ErrInvalidValue.Error()))
	}

	return errors
}
