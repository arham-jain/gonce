package config

// ConditionConfig a struct used to unmarsahal the json config file
// Contains info like thread count, rampup, execution time, etc.
type ConditionConfig struct {
	Threads                int64       `json:"threads"`
	RampUpTimeInSeconds    int64       `json:"rampUpTimeInSeconds"`
	ExecutionTimeInSeconds int64       `json:"executionTimeInSeconds"`
	HTTPRequest            HTTPRequest `json:"httpRequest"`
}
