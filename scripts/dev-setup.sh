#!/bin/bash
# Developer setup script for backstage-gen-cli
# Sets up the development environment

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
print_header() { echo -e "${CYAN}$1${NC}"; }

GO_VERSION_REQUIRED="1.21"

# Check Go installation
check_go() {
    print_header "Checking Go installation..."

    if ! command -v go &> /dev/null; then
        print_error "Go is not installed"
        echo "Please install Go ${GO_VERSION_REQUIRED} or later:"
        echo "  https://go.dev/doc/install"
        exit 1
    fi

    GO_VERSION=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+')
    print_status "Go version: $GO_VERSION"

    # Simple version check
    if [[ "$(printf '%s\n' "$GO_VERSION_REQUIRED" "$GO_VERSION" | sort -V | head -n1)" != "$GO_VERSION_REQUIRED" ]]; then
        print_warning "Go version $GO_VERSION may be too old. Recommended: $GO_VERSION_REQUIRED+"
    fi
}

# Check Make
check_make() {
    print_header "Checking Make..."

    if command -v make &> /dev/null; then
        print_status "Make is installed"
    else
        print_warning "Make is not installed (optional but recommended)"
        echo "Install with:"
        echo "  macOS: xcode-select --install"
        echo "  Ubuntu: sudo apt install build-essential"
    fi
}

# Check Git
check_git() {
    print_header "Checking Git..."

    if command -v git &> /dev/null; then
        print_status "Git is installed"
    else
        print_error "Git is not installed"
        exit 1
    fi
}

# Check Docker (optional)
check_docker() {
    print_header "Checking Docker (optional)..."

    if command -v docker &> /dev/null; then
        if docker info &> /dev/null; then
            print_status "Docker is installed and running"
        else
            print_warning "Docker is installed but not running"
        fi
    else
        print_warning "Docker is not installed (optional, needed for LLM features)"
    fi
}

# Download Go dependencies
setup_dependencies() {
    print_header "Downloading Go dependencies..."

    go mod download
    go mod tidy

    print_status "Dependencies downloaded"
}

# Build the project
build_project() {
    print_header "Building project..."

    if command -v make &> /dev/null; then
        make build
    else
        go build -o bin/backstage-gen-cli .
    fi

    print_status "Build complete: bin/backstage-gen-cli"
}

# Run tests
run_tests() {
    print_header "Running tests..."

    if command -v make &> /dev/null; then
        make test
    else
        go test -v ./...
    fi

    print_status "Tests passed"
}

# Run linters
run_lint() {
    print_header "Running linters..."

    if command -v make &> /dev/null; then
        make lint
    else
        go vet ./...
        go fmt ./...
    fi

    print_status "Linting complete"
}

# Install development tools
install_dev_tools() {
    print_header "Installing development tools..."

    # golangci-lint (optional but recommended)
    if ! command -v golangci-lint &> /dev/null; then
        print_status "Installing golangci-lint..."
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    else
        print_status "golangci-lint already installed"
    fi
}

# Create example config
create_example_config() {
    print_header "Creating example configuration..."

    if [ ! -f ".backstage-gen.yaml" ]; then
        ./bin/backstage-gen-cli config init 2>/dev/null || true
        print_status "Example config created"
    else
        print_status "Config already exists"
    fi
}

# Main setup
main() {
    echo ""
    echo -e "${CYAN}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║         backstage-gen-cli Developer Setup                  ║${NC}"
    echo -e "${CYAN}╚════════════════════════════════════════════════════════════╝${NC}"
    echo ""

    check_go
    check_git
    check_make
    check_docker

    echo ""
    setup_dependencies
    build_project
    run_tests
    run_lint

    echo ""
    print_status "Development environment ready!"
    echo ""
    echo "Quick commands:"
    echo "  make build     - Build the binary"
    echo "  make test      - Run tests"
    echo "  make lint      - Run linters"
    echo "  make clean     - Clean build artifacts"
    echo ""
    echo "Test the CLI:"
    echo "  ./bin/backstage-gen-cli --help"
    echo "  ./bin/backstage-gen-cli inspect"
    echo ""
}

# Parse arguments
SKIP_TESTS=false
SKIP_LINT=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --skip-tests)
            SKIP_TESTS=true
            shift
            ;;
        --skip-lint)
            SKIP_LINT=true
            shift
            ;;
        --help|-h)
            echo "Usage: $0 [options]"
            echo ""
            echo "Options:"
            echo "  --skip-tests  Skip running tests"
            echo "  --skip-lint   Skip running linters"
            echo "  --help, -h    Show this help"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Override functions if skipping
if [ "$SKIP_TESTS" = true ]; then
    run_tests() { print_warning "Skipping tests"; }
fi

if [ "$SKIP_LINT" = true ]; then
    run_lint() { print_warning "Skipping lint"; }
fi

main
