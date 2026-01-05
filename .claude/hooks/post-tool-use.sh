#!/bin/bash
# Post-tool-use hook: Auto-format Go code after edits
# Runs automatically after Edit/Write tools are used
# Ensures consistent code formatting without manual intervention

# Only run for Edit and Write tools
if [[ "$TOOL_NAME" != "Edit" && "$TOOL_NAME" != "Write" ]]; then
    exit 0
fi

# Only run for Go files
if [[ "$FILE_PATH" != *.go ]]; then
    exit 0
fi

# Check if file exists
if [[ ! -f "$FILE_PATH" ]]; then
    exit 0
fi

echo "🎨 Auto-formatting Go file: $FILE_PATH"

# Run go fmt
if ! go fmt "$FILE_PATH" >/dev/null 2>&1; then
    echo "⚠️  go fmt failed for $FILE_PATH"
    exit 0  # Don't fail the operation
fi

# Run goimports if available
if command -v goimports &> /dev/null; then
    if ! goimports -w "$FILE_PATH" >/dev/null 2>&1; then
        echo "⚠️  goimports failed for $FILE_PATH"
        exit 0
    fi
    echo "✓ Formatted and organized imports"
else
    echo "✓ Formatted (install goimports for import organization)"
fi

exit 0
