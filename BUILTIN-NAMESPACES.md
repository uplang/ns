# UP Built-in Namespaces

Core namespaces that are always available and provide essential functionality for UP templates.

## `$list` - List Generation & Manipulation

The `$list` namespace provides powerful list generation and transformation capabilities.

### `$list.generate(count, template)`

Generate a list by repeating a template with context.

**Signature:**
```up
$list.generate(count!int, template!any) -> list
```

**Context Available:**
- `$self.index` - Zero-based index (0, 1, 2, ...)
- `$self.number` - One-based number (1, 2, 3, ...)
- `$self.first` - Boolean, true if first item
- `$self.last` - Boolean, true if last item
- `$self.count` - Total number of items

**Example 1: Generate API calls**
```up
!use [api, list]

# Generate list of 3 API calls with sequential IDs
users $list.generate(3, $api.fetch_user($self.number))
```

**Result:**
```up
users [
  {name "Alice", id 1}  # $self.number = 1
  {name "Bob", id 2}    # $self.number = 2
  {name "Charlie", id 3} # $self.number = 3
]
```

**Example 2: Generate test data with sequence**
```up
!use [faker, id, list]

test_users $list.generate(5, {
  id $id.uuid
  sequence!int $self.number
  name $faker.name
  email $faker.email
  is_first!bool $self.first
  is_last!bool $self.last
})
```

**Result:**
```up
test_users [
  {
    id 550e8400-e29b-41d4-a716-446655440000
    sequence 1
    name "Alice Johnson"
    email "alice.j@example.com"
    is_first true
    is_last false
  }
  {
    id 550e8400-e29b-41d4-a716-446655440001
    sequence 2
    name "Bob Smith"
    email "bob.s@example.com"
    is_first false
    is_last false
  }
  # ... 3 more items
  {
    id 550e8400-e29b-41d4-a716-446655440004
    sequence 5
    name "Eve Wilson"
    email "eve.w@example.com"
    is_first false
    is_last true
  }
]
```

**Example 3: Generate range**
```up
!use [list]

# Generate sequential IDs
ids $list.generate(10, $self.number)
# Result: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

# Generate with offset
high_ids $list.generate(5, $math.add(100, $self.index))
# Result: [100, 101, 102, 103, 104]
```

### `$list.range(start, end, step)`

Generate a numeric range.

**Signature:**
```up
$list.range(start!int, end!int, step!int = 1) -> list<int>
```

**Examples:**
```up
!use [list]

numbers $list.range(1, 10)
# Result: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

evens $list.range(0, 20, 2)
# Result: [0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20]

countdown $list.range(10, 0, -1)
# Result: [10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0]
```

### `$list.map(list, template)`

Transform each item in a list using a template.

**Signature:**
```up
$list.map(list!list, template!any) -> list
```

**Context Available:**
- `$self.item` - Current item value
- `$self.index` - Zero-based index
- `$self.number` - One-based number
- All standard `$self` context

**Example:**
```up
!use [list]

vars {
  raw_ids [100, 200, 300]
}

# Transform each ID
user_ids $list.map($vars.raw_ids, user_$self.item)
# Result: ["user_100", "user_200", "user_300"]

# More complex transformation
enriched $list.map($vars.raw_ids, {
  id $self.item
  index $self.index
  label "ID-$self.item"
})
```

### `$list.filter(list, condition)`

Filter list items based on a condition.

**Signature:**
```up
$list.filter(list!list, condition!bool) -> list
```

**Context Available:**
- `$self.item` - Current item value
- All standard `$self` context

**Example:**
```up
!use [list]

vars {
  numbers [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
}

# Filter evens
evens $list.filter($vars.numbers, $math.mod($self.item, 2) == 0)
# Result: [2, 4, 6, 8, 10]
```

### `$list.repeat(value, count)`

Repeat a value N times.

**Signature:**
```up
$list.repeat(value!any, count!int) -> list
```

**Examples:**
```up
!use [list]

# Repeat simple value
zeros $list.repeat(0, 5)
# Result: [0, 0, 0, 0, 0]

# Repeat complex value
default_configs $list.repeat({
  enabled!bool true
  timeout!dur 30s
}, 3)
```

## `$self` - Context Variable

The `$self` variable provides context information within namespace function evaluations. It's automatically available whenever a template is being evaluated.

### Standard `$self` Properties

Available in all contexts:

```up
$self.index        # Zero-based index (0, 1, 2, ...)
$self.number       # One-based number (1, 2, 3, ...)
$self.first!bool   # True if first item
$self.last!bool    # True if last item
$self.count!int    # Total count of items
```

### Context-Specific Properties

#### In `$list.generate`:
```up
$self.index        # Current iteration (0-based)
$self.number       # Current iteration (1-based)
$self.first        # Is this the first iteration?
$self.last         # Is this the last iteration?
$self.count        # Total iterations
```

#### In `$list.map`:
```up
$self.item         # Current item from source list
$self.index        # Position in list (0-based)
$self.number       # Position (1-based)
$self.first        # Is this the first item?
$self.last         # Is this the last item?
$self.count        # Total items in list
```

#### In `$list.filter`:
```up
$self.item         # Current item being evaluated
$self.index        # Position in original list
```

### Using `$self` with Other Namespaces

**Example: Sequential IDs**
```up
!use [id, list]

# Generate UUIDs with sequence metadata
items $list.generate(3, {
  id $id.uuid
  sequence $self.number
  label "Item #$self.number of $self.count"
})
```

**Example: Conditional generation**
```up
!use [faker, list, random]

# Different data for first/last items
users $list.generate(5, {
  name $faker.name
  role $self.first ? admin : $self.last ? support : user
  priority!int $self.first ? 1 : $self.last ? 3 : 2
})
```

## `$math` - Mathematical Operations

Basic math operations for use in templates.

### Functions

```up
$math.add(a, b)           # Addition
$math.sub(a, b)           # Subtraction
$math.mul(a, b)           # Multiplication
$math.div(a, b)           # Division
$math.mod(a, b)           # Modulo
$math.pow(base, exp)      # Exponentiation
$math.sqrt(n)             # Square root
$math.abs(n)              # Absolute value
$math.min(a, b)           # Minimum
$math.max(a, b)           # Maximum
$math.round(n)            # Round to nearest integer
$math.floor(n)            # Round down
$math.ceil(n)             # Round up
```

**Examples:**
```up
!use [math, list]

# Calculate ranges
start 100
count 10
end $math.add($vars.start, $vars.count)

# Generate computed values
prices $list.generate(5, $math.mul($self.number, 9.99))
# Result: [9.99, 19.98, 29.97, 39.96, 49.95]

# Conditional math
discounts $list.generate(10, {
  item_id $self.number
  price!float 100.0
  discount!float $self.number > 5 ? 0.2 : 0.1
  final_price!float $math.mul(100.0, $math.sub(1, $self.discount))
})
```

## `$string` - String Operations

String manipulation functions.

### Functions

```up
$string.concat(a, b, ...)    # Concatenate strings
$string.upper(s)             # Uppercase
$string.lower(s)             # Lowercase
$string.trim(s)              # Trim whitespace
$string.replace(s, old, new) # Replace substring
$string.split(s, sep)        # Split into list
$string.join(list, sep)      # Join list into string
$string.length(s)            # String length
$string.substr(s, start, len) # Substring
$string.format(template, ...)  # Format string
```

**Examples:**
```up
!use [string, list, faker]

# Generate email addresses
emails $list.generate(3, $string.concat(
  $string.lower($faker.name),
  "@example.com"
))

# Format strings
labels $list.generate(5, $string.format(
  "Item #{} of {}",
  $self.number,
  $self.count
))
# Result: ["Item #1 of 5", "Item #2 of 5", ...]
```

## Type System for Namespace Functions

### Function Signatures in Schemas

**File: `list.up-schema`**

```up
namespace list

functions {
  generate {
    description "Generate a list by repeating a template"

    params {
      count!int {
        description "Number of items to generate"
        required!bool true
        min 0
        max 10000
      }

      template!template {
        description "Template to evaluate for each item"
        required!bool true
        context_provides [index, number, first, last, count]
      }
    }

    returns!list {
      description "Generated list"
      item_type any
    }

    examples [
      {
        call "$list.generate(3, $api.fetch($self.number))"
        result "[{data: 1}, {data: 2}, {data: 3}]"
      }
    ]
  }

  map {
    description "Transform each item using a template"

    params {
      list!list {
        description "Source list"
        required!bool true
      }

      template!template {
        description "Template to apply to each item"
        required!bool true
        context_provides [item, index, number, first, last, count]
      }
    }

    returns!list {
      description "Transformed list"
    }
  }

  range {
    description "Generate a numeric range"

    params {
      start!int {
        description "Start value (inclusive)"
        required!bool true
      }

      end!int {
        description "End value (inclusive)"
        required!bool true
      }

      step!int {
        description "Step size"
        default 1
        required!bool false
      }
    }

    returns!list<int> {
      description "List of numbers"
    }
  }
}

metadata {
  author "UP Core Team"
  builtin!bool true
  safe!bool true

  context_support {
    provides_self!bool true
    self_properties [index, number, first, last, count]
  }
}
```

### Type Annotations

#### `!template` Type

Indicates a parameter is a UP template expression that will be evaluated:

```up
template!template $faker.name
# The expression is not evaluated immediately
# It's evaluated in the context of the function
```

#### `!list<T>` Type

Generic list type with element type:

```up
numbers!list<int>
names!list<string>
users!list<object>
items!list<any>
```

#### Context Types

Functions can declare what context they provide:

```up
context_provides [index, number, item, first, last]
```

## Complete Example: Test Data Generation

```up
!use [list, faker, id, random, time, math]

# Generate 100 test users with realistic data
test_users $list.generate(100, {
  # Identity
  id $id.uuid
  user_id $self.number
  username $string.lower($string.concat(
    $faker.name,
    "_",
    $self.number
  ))

  # Profile
  profile {
    full_name $faker.name
    email $faker.email
    phone $faker.phone
    avatar_url "https://avatars.example.com/$self.user_id.jpg"
  }

  # Account status
  account {
    created_at $time.add($math.mul($self.index, -24h))
    status $self.first ? admin : $random.choice([active, pending, suspended])
    verified!bool $math.mod($self.number, 3) == 0
    login_count!int $random.int(0, 1000)
  }

  # Metadata
  metadata {
    sequence $self.number
    is_vip!bool $self.number <= 10
    tier $self.number <= 10 ? platinum : $self.number <= 50 ? gold : silver
  }
})

# Generate API endpoints with sequential IDs
api_tests $list.generate(20, {
  test_id test_$self.number
  method GET
  endpoint "/api/users/$self.number"
  expected_status!int 200
})

# Generate time series data
metrics $list.generate(24, {
  hour $self.index
  timestamp $time.add($math.mul($self.index, 1h))
  cpu_usage!float $math.add(50, $math.mul($random.float(-20, 20), $math.sin($self.index)))
  memory_usage!float $math.add(60, $math.mul(10, $math.cos($self.index)))
})

# Generate matrix/grid data
grid $list.generate(10, {
  row $self.number
  cells $list.generate(10, {
    col $self.number
    value $math.add($self.parent.row, $self.col)
  })
})
```

## Implementation Notes

### Template Evaluation

Templates are **lazy-evaluated** - they're not executed until the context is provided:

1. Parse template expression
2. For each iteration:
   - Set up `$self` context
   - Evaluate template with context
   - Collect result
3. Return list of results

### Nested Context

Templates can nest, with `$self` referring to the immediate context:

```up
$list.generate(3, {
  outer $self.number
  inner $list.generate(2, {
    outer_ref $self.parent.number  # Access parent context
    inner_val $self.number
  })
})
```

### Performance

- Templates are compiled once, evaluated many times
- Large lists (>1000 items) should show progress
- Consider streaming for very large lists

## Schema Validation

Built-in namespaces have type-checked parameters:

```up
# Type error: count must be int
$list.generate("5", $faker.name)  # ERROR

# Type error: template required
$list.generate(5)  # ERROR

# Valid
$list.generate(5, $faker.name)  # OK
```

## Summary

**`$list` namespace:**
- ✅ `generate(count, template)` - Repeat template N times
- ✅ `map(list, template)` - Transform each item
- ✅ `filter(list, condition)` - Filter items
- ✅ `range(start, end, step)` - Numeric ranges
- ✅ `repeat(value, count)` - Repeat value

**`$self` context:**
- ✅ Available in all template evaluations
- ✅ Provides index, number, first, last, count
- ✅ Context-specific properties (item, parent)
- ✅ Type-safe and documented

**Benefits:**
- 🎯 Declarative list generation
- 🔄 No imperative loops
- 📝 Self-documenting with context
- 🔒 Type-safe with schemas
- 🚀 Composable with all namespaces

**UP templates are declarative all the way down - even list generation!**

