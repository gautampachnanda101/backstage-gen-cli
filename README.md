# backstage-gen-cli

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/gautampachnanda101/backstage-gen-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/gautampachnanda101/backstage-gen-cli/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/gautampachnanda101/backstage-gen-cli?include_prereleases)](https://github.com/gautampachnanda101/backstage-gen-cli/releases)

A Go-based CLI tool for streamlining Backstage catalog integration. Automatically generates and validates `catalog-info.yaml` files by inspecting your repository.

## Features

- **Auto-detection**: Inspects repository to detect technology stack (20+ languages supported)
- **Cross-Platform**: Works on macOS, Linux, and Windows
- **AI-Powered**: Optional LLM integration for intelligent suggestions (via LiteLLM)
- **Interactive Wizard**: Step-by-step guided catalog generation
- **YAML Generation**: Creates properly formatted catalog files
- **Validation**: Validates against Backstage schema
- **Pre-commit Integration**: Can be used as a pre-commit hook
- **Fast & Portable**: Single binary with no dependencies

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap gautampachnanda101/backstage-gen-cli
brew install backstage-gen-cli
```

### Scoop (Windows)

```powershell
scoop bucket add backstage-gen-cli https://github.com/gautampachnanda101/scoop-backstage-gen-cli
scoop install backstage-gen-cli
```

### Chocolatey (Windows)

```powershell
choco install backstage-gen-cli
```

### Binary Download

```bash
# Linux/macOS
curl -L https://github.com/gautampachnanda101/backstage-gen-cli/releases/latest/download/backstage-gen-cli-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/') -o backstage-gen-cli
chmod +x backstage-gen-cli
sudo mv backstage-gen-cli /usr/local/bin/
```

### From Source

```bash
git clone https://github.com/gautampachnanda101/backstage-gen-cli.git
cd backstage-gen-cli
make build
sudo make install
```

See [Installation Guide](docs/installation.md) for more options.

## Quick Start

### Initialize (First Time Setup)

```bash
# Initialize with LiteLLM AI support
backstage-gen-cli init

# Or skip AI features
backstage-gen-cli init --skip-llm
```

### Generate catalog-info.yaml

```bash
# Interactive wizard with AI suggestions
backstage-gen-cli generate -i

# Auto-generate without prompts
backstage-gen-cli generate

# Preview with dry-run
backstage-gen-cli generate --dry-run

# Force overwrite
backstage-gen-cli generate --force
```

### Validate catalog

```bash
# Validate catalog-info.yaml
backstage-gen-cli lint

# Strict validation
backstage-gen-cli lint --strict
```

### Inspect repository

```bash
# Show detected information
backstage-gen-cli inspect

# Output as JSON
backstage-gen-cli inspect --json

# Verbose output
backstage-gen-cli inspect -v

# Debug output
backstage-gen-cli inspect -vvv
```

## Configuration

Create `.backstage-gen.yaml` in your project or home directory:

```yaml
# LLM Configuration (optional)
llm:
  provider: litellm
  base_url: http://localhost:4000
  model: ollama/llama2

# Organization settings
organization:
  name: "your-org"
  namespace: "default"

# Default values
defaults:
  owner: "platform-team"
  system: "infrastructure"
  lifecycle: "production"
```

See [Configuration Guide](docs/configuration.md) for all options.

## Commands

| Command | Description |
|---------|-------------|
| `generate` | Generate catalog-info.yaml from repository inspection |
| `inspect` | Show detected repository information |
| `lint` | Validate catalog files |
| `init` | Initialize with LLM support |
| `hooks` | Manage git hooks |
| `config` | Manage configuration |

### Global Flags

| Flag | Description |
|------|-------------|
| `--config` | Custom config file path |
| `-v, --verbose` | Increase verbosity (-v, -vv, -vvv) |
| `-q, --quiet` | Suppress non-essential output |
| `--help` | Show help |
| `--version` | Show version |

## Pre-commit Hook

```bash
# Install hook
backstage-gen-cli hooks install

# Uninstall hook
backstage-gen-cli hooks uninstall
```

## Documentation

- [Getting Started](docs/getting-started.md)
- [Installation Guide](docs/installation.md)
- [Configuration Guide](docs/configuration.md)
- [Architecture](docs/architecture.md)
- [Testing](docs/testing.md)
- [CI/CD Examples](docs/ci-cd-examples.md)
- [Contributing](docs/contributing.md)
- [MCP Integration](docs/mcp-integration.md)

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

See [Contributing Guide](docs/contributing.md) for development setup.

## License

MIT License - see [LICENSE](LICENSE) for details.
