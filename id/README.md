# UP Namespace: id

ID generation functions for UP templates.

## Installation

```bash
go install github.com/uplang/ns-id@latest
```

Or build from source:

```bash
cd ns/id
go build -o id
```

## Usage in UP Templates

```up
!use [id]

# Generate UUID
user_id $id.uuid

# Generate ULID (sortable)
request_id $id.ulid

# Generate NanoID (compact)
session_id $id.nanoid
short_id $id.nanoid(size=10)

# Generate Snowflake ID (distributed)
entity_id!int $id.snowflake(machine_id=5)
```

## Functions

### `uuid()`
Generates a UUID v4.

**Returns:** string (36 characters)

**Example:**
```up
id $id.uuid
# Result: 550e8400-e29b-41d4-a716-446655440000
```

### `ulid()`
Generates a ULID (Universally Unique Lexicographically Sortable Identifier).

**Returns:** string (26 characters)

**Features:**
- Timestamp-based sorting
- Case insensitive
- URL safe
- Compact

**Example:**
```up
id $id.ulid
# Result: 01ARZ3NDEKTSV4RRFFQ69G5FAV
```

### `nanoid(size?, alphabet?)`
Generates a NanoID.

**Parameters:**
- `size` (int, optional): Length (default: 21)
- `alphabet` (string, optional): Custom alphabet

**Returns:** string

**Examples:**
```up
id $id.nanoid
# Result: V1StGXR8_Z5jdHi6B-myT

short_id $id.nanoid(size=10)
# Result: V1StGXR8_Z

numeric_id $id.nanoid(size=12, alphabet="0123456789")
# Result: 123456789012
```

### `snowflake(machine_id?)`
Generates a Twitter Snowflake ID.

**Parameters:**
- `machine_id` (int, optional): Machine/node ID 0-1023 (default: 1)

**Returns:** int (64-bit)

**Features:**
- Time-ordered
- Decentralized
- No coordination needed
- High performance

**Example:**
```up
id!int $id.snowflake
# Result: 1450617600000000001

distributed_id!int $id.snowflake(machine_id=42)
# Result: 1450617600000042001
```

## Use Cases

- **UUID**: Standard unique identifiers, database keys
- **ULID**: Sortable IDs, log entries, time-series data
- **NanoID**: Compact IDs, URLs, short codes
- **Snowflake**: Distributed systems, high-throughput IDs

## Testing

Test directly:

```bash
echo '{"function":"uuid","params":{},"context":{}}' | ./id
echo '{"function":"nanoid","params":{"size":10},"context":{}}' | ./id
```

## Protocol

Follows the UP namespace plugin protocol (JSON stdin/stdout).

## License

MIT License

## Links

- [UP Language](https://uplang.org)
- [UP Specification](https://github.com/uplang/spec)
- [Schema](id.up-schema)
