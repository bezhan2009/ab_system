package errs

import (
	"ab_system/pkg/date_time"
)

type ApiErrorResponse struct {
	Code        ErrorCode              `json:"code"`
	Message     string                 `json:"message"`
	TraceId     string                 `json:"trace_id"`
	Timestamp   string                 `json:"timestamp"`
	Path        string                 `json:"path"`
	Details     map[string]interface{} `json:"details,omitempty"`
	FieldErrors []FieldError           `json:"fieldErrors,omitempty"`
}

type ValidationErrorResponse struct {
	Code        ErrorCode    `json:"code"`
	Message     string       `json:"message"`
	TraceId     string       `json:"trace_id"`
	Timestamp   string       `json:"timestamp"`
	Path        string       `json:"path"`
	FieldErrors []FieldError `json:"fieldErrors"`
}

type ErrorMapping struct {
	HttpStatus int                    `json:"http_status"`
	ErrorCode  ErrorCode              `json:"error_code"`
	Details    map[string]interface{} `json:"details"`
}

type ValidationErrors struct {
	FieldErrors []FieldError
}

func (e *ValidationErrors) Error() string {
	if len(e.FieldErrors) > 0 {
		return e.FieldErrors[0].Issue
	}

	return "ошибки валидации"
}

func NewValidationErrors(fieldErrors []FieldError) error {
	return &ValidationErrors{
		FieldErrors: fieldErrors,
	}
}

func NewErrorResponse(code ErrorCode, message, traceId, path string) *ApiErrorResponse {
	return &ApiErrorResponse{
		Code:      code,
		Message:   message,
		TraceId:   traceId,
		Timestamp: date_time.GetCurrentTimestamp(),
		Path:      path,
	}
}
