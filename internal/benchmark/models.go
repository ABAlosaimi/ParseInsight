package benchmark

import "time"

// TestConfig defines configuration for benchmark tests
type TestConfig struct {
	Iterations    int           `json:"iterations"`
	Concurrency   int           `json:"concurrency"`
	Libraries     []string      `json:"libraries"`
	MeasureMemory bool          `json:"measure_memory"`
	Duration      time.Duration `json:"duration"`
}

// HTTPMessage represents user's HTTP message to test
type HTTPMessage struct {
	Raw         string     `json:"raw"`
	MessageType string     `json:"message_type"` // "request" or "response"
	TestConfig  TestConfig `json:"test_config"`
}

// Constants for limits
const (
	MaxMessageSize = 1 * 1024 * 1024  // 1MB
	MaxIterations  = 10000000         // 10M
	MaxConcurrency = 100
	MaxDuration    = 30 * time.Second
	MinIterations  = 1
	MinConcurrency = 1
)
