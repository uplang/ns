# Dynamic Variable Namespaces

UP supports **dynamic variable namespaces** for generating content on-the-fly, useful for testing, mocking, and rapid prototyping. Dynamic namespaces provide functions for time, dates, random data, IDs, and more.

## Core Principle

**Static UP remains simple** - Dynamic namespaces are **opt-in via pragmas**. If you don't use them, your UP is just static data. If you do, the template processor generates dynamic values at processing time.

## Pragma Syntax

Declare which dynamic namespaces you need at the **top of the file** using `!use`:

```up
!use [time, random, id]

# Now you can use these namespaces
created_at $time.now
request_id $id.uuid
test_value $random.int(1, 100)
```

### Single Namespace

```up
!use time

timestamp $time.now
```

### Multiple Namespaces

```up
!use [time, random, id, faker]

# All four namespaces available
```

### Namespace Aliases

When you need multiple versions or implementations of the same namespace, use aliases with `as`:

```up
!use [time, time@1.0.0 as oldtime, github.com/other/time-ns as mytime]

# Use different time implementations
current $time.now
legacy $oldtime.now
custom $mytime.now
```

**Alias Syntax:**
```up
!use [namespace as alias]
!use [namespace@version as alias]
!use [url as alias]
```

**Examples:**
```up
# Multiple versions
!use [faker, faker@1.0.0 as faker1, faker@2.0.0 as faker2]

user1 {
  name $faker1.name
  email $faker1.email
}

user2 {
  name $faker2.name
  email $faker2.email
}
```

```up
# Different implementations
!use [
  github.com/uplang/ns-random as random,
  github.com/myorg/secure-random as securerandom
]

# Use appropriate random for each purpose
session_id $random.uuid
crypto_key $securerandom.bytes(size=32)
```

```up
# Avoid naming conflicts
!use [
  time,
  github.com/myorg/custom-time as mytime
]

created_at $time.now
custom_timestamp $mytime.special_format
```

## Namespace Reference

### `time` - Time and Timestamps

Current time in various formats:

```up
!use time

# Current time (default: RFC3339)
now $time.now
now_iso $time.now(iso8601)
now_unix $time.now(unix)
now_custom $time.now("2006-01-02 15:04:05")

# Relative times
yesterday $time.add(-24h)
next_week $time.add(168h)
last_month $time.add(-720h)

# Formatting
date_only $time.format(date)
time_only $time.format(time)
custom $time.format("Mon Jan 2 15:04:05 MST 2006")
```

**Available Functions:**
- `$time.now` - Current timestamp (RFC3339)
- `$time.now(format)` - Formatted current time
- `$time.add(duration)` - Time relative to now
- `$time.format(format)` - Format current time
- `$time.unix` - Unix timestamp
- `$time.unix_ms` - Unix milliseconds

### `date` - Calendar Dates

Date-specific operations:

```up
!use date

# Current date
today $date.today
today_iso $date.today(iso)

# Relative dates
yesterday $date.add(-1d)
next_week $date.add(7d)
last_year $date.add(-365d)

# Date components
year $date.year
month $date.month
day $date.day
weekday $date.weekday
```

**Available Functions:**
- `$date.today` - Current date (YYYY-MM-DD)
- `$date.add(days)` - Date relative to today
- `$date.year`, `$date.month`, `$date.day` - Components
- `$date.weekday` - Day of week (Monday, Tuesday, ...)

### `random` - Random Data Generation

Generate random values for testing:

```up
!use random

# Numbers
port $random.int(8000, 9000)
percentage $random.float(0.0, 100.0)
size $random.int(1, 1000)

# Strings
password $random.string(32)
hex $random.hex(16)
base64 $random.base64(24)

# Boolean
enabled $random.bool

# Lists (pick one)
environment $random.choice([dev, staging, prod])
region $random.choice([us-east-1, us-west-2, eu-west-1])
```

**Available Functions:**
- `$random.int(min, max)` - Random integer
- `$random.float(min, max)` - Random float
- `$random.string(length)` - Random alphanumeric
- `$random.hex(length)` - Random hex string
- `$random.base64(length)` - Random base64
- `$random.bool` - Random boolean
- `$random.choice([items])` - Pick random item from list
- `$random.uuid` - Alias for `$id.uuid`

### `id` - ID Generation

Generate various types of identifiers:

```up
!use id

# UUIDs
request_id $id.uuid
correlation_id $id.uuid4

# Other ID formats
short_id $id.short
nano_id $id.nano
ulid $id.ulid

# Custom formats
custom_id $id.format("REQ-{random.hex(8)}")
```

**Available Functions:**
- `$id.uuid` - UUID v4
- `$id.uuid4` - UUID v4 (explicit)
- `$id.short` - Short ID (8 chars)
- `$id.nano` - Nano ID (21 chars)
- `$id.ulid` - ULID
- `$id.sequential` - Sequential ID (session-scoped)

### `faker` - Realistic Fake Data

Generate realistic fake data for testing (inspired by [go-faker](https://github.com/go-faker/faker)):

```up
!use faker

# Personal information
user_name $faker.name
user_email $faker.email
user_phone $faker.phone

# Address
street $faker.street
city $faker.city
state $faker.state
zip $faker.zipcode
country $faker.country

# Internet
username $faker.username
domain $faker.domain
ipv4 $faker.ipv4
ipv6 $faker.ipv6
mac $faker.mac_address

# Lorem ipsum
sentence $faker.sentence
paragraph $faker.paragraph
text $faker.text(100)

# Business
company $faker.company
job_title $faker.job_title
buzzword $faker.buzzword

# Financial
credit_card $faker.credit_card
currency $faker.currency
price $faker.price(10.0, 1000.0)

# Dates and times
past_date $faker.date(past)
future_date $faker.date(future)
birthday $faker.date(birthday)
```

**Available Functions:**
- Personal: `name`, `email`, `phone`, `username`
- Address: `street`, `city`, `state`, `zipcode`, `country`
- Internet: `domain`, `ipv4`, `ipv6`, `url`, `mac_address`
- Text: `sentence`, `paragraph`, `text(length)`, `word`
- Business: `company`, `job_title`, `buzzword`
- Financial: `credit_card`, `currency`, `price(min, max)`
- Dates: `date(past)`, `date(future)`, `date(birthday)`

### `env` - Environment Variables

Access environment variables:

```up
!use env

# Access environment variables
home $env.HOME
user $env.USER
path $env.PATH

# With defaults
api_key $env.API_KEY(default-key)
port $env.PORT(8080)
```

**Available Functions:**
- `$env.VAR_NAME` - Get environment variable
- `$env.VAR_NAME(default)` - Get with default value

### `file` - File System Operations

Read file contents:

```up
!use file

# Read entire file
config $file.read(config.json)
template $file.read(template.html)

# Read with encoding
utf8_content $file.read(data.txt, utf8)
```

**Available Functions:**
- `$file.read(path)` - Read file contents
- `$file.exists(path)` - Check if file exists (boolean)

## Parameter Syntax

### Positional Parameters

```up
# Single parameter
$random.int(100)              # Random int from 0 to 100
$time.add(24h)                # 24 hours from now
$faker.text(50)               # 50 chars of text

# Multiple parameters
$random.int(10, 100)          # Random int between 10 and 100
$faker.price(10.0, 1000.0)    # Price between $10 and $1000
$time.format("2006-01-02")    # Custom format
```

### Named Parameters (Optional)

For clarity in complex cases:

```up
# Named parameters
$random.int(min=10, max=100)
$faker.price(min=10.0, max=1000.0)
$time.add(duration=24h)
```

### No Parameters

```up
# Just call it
$time.now
$date.today
$id.uuid
$faker.email
```

## Complete Examples

### Example 1: Test Configuration

```up
!use [time, random, id, faker]

# Test environment configuration
test_run {
  id $id.uuid
  started_at $time.now
  environment test

  database {
    host $faker.ipv4
    port $random.int(5432, 5532)
    name test_db_$random.hex(8)
    user test_user
    password $random.base64(32)
  }

  api {
    endpoint https://$faker.domain/api
    key api_$random.string(40)
    timeout!dur 30s
  }

  test_users [
    {
      id $id.uuid
      name $faker.name
      email $faker.email
      created_at $time.add(-30d)
    }
    {
      id $id.uuid
      name $faker.name
      email $faker.email
      created_at $time.add(-15d)
    }
    {
      id $id.uuid
      name $faker.name
      email $faker.email
      created_at $time.add(-7d)
    }
  ]
}
```

### Example 2: Mock API Responses

```up
!use [time, id, faker, random]

# Generate mock user for testing
users!list {
  user_$random.int(1000, 9999) {
    id $id.uuid
    username $faker.username
    email $faker.email
    full_name $faker.name
    phone $faker.phone

    address {
      street $faker.street
      city $faker.city
      state $faker.state
      zip $faker.zipcode
      country USA
    }

    account {
      created_at $time.add(-180d)
      last_login $time.add(-2h)
      verified!bool true
      balance!float $random.float(0.0, 10000.0)
      tier $random.choice([free, pro, enterprise])
    }

    metadata {
      ip_address $faker.ipv4
      user_agent Mozilla/5.0...
      session_id $id.short
    }
  }
}
```

### Example 3: Load Testing Data

```up
!use [random, faker, time, id]

# Generate load test scenario
load_test {
  scenario_id $id.uuid
  duration!dur 5m
  ramp_up!dur 30s

  virtual_users!int $random.int(100, 500)

  endpoints [
    {
      path /api/users
      method GET
      weight!int $random.int(40, 60)
      response_time!dur $random.int(50, 200)ms
    }
    {
      path /api/products
      method GET
      weight!int $random.int(20, 40)
      response_time!dur $random.int(100, 300)ms
    }
    {
      path /api/orders
      method POST
      weight!int $random.int(10, 20)
      response_time!dur $random.int(200, 500)ms
    }
  ]

  test_data {
    user_pool!int $random.int(1000, 5000)
    product_count!int $random.int(500, 2000)
    order_rate!int $random.int(10, 50)
  }
}
```

### Example 4: Development Seed Data

```up
!use [faker, id, time, random]

# Seed database with realistic data
seed_data {
  created_at $time.now
  version 1.0

  companies [
    {
      id $id.uuid
      name $faker.company
      domain $faker.domain
      email contact@$faker.domain
      phone $faker.phone

      address {
        street $faker.street
        city $faker.city
        state $faker.state
        country $faker.country
      }

      employees!int $random.int(10, 1000)
      founded_year!int $random.int(1990, 2023)
      revenue!float $random.float(100000, 10000000)
    }
  ]

  products!list {
    product_$random.hex(8) {
      id $id.uuid
      name $faker.buzzword Product
      description $faker.paragraph
      price!float $faker.price(9.99, 999.99)
      sku PROD-$random.string(10)
      in_stock!bool $random.bool
      quantity!int $random.int(0, 500)
      category $random.choice([electronics, clothing, home, sports])
    }
  }
}
```

## Implementation Notes

### Processing Order

1. **Parse UP document** - Parse static structure first
2. **Check for `!use` pragma** - Identify required namespaces
3. **Initialize namespaces** - Set up dynamic generators
4. **Resolve dynamic variables** - Replace `$namespace.function()` calls
5. **Resolve static variables** - Process `$vars.*` references
6. **Output final document** - Static UP with all values resolved

### Namespace Scoping

- **Per-document** - Each document has its own namespace instances
- **Deterministic when seeded** - Can provide seed for reproducible generation
- **Sequential IDs** - `$id.sequential` starts at 1 per document

### Seeding for Reproducibility

```up
!use [random, faker]
!seed 12345

# All random/faker calls will be deterministic with this seed
user_id $random.int(1000, 9999)  # Always same value with same seed
```

### Combining with Static Variables

```up
!use [time, id]

vars {
  environment production
  region us-west-2
  base_url https://$vars.environment.example.com
}

deployment {
  id $id.uuid
  environment $vars.environment
  region $vars.region
  url $vars.base_url
  deployed_at $time.now
  deployed_by $env.USER
}
```

## CLI Usage

### Generate with Dynamic Namespaces

```bash
# Process template with dynamic namespaces
up template process -i template.up -o output.up

# Generate with specific seed for reproducibility
up template process -i template.up --seed 12345 -o output.up

# Generate multiple instances
up template process -i template.up --count 10 -o outputs/
```

### Available Flags

- `--seed N` - Seed random generators for reproducibility
- `--count N` - Generate N instances
- `--env-file FILE` - Load environment variables from file
- `--no-dynamic` - Disable dynamic namespaces (error if `!use` present)

## Security Considerations

### Dangerous Namespaces

Some namespaces can be dangerous in production:

- `$file.*` - File system access
- `$env.*` - Environment variable access
- `$exec.*` - Command execution (if implemented)

**Best practices:**

1. **Whitelist namespaces** in production parsers
2. **Sandbox file access** - Limit to specific directories
3. **Validate inputs** - Sanitize all dynamic values
4. **Audit usage** - Log all dynamic namespace usage

### Recommended Safe List

```up
# Safe for production
!use [time, date, id]

# Use with caution
!use [random, faker]

# Dangerous - only in development
!use [env, file]
```

## Comparison with Other Systems

### vs Sprig (Helm)

**Sprig:**
```go-template
{{ now | date "2006-01-02" }}
{{ randAlphaNum 32 }}
{{ uuidv4 }}
```

**UP:**
```up
!use [time, random, id]

date $time.format("2006-01-02")
random $random.string(32)
id $id.uuid
```

**Advantages:**
- ✅ No template language mixing
- ✅ Explicit namespace declaration
- ✅ Pure UP syntax
- ✅ Static UP stays simple

### vs go-faker

**go-faker** ([example](https://github.com/go-faker/faker/blob/main/example_with_tags_test.go)):
```go
type Example struct {
    Name  string `faker:"name"`
    Email string `faker:"email"`
    Phone string `faker:"phone_number"`
}
```

**UP:**
```up
!use faker

example {
  name $faker.name
  email $faker.email
  phone $faker.phone
}
```

**Advantages:**
- ✅ No struct tags needed
- ✅ Direct use in config files
- ✅ Mix with static data easily
- ✅ Template generation on demand

## Design Philosophy

1. **Opt-in complexity** - Static UP is simple; dynamic is opt-in via `!use`
2. **Explicit over implicit** - Declare namespaces upfront
3. **Clean syntax** - `$namespace.function(params)` is obvious
4. **Parser-friendly** - Pragmas are easy to detect and handle
5. **Composable** - Mix dynamic and static variables freely
6. **Deterministic when needed** - Seed for reproducible tests
7. **Safe by default** - Dangerous namespaces require explicit opt-in

## Reserved Namespaces

The following namespace names are **reserved** for standard implementations:

- `time` - Time and timestamps
- `date` - Calendar dates
- `random` - Random data
- `id` - ID generation
- `faker` - Realistic fake data
- `env` - Environment variables
- `file` - File system access
- `http` - HTTP requests (future)
- `exec` - Command execution (future, dangerous)
- `vault` - Secret management (future)
- `aws`, `gcp`, `azure` - Cloud provider integrations (future)

**Custom namespaces** should use your own prefix: `myapp_*`, `company_*`, etc.

## Future Extensions

Possible future namespaces:

```up
# HTTP requests
!use http
api_response $http.get("https://api.example.com/status")

# Cloud resources
!use aws
instance_id $aws.ec2.instance_id
region $aws.region

# Secrets management
!use vault
db_password $vault.read("secret/database/password")

# Crypto operations
!use crypto
hash $crypto.sha256("content")
signature $crypto.sign(data, key)
```

---

**Dynamic namespaces make UP a powerful tool for testing, mocking, and rapid prototyping while keeping the core language simple and parseable.**

## Extensibility: Custom Namespaces

While the namespaces listed above are built-in or standard, UP's namespace system is **fully pluggable**. You can create custom namespaces using any executable (shell script, binary, etc.).

See **[NAMESPACE-PLUGINS.md](NAMESPACE-PLUGINS.md)** for complete documentation on:

- Creating custom namespace plugins
- Writing UP schemas to document your namespaces
- Plugin protocol (JSON-based stdin/stdout)
- Example plugins in bash, Python, Go, Rust
- Security and sandboxing

**Example: Custom namespace via shell script**

```bash
#!/bin/bash
# ./up-namespaces/myns

INPUT=$(cat)
FUNCTION=$(echo "$INPUT" | jq -r '.function')

case "$FUNCTION" in
  generate)
    echo '{"value": "custom content", "type": "string"}'
    ;;
esac
```

```up
!use [myns]

data $myns.generate
```

**UP schemas are written in UP syntax:**

```up
# myns.up-schema
namespace myns
version 1.0.0

functions {
  generate {
    description Generate custom content
    returns!string
  }
}
```

This makes UP infinitely extensible while keeping the core language simple and the extension mechanism transparent.

