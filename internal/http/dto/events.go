package dto

import (
	"ab_system/pkg/errs"
	"time"
)

type EventRequest struct {
	EventID    string         `json:"event_id"`
	EventType  string         `json:"event_type"`
	DecisionID string         `json:"decision_id"`
	UserID     string         `json:"user_id"`
	ClientTime *time.Time     `json:"client_time,omitempty"`
	Payload    map[string]any `json:"payload,omitempty"`
}

type EventBatchRequest []EventRequest

type EventResponse struct {
	Accepted  int               `json:"accepted"`
	Duplicate int               `json:"duplicate"`
	Rejected  int               `json:"rejected"`
	Errors    []errs.FieldError `json:"errors,omitempty"`
}
