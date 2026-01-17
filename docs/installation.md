# Installation Guide

## Quick Install

### macOS (Homebrew)

```bash
brew tap gautampachnanda101/backstage-gen-cli
brew install backstage-gen-cli
```

### Windows (Scoop)

```powershell
scoop bucket add backstage-gen-cli https://github.com/gautampachnanda101/scoop-backstage-gen-cli
scoop install backstage-gen-cli
```

### Windows (Chocolatey)

```powershell
choco install backstage-gen-cli
```

### Linux/macOS (Binary Download)

```bash
# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map architecture names
case $ARCH in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
esac

# Download latest release
curl -L "https://github.com/gautampachnanda101/backstage-gen-cli/releases/latest/download/backstage-gen-cli-${OS}-${ARCH}" -o backstage-gen-cli
chmod +x backstage-gen-cli
sudo mv backstage-gen-cli /usr/local/bin/
```

### Windows (PowerShell)

```powershell
# Download latest release
Invoke-WebRequest -Uri "https://github.com/gautampachnanda101/backstage-gen-cli/releases/latest/download/backstage-gen-cli-windows-amd64.exe" -OutFile "backstage-gen-cli.exe"

# Add to PATH or move to a directory in PATH
Move-Item backstage-gen-cli.exe C:\Windows\System32\
```

## Build from Source

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional)

### Steps

```bash
# Clone the repository
git clone https://github.com/gautampachnanda101/backstage-gen-cli.git
cd backstage-gen-cli

# Build
make build

# Or without make
go build -o bin/backstage-gen-cli .

# Install system-wide
sudo make install
# Or manually
sudo cp bin/backstage-gen-cli /usr/local/bin/
```

## Verify Installation

```bash
backstage-gen-cli --version
```

Expected output:
```
backstage-gen-cli version v0.1.0 (commit: abc1234, built: 2025-01-17T00:00:00Z)
```

## Platform-Specific Notes

### macOS

If you see a security warning on macOS:
1. Go to System Preferences > Security & Privacy > General
2. Click "Allow Anyway" for backstage-gen-cli
3. Or use: `xattr -d com.apple.quarantine /usr/local/bin/backstage-gen-cli`

### Linux

Ensure `/usr/local/bin` is in your PATH:
```bash
echo 'export PATH=$PATH:/usr/local/bin' >> ~/.bashrc
source ~/.bashrc
```

### Windows

Add the installation directory to your PATH environment variable if needed.

## Uninstalling

### Homebrew
```bash
brew uninstall backstage-gen-cli
brew untap gautampachnanda101/backstage-gen-cli
```

### Scoop
```powershell
scoop uninstall backstage-gen-cli
```

### Chocolatey
```powershell
choco uninstall backstage-gen-cli
```

### Binary
```bash
sudo rm /usr/local/bin/backstage-gen-cli
```

## Updating

### Homebrew
```bash
brew upgrade backstage-gen-cli
```

### Scoop
```powershell
scoop update backstage-gen-cli
```

### Chocolatey
```powershell
choco upgrade backstage-gen-cli
```

### Binary
Re-run the installation commands to download the latest version.
