package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
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
	case "int":
		result, resultType, err = handleInt(req.Params)
	case "float":
		result, resultType, err = handleFloat(req.Params)
	case "bool":
		result, resultType, err = handleBool(req.Params)
	case "choice":
		result, resultType, err = handleChoice(req.Params)
	case "bytes":
		result, resultType, err = handleBytes(req.Params)
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

func handleInt(params map[string]any) (any, string, error) {
	min := getInt64(params, "min", 0)
	max := getInt64(params, "max", 100)

	if min >= max {
		return nil, "", fmt.Errorf("min must be less than max")
	}

	diff := max - min
	n, err := rand.Int(rand.Reader, big.NewInt(diff))
	if err != nil {
		return nil, "", err
	}

	return min + n.Int64(), "int", nil
}

func handleFloat(params map[string]any) (any, string, error) {
	min := getFloat64(params, "min", 0.0)
	max := getFloat64(params, "max", 1.0)

	if min >= max {
		return nil, "", fmt.Errorf("min must be less than max")
	}

	// Generate random bytes and convert to float
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return nil, "", err
	}

	// Convert to uint64 then to float in range [0, 1)
	n := uint64(0)
	for _, v := range b {
		n = (n << 8) | uint64(v)
	}

	// Scale to [0, 1)
	f := float64(n) / float64(^uint64(0))

	// Scale to [min, max)
	result := min + f*(max-min)

	return result, "float", nil
}

func handleBool(params map[string]any) (any, string, error) {
	b := make([]byte, 1)
	if _, err := rand.Read(b); err != nil {
		return nil, "", err
	}

	return b[0]&1 == 1, "bool", nil
}

func handleChoice(params map[string]any) (any, string, error) {
	items, ok := params["items"].([]any)
	if !ok || len(items) == 0 {
		return nil, "", fmt.Errorf("items parameter required and must be non-empty list")
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(items))))
	if err != nil {
		return nil, "", err
	}

	return items[n.Int64()], "string", nil
}

func handleBytes(params map[string]any) (any, string, error) {
	size := int(getInt64(params, "size", 16))

	if size <= 0 || size > 1024 {
		return nil, "", fmt.Errorf("size must be between 1 and 1024")
	}

	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return nil, "", err
	}

	// Encode as hex string
	return fmt.Sprintf("%x", b), "string", nil
}

func getString(params map[string]any, key, defaultValue string) string {
	if v, ok := params[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultValue
}

func getInt64(params map[string]any, key string, defaultValue int64) int64 {
	if v, ok := params[key]; ok {
		switch val := v.(type) {
		case float64:
			return int64(val)
		case int64:
			return val
		case int:
			return int64(val)
		}
	}
	return defaultValue
}

func getFloat64(params map[string]any, key string, defaultValue float64) float64 {
	if v, ok := params[key]; ok {
		switch val := v.(type) {
		case float64:
			return val
		case int:
			return float64(val)
		case int64:
			return float64(val)
		}
	}
	return defaultValue
}

func sendResponse(value any, valueType string) {
	resp := Response{
		Value: value,
		Type:  valueType,
	}
	json.NewEncoder(os.Stdout).Encode(resp)
}

func sendError(message string) {
	resp := Response{
		Error: message,
	}
	json.NewEncoder(os.Stdout).Encode(resp)
	os.Exit(1)
}

