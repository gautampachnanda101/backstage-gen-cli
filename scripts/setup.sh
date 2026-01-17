#!/bin/bash
# Setup script for backstage-gen-cli
# Downloads and installs the latest release

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

print_status() { echo -e "${GREEN}[INFO]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARN]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1"; }

GITHUB_REPO="gautampachnanda101/backstage-gen-cli"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="backstage-gen-cli"

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case $ARCH in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac

    case $OS in
        linux|darwin)
            ;;
        mingw*|msys*|cygwin*)
            OS="windows"
            BINARY_NAME="backstage-gen-cli.exe"
            ;;
        *)
            print_error "Unsupported OS: $OS"
            exit 1
            ;;
    esac

    PLATFORM="${OS}-${ARCH}"
    print_status "Detected platform: $PLATFORM"
}

# Get latest release version
get_latest_version() {
    print_status "Fetching latest version..."
    VERSION=$(curl -s "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

    if [ -z "$VERSION" ]; then
        print_error "Failed to get latest version. Using 'latest' tag."
        VERSION="latest"
    else
        print_status "Latest version: $VERSION"
    fi
}

# Download binary
download_binary() {
    local url="https://github.com/${GITHUB_REPO}/releases/latest/download/${BINARY_NAME}-${PLATFORM}"
    local tmp_file=$(mktemp)

    print_status "Downloading from: $url"

    if command -v curl &> /dev/null; then
        curl -fsSL "$url" -o "$tmp_file"
    elif command -v wget &> /dev/null; then
        wget -q "$url" -O "$tmp_file"
    else
        print_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi

    if [ ! -f "$tmp_file" ] || [ ! -s "$tmp_file" ]; then
        print_error "Download failed"
        exit 1
    fi

    echo "$tmp_file"
}

# Install binary
install_binary() {
    local tmp_file=$1

    # Make executable
    chmod +x "$tmp_file"

    # Check if we need sudo
    if [ -w "$INSTALL_DIR" ]; then
        mv "$tmp_file" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        print_status "Installing to ${INSTALL_DIR} (requires sudo)"
        sudo mv "$tmp_file" "${INSTALL_DIR}/${BINARY_NAME}"
    fi

    print_status "Installed to: ${INSTALL_DIR}/${BINARY_NAME}"
}

# Verify installation
verify_installation() {
    if command -v "$BINARY_NAME" &> /dev/null; then
        print_status "Verifying installation..."
        VERSION=$("$BINARY_NAME" --version 2>/dev/null || echo "unknown")
        print_status "Installed version: $VERSION"
    else
        print_warning "Binary installed but not in PATH"
        echo "Add ${INSTALL_DIR} to your PATH:"
        echo "  export PATH=\$PATH:${INSTALL_DIR}"
    fi
}

# Optional: Install Docker
install_docker_hint() {
    if ! command -v docker &> /dev/null; then
        echo ""
        print_warning "Docker is not installed"
        echo "Docker is optional but required for LLM features (LiteLLM)"
        echo "Install Docker: https://docs.docker.com/get-docker/"
    fi
}

# Main installation
main() {
    echo ""
    echo -e "${CYAN}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║         backstage-gen-cli Installation                     ║${NC}"
    echo -e "${CYAN}╚════════════════════════════════════════════════════════════╝${NC}"
    echo ""

    detect_platform
    get_latest_version

    TMP_FILE=$(download_binary)
    install_binary "$TMP_FILE"
    verify_installation
    install_docker_hint

    echo ""
    print_status "Installation complete!"
    echo ""
    echo "Get started:"
    echo "  backstage-gen-cli --help"
    echo "  backstage-gen-cli inspect"
    echo "  backstage-gen-cli generate --dry-run"
    echo ""
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --install-dir)
            INSTALL_DIR="$2"
            shift 2
            ;;
        --help|-h)
            echo "Usage: $0 [options]"
            echo ""
            echo "Options:"
            echo "  --install-dir DIR  Installation directory (default: /usr/local/bin)"
            echo "  --help, -h         Show this help"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

main
