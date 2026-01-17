# Contributing Guide

Thank you for your interest in contributing to backstage-gen-cli!

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- Make
- Docker (optional, for LLM features)

### Development Setup

```bash
# Clone the repository
git clone https://github.com/gautampachnanda101/backstage-gen-cli.git
cd backstage-gen-cli

# Install dependencies
go mod download

# Build
make build

# Run tests
make test

# Run linters
make lint
```

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/my-feature
# or
git checkout -b fix/my-bug-fix
```

### 2. Make Changes

- Follow Go best practices and conventions
- Add tests for new functionality
- Update documentation as needed

### 3. Test Your Changes

```bash
# Run unit tests
make test

# Run linters
make lint

# Build the binary
make build

# Test manually
./bin/backstage-gen-cli inspect
./bin/backstage-gen-cli generate --dry-run
```

### 4. Commit Your Changes

We follow conventional commit messages:

```bash
git commit -m "feat: add new language detection for Kotlin"
git commit -m "fix: handle empty README files gracefully"
git commit -m "docs: update configuration examples"
git commit -m "test: add tests for validator package"
```

Prefixes:
- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `test:` - Test additions/changes
- `refactor:` - Code refactoring
- `chore:` - Build/tooling changes

### 5. Create Pull Request

1. Push your branch to GitHub
2. Create a Pull Request against `main`
3. Fill out the PR template
4. Wait for CI checks to pass
5. Request review

## Code Style

### Go Style

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` for formatting
- Use `go vet` for static analysis
- Keep functions focused and small
- Add comments for exported functions

### Project Structure

```
backstage-gen-cli/
├── cmd/           # CLI commands (Cobra)
├── pkg/           # Core packages
│   ├── config/    # Configuration handling
│   ├── detector/  # Repository detection
│   ├── generator/ # Catalog generation
│   ├── llm/       # LLM client
│   ├── output/    # Output/verbosity handling
│   ├── validator/ # Catalog validation
│   └── wizard/    # Interactive wizard
├── docs/          # Documentation
├── scripts/       # Helper scripts
├── packaging/     # Package manager files
└── main.go        # Entry point
```

## Testing Guidelines

### Unit Tests

- Test files go alongside source files (`*_test.go`)
- Use table-driven tests where appropriate
- Mock external dependencies
- Aim for meaningful coverage, not just numbers

Example:
```go
func TestDetectLanguages(t *testing.T) {
    tests := []struct {
        name     string
        files    []string
        expected []string
    }{
        {
            name:     "Go project",
            files:    []string{"go.mod"},
            expected: []string{"Go"},
        },
        // more cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

### Running Tests

```bash
# All tests
make test

# Specific package
go test -v ./pkg/detector/...

# With coverage
go test -cover ./...

# With race detection
go test -race ./...
```

## Documentation

- Update docs/ for user-facing changes
- Add code comments for complex logic
- Update README.md if needed
- Include examples in documentation

## Reporting Issues

### Bug Reports

Include:
- backstage-gen-cli version (`--version`)
- Operating system
- Steps to reproduce
- Expected vs actual behavior
- Any error messages

### Feature Requests

Include:
- Use case description
- Proposed solution
- Alternatives considered

## Pull Request Process

1. Ensure all tests pass
2. Update documentation
3. Add tests for new code
4. Follow commit message conventions
5. Keep PRs focused and reasonably sized
6. Respond to review feedback

## Release Process

Releases are automated via GitHub Actions when a tag is pushed:

```bash
git tag v0.2.0
git push origin v0.2.0
```

This triggers:
1. Cross-platform builds
2. SHA256 checksum generation
3. GitHub Release creation
4. Binary uploads

## Questions?

- Open a [GitHub Issue](https://github.com/gautampachnanda101/backstage-gen-cli/issues)
- Check existing issues and discussions

Thank you for contributing!
