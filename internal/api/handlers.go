package api

import (
	"encoding/json"
	"net/http"

	"github.com/ABAlosaimi/ParseInsight/internal/adapters"
	"github.com/ABAlosaimi/ParseInsight/internal/benchmark"
)

// Handler manages HTTP API endpoints
type Handler struct {
	runner   *benchmark.Runner
	registry *adapters.Registry
}

// NewHandler creates a new API handler
func NewHandler() *Handler {
	return &Handler{
		runner:   benchmark.NewRunner(),
		registry: adapters.NewRegistry(),
	}
}

// HandleBenchmark processes benchmark requests
func (h *Handler) HandleBenchmark(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req BenchmarkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Create HTTP message from request
	msg := benchmark.HTTPMessage{
		Raw: req.Message,
		MessageType: req.MessageType,
		TestConfig: benchmark.TestConfig{
		Iterations:  req.Iterations,
		Concurrency: req.Concurrency,
		Libraries:   req.Libraries,
		},
	}

	// Run benchmark
	result, err := h.runner.Run(msg)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Write response
	h.writeJSON(w, result, http.StatusOK)
}

// HandleLibraries returns available parser libraries
func (h *Handler) HandleLibraries(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := LibrariesResponse{
		Libraries: h.registry.Available(),
	}

	h.writeJSON(w, response, http.StatusOK)
}

// HandleIndex serves the main HTML page
func (h *Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "web/index.html")
}

// writeJSON writes JSON response
func (h *Handler) writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError writes error response
func (h *Handler) writeError(w http.ResponseWriter, message string, status int) {
	h.writeJSON(w, ErrorResponse{Error: message}, status)
}