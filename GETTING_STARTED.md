# Getting Started with backstage-gen

## Quick Setup

### Prerequisites
- Go 1.21 or later
- Make (optional but recommended)
- Git

### 1. Build the Binary

```bash
cd backstage-gen

# Download dependencies
go mod download

# Build
make build
# OR without make:
go build -o bin/backstage-gen .
```

### 2. Verify Installation

```bash
./bin/backstage-gen --version
```

### 3. Try It Out

```bash
# Inspect current directory
./bin/backstage-gen inspect

# Generate catalog (dry-run first)
./bin/backstage-gen generate --dry-run

# Generate actual file
./bin/backstage-gen generate

# Validate it
./bin/backstage-gen lint
```

### 4. Install System-Wide (Optional)

```bash
sudo make install
# Now you can use: backstage-gen
```

## Usage Examples

### Generate for Your Project

```bash
cd /path/to/your/project
backstage-gen generate
```

### Validate Existing Catalog

```bash
backstage-gen lint --strict
```

### Set Up Pre-commit Hook

```bash
backstage-gen hooks install
```

### Customize Configuration

Create `.backstage-gen.yaml`:

```yaml
organization:
  name: "myorg"
  namespace: "default"

defaults:
  owner: "platform-team"
  lifecycle: "production"
```

## What It Detects

- **Languages**: Go, Python, JavaScript, TypeScript, Java, Rust, Ruby, PHP, C#
- **Frameworks**: React, Vue, Django, Flask, Express, Spring Boot
- **Infrastructure**: Docker, Kubernetes, Helm, Terraform
- **Build Tools**: Maven, Gradle, npm, pip, cargo

## Common Workflows

### 1. New Project Setup

```bash
cd my-new-service
backstage-gen generate
backstage-gen hooks install
git add catalog-info.yaml
git commit -m "Add Backstage catalog"
```

### 2. Bulk Generation

```bash
for dir in services/*/; do
  cd "$dir"
  backstage-gen generate --force
  cd ../..
done
```

### 3. CI/CD Validation

```yaml
# .github/workflows/validate.yml
- name: Validate Catalog
  run: |
    curl -L URL/backstage-gen -o backstage-gen
    chmod +x backstage-gen
    ./backstage-gen lint --strict
```

## Troubleshooting

### Build Fails

```bash
# Clean and retry
rm -rf bin/
go mod tidy
go mod download
make build
```

### Import Errors

Update module path in all files:
```bash
# Replace yourusername with your GitHub username
find . -name "*.go" -exec sed -i 's/yourusername/YOUR_USERNAME/g' {} \;
```

## Next Steps

1. Read full README.md
2. Customize configuration
3. Try on your projects
4. Set up pre-commit hooks
5. Integrate with CI/CD

## Support

- Issues: GitHub Issues
- Docs: README.md
- Examples: See examples/ directory

Happy cataloging! 🚀
