package dto

import "ab_system/pkg/errs"

type DslValidateRequest struct {
	DslExpression string `json:"dslExpression" binding:"required,min=3,max=2000"`
}

type DslValidateResult struct {
	IsValid              bool                 `json:"isValid"`
	NormalizedExpression *string              `json:"normalizedExpression,omitempty"`
	Errors               []DslValidationError `json:"errors"`
}

type DslValidationError struct {
	Code     errs.ErrorCode `json:"code"`
	Message  string         `json:"message"`
	Position *int           `json:"position,omitempty"`
	Near     *string        `json:"near,omitempty"`
}
