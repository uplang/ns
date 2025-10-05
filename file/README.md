# UP Namespace: file

File system operations for UP templates.

## Installation

```bash
go install github.com/uplang/ns-file@latest
```

Or build from source:

```bash
cd ns/file
go build -o file
```

## Usage in UP Templates

```up
!use [file]

# Read file contents
config_data $file.read(path="config.txt")
template_content $file.read(path="template.html")

# Check file existence
has_config!bool $file.exists(path="/etc/app/config.yml")

# List directory contents
files $file.list(path="/data", pattern="*.json")

# Path manipulation
name $file.basename(path="/path/to/file.txt")
dir $file.dirname(path="/path/to/file.txt")
extension $file.ext(path="document.pdf")

# Join paths
full_path $file.join(parts=["${HOME}", ".config", "app", "settings.yml"])
```

## Functions

### `read(path)`
Reads file contents as string.

**Parameters:**
- `path` (string, required): File path

**Returns:** string

### `exists(path)`
Checks if file or directory exists.

**Parameters:**
- `path` (string, required): Path to check

**Returns:** bool

### `list(path, pattern?)`
Lists files in a directory.

**Parameters:**
- `path` (string, required): Directory path
- `pattern` (string, optional): Glob pattern (default: "*")

**Returns:** list of strings

### `basename(path)`
Returns the base name of a path.

**Parameters:**
- `path` (string, required): File path

**Returns:** string (filename)

### `dirname(path)`
Returns the directory part of a path.

**Parameters:**
- `path` (string, required): File path

**Returns:** string (directory)

### `ext(path)`
Returns the file extension.

**Parameters:**
- `path` (string, required): File path

**Returns:** string (includes the dot)

### `join(parts)`
Joins path components.

**Parameters:**
- `parts` (list, required): Path components

**Returns:** string (joined path)

## Security Note

File operations are restricted to paths accessible by the process. Paths are **not validated** for security - use with caution in untrusted contexts.

## Testing

```bash
echo '{"function":"read","params":{"path":"test.txt"},"context":{}}' | ./file
echo '{"function":"exists","params":{"path":"/tmp"},"context":{}}' | ./file
echo '{"function":"list","params":{"path":"."},"context":{}}' | ./file
```

## License

MIT License

## Links

- [UP Language](https://uplang.org)
- [UP Specification](https://github.com/uplang/spec)
- [Schema](file.up-schema)
