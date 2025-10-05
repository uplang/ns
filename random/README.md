# UP Namespace: random

Random value generation functions for UP templates.

## Installation

```bash
go install github.com/uplang/ns-random@latest
```

Or build from source:

```bash
cd ns/random
go build -o random
```

## Usage in UP Templates

```up
!use [random]

# Random integers
score!int $random.int(min=0, max=100)
dice!int $random.int(min=1, max=6)

# Random floats
temperature!float $random.float(min=-10.0, max=40.0)
probability!float $random.float

# Random boolean
enabled!bool $random.bool

# Random choice
color $random.choice(items=["red", "green", "blue"])

# Random bytes
token $random.bytes(size=32)
```

## Functions

### `int(min?, max?)`
Generates a random integer.

**Parameters:**
- `min` (int, optional): Minimum value (default: 0)
- `max` (int, optional): Maximum value (default: 100)

**Returns:** int

### `float(min?, max?)`
Generates a random floating-point number.

**Parameters:**
- `min` (float, optional): Minimum value (default: 0.0)
- `max` (float, optional): Maximum value (default: 1.0)

**Returns:** float

### `bool()`
Generates a random boolean.

**Returns:** bool

### `choice(items)`
Selects a random element from a list.

**Parameters:**
- `items` (list, required): List to choose from

**Returns:** any (type of selected item)

### `bytes(size?)`
Generates random bytes as hex string.

**Parameters:**
- `size` (int, optional): Number of bytes (default: 16)

**Returns:** string (hex-encoded)

## Security Note

This namespace uses Go's `math/rand` for random generation, which is **not cryptographically secure**. 

For cryptographic purposes (passwords, keys, tokens), use a dedicated secure random namespace.

## Testing

```bash
echo '{"function":"int","params":{"min":1,"max":10},"context":{}}' | ./random
echo '{"function":"choice","params":{"items":["a","b","c"]},"context":{}}' | ./random
```

## License

MIT License

## Links

- [UP Language](https://uplang.org)
- [UP Specification](https://github.com/uplang/spec)
- [Schema](random.up-schema)
