# UP Namespace: env

Environment variable access and manipulation for UP templates.

## Installation

```bash
go install github.com/uplang/ns-env@latest
```

Or build from source:

```bash
cd ns/env
go build -o env
```

## Usage in UP Templates

```up
!use [env]

# Get environment variables
database {
  host $env.get(key="DB_HOST", default="localhost")
  port!int $env.get(key="DB_PORT", default="5432")
  user $env.get(key="DB_USER")
}

# Check if variable exists
has_token!bool $env.has(key="API_TOKEN")

# List variables with prefix
db_vars $env.list(prefix="DB_")

# Expand template strings
config_path $env.expand(template="${HOME}/.config/app")
connection $env.expand(template="$DB_USER@$DB_HOST:$DB_PORT")
```

## Functions

### `get(key, default?)`
Gets an environment variable value.

**Parameters:**
- `key` (string, required): Variable name
- `default` (string, optional): Default value if not set

**Returns:** string

**Best Practice:** Always provide defaults for optional config.

### `has(key)`
Checks if an environment variable exists.

**Parameters:**
- `key` (string, required): Variable name

**Returns:** bool

### `list(prefix?)`
Lists all environment variables.

**Parameters:**
- `prefix` (string, optional): Filter by prefix

**Returns:** list of strings

### `expand(template)`
Expands environment variables in a string.

**Parameters:**
- `template` (string, required): String with $VAR or ${VAR} references

**Returns:** string

**Supports:**
- `$VAR` syntax
- `${VAR}` syntax
- Undefined variables left as-is

## Best Practices

1. **Always provide defaults** for `get()`:
   ```up
   port $env.get(key="PORT", default="8080")
   ```

2. **Check required variables** with `has()`:
   ```up
   has_api_key!bool $env.has(key="API_KEY")
   ```

3. **Use `list()` for discovery**:
   ```up
   db_config $env.list(prefix="DB_")
   ```

4. **Use `expand()` for complex templates**:
   ```up
   path $env.expand(template="${HOME}/data/${APP_ENV}")
   ```

## Testing

```bash
export DB_HOST=localhost
export DB_PORT=5432
echo '{"function":"get","params":{"key":"DB_HOST"},"context":{}}' | ./env
echo '{"function":"list","params":{"prefix":"DB_"},"context":{}}' | ./env
```

## License

MIT License

## Links

- [UP Language](https://uplang.org)
- [UP Specification](https://github.com/uplang/spec)
- [Schema](env.up-schema)
