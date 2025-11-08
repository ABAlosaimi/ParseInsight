package api

// BenchmarkRequest represents the API request for benchmarking
type BenchmarkRequest struct {
	Message     string   `json:"message"`
	MessageType string   `json:"message_type,omitempty"`
	Iterations  int      `json:"iterations,omitempty"`
	Concurrency int      `json:"concurrency,omitempty"`
	Libraries   []string `json:"libraries,omitempty"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// LibrariesResponse contains available libraries
type LibrariesResponse struct {
	Libraries []string `json:"libraries"`
}
