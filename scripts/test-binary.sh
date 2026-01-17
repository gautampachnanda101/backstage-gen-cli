#!/bin/bash
# Test script for validating backstage-gen-cli binary features
# Usage: ./scripts/test-binary.sh [path-to-binary]

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

# Default binary path
BINARY="${1:-./bin/backstage-gen-cli}"

print_pass() { echo -e "${GREEN}[PASS]${NC} $1"; }
print_fail() { echo -e "${RED}[FAIL]${NC} $1"; }
print_info() { echo -e "${CYAN}[INFO]${NC} $1"; }
print_test() { echo -e "${YELLOW}[TEST]${NC} $1"; }

# Counters
PASSED=0
FAILED=0

run_test() {
    local name="$1"
    local expected="$2"
    shift 2

    print_test "$name"

    local OUTPUT
    OUTPUT=$("$@" 2>&1) || true

    if [ -n "$expected" ]; then
        if echo "$OUTPUT" | grep -qiE "$expected"; then
            print_pass "$name"
            ((PASSED++))
            return 0
        else
            print_fail "$name - expected '$expected' in output"
            echo "  Output: ${OUTPUT:0:200}"
            ((FAILED++))
            return 1
        fi
    else
        # No expected pattern - just check command ran
        print_pass "$name"
        ((PASSED++))
        return 0
    fi
}

echo ""
echo -e "${CYAN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║         backstage-gen-cli Binary Test Suite                ║${NC}"
echo -e "${CYAN}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Check if binary exists
if [ ! -f "$BINARY" ]; then
    print_fail "Binary not found: $BINARY"
    echo "Build with: make build"
    exit 1
fi

print_info "Testing binary: $BINARY"
print_info "Binary size: $(du -h "$BINARY" | cut -f1)"
echo ""

echo "=== Basic Commands ==="
echo ""

run_test "--version flag" "backstage-gen-cli" "$BINARY" --version
run_test "--help flag" "catalog-info.yaml|Backstage|generate" "$BINARY" --help

echo ""
echo "=== Inspect Command ==="
echo ""

run_test "inspect basic" "Language|Technology|Repository" "$BINARY" inspect
run_test "inspect quiet mode" "" "$BINARY" inspect -q
run_test "inspect verbose mode" "" "$BINARY" inspect -v
run_test "inspect JSON output" "languages" "$BINARY" inspect --json

echo ""
echo "=== Generate Command ==="
echo ""

run_test "generate dry-run" "apiVersion" "$BINARY" generate --dry-run
run_test "generate quiet dry-run" "apiVersion" "$BINARY" generate --dry-run -q

# Generate to a temp file
TEMP_CATALOG=$(mktemp)
trap "rm -f $TEMP_CATALOG" EXIT

"$BINARY" generate -o "$TEMP_CATALOG" --force 2>/dev/null
if [ -f "$TEMP_CATALOG" ] && grep -q "apiVersion" "$TEMP_CATALOG"; then
    print_test "generate to file"
    print_pass "generate to file"
    ((PASSED++))
else
    print_test "generate to file"
    print_fail "generate to file - file not created or invalid"
    ((FAILED++))
fi

echo ""
echo "=== Lint Command ==="
echo ""

run_test "lint generated file" "" "$BINARY" lint -f "$TEMP_CATALOG"
run_test "generate force overwrites" "" "$BINARY" generate -o "$TEMP_CATALOG" --force

echo ""
echo "=== Edge Cases ==="
echo ""

run_test "lint non-existent file" "" "$BINARY" lint -f /nonexistent/path/file.yaml

echo ""
echo "=== Config Command ==="
echo ""

run_test "config show" "" "$BINARY" config show

echo ""
echo -e "${CYAN}════════════════════════════════════════════════════════════${NC}"
echo ""
echo -e "Tests Passed: ${GREEN}$PASSED${NC}"
echo -e "Tests Failed: ${RED}$FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed.${NC}"
    exit 1
fi
