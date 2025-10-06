package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Request represents the JSON input from the UP template engine
type Request struct {
	Function string         `json:"function"`
	Params   map[string]any `json:"params"`
	Context  map[string]any `json:"context"`
}

// Response represents the JSON output to the UP template engine
type Response struct {
	Value any    `json:"value"`
	Type  string `json:"type"`
	Error string `json:"error,omitempty"`
}

func main() {
	var req Request
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		sendError(fmt.Sprintf("Invalid request: %v", err))
		return
	}

	var result any
	var resultType string
	var err error

	switch req.Function {
	case "now":
		result, resultType, err = handleNow(req.Params)
	case "unix":
		result, resultType, err = handleUnix(req.Params)
	case "format":
		result, resultType, err = handleFormat(req.Params)
	case "parse":
		result, resultType, err = handleParse(req.Params)
	case "add":
		result, resultType, err = handleAdd(req.Params)
	case "sub":
		result, resultType, err = handleSub(req.Params)
	case "since":
		result, resultType, err = handleSince(req.Params)
	case "until":
		result, resultType, err = handleUntil(req.Params)
	default:
		sendError(fmt.Sprintf("Unknown function: %s", req.Function))
		return
	}

	if err != nil {
		sendError(err.Error())
		return
	}

	sendResponse(result, resultType)
}

// handleNow returns the current time
func handleNow(params map[string]any) (any, string, error) {
	format := getString(params, "format", time.RFC3339)
	return time.Now().Format(format), "ts", nil
}

// handleUnix returns the current Unix timestamp
func handleUnix(params map[string]any) (any, string, error) {
	return time.Now().Unix(), "int", nil
}

// handleFormat formats a time string
func handleFormat(params map[string]any) (any, string, error) {
	timeStr := getString(params, "time", "")
	if timeStr == "" {
		return nil, "", fmt.Errorf("time parameter required")
	}

	format := getString(params, "format", time.RFC3339)
	inputFormat := getString(params, "input_format", time.RFC3339)

	t, err := time.Parse(inputFormat, timeStr)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse time: %v", err)
	}

	return t.Format(format), "string", nil
}

// handleParse parses a time string
func handleParse(params map[string]any) (any, string, error) {
	timeStr := getString(params, "time", "")
	if timeStr == "" {
		return nil, "", fmt.Errorf("time parameter required")
	}

	format := getString(params, "format", time.RFC3339)
	t, err := time.Parse(format, timeStr)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse time: %v", err)
	}

	return t.Format(time.RFC3339), "ts", nil
}

// handleAdd adds duration to a time
func handleAdd(params map[string]any) (any, string, error) {
	timeStr := getString(params, "time", "")
	if timeStr == "" {
		// Use current time if not specified
		timeStr = time.Now().Format(time.RFC3339)
	}

	durationStr := getString(params, "duration", "")
	if durationStr == "" {
		return nil, "", fmt.Errorf("duration parameter required")
	}

	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse time: %v", err)
	}

	d, err := time.ParseDuration(durationStr)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse duration: %v", err)
	}

	return t.Add(d).Format(time.RFC3339), "ts", nil
}

// handleSub subtracts duration from a time
func handleSub(params map[string]any) (any, string, error) {
	timeStr := getString(params, "time", "")
	if timeStr == "" {
		timeStr = time.Now().Format(time.RFC3339)
	}

	durationStr := getString(params, "duration", "")
	if durationStr == "" {
		return nil, "", fmt.Errorf("duration parameter required")
	}

	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse time: %v", err)
	}

	d, err := time.ParseDuration(durationStr)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse duration: %v", err)
	}

	return t.Add(-d).Format(time.RFC3339), "ts", nil
}

// handleSince returns duration since a time
func handleSince(params map[string]any) (any, string, error) {
	timeStr := getString(params, "time", "")
	if timeStr == "" {
		return nil, "", fmt.Errorf("time parameter required")
	}

	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse time: %v", err)
	}

	return time.Since(t).String(), "dur", nil
}

// handleUntil returns duration until a time
func handleUntil(params map[string]any) (any, string, error) {
	timeStr := getString(params, "time", "")
	if timeStr == "" {
		return nil, "", fmt.Errorf("time parameter required")
	}

	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse time: %v", err)
	}

	return time.Until(t).String(), "dur", nil
}

// Helper functions

func getString(params map[string]any, key, defaultValue string) string {
	if v, ok := params[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultValue
}

func sendResponse(value any, valueType string) {
	resp := Response{
		Value: value,
		Type:  valueType,
	}
	if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to encode response: %v\n", err)
		os.Exit(1)
	}
}

func sendError(message string) {
	resp := Response{
		Error: message,
	}
	if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to encode error response: %v\n", err)
	}
	os.Exit(1)
}
