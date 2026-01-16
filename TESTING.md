# Testing Documentation

## GitHub Actions CI/CD Testing

Our comprehensive CI/CD pipeline tests the CLI against real-world open source projects to ensure reliability and prevent regressions.

### Test Coverage

#### 1. Make Tasks Testing
- ✅ `make build` - Builds the binary successfully
- ✅ `make test` - Runs all Go tests
- ✅ `make lint` - Runs linters (go vet, go fmt)
- ✅ `make clean` - Cleans build artifacts

#### 2. Cross-Platform Builds
Builds binaries for all major platforms:
- ✅ Linux (amd64, arm64)
- ✅ macOS (amd64, arm64)
- ✅ Windows (amd64)

All binaries include version info and are uploaded as artifacts.

#### 3. Real-World Project Testing

We test against 10+ diverse open source projects to ensure the CLI works across different:
- Programming languages
- Framework types
- Repository structures
- README formats

**Projects Tested:**

| Project | Language | Type | What We Validate |
|---------|----------|------|------------------|
| **kubernetes/kubernetes** | Go | Orchestration | Go detection, large repos |
| **django/django** | Python | Web Framework | Python, framework detection |
| **pallets/flask** | Python | Microframework | Description parsing, clean text |
| **facebook/react** | JavaScript | UI Library | JS/TS detection, monorepos |
| **vercel/next.js** | JavaScript/TypeScript | Framework | Mixed language detection |
| **vuejs/core** | JavaScript/TypeScript | Framework | Vue detection, annotations |
| **spring-projects/spring-boot** | Java | Framework | Java, Maven detection |
| **rust-lang/rust** | Rust | Compiler | Rust detection, complex repos |
| **rails/rails** | Ruby | Framework | Ruby, Rails detection |
| **grafana/grafana** | Go/TypeScript | Multi-lang | Multi-language handling |

#### 4. Command Testing

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
- Contains appropriate tags

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
- Verifies metadata structure

#### 5. Quality Checks

**Language Detection Accuracy:**
- Verifies detected language matches expected
- Ensures primary language is correctly identified
- Validates multi-language project handling

**Annotation Quality:**
- `backstage.io/source-location` - Git URL
- `backstage.io/techdocs-ref` - Documentation reference
- `github.com/project-slug` - Extracted from remote
- No non-standard annotations

**Description Quality:**
- No HTML tags or markup
- No image/badge URLs
- Plain text only
- Reasonable length (30-200 chars)
- Meaningful content

**Tag Appropriateness:**
- Language tags present
- Framework tags when applicable
- Infrastructure tags (docker, k8s) when detected
- No duplicate or empty tags

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

#### Test Specific Features

**Clean Description Parsing:**
```bash
# Should extract clean text without HTML/markdown
backstage-gen-cli inspect | grep "Description:"
```

**Annotation Validation:**
```bash
# Should include standard annotations only
backstage-gen-cli generate --dry-run | grep "annotations:"
```

**Multi-Language Detection:**
```bash
# Test on a multi-language repo
cd /path/to/mixed/language/repo
backstage-gen-cli inspect | grep "Language:"
```

### Regression Testing

The GitHub Actions workflow automatically runs on every push to `main` and on all pull requests. This ensures:

1. **No Breaking Changes** - All commands work as expected
2. **Consistent Output** - Generated catalogs remain valid
3. **Cross-Platform Compatibility** - Builds work on all platforms
4. **Real-World Validation** - Works with actual OSS projects

### Test Artifacts

Each CI run produces:
- Cross-platform binaries (7 days retention)
- Release binaries with checksums (30 days retention)
- Test summary reports
- Language detection accuracy reports

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
3. Commit and push - CI will automatically test

### Known Limitations

- Description parsing may fall back to generic for repos with complex READMEs
- Some edge cases in multi-line HTML removal
- LLM features require local setup (not tested in CI)

## Manual Testing Checklist

Before releasing, manually verify:

- [ ] Interactive wizard works (`-i` flag)
- [ ] Config file creation (`backstage-gen-cli config init`)
- [ ] Hook installation (`backstage-gen-cli hooks install`)
- [ ] Version display (`--version`)
- [ ] Help text is clear (`--help`)
- [ ] Binary naming is consistent
- [ ] All annotations are standard Backstage
- [ ] Descriptions are clean text
- [ ] Tags are meaningful and lowercase
