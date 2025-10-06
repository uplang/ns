package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Request struct {
	Function string         `json:"function"`
	Params   map[string]any `json:"params"`
	Context  map[string]any `json:"context"`
}

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
	case "get":
		result, resultType, err = handleGet(req.Params)
	case "has":
		result, resultType, err = handleHas(req.Params)
	case "list":
		result, resultType, err = handleList(req.Params)
	case "expand":
		result, resultType, err = handleExpand(req.Params)
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

func handleGet(params map[string]any) (any, string, error) {
	key := getString(params, "key", "")
	if key == "" {
		return nil, "", fmt.Errorf("key parameter required")
	}

	defaultValue := getString(params, "default", "")
	value := os.Getenv(key)

	if value == "" && defaultValue != "" {
		return defaultValue, "string", nil
	}

	return value, "string", nil
}

func handleHas(params map[string]any) (any, string, error) {
	key := getString(params, "key", "")
	if key == "" {
		return nil, "", fmt.Errorf("key parameter required")
	}

	_, exists := os.LookupEnv(key)
	return exists, "bool", nil
}

func handleList(params map[string]any) (any, string, error) {
	prefix := getString(params, "prefix", "")

	env := os.Environ()
	result := make(map[string]string)

	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			if prefix == "" || strings.HasPrefix(key, prefix) {
				result[key] = parts[1]
			}
		}
	}

	return result, "block", nil
}

func handleExpand(params map[string]any) (any, string, error) {
	text := getString(params, "text", "")
	if text == "" {
		return nil, "", fmt.Errorf("text parameter required")
	}

	expanded := os.ExpandEnv(text)
	return expanded, "string", nil
}

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
