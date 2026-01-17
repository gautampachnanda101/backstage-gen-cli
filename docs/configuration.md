# Configuration Guide

## Configuration File

backstage-gen-cli uses a YAML configuration file. It searches for configuration in this order:

1. `--config` flag (explicit path)
2. `.backstage-gen.yaml` (current directory)
3. `~/.backstage-gen.yaml` (home directory)
4. `~/.config/backstage-gen/config.yaml`

## Creating Configuration

### Quick Setup

```bash
# Create example config in home directory
backstage-gen-cli config init

# Or use init command for guided setup
backstage-gen-cli init
```

### Manual Creation

Create `.backstage-gen.yaml` in your project or home directory:

```yaml
# LLM Configuration (for AI-powered suggestions)
llm:
  provider: litellm
  base_url: http://localhost:4000
  model: ollama/llama2
  timeout: 30

# Organization settings
organization:
  name: my-org
  namespace: default

# Default values for generated catalogs
defaults:
  owner: platform-team
  system: ""
  lifecycle: production

  # Custom annotations to add to all catalogs
  annotations:
    company.io/team: "platform"
    company.io/slack-channel: "#backstage"

  # Custom tags to add to all catalogs
  tags:
    - internal
    - company-standard
```

## Configuration Options

### LLM Settings

| Option | Description | Default |
|--------|-------------|---------|
| `provider` | LLM provider (litellm, ollama, openrouter) | `litellm` |
| `base_url` | Provider API URL | `http://localhost:4000` |
| `model` | Model to use | `ollama/llama2` |
| `api_key` | API key (if required) | - |
| `timeout` | Request timeout in seconds | `30` |

### Organization Settings

| Option | Description | Default |
|--------|-------------|---------|
| `name` | Organization name | - |
| `namespace` | Backstage namespace | `default` |

### Defaults

| Option | Description | Default |
|--------|-------------|---------|
| `owner` | Default component owner | `platform-team` |
| `system` | Default system name | - |
| `lifecycle` | Default lifecycle stage | `production` |
| `annotations` | Custom annotations map | - |
| `tags` | Custom tags list | - |

## LLM Provider Examples

### LiteLLM with Docker (Recommended)

```yaml
llm:
  provider: litellm
  base_url: http://localhost:4000
  model: ollama/llama2
  timeout: 30
```

Start LiteLLM:
```bash
docker run -d --name litellm -p 4000:4000 ghcr.io/berriai/litellm:main-latest
```

### Direct Ollama

```yaml
llm:
  provider: ollama
  base_url: http://localhost:11434
  model: llama2
  timeout: 30
```

### OpenRouter (Cloud)

```yaml
llm:
  provider: openrouter
  base_url: https://openrouter.ai
  model: mistralai/mistral-7b-instruct:free
  api_key: your-api-key
  timeout: 30
```

### OpenAI via LiteLLM

```yaml
llm:
  provider: litellm
  base_url: http://localhost:4000
  model: gpt-3.5-turbo
  api_key: your-openai-key
```

### Anthropic Claude via LiteLLM

```yaml
llm:
  provider: litellm
  base_url: http://localhost:4000
  model: claude-2
  api_key: your-anthropic-key
```

## Environment Variables

Configuration can also be set via environment variables with the `BACKSTAGE_GEN_` prefix:

```bash
export BACKSTAGE_GEN_LLM_PROVIDER=ollama
export BACKSTAGE_GEN_LLM_BASE_URL=http://localhost:11434
export BACKSTAGE_GEN_DEFAULTS_OWNER=my-team
```

## Command-Line Flags

Global flags available on all commands:

| Flag | Description |
|------|-------------|
| `--config` | Specify config file path |
| `-v, --verbose` | Increase verbosity (-v, -vv, -vvv) |
| `-q, --quiet` | Suppress non-essential output |

## Project-Level Configuration

You can have different configurations per project by placing `.backstage-gen.yaml` in each project directory. Project-level config takes precedence over home directory config.

Example project config:
```yaml
# Project-specific overrides
defaults:
  owner: frontend-team
  system: web-platform
  lifecycle: experimental
  tags:
    - frontend
    - react
```

## Validation

Verify your configuration is being loaded:

```bash
# Verbose mode shows config file path
backstage-gen-cli inspect -v

# Debug mode shows all config details
backstage-gen-cli inspect -vvv
```
