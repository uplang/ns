package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/google/uuid"
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
	case "uuid", "uuid4":
		result, resultType, err = handleUUID()
	case "ulid":
		result, resultType, err = handleULID()
	case "nanoid":
		result, resultType, err = handleNanoID(req.Params)
	case "snowflake":
		result, resultType, err = handleSnowflake(req.Params)
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

func handleUUID() (any, string, error) {
	return uuid.New().String(), "uuid", nil
}

func handleULID() (any, string, error) {
	// Simple ULID-like implementation (timestamp + random)
	// In production, use github.com/oklog/ulid
	return fmt.Sprintf("%013x%013x", getTimestamp(), getRandomHex(13)), "string", nil
}

func handleNanoID(params map[string]any) (any, string, error) {
	size := getInt(params, "size", 21)
	alphabet := getString(params, "alphabet", "_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	result := make([]byte, size)
	for i := 0; i < size; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return nil, "", err
		}
		result[i] = alphabet[n.Int64()]
	}

	return string(result), "string", nil
}

func handleSnowflake(params map[string]any) (any, string, error) {
	// Simplified Snowflake ID (timestamp + worker + sequence)
	timestamp := getTimestamp()
	worker := getInt64(params, "worker", 0) & 0x3FF                        // 10 bits
	sequence := getInt64(params, "sequence", getRandomInt64(4096)) & 0xFFF // 12 bits

	id := (timestamp << 22) | (worker << 12) | sequence
	return id, "int", nil
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

func getTimestamp() int64 {
	return time.Now().UnixMilli()
}

func getRandomHex(n int) int64 {
	max := big.NewInt(1 << uint(n*4))
	r, _ := rand.Int(rand.Reader, max)
	return r.Int64()
}

func getRandomInt64(max int64) int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(max))
	return n.Int64()
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
