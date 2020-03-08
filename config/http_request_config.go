package config

// HTTPRequest struct to unmarshal data to make the http call required for the test
type HTTPRequest struct {
	Method             string            `json:"method"`
	URLWithEndpoint    string            `json:"urlWithEndpoint"`
	PayLoad            interface{}       `json:"payload"`
	Headers            map[string]string `json:"headers"`
	QueryParams        map[string]string `json:"queryParams"`
	SuccessStatusCodes map[int]bool      `json:"successStatusCodes"`
}
