# Getting Started with backstage-gen-cli

## Quick Setup

### Prerequisites
- Go 1.21 or later
- Make (optional but recommended)
- Git

### 1. Build the Binary

```bash
cd backstage-gen-cli

# Download dependencies
go mod download

# Build
make build
# OR without make:
go build -o bin/backstage-gen-cli .
```

### 2. Verify Installation

```bash
./bin/backstage-gen-cli --version
```

### 3. Try It Out

```bash
# Inspect current directory
./bin/backstage-gen-cli inspect

# Generate catalog (dry-run first)
./bin/backstage-gen-cli generate --dry-run

# Generate actual file
./bin/backstage-gen-cli generate

# Validate it
./bin/backstage-gen-cli lint
```

### 4. Install System-Wide (Optional)

```bash
sudo make install
# Now you can use: backstage-gen-cli
```

## Usage Examples

### Generate for Your Project

```bash
cd /path/to/your/project
backstage-gen-cli generate
```

### Validate Existing Catalog

```bash
backstage-gen-cli lint --strict
```

### Set Up Pre-commit Hook

```bash
backstage-gen-cli hooks install
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

- **Languages**: Go, Python, JavaScript, TypeScript, Java, Rust, Ruby, PHP, C#, and more
- **Frameworks**: React, Vue, Django, Flask, Express, Spring Boot, Rails, and more
- **Infrastructure**: Docker, Kubernetes, Helm, Terraform
- **Build Tools**: Maven, Gradle, npm, pip, cargo, Make

## Common Workflows

### 1. New Project Setup

```bash
cd my-new-service
backstage-gen-cli generate
backstage-gen-cli hooks install
git add catalog-info.yaml
git commit -m "Add Backstage catalog"
```

### 2. Bulk Generation

```bash
for dir in services/*/; do
  cd "$dir"
  backstage-gen-cli generate --force
  cd ../..
done
```

### 3. CI/CD Validation

```yaml
# .github/workflows/validate.yml
- name: Validate Catalog
  run: |
    curl -L https://github.com/gautampachnanda101/backstage-gen-cli/releases/latest/download/backstage-gen-cli-linux-amd64 -o backstage-gen-cli
    chmod +x backstage-gen-cli
    ./backstage-gen-cli lint --strict
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

1. Read the full [README](../README.md)
2. Check the [configuration guide](configuration.md)
3. Set up [CI/CD integration](ci-cd-examples.md)
4. Learn about [testing](testing.md)

## Support

- Issues: [GitHub Issues](https://github.com/gautampachnanda101/backstage-gen-cli/issues)
- Docs: See other files in this directory
