package benchmark

import (
	"fmt"
	"sync"
	"time"

	"github.com/ABAlosaimi/ParseInsight/internal/adapters"
)

// Runner executes benchmarks across multiple parsers
type Runner struct {
	registry  *adapters.Registry
	validator *Validator
}

// NewRunner creates a new benchmark runner
func NewRunner() *Runner {
	return &Runner{
		registry:  adapters.NewRegistry(),
		validator: NewValidator(),
	}
}

// Run executes benchmarks for the given HTTP message
func (r *Runner) Run(msg HTTPMessage) (*Result, error) {
	// Validate message
	if err := r.validator.ValidateHTTPMessage(&msg); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Validate config
	if err := r.validator.ValidateTestConfig(&msg.TestConfig); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	// Run benchmarks for each library
	results := make([]adapters.BenchmarkResult, 0, len(msg.TestConfig.Libraries))

	for _, libName := range msg.TestConfig.Libraries {
		adapter, err := r.registry.Get(libName, msg.MessageType)
		if err != nil {
			results = append(results, adapters.BenchmarkResult{
				Library: libName,
				Success: false,
				Error:   err.Error(),
			})
			continue
		}

		var benchResult adapters.BenchmarkResult

		if msg.TestConfig.Concurrency > 1 {
			benchResult = r.runConcurrent(adapter, []byte(msg.Raw), msg.TestConfig)
		} else {
			benchResult = adapter.BenchmarkParse([]byte(msg.Raw), msg.TestConfig.Iterations)
		}

		results = append(results, benchResult)
	}

	// Determine winner
	r.markWinner(results)

	// Generate recommendation
	recommendation := r.generateRecommendation(results)

	return &Result{
		Results:        results,
		Recommendation: recommendation,
		MessageType:    msg.MessageType,
	}, nil
}

// runConcurrent executes benchmark with concurrent workers
func (r *Runner) runConcurrent(adapter adapters.ParserAdapter, raw []byte, config TestConfig) adapters.BenchmarkResult {
	var wg sync.WaitGroup
	workerResults := make([]adapters.BenchmarkResult, config.Concurrency)

	iterPerWorker := config.Iterations / config.Concurrency

	for i := 0; i < config.Concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			workerResults[workerID] = adapter.BenchmarkParse(raw, iterPerWorker)
		}(i)
	}

	wg.Wait()

	// Aggregate results
	return r.aggregateResults(adapter.Name(), workerResults)
}

// aggregateResults combines results from multiple workers
func (r *Runner) aggregateResults(libName string, results []adapters.BenchmarkResult) adapters.BenchmarkResult {
	aggregate := adapters.BenchmarkResult{
		Library: libName,
		Success: true,
	}

	var totalOps float64
	var totalMemory uint64
	var totalAllocs uint64

	for _, res := range results {
		if !res.Success {
			aggregate.Success = false
			aggregate.Error = res.Error
			return aggregate
		}

		aggregate.TotalTime += res.TotalTime
		totalOps += res.OpsPerSecond
		totalMemory += res.MemoryAllocated
		totalAllocs += res.AllocsPerOp
	}

	count := len(results)
	if count > 0 {
		aggregate.OpsPerSecond = totalOps / float64(count)
		aggregate.MemoryAllocated = totalMemory / uint64(count)
		aggregate.AllocsPerOp = totalAllocs / uint64(count)
		aggregate.AvgTimePerParse = aggregate.TotalTime / time.Duration(count)
	}

	return aggregate
}

// markWinner marks the fastest successful result as winner
func (r *Runner) markWinner(results []adapters.BenchmarkResult) {
	var maxOps float64
	var winnerIdx = -1

	for i, res := range results {
		if res.Success && res.OpsPerSecond > maxOps {
			maxOps = res.OpsPerSecond
			winnerIdx = i
		}
	}

	if winnerIdx >= 0 {
		results[winnerIdx].Winner = true
	}
}

// generateRecommendation creates a human-readable recommendation
func (r *Runner) generateRecommendation(results []adapters.BenchmarkResult) string {
	if len(results) == 0 {
		return "No results available"
	}

	var winner *adapters.BenchmarkResult
	var runner *adapters.BenchmarkResult

	// Find winner and runner-up
	for i := range results {
		if !results[i].Success {
			continue
		}
		if results[i].Winner {
			winner = &results[i]
		} else if runner == nil || results[i].OpsPerSecond > runner.OpsPerSecond {
			runner = &results[i]
		}
	}

	if winner == nil {
		return "All parsers failed"
	}

	if runner == nil || len(results) == 1 {
		return fmt.Sprintf("%s completed successfully with %.0f ops/sec",
			winner.Library, winner.OpsPerSecond)
	}

	speedup := winner.OpsPerSecond / runner.OpsPerSecond
	return fmt.Sprintf("%s is %.2fx faster than %s (%.0f vs %.0f ops/sec)",
		winner.Library, speedup, runner.Library, winner.OpsPerSecond, runner.OpsPerSecond)
}

// Result contains benchmark results and analysis
type Result struct {
	Results        []adapters.BenchmarkResult `json:"results"`
	Recommendation string                      `json:"recommendation"`
	MessageType    string                      `json:"message_type"`
}
