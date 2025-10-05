# Greeting Namespace - Shell Script Example

**⭐ This is a reference implementation showing how ANY executable can become an UP namespace!**

This directory contains a **shell script example** demonstrating UP's pluggable namespace system. The `greeting` namespace is a simple bash script that shows how to implement the JSON stdin/stdout protocol that powers UP's dynamic namespaces.

## Why This Example Matters

This example demonstrates that UP namespaces are **language-agnostic**:
- ✅ **ANY** executable can be a namespace (shell scripts, Python, Go, Rust, Ruby, etc.)
- ✅ Simple **JSON-based protocol** (read from stdin, write to stdout)
- ✅ No complex APIs or SDKs required
- ✅ Easy to understand and extend

**Use this as a template** for creating your own custom namespaces!

## What Are Namespace Plugins?

Namespace plugins are **executables** (shell scripts, binaries, programs in any language) that provide dynamic functionality to UP documents. They follow a simple JSON-based protocol for communication.

### Files

- **`greeting`** - Executable bash script (the plugin)
- **`greeting.up-schema`** - UP schema describing the namespace

### Testing the Plugin

```bash
# Test hello function
echo '{"function":"hello","params":{"positional":["Alice"]}}' | ./greeting

# Test goodbye function
echo '{"function":"goodbye","params":{"positional":["Bob"]}}' | ./greeting

# Test with named parameters
echo '{"function":"hello","params":{"named":{"name":"Charlie","excited":"false"}}}' | ./greeting
```

### Using in UP

```up
!use [greeting]

# Call functions
message1 $greeting.hello(Alice)
message2 $greeting.goodbye(Bob)
message3 $greeting.wave(Charlie)
```

## The JSON stdin/stdout Protocol

**This is the magic!** The protocol is dead simple:

1. **Read JSON from stdin** - UP sends function calls as JSON
2. **Process the request** - Your executable does whatever it needs to do
3. **Write JSON to stdout** - Return the result as JSON
4. **Exit** - That's it!

No libraries required. No frameworks. Just JSON in, JSON out.

### Input (stdin)

JSON with function name and parameters:

```json
{
  "function": "hello",
  "params": {
    "positional": ["Alice"],
    "named": {
      "excited": "true"
    }
  },
  "context": {
    "seed": 12345,
    "file": "/path/to/file.up",
    "line": 42
  }
}
```

### Output (stdout)

JSON with result:

```json
{
  "value": "Hello, Alice!",
  "type": "string"
}
```

Or for errors:

```json
{
  "error": "Unknown function: invalid",
  "code": "UNKNOWN_FUNCTION"
}
```

## Schema Files

Each plugin should have a corresponding `.up-schema` file using UP syntax to describe:

- Available functions
- Parameters and their types
- Return types
- Examples
- Metadata (author, version, requirements)

See `greeting.up-schema` for a complete example.

## Creating Your Own Plugin

1. **Create executable** - Any language, just needs to be executable
2. **Implement protocol** - Read JSON from stdin, write JSON to stdout
3. **Write schema** - Document your namespace in `.up-schema` file
4. **Make executable** - `chmod +x your-plugin`
5. **Place in directory** - `./up-namespaces/`, `~/.up/namespaces/`, or `/usr/local/up/namespaces/`

### Minimal Example

```bash
#!/bin/bash
# my-plugin

INPUT=$(cat)
FUNCTION=$(echo "$INPUT" | jq -r '.function')

case "$FUNCTION" in
  greet)
    echo '{"value": "Hello from my plugin!", "type": "string"}'
    ;;
  *)
    echo '{"error": "Unknown function", "code": "UNKNOWN_FUNCTION"}'
    exit 1
    ;;
esac
```

## Plugin Locations

UP searches for plugins in this order:

1. `./up-namespaces/{namespace}` - Project-local
2. `~/.up/namespaces/{namespace}` - User-installed
3. `/usr/local/up/namespaces/{namespace}` - System-wide
4. Built-in - Internal implementation (time, date, id, random, faker)

First match wins.

## Security

- Plugins run as separate processes
- Can be sandboxed
- Dangerous namespaces (file, exec, env) should require explicit opt-in
- Review plugin code before running
- Use schemas to understand what a plugin does

## More Information

See **[NAMESPACE-PLUGINS.md](../NAMESPACE-PLUGINS.md)** for complete documentation including:
- Detailed protocol specification
- Examples in multiple languages (Python, Go, Rust)
- Schema system details
- Security considerations
- CLI tools for managing plugins

## Contributing

To contribute a plugin:

1. Create plugin and schema
2. Test thoroughly
3. Document in schema
4. Submit PR with both files
5. Include examples and tests

Plugins should be:
- Well-documented with schemas
- Safe by default
- Follow the JSON protocol
- Include error handling
- Have clear examples

