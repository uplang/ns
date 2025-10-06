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
	case "upper":
		result, resultType, err = handleUpper(req.Params)
	case "lower":
		result, resultType, err = handleLower(req.Params)
	case "title":
		result, resultType, err = handleTitle(req.Params)
	case "trim":
		result, resultType, err = handleTrim(req.Params)
	case "trimPrefix":
		result, resultType, err = handleTrimPrefix(req.Params)
	case "trimSuffix":
		result, resultType, err = handleTrimSuffix(req.Params)
	case "split":
		result, resultType, err = handleSplit(req.Params)
	case "join":
		result, resultType, err = handleJoin(req.Params)
	case "replace":
		result, resultType, err = handleReplace(req.Params)
	case "replaceAll":
		result, resultType, err = handleReplaceAll(req.Params)
	case "contains":
		result, resultType, err = handleContains(req.Params)
	case "hasPrefix":
		result, resultType, err = handleHasPrefix(req.Params)
	case "hasSuffix":
		result, resultType, err = handleHasSuffix(req.Params)
	case "slice":
		result, resultType, err = handleSlice(req.Params)
	case "repeat":
		result, resultType, err = handleRepeat(req.Params)
	case "reverse":
		result, resultType, err = handleReverse(req.Params)
	case "length":
		result, resultType, err = handleLength(req.Params)
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

func handleUpper(params map[string]any) (any, string, error) {
	s := getString(params, "s", "")
	if s == "" {
		return nil, "", fmt.Errorf("s parameter required")
	}
	return strings.ToUpper(s), "string", nil
}

func handleLower(params map[string]any) (any, string, error) {
	s := getString(params, "s", "")
	if s == "" {
		return nil, "", fmt.Errorf("s parameter required")
	}
	return strings.ToLower(s), "string", nil
}

func handleTitle(params map[string]any) (any, string, error) {
	s := getString(params, "s", "")
	if s == "" {
		return nil, "", fmt.Errorf("s parameter required")
	}
	return strings.Title(s), "string", nil
}

func handleTrim(params map[string]any) (any, string, error) {
	s := getString(params, "s", "")
	if s == "" {
		return nil, "", fmt.Errorf("s parameter required")
	}
	cutset := getString(params, "cutset", " \t\n\r")
	return strings.Trim(s, cutset), "string", nil
}

func handleTrimPrefix(params map[string]any) (any, string, error) {
	s := getString(params, "s", "")
	if s == "" {
		return nil, "", fmt.Errorf("s parameter required")
	}
	prefix := getString(params, "prefix", "")
	if prefix == "" {
		return nil, "", fmt.Errorf("prefix parameter required")
	}
	return strings.TrimPrefix(s, prefix), "string", nil
}

func handleTrimSuffix(params map[string]any) (any, string, error) {
	s := getString(params, "s", "")
	if s == "" {
		return nil, "", fmt.Errorf("s parameter required")
	}
	suffix := getString(params, "suffix", "")
	if suffix == "" {
		return nil, "", fmt.Errorf("suffix parameter required")
	}
	return strings.TrimSuffix(s, suffix), "string", nil
}

func handleSplit(params map[string]any) (any, string, error) {
	s := getString(params, "s", "")
	if s == "" {
		return nil, "", fmt.Errorf("s parameter required")
	}
	sep := getString(params, "sep", ",")
	parts := strings.Split(s, sep)
	return parts, "list", nil
}

func handleJoin(params map[string]any) (any, string, error) {
	items, ok := params["items"].([]any)
	if !ok {
		return nil, "", fmt.Errorf("items parameter required and must be a list")
	}

	sep := getString(params, "sep", ",")

	strItems := make([]string, len(items))
	for i, item := range items {
		strItems[i] = fmt.Sprint(item)
	}

	return strings.Join(strItems, sep), "string", nil
}

func handleReplace(params map[string]any) (any, string, error) {
	s := getString(params, "s", "")
	if s == "" {
		return nil, "", fmt.Errorf("s parameter required")
	}
	old := getString(params, "old", "")
	if old == "" {
		return nil, "", fmt.Errorf("old parameter required")
	}
	new := getString(params, "new", "")

	n := getInt(params, "n", 1)
	return strings.Replace(s, old, new, n), "string", nil
}

func handleReplaceAll(params map[string]any) (any, string, error) {
	s := getString(params, "s", "")
	if s == "" {
		return nil, "", fmt.Errorf("s parameter required")
	}
	old := getString(params, "old", "")
	if old == "" {
		return nil, "", fmt.Errorf("old parameter required")
	}
	new := getString(params, "new", "")

	return strings.ReplaceAll(s, old, new), "string", nil
}

func handleContains(params map[string]any) (any, string, error) {
	s := getString(params, "s", "")
	if s == "" {
		return nil, "", fmt.Errorf("s parameter required")
	}
	substr := getString(params, "substr", "")
	if substr == "" {
		return nil, "", fmt.Errorf("substr parameter required")
	}
	return strings.Contains(s, substr), "bool", nil
}

func handleHasPrefix(params map[string]any) (any, string, error) {
	s := getString(params, "s", "")
	if s == "" {
		return nil, "", fmt.Errorf("s parameter required")
	}
	prefix := getString(params, "prefix", "")
	if prefix == "" {
		return nil, "", fmt.Errorf("prefix parameter required")
	}
	return strings.HasPrefix(s, prefix), "bool", nil
}

func handleHasSuffix(params map[string]any) (any, string, error) {
	s := getString(params, "s", "")
	if s == "" {
		return nil, "", fmt.Errorf("s parameter required")
	}
	suffix := getString(params, "suffix", "")
	if suffix == "" {
		return nil, "", fmt.Errorf("suffix parameter required")
	}
	return strings.HasSuffix(s, suffix), "bool", nil
}

func handleSlice(params map[string]any) (any, string, error) {
	s := getString(params, "s", "")
	if s == "" {
		return nil, "", fmt.Errorf("s parameter required")
	}

	start := getInt(params, "start", 0)
	end := getInt(params, "end", len(s))

	if start < 0 {
		start = 0
	}
	if end > len(s) {
		end = len(s)
	}
	if start > end {
		start = end
	}

	return s[start:end], "string", nil
}

func handleRepeat(params map[string]any) (any, string, error) {
	s := getString(params, "s", "")
	if s == "" {
		return nil, "", fmt.Errorf("s parameter required")
	}

	count := getInt(params, "count", 1)
	if count < 0 {
		return nil, "", fmt.Errorf("count must be non-negative")
	}
	if count > 10000 {
		return nil, "", fmt.Errorf("count too large (max 10000)")
	}

	return strings.Repeat(s, count), "string", nil
}

func handleReverse(params map[string]any) (any, string, error) {
	s := getString(params, "s", "")
	if s == "" {
		return nil, "", fmt.Errorf("s parameter required")
	}

	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes), "string", nil
}

func handleLength(params map[string]any) (any, string, error) {
	s := getString(params, "s", "")
	if s == "" {
		return nil, "", fmt.Errorf("s parameter required")
	}

	return len(s), "int", nil
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

