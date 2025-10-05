# UP Namespace: math

Mathematical operations for UP templates.

## Installation

```bash
go install github.com/uplang/ns-math@latest
```

Or build from source:

```bash
cd ns/math
go build -o math
```

## Usage in UP Templates

```up
!use [math]

# Basic arithmetic
total $math.add(a=10, b=20)
difference $math.sub(a=100, b=25)
product $math.mul(a=5, b=7)
quotient $math.div(a=100, b=4)
remainder!int $math.mod(a=17, b=5)

# Power and roots
squared $math.pow(x=5, y=2)
cubed $math.pow(x=3, y=3)
root $math.sqrt(x=16)

# Absolute value
distance $math.abs(x=-42)

# Min/Max
minimum $math.min(a=5, b=10)
maximum $math.max(a=5, b=10)

# Rounding
ceiling $math.ceil(x=4.2)
floor $math.floor(x=4.8)
rounded $math.round(x=4.5)
```

## Functions

### Basic Operations

- `add(a, b)` - Addition
- `sub(a, b)` - Subtraction
- `mul(a, b)` - Multiplication
- `div(a, b)` - Division (b must be non-zero)
- `mod(a, b)` - Modulo/remainder (b must be non-zero)

### Advanced Operations

- `pow(x, y)` - Power (x^y)
- `sqrt(x)` - Square root (x must be non-negative)
- `abs(x)` - Absolute value

### Comparison

- `min(a, b)` - Minimum of two numbers
- `max(a, b)` - Maximum of two numbers

### Rounding

- `ceil(x)` - Round up to nearest integer
- `floor(x)` - Round down to nearest integer
- `round(x)` - Round to nearest integer

## Use Cases

- **Configuration calculations**: Dynamic port assignments, timeouts
- **Data transformations**: Scaling values, unit conversions
- **Business logic**: Pricing, discounts, quotas
- **Test data**: Random ranges, distribution calculations

## Testing

```bash
echo '{"function":"add","params":{"a":10,"b":20},"context":{}}' | ./math
echo '{"function":"sqrt","params":{"x":16},"context":{}}' | ./math
echo '{"function":"round","params":{"x":4.7},"context":{}}' | ./math
```

## License

MIT License

## Links

- [UP Language](https://uplang.org)
- [UP Specification](https://github.com/uplang/spec)
- [Schema](math.up-schema)
