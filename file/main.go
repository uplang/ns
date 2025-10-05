package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	case "read":
		result, resultType, err = handleRead(req.Params)
	case "exists":
		result, resultType, err = handleExists(req.Params)
	case "list":
		result, resultType, err = handleList(req.Params)
	case "basename":
		result, resultType, err = handleBasename(req.Params)
	case "dirname":
		result, resultType, err = handleDirname(req.Params)
	case "ext":
		result, resultType, err = handleExt(req.Params)
	case "join":
		result, resultType, err = handleJoin(req.Params)
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

func handleRead(params map[string]any) (any, string, error) {
	path := getString(params, "path", "")
	if path == "" {
		return nil, "", fmt.Errorf("path parameter required")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read file: %v", err)
	}

	return string(data), "string", nil
}

func handleExists(params map[string]any) (any, string, error) {
	path := getString(params, "path", "")
	if path == "" {
		return nil, "", fmt.Errorf("path parameter required")
	}

	_, err := os.Stat(path)
	return err == nil, "bool", nil
}

func handleList(params map[string]any) (any, string, error) {
	dir := getString(params, "dir", ".")
	pattern := getString(params, "pattern", "*")

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read directory: %v", err)
	}

	var result []string
	for _, entry := range entries {
		name := entry.Name()
		if matched, _ := filepath.Match(pattern, name); matched {
			result = append(result, name)
		}
	}

	return result, "list", nil
}

func handleBasename(params map[string]any) (any, string, error) {
	path := getString(params, "path", "")
	if path == "" {
		return nil, "", fmt.Errorf("path parameter required")
	}

	return filepath.Base(path), "string", nil
}

func handleDirname(params map[string]any) (any, string, error) {
	path := getString(params, "path", "")
	if path == "" {
		return nil, "", fmt.Errorf("path parameter required")
	}

	return filepath.Dir(path), "string", nil
}

func handleExt(params map[string]any) (any, string, error) {
	path := getString(params, "path", "")
	if path == "" {
		return nil, "", fmt.Errorf("path parameter required")
	}

	ext := filepath.Ext(path)
	// Remove leading dot
	if strings.HasPrefix(ext, ".") {
		ext = ext[1:]
	}

	return ext, "string", nil
}

func handleJoin(params map[string]any) (any, string, error) {
	parts, ok := params["parts"].([]any)
	if !ok || len(parts) == 0 {
		return nil, "", fmt.Errorf("parts parameter required and must be non-empty list")
	}

	strParts := make([]string, len(parts))
	for i, p := range parts {
		strParts[i] = fmt.Sprint(p)
	}

	return filepath.Join(strParts...), "string", nil
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
	json.NewEncoder(os.Stdout).Encode(resp)
}

func sendError(message string) {
	resp := Response{
		Error: message,
	}
	json.NewEncoder(os.Stdout).Encode(resp)
	os.Exit(1)
}

