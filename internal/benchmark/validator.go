package benchmark

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

// Validator validates HTTP messages and configurations
type Validator struct{}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateHTTPMessage validates the HTTP message and determines its type
func (v *Validator) ValidateHTTPMessage(msg *HTTPMessage) error {
	if msg.Raw == "" {
		return fmt.Errorf("HTTP message cannot be empty")
	}

	if len(msg.Raw) > MaxMessageSize {
		return fmt.Errorf("message size exceeds maximum of %d bytes", MaxMessageSize)
	}

	// Auto-detect message type if not specified
	if msg.MessageType == "" {
		msgType, err := v.detectMessageType(msg.Raw)
		if err != nil {
			return err
		}
		msg.MessageType = msgType
	}

	// Validate message type
	if msg.MessageType != "request" && msg.MessageType != "response" {
		return fmt.Errorf("message_type must be 'request' or 'response'")
	}

	// Try to parse with net/http to validate format
	if err := v.validateFormat(msg.Raw, msg.MessageType); err != nil {
		return fmt.Errorf("invalid HTTP message format: %v", err)
	}

	return nil
}

// detectMessageType auto-detects if message is request or response
func (v *Validator) detectMessageType(raw string) (string, error) {
	lines := strings.Split(raw, "\n")
	if len(lines) == 0 {
		return "", fmt.Errorf("empty message")
	}

	firstLine := strings.TrimSpace(lines[0])

	// Check if it's a request (starts with method)
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "CONNECT", "TRACE"}
	for _, method := range methods {
		if strings.HasPrefix(firstLine, method+" ") {
			return "request", nil
		}
	}

	// Check if it's a response (starts with HTTP/)
	if strings.HasPrefix(firstLine, "HTTP/") {
		return "response", nil
	}

	return "", fmt.Errorf("unable to detect message type from first line: %s", firstLine)
}

// validateFormat validates HTTP message format using net/http
func (v *Validator) validateFormat(raw, messageType string) error {
	reader := bufio.NewReader(bytes.NewReader([]byte(raw)))

	if messageType == "request" {
		_, err := http.ReadRequest(reader)
		return err
	}

	_, err := http.ReadResponse(reader, nil)
	return err
}

// ValidateTestConfig validates benchmark configuration
func (v *Validator) ValidateTestConfig(config *TestConfig) error {
	// Set defaults
	if config.Iterations == 0 {
		config.Iterations = 10000
	}
	if config.Concurrency == 0 {
		config.Concurrency = 1
	}
	if len(config.Libraries) == 0 {
		config.Libraries = []string{"net/http", "fasthttp"}
	}

	// Validate limits
	if config.Iterations < MinIterations || config.Iterations > MaxIterations {
		return fmt.Errorf("iterations must be between %d and %d", MinIterations, MaxIterations)
	}

	if config.Concurrency < MinConcurrency || config.Concurrency > MaxConcurrency {
		return fmt.Errorf("concurrency must be between %d and %d", MinConcurrency, MaxConcurrency)
	}

	if config.Duration > MaxDuration {
		return fmt.Errorf("duration must not exceed %v", MaxDuration)
	}

	return nil
}
