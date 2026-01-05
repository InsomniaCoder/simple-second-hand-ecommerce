# Claude Code Automation Guide

This project demonstrates comprehensive Claude Code automation following Boris's best practices.

## Overview

All automation is configured in the `.claude/` directory:
- **CLAUDE.md** - Team-shared conventions and standards
- **settings.json** - Permissions and hook configuration
- **commands/** - Slash commands for workflows
- **agents/** - Specialized subagents
- **hooks/** - Automation hooks

---

## Slash Commands

### /commit-push-pr
**Purpose**: Complete git workflow with validation

**What it does**:
1. Runs `go fmt ./...` to check formatting
2. Runs `go test ./...` with coverage check (≥80%)
3. Runs `go vet ./...` for static analysis
4. Verifies you're on a feature branch (not main/master)
5. Generates conventional commit message from changes
6. Commits, pushes, and creates PR (if gh CLI available)

**Usage**:
```bash
/commit-push-pr

# Skip tests (not recommended)
/commit-push-pr --skip-tests

# Skip PR creation
/commit-push-pr --skip-pr
```

**When to use**: Ready to commit your changes

---

### /test-and-verify
**Purpose**: Run tests with quality analysis

**What it does**:
1. Runs `go test -v -race -cover ./...`
2. Generates coverage report
3. Runs `go vet` for static analysis
4. Activates `go-test-verifier` subagent for detailed analysis
5. Provides recommendations for missing coverage

**Usage**:
```bash
/test-and-verify

# With verbose output
/test-and-verify --verbose

# Generate HTML coverage
/test-and-verify --coverage-html
```

**When to use**: Before committing, after implementing features

---

### /format-all
**Purpose**: Format all Go files

**What it does**:
1. Runs `go fmt ./...`
2. Runs `goimports -w .` (if available)
3. Reports formatted files

**Usage**:
```bash
/format-all

# Check without applying
/format-all --check-only
```

**When to use**: Before committing, after refactoring

---

### /api-test
**Purpose**: Test all API endpoints

**What it does**:
1. Health check
2. Tests all CRUD operations (POST, GET, PUT, DELETE)
3. Tests error handling (400, 404, 500)
4. Validates responses
5. Generates test report

**Usage**:
```bash
/api-test

# Custom URL
/api-test --url http://localhost:3000

# Verbose output
/api-test --verbose
```

**When to use**: After API changes, before deployment

---

## Subagents

### go-test-verifier
**Purpose**: Deep test coverage and quality analysis

**Activated by**: `/test-and-verify` command

**What it does**:
- Analyzes coverage gaps with risk assessment
- Checks test quality (table-driven, naming, assertions)
- Generates specific test recommendations
- Provides implementation examples
- Uses Sequential MCP for complex analysis

**Example output**:
```
Coverage: 88.7% ✓
Quality Metrics:
  ✓ Table-Driven Tests: 80.5%
  ✓ Proper Naming: 97.4%
  ✓ Race-Free: Yes

Recommendations:
  • Add edge case tests for price boundaries
  • Test error handling paths in service.go:145
```

---

### go-code-simplifier
**Purpose**: Reduce code complexity

**Triggers**: After implementation, code reviews

**What it does**:
- Calculates cyclomatic complexity
- Identifies duplication and deep nesting
- Suggests refactorings with before/after examples
- Applies Go idioms (early returns, error wrapping)
- Ensures tests pass after changes

**Example**:
```
Function: UpdateItem
  Complexity: 12 (HIGH - target: <10)
  Lines: 87 (target: <50)

Recommendations:
  1. Extract validation → -4 complexity
  2. Use early returns → -2 complexity
  3. Extract update logic → -3 complexity
```

---

### api-tester
**Purpose**: Comprehensive API testing

**Triggers**: API implementation, integration testing

**What it does**:
- Tests all CRUD operations
- Validates status codes and JSON responses
- Tests error handling
- Verifies data persistence
- Generates detailed test report

---

## Hooks

### post-tool-use.sh
**Purpose**: Auto-format Go files after edits

**Runs**: Automatically after Edit/Write tools on `.go` files

**What it does**:
1. Runs `go fmt` on the modified file
2. Runs `goimports` if available
3. Silent operation (doesn't block workflow)

**Example**:
```
Edit internal/service/item_service.go
  ↓
🎨 Auto-formatting Go file: internal/service/item_service.go
✓ Formatted and organized imports
```

---

### pre-commit.sh
**Purpose**: Validate code before commits

**Runs**: Automatically before `git commit`

**What it does**:
1. Checks if code is formatted (`go fmt`)
2. Runs tests (`go test -short ./...`)
3. Runs `go vet` for warnings
4. Checks for TODOs and debug prints (warnings only)
5. **Blocks commit if validation fails**

**Example**:
```
git commit -m "Add feature"
  ↓
🔍 Running pre-commit validation...
✓ All files properly formatted
✓ All tests passed
✓ No vet warnings

✅ All validation checks passed
Proceeding with commit...
```

**Bypass** (not recommended):
```bash
git commit --no-verify
```

---

## Daily Workflow

### Morning: Start Development
```bash
1. Check project status
   git status

2. Verify tests pass
   /test-and-verify

3. Review project conventions
   cat .claude/CLAUDE.md
```

### Development: Implement Feature
```bash
4. Write code
   # Files auto-format on save (post-tool-use hook)

5. Write tests alongside
   # Table-driven tests in *_test.go files

6. Verify locally
   /test-and-verify
```

### Commit: Ready to Push
```bash
7. Run commit workflow
   /commit-push-pr

   # This automatically:
   # - Formats code
   # - Runs tests
   # - Generates commit message
   # - Pushes and creates PR
```

### Review: Code Quality
```bash
8. Review PR
   # Check coverage report
   # Verify tests pass in CI
   # Review generated commit message

9. Merge when approved
```

---

## Configuration

### Pre-Allowed Commands
From `.claude/settings.json`:
```json
{
  "permissions": {
    "allow": [
      "Bash(go build*)",
      "Bash(go test*)",
      "Bash(go fmt*)",
      "Bash(git *)",
      "Bash(make*)"
    ]
  }
}
```

These commands run without permission prompts.

### Testing Configuration
```json
{
  "testing": {
    "minCoverage": 80,
    "testCommand": "go test -v -race -cover ./..."
  }
}
```

### Hook Registration
```json
{
  "hooks": {
    "postToolUse": ".claude/hooks/post-tool-use.sh",
    "preCommit": ".claude/hooks/pre-commit.sh"
  }
}
```

---

## MCP Integration

### Sequential MCP
**Used for**: Complex reasoning and analysis

**Examples**:
- Test coverage gap analysis (`go-test-verifier`)
- Code complexity analysis (`go-code-simplifier`)
- Commit message generation (`/commit-push-pr`)

### Context7 MCP
**Used for**: Official Go documentation and patterns

**Examples**:
- "What's the Go idiom for error handling?"
- "How to structure HTTP middleware?"
- "Best practices for concurrent code?"

### Serena MCP
**Used for**: Project memory and session persistence

**Commands**:
- `/sc:load` - Load project context
- `/sc:save` - Save session state
- Memory-driven development workflows

---

## Extending Automation

### Adding New Slash Commands
1. Create `.claude/commands/your-command.md`
2. Follow the frontmatter format:
   ```yaml
   ---
   name: your-command
   description: "What it does"
   category: workflow|testing|utility
   ---
   ```
3. Document usage and behavior

### Adding New Subagents
1. Create `.claude/agents/your-agent.md`
2. Define:
   - Purpose and triggers
   - Behavioral mindset
   - Execution flow
   - Boundaries (will/won't do)

### Adding New Hooks
1. Create `.claude/hooks/your-hook.sh`
2. Make executable: `chmod +x`
3. Register in `.claude/settings.json`

---

## Troubleshooting

### Slash Commands Not Working
- Check `.claude/commands/` directory exists
- Verify frontmatter format is correct
- Restart Claude Code session

### Hooks Not Running
- Check `.claude/settings.json` registration
- Verify hooks are executable (`chmod +x`)
- Check hook script has no syntax errors

### Tests Failing
- Run `/test-and-verify` for detailed analysis
- Check `go test -v ./...` output
- Review `.claude/CLAUDE.md` for standards

### Format Issues
- Run `/format-all` to format all files
- Install `goimports`: `go install golang.org/x/tools/cmd/goimports@latest`
- Check `go fmt ./...` output

---

## Success Metrics

Your automation is working if:
- ✅ Files auto-format after edits
- ✅ Pre-commit hook blocks bad commits
- ✅ `/commit-push-pr` successfully validates
- ✅ `/test-and-verify` provides clear reports
- ✅ Coverage maintained at ≥80%
- ✅ No manual formatting needed
- ✅ Conventional commits generated automatically

---

## Resources

- **Project Conventions**: `.claude/CLAUDE.md`
- **API Documentation**: `docs/API.md`
- **README**: `README.md`
- **Boris's Blog Post**: Original Claude Code best practices
