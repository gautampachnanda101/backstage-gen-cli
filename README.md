# backstage-gen

A Go-based CLI tool for streamlining Backstage catalog integration. Automatically generates and validates `catalog-info.yaml` files by inspecting your repository.

## Features

- 🔍 **Auto-detection**: Inspects repository to detect technology stack
- 📝 **YAML Generation**: Creates properly formatted catalog files
- ✅ **Validation**: Validates against Backstage schema
- 🪝 **Pre-commit Integration**: Can be used as a pre-commit hook
- 🚀 **Fast & Portable**: Single binary with no dependencies

## Installation

### From Source

```bash
git clone https://github.com/yourusername/backstage-gen.git
cd backstage-gen
make build
sudo make install
```

## Quick Start

### Generate catalog-info.yaml

```bash
# Generate in current directory
backstage-gen generate

# Generate with dry-run
backstage-gen generate --dry-run

# Force overwrite
backstage-gen generate --force
```

### Validate catalog

```bash
# Validate catalog-info.yaml
backstage-gen lint

# Strict validation
backstage-gen lint --strict
```

### Inspect repository

```bash
# Show detected information
backstage-gen inspect

# Output as JSON
backstage-gen inspect --json
```

## Pre-commit Hook

```bash
# Install hook
backstage-gen hooks install

# Uninstall hook
backstage-gen hooks uninstall
```

## Configuration

Create `.backstage-gen.yaml`:

```yaml
organization:
  name: "your-org"
  namespace: "default"

defaults:
  owner: "platform-team"
  system: "infrastructure"
  lifecycle: "production"
```

## Commands

### `generate`
Generate catalog-info.yaml from repository inspection.

**Flags:**
- `-o, --output` - Output file path
- `-t, --template` - Template to use
- `-i, --interactive` - Interactive mode
- `--dry-run` - Show without writing
- `--force` - Overwrite existing

### `lint`
Validate catalog files.

**Flags:**
- `-f, --file` - File to validate
- `--strict` - Treat warnings as errors
- `--format` - Output format

### `inspect`
Show detected repository information.

**Flags:**
- `--json` - Output as JSON
- `--verbose` - Detailed output

### `hooks`
Manage git hooks.

**Usage:**
- `backstage-gen hooks install`
- `backstage-gen hooks uninstall`

## Examples

```bash
# Complete workflow
backstage-gen inspect
backstage-gen generate --dry-run
backstage-gen generate
backstage-gen lint --strict
backstage-gen hooks install
```

## Development

```bash
# Build
make build

# Run tests
make test

# Run linters
make lint

# Clean
make clean
```

## License

MIT License
