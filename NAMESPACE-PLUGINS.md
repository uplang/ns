# UP Namespace Plugins

UP's dynamic namespace system is **pluggable** - any executable (shell script, binary, program) can provide a namespace. Namespaces are documented using **UP schemas** - UP describing UP.

## Core Concept

```up
!use [myns]

# Calls executable: ./up-namespaces/myns or ~/.up/namespaces/myns
value $myns.generate(param1, param2)
```

The executable receives the function call, processes it, and returns the value. Simple, universal, language-agnostic.

## Plugin Locations

UP searches for namespace executables in order:

1. **Local project**: `./up-namespaces/{namespace}`
2. **User directory**: `~/.up/namespaces/{namespace}`
3. **System directory**: `/usr/local/up/namespaces/{namespace}`
4. **Built-in**: Internal implementation (time, date, id, random, faker, env, file)

First match wins.

## Executable Protocol

### Input (stdin)

The executable receives **JSON** on stdin with the function call:

```json
{
  "function": "generate",
  "params": {
    "positional": ["value1", "value2"],
    "named": {
      "key": "value"
    }
  },
  "context": {
    "seed": 12345,
    "file": "/path/to/document.up",
    "line": 42
  }
}
```

### Output (stdout)

The executable returns **JSON** on stdout with the result:

```json
{
  "value": "generated content",
  "type": "string"
}
```

Or for errors:

```json
{
  "error": "Function not found: invalid_function",
  "code": "FUNCTION_NOT_FOUND"
}
```

### Exit Codes

- `0` - Success
- `1` - Error (with error JSON on stdout)
- `2` - Invalid input

## Simple Shell Script Example

**File: `./up-namespaces/greeting`**

```bash
#!/bin/bash

# Read input
INPUT=$(cat)

# Parse function name (using jq)
FUNCTION=$(echo "$INPUT" | jq -r '.function')
PARAMS=$(echo "$INPUT" | jq -r '.params.positional[0] // "World"')

case "$FUNCTION" in
  hello)
    echo "{\"value\": \"Hello, $PARAMS!\", \"type\": \"string\"}"
    ;;
  goodbye)
    echo "{\"value\": \"Goodbye, $PARAMS!\", \"type\": \"string\"}"
    ;;
  *)
    echo "{\"error\": \"Unknown function: $FUNCTION\", \"code\": \"UNKNOWN_FUNCTION\"}"
    exit 1
    ;;
esac
```

Make it executable:
```bash
chmod +x ./up-namespaces/greeting
```

**Usage in UP:**

```up
!use [greeting]

message $greeting.hello(Alice)
farewell $greeting.goodbye(Bob)
```

**Output:**
```up
message Hello, Alice!
farewell Goodbye, Bob!
```

## UP Schema System

Each namespace is documented using a **`.up-schema`** file that describes its interface.

### Schema Location

Schemas are co-located with executables:

- `./up-namespaces/greeting` (executable)
- `./up-namespaces/greeting.up-schema` (schema)

### Schema Syntax

UP schemas use UP syntax to describe namespaces:

**File: `./up-namespaces/greeting.up-schema`**

```up
# UP Namespace Schema
# Describes the greeting namespace

namespace greeting

version 1.0.0
description Simple greeting generator for testing

# Define available functions
functions {
  hello {
    description Generate a hello greeting

    params {
      name!string {
        description Name to greet
        default World
        required!bool false
      }

      excited!bool {
        description Add exclamation
        default true
        required!bool false
      }
    }

    returns!string {
      description Greeting message
      example Hello, Alice!
    }

    examples [
      {
        call $greeting.hello(Alice)
        result Hello, Alice!
      }
      {
        call $greeting.hello(Bob, excited=false)
        result Hello, Bob
      }
    ]
  }

  goodbye {
    description Generate a goodbye message

    params {
      name!string {
        description Name to say goodbye to
        required!bool true
      }
    }

    returns!string {
      description Farewell message
    }

    examples [
      {
        call $greeting.goodbye(Alice)
        result Goodbye, Alice!
      }
    ]
  }
}

# Namespace metadata
metadata {
  author John Doe
  repository https://github.com/example/up-greeting
  license MIT

  tags [greeting, test, example]

  requires {
    up_version >= 1.0.0
    executables [bash, jq]
  }
}
```

## Built-in Namespace Schemas

Even built-in namespaces are documented with schemas:

**File: `~/.up/namespaces/time.up-schema`**

```up
namespace time

version 1.0.0
description Time and timestamp operations

functions {
  now {
    description Current timestamp

    params {
      format!string {
        description Time format
        default rfc3339
        options [rfc3339, iso8601, unix, unix_ms]
        required!bool false
      }
    }

    returns!string {
      description Formatted timestamp
      example 2025-10-05T14:30:00Z
    }

    examples [
      {
        call $time.now
        result 2025-10-05T14:30:00Z
      }
      {
        call $time.now(unix)
        result 1728139800
      }
    ]
  }

  add {
    description Add duration to current time

    params {
      duration!dur {
        description Duration to add (can be negative)
        required!bool true
        example 24h
      }

      format!string {
        description Output format
        default rfc3339
        required!bool false
      }
    }

    returns!string {
      description Future/past timestamp
    }

    examples [
      {
        call $time.add(24h)
        result 2025-10-06T14:30:00Z
      }
      {
        call $time.add(-7d)
        result 2025-09-28T14:30:00Z
      }
    ]
  }

  format {
    description Format current time with custom format

    params {
      pattern!string {
        description Go time format pattern
        required!bool true
        example 2006-01-02 15:04:05
      }
    }

    returns!string {
      description Formatted time
    }
  }
}

metadata {
  author UP Core Team
  builtin!bool true
  safe!bool true
}
```

## Complex Example: Random Namespace

**File: `./up-namespaces/random.up-schema`**

```up
namespace random

version 1.0.0
description Random data generation for testing

functions {
  int {
    description Generate random integer

    params {
      min!int {
        description Minimum value (inclusive)
        default 0
        required!bool false
      }

      max!int {
        description Maximum value (exclusive)
        required!bool true
      }
    }

    returns!int {
      description Random integer
    }

    examples [
      {
        call $random.int(100)
        result 42
        note Returns integer between 0 and 99
      }
      {
        call $random.int(10, 20)
        result 15
        note Returns integer between 10 and 19
      }
    ]
  }

  float {
    description Generate random floating-point number

    params {
      min!float {
        description Minimum value
        default 0.0
        required!bool false
      }

      max!float {
        description Maximum value
        required!bool true
      }

      precision!int {
        description Decimal places
        default 2
        required!bool false
      }
    }

    returns!float {
      description Random float
    }
  }

  string {
    description Generate random string

    params {
      length!int {
        description String length
        required!bool true
      }

      charset!string {
        description Character set to use
        default alphanumeric
        options [alphanumeric, alpha, numeric, hex, base64]
        required!bool false
      }
    }

    returns!string {
      description Random string
    }

    examples [
      {
        call $random.string(16)
        result aB3xK9mL2pQr8sT1
      }
      {
        call $random.string(8, charset=hex)
        result a3b8c9d2
      }
    ]
  }

  bool {
    description Generate random boolean

    params {
      probability!float {
        description Probability of true (0.0 to 1.0)
        default 0.5
        required!bool false
      }
    }

    returns!bool {
      description Random boolean
    }
  }

  choice {
    description Pick random item from list

    params {
      items!list {
        description List to choose from
        required!bool true
      }
    }

    returns {
      description Random item from list
      type any
    }

    examples [
      {
        call $random.choice([red, green, blue])
        result green
      }
    ]
  }
}

metadata {
  author UP Core Team
  builtin!bool true
  safe!bool true
  seeded!bool true
  note Random functions respect the --seed flag for reproducibility
}
```

## Schema Validation

UP schemas are themselves UP documents with a specific structure. The meta-schema:

**File: `namespace-schema.up-schema`** (meta!)

```up
# Meta-schema: Schema for namespace schemas
namespace namespace-schema

version 1.0.0
description Schema for describing UP namespace plugins

schema {
  namespace!string {
    description Namespace identifier
    required!bool true
    pattern ^[a-z][a-z0-9_]*$
  }

  version!string {
    description Semantic version
    required!bool true
    pattern ^\d+\.\d+\.\d+$
  }

  description!string {
    description Human-readable description
    required!bool true
  }

  functions!block {
    description Available functions
    required!bool true

    schema {
      function_name {
        description!string Function description

        params!block {
          description Parameters accepted by function

          schema {
            param_name {
              description!string Parameter description
              type!string Parameter type
              default Value if not provided
              required!bool Is parameter required
              options!list Valid values
              pattern!string Validation regex
            }
          }
        }

        returns {
          description!string Return value description
          type!string Return type
          example Example return value
        }

        examples!list {
          description Usage examples

          schema {
            call!string Function call syntax
            result Value returned
            note!string Additional notes
          }
        }
      }
    }
  }

  metadata!block {
    description Namespace metadata

    schema {
      author!string Schema author
      repository!string Source repository
      license!string License type
      builtin!bool Built-in namespace
      safe!bool Safe for production
      tags!list Searchable tags

      requires!block {
        description Requirements

        schema {
          up_version!string Minimum UP version
          executables!list Required executables
          env_vars!list Required environment variables
        }
      }
    }
  }
}
```

## Python Example

**File: `./up-namespaces/math`**

```python
#!/usr/bin/env python3

import json
import sys
import math

def main():
    # Read input
    input_data = json.load(sys.stdin)

    function = input_data.get('function')
    params = input_data.get('params', {})
    positional = params.get('positional', [])
    named = params.get('named', {})

    try:
        if function == 'sqrt':
            value = float(positional[0])
            result = math.sqrt(value)
            print(json.dumps({"value": result, "type": "float"}))

        elif function == 'pow':
            base = float(positional[0])
            exp = float(positional[1]) if len(positional) > 1 else 2.0
            result = math.pow(base, exp)
            print(json.dumps({"value": result, "type": "float"}))

        elif function == 'abs':
            value = float(positional[0])
            result = abs(value)
            print(json.dumps({"value": result, "type": "float"}))

        else:
            print(json.dumps({
                "error": f"Unknown function: {function}",
                "code": "UNKNOWN_FUNCTION"
            }))
            sys.exit(1)

    except Exception as e:
        print(json.dumps({
            "error": str(e),
            "code": "EXECUTION_ERROR"
        }))
        sys.exit(1)

if __name__ == '__main__':
    main()
```

**Schema: `./up-namespaces/math.up-schema`**

```up
namespace math

version 1.0.0
description Mathematical operations

functions {
  sqrt {
    description Square root

    params {
      value!float {
        description Number to take square root of
        required!bool true
      }
    }

    returns!float {
      description Square root result
    }

    examples [
      {
        call $math.sqrt(16)
        result 4.0
      }
    ]
  }

  pow {
    description Power/exponentiation

    params {
      base!float {
        description Base number
        required!bool true
      }

      exponent!float {
        description Exponent
        default 2.0
        required!bool false
      }
    }

    returns!float {
      description Result of base^exponent
    }

    examples [
      {
        call $math.pow(2, 8)
        result 256.0
      }
      {
        call $math.pow(5)
        result 25.0
        note Defaults to squaring
      }
    ]
  }

  abs {
    description Absolute value

    params {
      value!float {
        description Number
        required!bool true
      }
    }

    returns!float {
      description Absolute value
    }
  }
}

metadata {
  author Example Team
  license MIT
  safe!bool true

  requires {
    executables [python3]
  }
}
```

## Go Binary Example

**File: `./up-namespaces/hash/main.go`**

```go
package main

import (
    "crypto/md5"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "os"
)

type Input struct {
    Function string `json:"function"`
    Params   struct {
        Positional []string          `json:"positional"`
        Named      map[string]string `json:"named"`
    } `json:"params"`
}

type Output struct {
    Value string `json:"value"`
    Type  string `json:"type"`
}

type ErrorOutput struct {
    Error string `json:"error"`
    Code  string `json:"code"`
}

func main() {
    var input Input
    if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
        sendError("Invalid input JSON", "INVALID_INPUT")
        return
    }

    if len(input.Params.Positional) == 0 {
        sendError("Missing input string", "MISSING_PARAM")
        return
    }

    data := input.Params.Positional[0]

    switch input.Function {
    case "md5":
        hash := md5.Sum([]byte(data))
        sendResult(hex.EncodeToString(hash[:]))

    case "sha256":
        hash := sha256.Sum256([]byte(data))
        sendResult(hex.EncodeToString(hash[:]))

    default:
        sendError(fmt.Sprintf("Unknown function: %s", input.Function), "UNKNOWN_FUNCTION")
    }
}

func sendResult(value string) {
    json.NewEncoder(os.Stdout).Encode(Output{
        Value: value,
        Type:  "string",
    })
}

func sendError(message, code string) {
    json.NewEncoder(os.Stdout).Encode(ErrorOutput{
        Error: message,
        Code:  code,
    })
    os.Exit(1)
}
```

Build and install:
```bash
go build -o ~/.up/namespaces/hash ./up-namespaces/hash/
```

**Usage:**

```up
!use [hash]

config_hash $hash.md5(my-config-data)
checksum $hash.sha256($vars.content)
```

## CLI Tools for Schemas

```bash
# List available namespaces
up namespace list

# Show namespace schema
up namespace show time

# Validate a namespace schema
up namespace validate ./up-namespaces/greeting.up-schema

# Test a namespace plugin
up namespace test greeting --function hello --params "Alice"

# Install namespace from repository
up namespace install https://github.com/example/up-weather

# Create namespace template
up namespace init myns
```

## Discovery and Registration

**User namespace registry: `~/.up/namespaces/registry.up`**

```up
# Local namespace registry

namespaces {
  greeting {
    path ~/.up/namespaces/greeting
    schema ~/.up/namespaces/greeting.up-schema
    version 1.0.0
    enabled!bool true
  }

  math {
    path ~/.up/namespaces/math
    schema ~/.up/namespaces/math.up-schema
    version 1.0.0
    enabled!bool true
  }

  weather {
    path ~/.up/namespaces/weather
    schema ~/.up/namespaces/weather.up-schema
    version 2.1.0
    enabled!bool true
    requires {
      api_key!env WEATHER_API_KEY
    }
  }
}
```

## Security Model

### Sandboxing

Namespaces should run in a sandbox:

```up
# Namespace security policy
namespace_policy {
  greeting {
    allowed!bool true
    sandbox!bool false
    safe!bool true
  }

  file {
    allowed!bool true
    sandbox!bool true
    safe!bool false

    restrictions {
      read_only!bool true
      allowed_paths [
        ./config
        ./data
      ]
      blocked_paths [
        ~/.ssh
        /etc
      ]
    }
  }

  exec {
    allowed!bool false
    safe!bool false
    note Dangerous - disabled by default
  }
}
```

### Capability Declarations

Namespaces declare what they need:

```up
# In namespace schema
metadata {
  capabilities {
    network!bool false
    filesystem!bool true
    environment!bool false

    filesystem_access {
      read!list [./config, ./data]
      write!list []
    }
  }
}
```

## Performance Considerations

### Caching

```up
# Namespace calls can be cached
!use [expensive]
!cache expensive.compute

# First call executes
result1 $expensive.compute(input)

# Second call with same input returns cached value
result2 $expensive.compute(input)
```

### Batch Operations

```up
# Multiple calls can be batched
!use [api]
!batch api

users [
  $api.fetch_user(1)
  $api.fetch_user(2)
  $api.fetch_user(3)
]

# Sent as single batch request to namespace
```

## Security and Versioning

Namespace plugins require strong security guarantees. See **[NAMESPACE-SECURITY.md](NAMESPACE-SECURITY.md)** for complete documentation on:

- **Version pinning** - `!use [greeting@1.2.3]` or `!use [greeting^1.2.0]`
- **Hash verification** - SHA-256 integrity checking via lock files
- **Cryptographic signatures** - Ed25519/RSA signing for authenticity
- **Trusted signers** - Public key management and trust levels
- **Security policies** - `.up-security` configuration files
- **Audit logging** - Track security events
- **Registry integration** - Verified distribution

**Quick example:**

```up
# Exact version pinning
!use [greeting@1.2.3, faker@3.1.4]

# Lock file ensures reproducibility
# up-namespaces.lock contains hashes
# Verification happens automatically
```

## Summary

**Plugin System:**
- ✅ Any executable can be a namespace
- ✅ Simple JSON protocol (stdin/stdout)
- ✅ Language-agnostic (bash, python, go, rust, etc.)
- ✅ Located in standard directories

**Schema System:**
- ✅ Schemas written in UP syntax
- ✅ Self-documenting namespaces
- ✅ Validates params and returns
- ✅ Provides examples and metadata

**Security:**
- 🔒 Cryptographic verification
- 📌 Version pinning and lock files
- 🛡️ Signature verification
- 🔍 Audit logging
- ⚙️ Flexible security policies

**Benefits:**
- 🚀 Extensible without modifying UP core
- 📚 Self-documenting with `.up-schema` files
- 🔒 Secure and reproducible
- 🌍 Language-agnostic plugins
- 🎯 Simple shell scripts to complex services

**UP describes UP - schemas all the way down!**

