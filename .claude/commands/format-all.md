---
name: format-all
description: "Format and organize imports for all Go files"
category: utility
complexity: simple
mcp-servers: []
personas: []
---

# /format-all - Code Formatting

## Purpose
Format all Go files and organize imports consistently across the project. Ensures code style compliance before commits.

## Triggers
- Before committing changes
- After refactoring sessions
- Code style standardization needs
- Manual formatting requests

## Usage
```bash
/format-all [--check-only] [--verbose]
```

### Flags
- `--check-only`: Check if formatting needed without applying changes
- `--verbose`: Show all formatted files

---

## Execution Steps

### Step 1: Go Format
```bash
go fmt ./...
```

**What it does**:
- Formats all `.go` files recursively
- Applies standard Go formatting rules
- Fixes indentation, spacing, alignment
- Returns list of files that were reformatted

### Step 2: Organize Imports (if goimports available)
```bash
goimports -w .
```

**What it does**:
- Organizes imports in standard order:
  1. Standard library
  2. External packages
  3. Internal packages
- Removes unused imports
- Adds missing imports

---

## Example Output

### Standard Execution
```
User: /format-all

Claude:
🎨 Formatting Go code...

Running: go fmt ./...
✓ Formatted 3 files:
  • internal/service/item_service.go
  • internal/handlers/handlers.go
  • cmd/api/main.go

Running: goimports -w .
✓ Organized imports in 5 files:
  • internal/repository/memory.go
  • internal/service/item_service.go
  • internal/handlers/handlers.go
  • internal/handlers/routes.go
  • cmd/api/main.go

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ Formatting complete

Total files processed: 8
Files formatted: 3
Imports organized: 5

All files now conform to Go standards.
```

### Check-Only Mode
```
User: /format-all --check-only

Claude:
🔍 Checking code formatting...

Files needing formatting:
  • internal/service/item_service.go
  • internal/handlers/handlers.go
  • cmd/api/main.go

Run without --check-only to apply formatting.
```

### Already Formatted
```
User: /format-all

Claude:
🎨 Formatting Go code...

Running: go fmt ./...
✓ All files already formatted

Running: goimports -w .
✓ All imports properly organized

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ No changes needed

All files conform to Go standards.
```

---

## Tool Coordination

### Tools Used
- **Bash**: Execute go fmt and goimports
- **Glob**: Find all .go files
- **Read** (optional): Verify formatting changes

### Integration
- **post-tool-use hook**: Auto-formats files after edits
- **pre-commit hook**: Verifies formatting before commit
- **/commit-push-pr**: Checks formatting as part of validation

---

## Configuration

### Install goimports (optional but recommended)
```bash
go install golang.org/x/tools/cmd/goimports@latest
```

### IDE Integration
Most IDEs support format-on-save:
- **VS Code**: `"editor.formatOnSave": true`
- **GoLand**: Settings → Tools → File Watchers → go fmt
- **Vim**: Use `vim-go` with `:GoFmt`

---

## Boundaries

### Will Do ✓
- Format all Go files with `go fmt`
- Organize imports with `goimports`
- Report formatted files
- Check formatting without applying changes

### Will Not Do ✗
- Format non-Go files
- Modify code logic or behavior
- Add/remove code beyond imports
- Override custom formatting in special cases

---

## Quick Reference

```bash
# Format all files
/format-all

# Check without applying
/format-all --check-only

# Verbose output
/format-all --verbose

# Format specific package
go fmt ./internal/service

# Format single file
go fmt internal/service/item_service.go
```

---

## Success Criteria

- ✅ All Go files formatted
- ✅ Imports organized
- ✅ No formatting errors
- ✅ Ready for commit
