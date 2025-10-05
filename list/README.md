# UP Namespace: list

List generation and manipulation functions for UP templates.

## Installation

```bash
go install github.com/uplang/ns-list@latest
```

Or build from source:

```bash
cd ns/list
go build -o list
```

## Usage in UP Templates

```up
!use [list, id]

# Generate list from template
users $list.generate(count=3, template={
  id $id.uuid
  name "User ${self.number}"
  active!bool true
})

# Join list items
tags ["go", "rust", "python"]
tags_string $list.join(items=$tags, separator=", ")

# Slice a list
numbers [1, 2, 3, 4, 5]
first_three $list.slice(items=$numbers, start=0, end=3)

# List operations
items ["a", "b", "c"]
length!int $list.length(items=$items)
has_b!bool $list.contains(items=$items, value="b")
index_b!int $list.index(items=$items, value="b")
```

## Functions

### `generate(count, template)`
Generates a list by repeating a template.

**Parameters:**
- `count` (int, required): Number of items
- `template` (any, required): Template to repeat

**Returns:** list

**Context variables** ($self):
- `$self.number`: 1-based index
- `$self.index`: 0-based index
- `$self.first`: true if first item
- `$self.last`: true if last item
- `$self.count`: total count

### `join(items, separator?)`
Joins list items into a string.

**Parameters:**
- `items` (list, required): List to join
- `separator` (string, optional): Separator (default: ", ")

**Returns:** string

### `slice(items, start?, end?)`
Extracts a portion of a list.

**Parameters:**
- `items` (list, required): Source list
- `start` (int, optional): Start index (default: 0)
- `end` (int, optional): End index (default: length)

**Returns:** list

### `length(items)`
Returns the length of a list.

**Parameters:**
- `items` (list, required): List to measure

**Returns:** int

### `contains(items, value)`
Checks if a list contains a value.

**Parameters:**
- `items` (list, required): List to search
- `value` (any, required): Value to find

**Returns:** bool

### `index(items, value)`
Returns the index of a value in a list.

**Parameters:**
- `items` (list, required): List to search
- `value` (any, required): Value to find

**Returns:** int (-1 if not found)

## Testing

```bash
echo '{"function":"generate","params":{"count":3,"template":"item"},"context":{}}' | ./list
echo '{"function":"join","params":{"items":["a","b","c"],"separator":", "},"context":{}}' | ./list
```

## License

MIT License

## Links

- [UP Language](https://uplang.org)
- [UP Specification](https://github.com/uplang/spec)
- [Schema](list.up-schema)
