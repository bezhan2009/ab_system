package dto

type DecideRequest struct {
	UserID  string         `json:"user_id" binding:"required"`
	Flags   []string       `json:"flags" binding:"required,min=1"`
	Payload map[string]any `json:"payload"`
}

type FlagValue struct {
	FlagKey      string `json:"flag_key"`
	Value        string `json:"value"`
	ExperimentID string `json:"experiment_id,omitempty"`
	VariantID    string `json:"variant_id,omitempty"`
}

type DecideResponse struct {
	DecisionID string       `json:"decision_id"`
	Flags      []*FlagValue `json:"flags"`
}
