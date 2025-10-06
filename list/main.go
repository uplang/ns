package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Request struct {
	Function string                 `json:"function"`
	Params   map[string]any `json:"params"`
	Context  map[string]any `json:"context"`
}

type Response struct {
	Value any `json:"value"`
	Type  string      `json:"type"`
	Error string      `json:"error,omitempty"`
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
	case "generate":
		result, resultType, err = handleGenerate(req.Params, req.Context)
	case "join":
		result, resultType, err = handleJoin(req.Params)
	case "slice":
		result, resultType, err = handleSlice(req.Params)
	case "length":
		result, resultType, err = handleLength(req.Params)
	case "contains":
		result, resultType, err = handleContains(req.Params)
	case "index":
		result, resultType, err = handleIndex(req.Params)
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

func handleGenerate(params map[string]any, context map[string]any) (any, string, error) {
	count := getInt(params, "count", 0)
	if count <= 0 {
		return nil, "", fmt.Errorf("count parameter required and must be positive")
	}

	template, ok := params["template"]
	if !ok {
		return nil, "", fmt.Errorf("template parameter required")
	}

	result := make([]any, count)
	for i := 0; i < count; i++ {
		// Create $self context
		self := map[string]any{
			"number": i + 1,
			"index":  i,
			"count":  count,
			"first":  i == 0,
			"last":   i == count-1,
		}

		// Clone template and inject $self
		item := cloneWithContext(template, self)
		result[i] = item
	}

	return result, "list", nil
}

func handleJoin(params map[string]any) (any, string, error) {
	items, ok := params["items"].([]any)
	if !ok {
		return nil, "", fmt.Errorf("items parameter required and must be a list")
	}

	separator := getString(params, "separator", ",")

	strItems := make([]string, len(items))
	for i, item := range items {
		strItems[i] = fmt.Sprint(item)
	}

	return strings.Join(strItems, separator), "string", nil
}

func handleSlice(params map[string]any) (any, string, error) {
	items, ok := params["items"].([]any)
	if !ok {
		return nil, "", fmt.Errorf("items parameter required and must be a list")
	}

	start := getInt(params, "start", 0)
	end := getInt(params, "end", len(items))

	if start < 0 {
		start = 0
	}
	if end > len(items) {
		end = len(items)
	}
	if start > end {
		start = end
	}

	return items[start:end], "list", nil
}

func handleLength(params map[string]any) (any, string, error) {
	items, ok := params["items"].([]any)
	if !ok {
		return nil, "", fmt.Errorf("items parameter required and must be a list")
	}

	return len(items), "int", nil
}

func handleContains(params map[string]any) (any, string, error) {
	items, ok := params["items"].([]any)
	if !ok {
		return nil, "", fmt.Errorf("items parameter required and must be a list")
	}

	value, ok := params["value"]
	if !ok {
		return nil, "", fmt.Errorf("value parameter required")
	}

	for _, item := range items {
		if fmt.Sprint(item) == fmt.Sprint(value) {
			return true, "bool", nil
		}
	}

	return false, "bool", nil
}

func handleIndex(params map[string]any) (any, string, error) {
	items, ok := params["items"].([]any)
	if !ok {
		return nil, "", fmt.Errorf("items parameter required and must be a list")
	}

	index := getInt(params, "index", 0)

	if index < 0 || index >= len(items) {
		return nil, "", fmt.Errorf("index out of range")
	}

	return items[index], "string", nil
}

func cloneWithContext(template any, self map[string]any) any {
	// Simple clone - in real implementation would process $self references
	return template
}

func getString(params map[string]any, key, defaultValue string) string {
	if v, ok := params[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultValue
}

func getInt(params map[string]any, key string, defaultValue int) int {
	if v, ok := params[key]; ok {
		switch val := v.(type) {
		case float64:
			return int(val)
		case int:
			return val
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

