package adapters

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

// NetHTTPAdapter implements ParserAdapter for net/http
type NetHTTPAdapter struct {
	messageType string // "request" or "response"
}

// NewNetHTTPAdapter creates a new net/http adapter
func NewNetHTTPAdapter(messageType string) *NetHTTPAdapter {
	return &NetHTTPAdapter{messageType: messageType}
}

// Name returns the adapter name
func (a *NetHTTPAdapter) Name() string {
	return "net/http"
}

// Parse parses an HTTP message using net/http
func (a *NetHTTPAdapter) Parse(raw []byte) error {
	reader := bufio.NewReader(bytes.NewReader(raw))

	if a.messageType == "request" {
		_, err := http.ReadRequest(reader)
		return err
	}

	_, err := http.ReadResponse(reader, nil)
	return err
}

// BenchmarkParse runs benchmark for net/http parser
func (a *NetHTTPAdapter) BenchmarkParse(raw []byte, iterations int) BenchmarkResult {
	result := BenchmarkResult{
		Library: a.Name(),
		Success: true,
	}

	// Clean GC before benchmark
	runtime.GC()

	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	start := time.Now()

	for i := 0; i < iterations; i++ {
		if err := a.Parse(raw); err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("Parse error: %v", err)
			return result
		}
	}

	duration := time.Since(start)

	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	result.TotalTime = duration
	result.AvgTimePerParse = duration / time.Duration(iterations)
	result.OpsPerSecond = float64(iterations) / duration.Seconds()
	result.MemoryAllocated = memAfter.TotalAlloc - memBefore.TotalAlloc

	if iterations > 0 {
		result.AllocsPerOp = (memAfter.Mallocs - memBefore.Mallocs) / uint64(iterations)
	}

	return result
}
