# Homebrew Tap for backstage-gen-cli

This directory contains the Homebrew formula for backstage-gen-cli.

## Setting Up the Tap Repository

To distribute via Homebrew, create a separate repository named `homebrew-backstage-gen-cli`:

1. Create repository: `https://github.com/gautampachnanda101/homebrew-backstage-gen-cli`

2. Copy the formula:
   ```bash
   cp backstage-gen-cli.rb /path/to/homebrew-backstage-gen-cli/Formula/
   ```

3. Update SHA256 checksums:
   ```bash
   # Get checksums from release assets
   curl -L https://github.com/gautampachnanda101/backstage-gen-cli/releases/download/v0.1.0-alpha.1/backstage-gen-cli-darwin-arm64.sha256
   ```

4. Push to GitHub

## Installing

Users can then install with:

```bash
brew tap gautampachnanda101/backstage-gen-cli
brew install backstage-gen-cli
```

## Updating the Formula

After each release:

1. Update `version` in the formula
2. Download new binaries and compute SHA256:
   ```bash
   shasum -a 256 backstage-gen-cli-darwin-amd64
   shasum -a 256 backstage-gen-cli-darwin-arm64
   shasum -a 256 backstage-gen-cli-linux-amd64
   shasum -a 256 backstage-gen-cli-linux-arm64
   ```
3. Update SHA256 placeholders in the formula
4. Commit and push

## Automated Updates

Consider using GitHub Actions to automatically update the formula on new releases:

```yaml
# .github/workflows/update-homebrew.yml
name: Update Homebrew Formula

on:
  release:
    types: [published]

jobs:
  update:
    runs-on: ubuntu-latest
    steps:
      - name: Update formula
        run: |
          # Script to update formula with new version and checksums
```

## Testing Locally

```bash
# Test the formula
brew install --build-from-source ./backstage-gen-cli.rb

# Audit the formula
brew audit --strict ./backstage-gen-cli.rb
```
