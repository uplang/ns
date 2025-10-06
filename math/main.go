package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
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
	case "add":
		result, resultType, err = handleAdd(req.Params)
	case "sub":
		result, resultType, err = handleSub(req.Params)
	case "mul":
		result, resultType, err = handleMul(req.Params)
	case "div":
		result, resultType, err = handleDiv(req.Params)
	case "mod":
		result, resultType, err = handleMod(req.Params)
	case "pow":
		result, resultType, err = handlePow(req.Params)
	case "sqrt":
		result, resultType, err = handleSqrt(req.Params)
	case "abs":
		result, resultType, err = handleAbs(req.Params)
	case "min":
		result, resultType, err = handleMin(req.Params)
	case "max":
		result, resultType, err = handleMax(req.Params)
	case "ceil":
		result, resultType, err = handleCeil(req.Params)
	case "floor":
		result, resultType, err = handleFloor(req.Params)
	case "round":
		result, resultType, err = handleRound(req.Params)
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

func handleAdd(params map[string]any) (any, string, error) {
	a := getFloat64(params, "a", 0)
	b := getFloat64(params, "b", 0)
	return a + b, "float", nil
}

func handleSub(params map[string]any) (any, string, error) {
	a := getFloat64(params, "a", 0)
	b := getFloat64(params, "b", 0)
	return a - b, "float", nil
}

func handleMul(params map[string]any) (any, string, error) {
	a := getFloat64(params, "a", 0)
	b := getFloat64(params, "b", 0)
	return a * b, "float", nil
}

func handleDiv(params map[string]any) (any, string, error) {
	a := getFloat64(params, "a", 0)
	b := getFloat64(params, "b", 0)
	if b == 0 {
		return nil, "", fmt.Errorf("division by zero")
	}
	return a / b, "float", nil
}

func handleMod(params map[string]any) (any, string, error) {
	a := getInt64(params, "a", 0)
	b := getInt64(params, "b", 0)
	if b == 0 {
		return nil, "", fmt.Errorf("modulo by zero")
	}
	return a % b, "int", nil
}

func handlePow(params map[string]any) (any, string, error) {
	base := getFloat64(params, "base", 0)
	exponent := getFloat64(params, "exponent", 0)
	return math.Pow(base, exponent), "float", nil
}

func handleSqrt(params map[string]any) (any, string, error) {
	x := getFloat64(params, "x", 0)
	if x < 0 {
		return nil, "", fmt.Errorf("cannot take square root of negative number")
	}
	return math.Sqrt(x), "float", nil
}

func handleAbs(params map[string]any) (any, string, error) {
	x := getFloat64(params, "x", 0)
	return math.Abs(x), "float", nil
}

func handleMin(params map[string]any) (any, string, error) {
	values, ok := params["values"].([]any)
	if !ok || len(values) == 0 {
		return nil, "", fmt.Errorf("values parameter required and must be non-empty list")
	}

	min := toFloat64(values[0])
	for _, v := range values[1:] {
		val := toFloat64(v)
		if val < min {
			min = val
		}
	}

	return min, "float", nil
}

func handleMax(params map[string]any) (any, string, error) {
	values, ok := params["values"].([]any)
	if !ok || len(values) == 0 {
		return nil, "", fmt.Errorf("values parameter required and must be non-empty list")
	}

	max := toFloat64(values[0])
	for _, v := range values[1:] {
		val := toFloat64(v)
		if val > max {
			max = val
		}
	}

	return max, "float", nil
}

func handleCeil(params map[string]any) (any, string, error) {
	x := getFloat64(params, "x", 0)
	return math.Ceil(x), "float", nil
}

func handleFloor(params map[string]any) (any, string, error) {
	x := getFloat64(params, "x", 0)
	return math.Floor(x), "float", nil
}

func handleRound(params map[string]any) (any, string, error) {
	x := getFloat64(params, "x", 0)
	return math.Round(x), "float", nil
}

func getFloat64(params map[string]any, key string, defaultValue float64) float64 {
	if v, ok := params[key]; ok {
		return toFloat64(v)
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

func toFloat64(v any) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	case int64:
		return float64(val)
	default:
		return 0
	}
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
