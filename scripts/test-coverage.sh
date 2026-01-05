#!/bin/bash
# Generate test coverage report with HTML visualization

echo "📊 Generating coverage report..."

# Run tests with coverage
go test -coverprofile=coverage.out ./...

if [ $? -ne 0 ]; then
    echo "❌ Tests failed"
    exit 1
fi

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# Display function-level coverage
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📈 Coverage by Function"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
go tool cover -func=coverage.out

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Coverage report generated"
echo "   • coverage.out   - Coverage data"
echo "   • coverage.html  - HTML visualization"
echo ""
echo "Open coverage.html in your browser to view detailed report"
