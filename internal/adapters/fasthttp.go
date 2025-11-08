package adapters

import (
	"bufio"
	"bytes"
	"fmt"
	"runtime"
	"time"

	"github.com/valyala/fasthttp"
)

// FastHTTPAdapter implements ParserAdapter for fasthttp
type FastHTTPAdapter struct {
	messageType string // "request" or "response"
}

// NewFastHTTPAdapter creates a new fasthttp adapter
func NewFastHTTPAdapter(messageType string) *FastHTTPAdapter {
	return &FastHTTPAdapter{messageType: messageType}
}

// Name returns the adapter name
func (a *FastHTTPAdapter) Name() string {
	return "fasthttp"
}

// Parse parses an HTTP message using fasthttp
func (a *FastHTTPAdapter) Parse(raw []byte) error {
	if a.messageType == "request" {
		var req fasthttp.Request
		err := req.ReadLimitBody(bufio.NewReader(bytes.NewReader(raw)), len(raw))
		return err
	}

	var resp fasthttp.Response
	err := resp.ReadLimitBody(bufio.NewReader(bytes.NewReader(raw)), len(raw))
	return err
}

// BenchmarkParse runs benchmark for fasthttp parser
func (a *FastHTTPAdapter) BenchmarkParse(raw []byte, iterations int) BenchmarkResult {
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
