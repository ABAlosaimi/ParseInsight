package adapters

import "time"

// ParserAdapter defines the interface for HTTP parser implementations
type ParserAdapter interface {
	Name() string
	Parse(raw []byte) error
	BenchmarkParse(raw []byte, iterations int) BenchmarkResult
}

// BenchmarkResult contains performance metrics for a parser
type BenchmarkResult struct {
	Library         string        `json:"library"`
	TotalTime       time.Duration `json:"total_time"`
	AvgTimePerParse time.Duration `json:"avg_time_per_parse"`
	OpsPerSecond    float64       `json:"ops_per_second"`
	MemoryAllocated uint64        `json:"memory_allocated"`
	AllocsPerOp     uint64        `json:"allocs_per_op"`
	Success         bool          `json:"success"`
	Error           string        `json:"error,omitempty"`
	Winner          bool          `json:"winner,omitempty"`
}
