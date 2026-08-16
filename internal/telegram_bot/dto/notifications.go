package dto


type NotifyRequest struct {
	EventType    string `json:"event_type" binding:"required"`
	ExperimentID string `json:"experiment_id"`
	Experiment   string `json:"experiment"`
	FlagKey      string `json:"flag_key"`
	UserID       string `json:"user_id"`
	Status       string `json:"status"`
	Metric       string `json:"metric"`
	Threshold    string `json:"threshold"`
	Value        string `json:"value"`
	Message      string `json:"message"`
	ReportURL    string `json:"report_url"`
	Timestamp    string `json:"timestamp"`
}