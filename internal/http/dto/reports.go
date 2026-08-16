package dto

import "time"

type MetricValue struct {
	MetricID    string  `json:"metric_id"`
	MetricTitle string  `json:"metric_title"`
	Value       float64 `json:"value"`
	Unit        string  `json:"unit"`
}

type VariantReport struct {
	VariantID    string        `json:"variant_id"`
	VariantName  string        `json:"variant_name"`
	MetricValues []MetricValue `json:"metric_values"`
}

type ExperimentReport struct {
	ExperimentID string          `json:"experiment_id"`
	From         time.Time       `json:"from"`
	To           time.Time       `json:"to"`
	Variants     []VariantReport `json:"variants"`
}
