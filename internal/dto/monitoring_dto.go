package dto

import "time"

type RainfallHistoryResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Hourly    float64   `json:"hourly"`
	Daily     float64   `json:"daily"`
}
