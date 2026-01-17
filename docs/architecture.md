# Architecture Overview

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    CLI (cmd/*.go)                            │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ │
│  │ inspect │ │generate │ │  lint   │ │  init   │ │  hooks  │ │
│  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘ │
└───────┼──────────┼──────────┼──────────┼──────────┼─────────┘
        │          │          │          │          │
┌───────▼──────────▼──────────▼──────────▼──────────▼─────────┐
│                    Core Packages (pkg/)                      │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐        │
│  │ detector │ │generator │ │validator │ │  config  │        │
│  └────┬─────┘ └────┬─────┘ └──────────┘ └──────────┘        │
│       │            │                                         │
│  ┌────▼────────────▼────┐  ┌──────────┐  ┌──────────┐       │
│  │         llm          │  │  wizard  │  │  output  │       │
│  └──────────────────────┘  └──────────┘  └──────────┘       │
└─────────────────────────────────────────────────────────────┘
```

## Package Responsibilities

### cmd/ - CLI Commands

Uses [Cobra](https://github.com/spf13/cobra) for command-line interface.

| File | Command | Description |
|------|---------|-------------|
| `root.go` | - | Global flags, config initialization |
| `inspect.go` | `inspect` | Display detected repository info |
| `generate.go` | `generate` | Generate catalog-info.yaml |
| `lint.go` | `lint` | Validate catalog files |
| `init.go` | `init` | Initialize with LiteLLM |
| `hooks.go` | `hooks` | Manage git hooks |
| `config.go` | `config` | Configuration management |

### pkg/detector - Repository Detection

Analyzes repository to detect:
- Programming languages (17+ supported)
- Frameworks (React, Django, Spring Boot, etc.)
- Build tools (Make, npm, Maven, etc.)
- Infrastructure (Docker, Kubernetes, Helm, Terraform)
- Repository type (service, library, website, resource)

**Detection Strategy (Hybrid):**
1. Fast pattern matching (root-level files like `go.mod`, `package.json`)
2. LLM disambiguation for multi-language projects
3. Extension scanning as fallback

### pkg/generator - Catalog Generation

Generates Backstage catalog entities:
- Component (service, library, website)
- Resource (infrastructure)

Features:
- Automatic annotation generation
- Tag generation from detected tech stack
- LLM-enhanced descriptions and domain analysis

### pkg/validator - Catalog Validation

Validates catalog files against:
- Backstage schema requirements
- Required field checks
- Field format validation
- Best practice recommendations

### pkg/llm - LLM Integration

Supports multiple LLM providers:
- LiteLLM (recommended, supports 100+ models)
- Ollama (local)
- OpenRouter (cloud)

Used for:
- Language detection disambiguation
- Domain analysis
- Description enhancement
- Tag suggestions

### pkg/config - Configuration

Handles configuration from:
- YAML config files
- Environment variables
- Command-line flags

Uses [Viper](https://github.com/spf13/viper) for configuration management.

### pkg/output - Verbosity Control

Manages output verbosity levels:
- Level 0: Quiet (errors only)
- Level 1: Normal (default)
- Level 2: Verbose (extra details)
- Level 3: Debug (all debug info)

### pkg/wizard - Interactive Wizard

Provides interactive prompts for:
- Component name
- Description
- Owner selection
- Lifecycle stage
- Tags

## Data Flow

### Inspect Command

```
User runs: backstage-gen-cli inspect
    │
    ▼
┌─────────────────────────────────────┐
│ Load configuration                   │
│ (config.Load())                     │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│ Initialize LLM client (optional)    │
│ (llm.NewClient())                   │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│ Create detector                      │
│ (detector.New() or NewWithLLM())    │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│ Detect repository info              │
│ - Languages (hybrid detection)      │
│ - Frameworks                         │
│ - Build tools                        │
│ - Infrastructure                     │
│ - Git info                           │
│ - Description (from README)         │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│ Output results                       │
│ (formatted or JSON)                 │
└─────────────────────────────────────┘
```

### Generate Command

```
User runs: backstage-gen-cli generate
    │
    ▼
┌─────────────────────────────────────┐
│ Load configuration + detect repo    │
│ (same as inspect)                   │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│ [Optional] Run interactive wizard   │
│ (wizard.Run())                      │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│ Generate catalog                     │
│ (generator.Generate())              │
│ - Create metadata                    │
│ - Add labels                         │
│ - Add annotations                    │
│ - Generate tags                      │
│ - [Optional] LLM enhancement        │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│ Marshal to YAML and write file      │
│ (or dry-run preview)                │
└─────────────────────────────────────┘
```

## Key Design Decisions

### 1. Hybrid Language Detection

Pattern-first detection is fast and accurate for single-language projects. LLM is only used when multiple languages are detected to determine the primary language.

### 2. Optional LLM

All core functionality works without LLM. LLM enhances results but is not required.

### 3. Configuration Layering

Configuration merges from multiple sources (files, env vars, flags) with clear precedence.

### 4. Verbosity Levels

Consistent verbosity handling across all commands via the output package.

## Dependencies

| Dependency | Purpose |
|------------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `github.com/spf13/viper` | Configuration |
| `github.com/fatih/color` | Colored output |
| `github.com/go-git/go-git/v5` | Git operations |
| `gopkg.in/yaml.v3` | YAML parsing |
