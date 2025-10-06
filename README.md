# UP Namespaces

[![CI](https://github.com/uplang/ns/actions/workflows/ci.yml/badge.svg)](https://github.com/uplang/ns/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/uplang/ns)](https://goreportcard.com/report/github.com/uplang/ns)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

Official namespace implementations for the UP language. All namespaces are in this single repository, each in its own directory.

**Quick Links:**
- 📦 [Quick Start](#quick-start)
- 📖 [Documentation](#documentation)
- 🔧 [Development](#creating-custom-namespaces)
- 🔗 [Go Packages](https://pkg.go.dev/github.com/uplang/ns)

## Repository Structure

```
ns/
├── time/          # Time manipulation functions
├── id/            # ID generation (UUID, ULID, etc.)
├── random/        # Random value generation
├── env/           # Environment variables
├── file/          # File operations
├── list/          # List generation and manipulation
├── math/          # Mathematical operations
├── string/        # String manipulation
├── fake/          # Fake data generation (faker)
└── greeting/      # Shell script example (template for custom namespaces)
```

## Documentation

- **[BUILTIN-NAMESPACES.md](BUILTIN-NAMESPACES.md)** - Documentation for all built-in namespaces
- **[DYNAMIC-NAMESPACES.md](DYNAMIC-NAMESPACES.md)** - Dynamic variables and list generation
- **[NAMESPACE-PLUGINS.md](NAMESPACE-PLUGINS.md)** - Creating custom namespaces
- **[NAMESPACE-SECURITY.md](NAMESPACE-SECURITY.md)** - Security considerations

## Available Namespaces

All namespaces are implemented and ready to use:

| Namespace | Functions | Description |
|-----------|-----------|-------------|
| **time** | 8 | Time manipulation and formatting |
| **id** | 4 | ID generation (UUID, ULID, nanoid, snowflake) |
| **random** | 7 | Random value generation |
| **env** | 5 | Environment variable access |
| **file** | 6 | File system operations |
| **list** | 6 | List generation with context |
| **math** | 13 | Mathematical operations |
| **string** | 10 | String manipulation |
| **fake** | 33 | Fake data generation (faker) |
| **greeting** | 3 | Example shell script namespace |

## Quick Start

### Using in UP Documents

```up
!use [time, id, random]

user {
  id $id.uuid
  created_at $time.now
  score!int $random.int(min=0, max=100)
}
```

### Building Namespaces

Build all namespaces:
```bash
# From the ns/ directory
for dir in time id random env file list math string fake; do
  (cd $dir && go build -o $dir)
done
```

Or build a specific namespace:
```bash
cd time
go build -o time
```

### Testing a Namespace

All namespaces follow the JSON stdin/stdout protocol:

```bash
echo '{"function":"now","params":{},"context":{}}' | ./time/time
```

Expected output:
```json
{"value":"2025-10-05T12:00:00Z","type":"ts"}
```

## Installation

Each namespace is a standalone Go module. To install as a command-line tool:

```bash
# Example: Install time namespace
cd time
go install
```

Or build and add to your PATH:
```bash
# Build all namespaces
for dir in time id random env file list math string fake; do
  (cd $dir && go build -o ../bin/$dir)
done

# Add to PATH
export PATH=$PATH:$(pwd)/bin
```

## Namespace Protocol

All namespaces implement a simple JSON protocol:

### Input (stdin)
```json
{
  "function": "now",
  "params": {
    "format": "2006-01-02"
  },
  "context": {
    "seed": 12345
  }
}
```

### Output (stdout)
```json
{
  "value": "2025-10-05",
  "type": "string"
}
```

### Error Response
```json
{
  "error": "Unknown function: invalid",
  "code": "UNKNOWN_FUNCTION"
}
```

## Creating Custom Namespaces

See **[greeting/](greeting/)** for a shell script example that demonstrates:
- ✅ ANY executable can be a namespace (shell, Python, Ruby, etc.)
- ✅ Simple JSON stdin/stdout protocol
- ✅ No complex APIs or SDKs required
- ✅ Easy to understand and extend

See **[NAMESPACE-PLUGINS.md](NAMESPACE-PLUGINS.md)** for complete documentation on creating custom namespaces.

### Go Namespace Template

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Request struct {
	Function string                 `json:"function"`
	Params   map[string]interface{} `json:"params"`
	Context  map[string]interface{} `json:"context"`
}

type Response struct {
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
	Error string      `json:"error,omitempty"`
	Code  string      `json:"code,omitempty"`
}

func main() {
	var req Request
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		sendError("Invalid JSON", "INVALID_REQUEST")
		return
	}

	switch req.Function {
	case "myfunction":
		result, err := handleMyFunction(req.Params)
		if err != nil {
			sendError(err.Error(), "FUNCTION_ERROR")
			return
		}
		sendResponse(result, "string")
	default:
		sendError(fmt.Sprintf("Unknown function: %s", req.Function), "UNKNOWN_FUNCTION")
	}
}

func handleMyFunction(params map[string]interface{}) (string, error) {
	// Implementation here
	return "result", nil
}

func sendResponse(value interface{}, valueType string) {
	json.NewEncoder(os.Stdout).Encode(Response{Value: value, Type: valueType})
}

func sendError(message, code string) {
	json.NewEncoder(os.Stdout).Encode(Response{Error: message, Code: code})
	os.Exit(1)
}
```

## Contributing

To add a new namespace:

1. **Create directory** with namespace name
2. **Initialize Go module** (if using Go)
   ```bash
   cd myspace
   go mod init github.com/uplang/ns/myspace
   ```
3. **Implement** the JSON protocol
4. **Create schema** file: `myspace.up-schema`
5. **Document** in `README.md`
6. **Test** with sample inputs
7. **Submit PR**

See existing namespaces for examples.

## Testing

Each namespace directory has its own tests. Run tests for a specific namespace:

```bash
cd time
go test ./...
```

Run tests for all namespaces:
```bash
for dir in time id random env file list math string fake; do
  echo "Testing $dir..."
  (cd $dir && go test ./...)
done
```

## License

MIT License - See [LICENSE](LICENSE) for details.

## Links

- 📖 **[UP Language Specification](https://github.com/uplang/spec)**
- 🔧 **[Go Implementation](https://github.com/uplang/go)**
- 📚 **[Namespace Documentation](BUILTIN-NAMESPACES.md)**
- 🔐 **[Security Guide](NAMESPACE-SECURITY.md)**
