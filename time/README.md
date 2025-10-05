# UP Namespace: time

Time manipulation and formatting functions for UP templates.

## Installation

```bash
go install github.com/uplang/ns-time@latest
```

Or build from source:

```bash
git clone https://github.com/uplang/ns-time.git
cd ns-time
go build -o time
```

## Usage in UP Templates

```up
!use [time]

# Current time
created_at $time.now

# Unix timestamp
timestamp!int $time.unix

# Format time
date $time.format(time="2025-10-05T12:00:00Z", format="2006-01-02")

# Add duration
future_time $time.add(duration="24h")

# Duration since
elapsed $time.since(time="2025-10-05T11:00:00Z")
```

## Functions

### `now(format?)`
Returns the current time.

**Parameters:**
- `format` (optional): Time format (default: RFC3339)

**Returns:** timestamp

**Example:**
```up
current $time.now
formatted $time.now(format="2006-01-02 15:04:05")
```

### `unix()`
Returns the current Unix timestamp in seconds.

**Returns:** int

**Example:**
```up
timestamp!int $time.unix
```

### `format(time, format?, input_format?)`
Formats a time string.

**Parameters:**
- `time` (required): Time to format
- `format` (optional): Output format (default: RFC3339)
- `input_format` (optional): Input format (default: RFC3339)

**Returns:** string

**Example:**
```up
date $time.format(time="2025-10-05T12:00:00Z", format="2006-01-02")
```

### `parse(time, format?)`
Parses a time string.

**Parameters:**
- `time` (required): Time string to parse
- `format` (optional): Input format (default: RFC3339)

**Returns:** timestamp

**Example:**
```up
parsed $time.parse(time="2025-10-05", format="2006-01-02")
```

### `add(time?, duration)`
Adds duration to a time.

**Parameters:**
- `time` (optional): Time to add to (default: now)
- `duration` (required): Duration to add (e.g., "1h", "30m", "24h")

**Returns:** timestamp

**Example:**
```up
tomorrow $time.add(duration="24h")
next_week $time.add(time="2025-10-05T12:00:00Z", duration="168h")
```

### `sub(time?, duration)`
Subtracts duration from a time.

**Parameters:**
- `time` (optional): Time to subtract from (default: now)
- `duration` (required): Duration to subtract

**Returns:** timestamp

**Example:**
```up
yesterday $time.sub(duration="24h")
```

### `since(time)`
Returns duration since a time.

**Parameters:**
- `time` (required): Reference time

**Returns:** duration

**Example:**
```up
elapsed $time.since(time="2025-10-05T11:00:00Z")
```

### `until(time)`
Returns duration until a time.

**Parameters:**
- `time` (required): Reference time

**Returns:** duration

**Example:**
```up
remaining $time.until(time="2025-12-31T23:59:59Z")
```

## Time Formats

Common Go time format patterns:

- `RFC3339`: `2006-01-02T15:04:05Z07:00` (default)
- `2006-01-02`: Date only
- `15:04:05`: Time only
- `2006-01-02 15:04:05`: Date and time
- `Mon Jan 2 15:04:05 MST 2006`: Full format

## Testing

Test the namespace directly:

```bash
echo '{"function":"now","params":{},"context":{}}' | ./time
```

Expected output:
```json
{"value":"2025-10-05T12:00:00Z","type":"ts"}
```

## Protocol

This namespace follows the UP namespace plugin protocol:

**Input (stdin):**
```json
{
  "function": "now",
  "params": {},
  "context": {}
}
```

**Output (stdout):**
```json
{
  "value": "2025-10-05T12:00:00Z",
  "type": "ts"
}
```

## License

MIT License - see [LICENSE](../LICENSE) for details.

## Links

- [UP Language](https://uplang.org)
- [UP Specification](https://github.com/uplang/spec)
- [Namespace Documentation](https://github.com/uplang/spec/blob/main/NAMESPACE-PLUGINS.md)

