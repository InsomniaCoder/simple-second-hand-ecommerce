#!/bin/bash
# Pre-commit hook: Validate code before allowing commit
# Ensures code quality and prevents broken code from entering repository
# Blocks commit if validation fails

set -e  # Exit on any error

echo "🔍 Running pre-commit validation..."

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Track validation status
VALIDATION_PASSED=true

# 1. Format Check
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 Checking code formatting..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check if any files need formatting
UNFORMATTED=$(gofmt -l . 2>/dev/null | grep -v vendor || true)

if [[ -n "$UNFORMATTED" ]]; then
    echo -e "${RED}✗ Code needs formatting${NC}"
    echo ""
    echo "Files needing formatting:"
    echo "$UNFORMATTED"
    echo ""
    echo "Run: go fmt ./..."
    echo "Or: /format-all"
    VALIDATION_PASSED=false
else
    echo -e "${GREEN}✓ All files properly formatted${NC}"
fi

# 2. Run Tests
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🧪 Running tests..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if go test -short ./... > /dev/null 2>&1; then
    echo -e "${GREEN}✓ All tests passed${NC}"
else
    echo -e "${RED}✗ Tests failed${NC}"
    echo ""
    echo "Run: go test -v ./..."
    echo "Or: /test-and-verify"
    VALIDATION_PASSED=false
fi

# 3. Run go vet
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔬 Running static analysis..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if go vet ./... > /dev/null 2>&1; then
    echo -e "${GREEN}✓ No vet warnings${NC}"
else
    echo -e "${RED}✗ Go vet found issues${NC}"
    echo ""
    echo "Run: go vet ./..."
    VALIDATION_PASSED=false
fi

# 4. Check for common issues
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔍 Checking for common issues..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check for TODO in production code (excluding test files and comments explaining TODOs)
TODOS=$(grep -r "TODO" --include="*.go" --exclude="*_test.go" . | grep -v "claude/hooks" | grep -v "^#" || true)
if [[ -n "$TODOS" ]]; then
    echo -e "${YELLOW}⚠️  TODOs found in production code${NC}"
    echo ""
    echo "$TODOS"
    echo ""
    echo "Consider resolving TODOs before committing"
    # Don't block commit for TODOs, just warn
fi

# Check for fmt.Println (use proper logging instead)
DEBUG_PRINTS=$(grep -r "fmt.Println" --include="*.go" --exclude="*_test.go" . | grep -v "claude/hooks" || true)
if [[ -n "$DEBUG_PRINTS" ]]; then
    echo -e "${YELLOW}⚠️  Debug print statements found${NC}"
    echo ""
    echo "$DEBUG_PRINTS"
    echo ""
    echo "Consider using proper logging (log.Printf)"
    # Don't block commit, just warn
fi

echo -e "${GREEN}✓ Common issues check complete${NC}"

# Final result
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [[ "$VALIDATION_PASSED" == true ]]; then
    echo -e "${GREEN}✅ All validation checks passed${NC}"
    echo ""
    echo "Proceeding with commit..."
    exit 0
else
    echo -e "${RED}❌ Validation failed${NC}"
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "🚫 Commit blocked"
    echo ""
    echo "Fix the issues above and try again."
    echo ""
    echo "Quick fixes:"
    echo "  • Format code:     /format-all"
    echo "  • Run tests:       /test-and-verify"
    echo "  • Fix vet issues:  go vet ./..."
    echo ""
    echo "Or bypass validation (NOT recommended):"
    echo "  git commit --no-verify"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    exit 1
fi
