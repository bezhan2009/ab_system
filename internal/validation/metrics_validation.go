package validation

import (
	"ab_system/internal/http/dto"
	"ab_system/pkg/errs"
)

func ValidateMetricCreate(metric dto.Metric) []errs.FieldError {
	var errors []errs.FieldError

	if metric.Title == "" {
		errors = append(errors, errs.NewFieldError("title", "обязательное поле"))
	} else if len(metric.Type) > 255 {
		errors = append(errors, errs.NewFieldError("title", "не более 255 символов"))
	}

	if metric.Type == "" {
		errors = append(errors, errs.NewFieldError("type", "обязательное поле"))
	} else {
		switch metric.Type {
		case "counter":
			if metric.CounterEventType == "" {
				errors = append(errors, errs.NewFieldError("counter_event_type", "обязательное поле для типа counter"))
			}
		case "ratio":
			if metric.NumeratorEventType == "" {
				errors = append(errors, errs.NewFieldError("numerator_event_type", "обязательное поле для типа ratio"))
			}
			if metric.DenominatorEventType == "" {
				errors = append(errors, errs.NewFieldError("denominator_event_type", "обязательное поле для типа ratio"))
			}
		case "histogram":
			if metric.HistogramEventType == "" {
				errors = append(errors, errs.NewFieldError("histogram_event_type", "обязательное поле для типа histogram"))
			}
			if metric.HistogramField == "" {
				errors = append(errors, errs.NewFieldError("histogram_field", "обязательное поле для типа histogram"))
			}
		default:
			errors = append(errors, errs.NewFieldError("type", "допустимые типы: counter, ratio, histogram"))
		}
	}

	if len(metric.Description) > 1000 {
		errors = append(errors, errs.NewFieldError("description", "не более 1000 символов"))
	}
	if len(metric.Unit) > 50 {
		errors = append(errors, errs.NewFieldError("unit", "не более 50 символов"))
	}

	return errors
}

func ValidateMetricUpdate(metric dto.Metric) []errs.FieldError {
	var errors []errs.FieldError

	if metric.Title != "" && len(metric.Title) > 255 {
		errors = append(errors, errs.NewFieldError("title", "не более 255 символов"))
	}

	if metric.Type != "" {
		switch metric.Type {
		case "counter":
			if metric.CounterEventType == "" {
				errors = append(errors, errs.NewFieldError("counter_event_type", "обязательное поле для типа counter"))
			}
		case "ratio":
			if metric.NumeratorEventType == "" {
				errors = append(errors, errs.NewFieldError("numerator_event_type", "обязательное поле для типа ratio"))
			}
			if metric.DenominatorEventType == "" {
				errors = append(errors, errs.NewFieldError("denominator_event_type", "обязательное поле для типа ratio"))
			}
		case "histogram":
			if metric.HistogramEventType == "" {
				errors = append(errors, errs.NewFieldError("histogram_event_type", "обязательное поле для типа histogram"))
			}
			if metric.HistogramField == "" {
				errors = append(errors, errs.NewFieldError("histogram_field", "обязательное поле для типа histogram"))
			}
		default:
			errors = append(errors, errs.NewFieldError("type", "допустимые типы: counter, ratio, histogram"))
		}
	}

	if len(metric.Description) > 1000 {
		errors = append(errors, errs.NewFieldError("description", "не более 1000 символов"))
	}
	if len(metric.Unit) > 50 {
		errors = append(errors, errs.NewFieldError("unit", "не более 50 символов"))
	}

	return errors
}
