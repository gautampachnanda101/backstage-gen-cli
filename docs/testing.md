# Testing Documentation

## Overview

backstage-gen-cli uses a multi-layered testing approach:
1. Unit tests for core packages
2. Integration tests against real OSS projects
3. Cross-platform build verification

## Unit Tests

### Running Unit Tests

```bash
# Run all tests
make test

# Run with verbose output
go test -v ./...

# Run specific package tests
go test -v ./pkg/detector/...
go test -v ./pkg/llm/...
go test -v ./pkg/config/...
go test -v ./pkg/generator/...
go test -v ./pkg/validator/...

# Run with race detection
go test -race ./...

# Run with coverage
go test -cover ./...
```

### Test Coverage by Package

| Package | Coverage | Key Tests |
|---------|----------|-----------|
| `pkg/detector` | Language, framework, build tool detection |
| `pkg/llm` | LLM client, mock server tests |
| `pkg/config` | Config loading, defaults |
| `pkg/generator` | Catalog generation, annotations |
| `pkg/validator` | Schema validation, field checks |
| `pkg/output` | Verbosity levels, quiet mode |

## GitHub Actions CI/CD Testing

Our comprehensive CI/CD pipeline tests the CLI against real-world open source projects.

### Test Coverage

#### 1. Make Tasks Testing
- `make build` - Builds the binary successfully
- `make test` - Runs all Go tests
- `make lint` - Runs linters (go vet, go fmt)
- `make clean` - Cleans build artifacts

#### 2. Cross-Platform Builds

Builds binaries for all major platforms:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

#### 3. Real-World Project Testing

We test against 10+ diverse open source projects:

| Project | Language | Type | What We Validate |
|---------|----------|------|------------------|
| **kubernetes/kubernetes** | Go | Orchestration | Go detection, large repos |
| **django/django** | Python | Web Framework | Python, framework detection |
| **pallets/flask** | Python | Microframework | Description parsing |
| **facebook/react** | JavaScript | UI Library | JS/TS detection |
| **vercel/next.js** | JavaScript/TypeScript | Framework | Mixed language detection |
| **vuejs/core** | JavaScript/TypeScript | Framework | Vue detection |
| **spring-projects/spring-boot** | Java | Framework | Java, Maven detection |
| **rust-lang/rust** | Rust | Compiler | Rust detection |
| **rails/rails** | Ruby | Framework | Ruby, Rails detection |
| **grafana/grafana** | Go/TypeScript | Multi-lang | Multi-language handling |

### Command Testing

For each project, we test:

**Inspect Command:**
```bash
backstage-gen-cli inspect
```
- Detects correct primary language
- Identifies frameworks
- Extracts build tools
- Parses Git information

**Generate Command (Dry-Run):**
```bash
backstage-gen-cli generate --dry-run
```
- Produces valid YAML
- Includes standard Backstage annotations
- Has clean, readable descriptions

**Generate Command (Actual):**
```bash
backstage-gen-cli generate -o catalog-test.yaml
```
- Creates valid catalog file
- Proper file permissions
- Correct YAML structure

**Lint Command:**
```bash
backstage-gen-cli lint -f catalog-test.yaml
```
- Validates against Backstage schema
- Checks required fields
- Ensures annotation format

### Running Tests Locally

#### Basic Test Suite
```bash
make test
make lint
make build
```

#### Test Against Real Project
```bash
# Clone a test project
git clone https://github.com/pallets/flask.git /tmp/flask-test
cd /tmp/flask-test

# Run inspection
backstage-gen-cli inspect

# Generate and validate
backstage-gen-cli generate --dry-run
backstage-gen-cli generate -o catalog-info.yaml
backstage-gen-cli lint
```

#### Test Verbosity Flags
```bash
# Quiet mode - no output
backstage-gen-cli inspect -q

# Verbose mode - extra details
backstage-gen-cli inspect -v

# Debug mode - all debug info
backstage-gen-cli inspect -vvv
```

### Viewing Test Results

1. **GitHub Actions Tab**: https://github.com/gautampachnanda101/backstage-gen-cli/actions
2. **Test Summary**: Available in each workflow run
3. **Artifacts**: Download binaries from successful runs

### Adding New Test Projects

To add a new project to test against:

1. Edit `.github/workflows/ci.yml`
2. Add to the `matrix.project` array:
```yaml
- repo: owner/repo-name
  branch: main
  expected_lang: language
  description: "Brief description"
```

## Manual Testing Checklist

Before releasing, manually verify:

- [ ] Interactive wizard works (`-i` flag)
- [ ] Config file creation (`backstage-gen-cli config init`)
- [ ] Hook installation (`backstage-gen-cli hooks install`)
- [ ] Version display (`--version`)
- [ ] Help text is clear (`--help`)
- [ ] Quiet mode suppresses output (`-q`)
- [ ] Verbose mode shows extra info (`-v`, `-vv`, `-vvv`)
- [ ] All annotations are standard Backstage
- [ ] Descriptions are clean text
- [ ] Tags are meaningful and lowercase
