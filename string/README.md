# UP Namespace: string

String manipulation functions for UP templates.

## Installation

```bash
go install github.com/uplang/ns-string@latest
```

Or build from source:

```bash
cd ns/string
go build -o string
```

## Usage in UP Templates

```up
!use [string]

# Case conversion
title $string.upper(s="hello world")
lowercase $string.lower(s="HELLO WORLD")
titlecase $string.title(s="hello world")

# Trimming
clean $string.trim(s="  hello  ")
no_prefix $string.trimPrefix(s="https://example.com", prefix="https://")
no_suffix $string.trimSuffix(s="file.txt", suffix=".txt")

# Splitting and joining
words $string.split(s="a,b,c", sep=",")
csv $string.join(items=["a", "b", "c"], sep=",")

# Replacement
fixed $string.replace(s="hello world", old="world", new="there")
all_replaced $string.replaceAll(s="foo foo foo", old="foo", new="bar")

# Checking
has_hello!bool $string.contains(s="hello world", substr="hello")
starts!bool $string.hasPrefix(s="https://example.com", prefix="https://")
ends!bool $string.hasSuffix(s="file.txt", suffix=".txt")

# Manipulation
substr $string.slice(s="hello", start=0, end=2)
repeated $string.repeat(s="ha", count=3)
reversed $string.reverse(s="hello")
len!int $string.length(s="hello")
```

## Functions

### Case Conversion

- `upper(s)` - Convert to uppercase
- `lower(s)` - Convert to lowercase
- `title(s)` - Convert to title case

### Trimming

- `trim(s, cutset?)` - Remove leading/trailing characters
- `trimPrefix(s, prefix)` - Remove prefix
- `trimSuffix(s, suffix)` - Remove suffix

### Splitting and Joining

- `split(s, sep?)` - Split into list (default sep: ",")
- `join(items, sep?)` - Join list into string (default sep: ",")

### Replacement

- `replace(s, old, new, n?)` - Replace n occurrences (default: 1)
- `replaceAll(s, old, new)` - Replace all occurrences

### Checking

- `contains(s, substr)` - Check if contains substring
- `hasPrefix(s, prefix)` - Check if starts with prefix
- `hasSuffix(s, suffix)` - Check if ends with suffix

### Manipulation

- `slice(s, start?, end?)` - Extract substring
- `repeat(s, count?)` - Repeat string n times (max: 10000)
- `reverse(s)` - Reverse string
- `length(s)` - Get string length in bytes

## Common Patterns

### URL Processing
```up
!use [string]

url "https://api.example.com/v1/users"
protocol $string.split(s=$url, sep="://")[0]
clean_url $string.trimPrefix(s=$url, prefix="https://")
```

### Text Formatting
```up
!use [string]

name "alice"
formatted $string.title(s=$name)
uppercased $string.upper(s=$name)
```

### Path Manipulation
```up
!use [string]

path "/path/to/file.txt"
filename $string.split(s=$path, sep="/")[-1]
extension $string.split(s=$filename, sep=".")[- 1]
without_ext $string.trimSuffix(s=$filename, suffix=$extension)
```

## Testing

```bash
echo '{"function":"upper","params":{"s":"hello"},"context":{}}' | ./string
echo '{"function":"split","params":{"s":"a,b,c","sep":","},"context":{}}' | ./string
echo '{"function":"reverse","params":{"s":"hello"},"context":{}}' | ./string
```

## License

MIT License

## Links

- [UP Language](https://uplang.org)
- [UP Specification](https://github.com/uplang/spec)
- [Schema](string.up-schema)
