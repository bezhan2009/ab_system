package validation

import (
	"ab_system/internal/domain/models"
	"ab_system/pkg/errs"
)

func ValidateGuardrailConfig(binding *models.ExperimentMetric) (err []errs.FieldError) {
	if binding.Threshold == 0 {
		err = append(err, errs.NewFieldError("threshold", errs.ErrThresholdRequired.Error()))
	}

	if binding.Operator == "" {
		err = append(err, errs.NewFieldError("operator", errs.ErrOperatorRequired.Error()))
	}

	switch binding.Operator {
	case ">", ">=", "<", "<=":
		// ok
	default:
		err = append(err, errs.NewFieldError("operator", errs.ErrInvalidOperator.Error()))
	}

	if binding.WindowMin <= 0 {
		err = append(err, errs.NewFieldError("window_min", errs.ErrWindowMinRequired.Error()))
	}

	if binding.Action == "" {
		err = append(err, errs.NewFieldError("action", errs.ErrActionRequired.Error()))
	}

	switch binding.Action {
	case "pause", "rollback":
		// ok
	default:
		err = append(err, errs.NewFieldError("action", errs.ErrInvalidAction.Error()))
	}

	return err
}
