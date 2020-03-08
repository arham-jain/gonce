package models

// Metrics struct to populate the final metrics post CurrentMetric evaluation
type Metrics struct {
	AverageResponseTime int64 `json:"average_response_time"`
	PeakResponseTime    int64 `json:"peak_response_time"`
	TotalRequests       int64 `json:"total_requests"`
	ErrorCount          int64 `json:"error_count"`
}
