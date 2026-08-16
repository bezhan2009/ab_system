package errs

import "fmt"

type FieldError struct {
	Field         string      `json:"field"`
	Issue         string      `json:"issue"`
	RejectedValue interface{} `json:"rejectedValue,omitempty"`
}

func (f FieldError) Error() string {
	return fmt.Sprintf("Field %s is %s", f.Field, f.RejectedValue)
}

func NewFieldError(field, issue string, rejectedValue ...interface{}) FieldError {
	fe := FieldError{
		Field: field,
		Issue: issue,
	}

	if len(rejectedValue) > 0 {
		fe.RejectedValue = rejectedValue[0]
	}

	return fe
}

type MultiValidationError struct {
	FieldErrors []FieldError `json:"fieldErrors"`
	message     string
}

func (e *MultiValidationError) Error() string {
	return e.message
}

func NewMultiValidationError(fieldErrors []FieldError) error {
	return &MultiValidationError{
		FieldErrors: fieldErrors,
		message:     "Некоторые поля не прошли валидацию",
	}
}

type ValidationError struct {
	Field string
	Issue string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("ошибка валидации с полем '%s': %s", e.Field, e.Issue)
}

func NewValidationError(field, issue string) error {
	return &ValidationError{
		Field: field,
		Issue: issue,
	}
}

func GetValidationErrorDetails(err error) (field, issue string, ok bool) {
	if vErr, ok := err.(*ValidationError); ok {
		return vErr.Field, vErr.Issue, true
	}
	return "", "", false
}
