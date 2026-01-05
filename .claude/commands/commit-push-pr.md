---
name: commit-push-pr
description: "Automated git workflow with comprehensive pre-checks for Go projects"
category: workflow
complexity: enhanced
mcp-servers: []
personas: []
---

# /commit-push-pr - Automated Git Workflow

## Purpose
Orchestrates complete git workflow with validation, intelligent commit message generation, and PR creation. Implements Boris's best practice of pre-commit validation with automation.

## Triggers
- Ready to commit and push changes
- Need intelligent conventional commit message generation
- Want automated PR creation
- Require comprehensive pre-commit validation

## Usage
```bash
/commit-push-pr [--skip-tests] [--skip-pr] [--force]
```

### Flags
- `--skip-tests`: Skip test execution (not recommended)
- `--skip-pr`: Skip PR creation
- `--force`: Override validation failures (dangerous)

---

## Execution Flow

### Phase 1: Pre-Commit Validation

#### 1.1 Format Check
```bash
go fmt ./...
```
**Purpose**: Ensure code is consistently formatted
**Failure**: Code needs formatting (must run `go fmt`)

#### 1.2 Test Execution
```bash
go test -race -cover ./...
```
**Purpose**: All tests must pass before commit
**Requirements**:
- All tests passing
- No race conditions
- Coverage ≥80%

#### 1.3 Static Analysis
```bash
go vet ./...
```
**Purpose**: Catch common Go mistakes
**Failure**: Vet warnings must be fixed

#### 1.4 Git Status Check
```bash
git status
git branch --show-current
```
**Purpose**: Verify repository state
**Requirements**:
- Must be on feature branch (not main/master)
- Check for untracked files
- Warn about large file additions

---

### Phase 2: Commit Creation

#### 2.1 Analyze Changes
```bash
# Check if anything is staged
git diff --cached --name-only

# If nothing staged, stage all changes (with confirmation)
if [ -z "$(git diff --cached --name-only)" ]; then
    echo "No changes staged. Stage all changes? (yes/no)"
    # If yes: git add .
fi
```

#### 2.2 Generate Commit Message

**Analysis Process**:
1. **Determine Type**: Analyze changed files
   - `internal/models/` → type: feat/refactor
   - `*_test.go` → type: test
   - `docs/` → type: docs
   - Bug fixes → type: fix

2. **Identify Scope**: Extract from file paths
   - `internal/handlers/` → scope: api
   - `internal/service/` → scope: service
   - `internal/repository/` → scope: repository

3. **Generate Subject**: Concise description
   - Imperative mood ("add", not "added")
   - Lowercase after colon
   - No period at end
   - Max 72 characters

4. **Create Body** (optional):
   - Explain "why" not "what"
   - List significant changes
   - Reference issues

**Conventional Commit Format**:
```
<type>(<scope>): <subject>

<body>

<footer>
```

**Example Generated Message**:
```
feat(api): add item update endpoint

Implements PUT /items/{id} with validation and
partial update support. Includes comprehensive
tests and error handling.

Closes #123
```

#### 2.3 Commit Execution
```bash
git commit -m "<generated message>"
```

---

### Phase 3: Push & PR Creation

#### 3.1 Push to Remote
```bash
git push origin <current-branch>
```

#### 3.2 Create Pull Request (if gh CLI available)
```bash
# Check if gh CLI is installed
if command -v gh &> /dev/null; then
    gh pr create \
        --title "<commit subject>" \
        --body "<commit body>"
else
    echo "⚠️  gh CLI not found. Skipping PR creation."
    echo "Create PR manually at: <repository URL>"
fi
```

---

## Example Execution

### Scenario: Adding new API endpoint

```
User: /commit-push-pr

Claude:
🔍 Running pre-commit validation...

📋 Format Check
✓ Code is properly formatted (no changes needed)

🧪 Test Execution
Running: go test -race -cover ./...
✓ All tests passed (132/132)
✓ No race conditions detected
✓ Coverage: 87.3% (target: 80%)

🔬 Static Analysis
Running: go vet ./...
✓ No issues found

📊 Git Status
Current branch: feat/add-update-endpoint
✓ On feature branch (not main/master)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📝 Analyzing changes...

Modified files:
  • internal/handlers/handlers.go
  • internal/handlers/handlers_test.go
  • internal/handlers/routes.go
  • docs/API.md

Change analysis:
  Type: feat (new feature added)
  Scope: api (HTTP handlers modified)
  Files: 4 modified, 187 lines added, 12 lines removed

💬 Generated commit message:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
feat(api): add item update endpoint

Implements PUT /items/{id} endpoint with validation
and partial update support. Includes comprehensive
tests and API documentation.

Changes:
- Add UpdateItem handler with request validation
- Implement partial update logic in service layer
- Add integration tests for update scenarios
- Update API documentation with examples
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🤔 Proceed with commit and push? (yes/no/edit)

User: yes

✓ Changes committed: 3f2a8b9
✓ Pushed to origin/feat/add-update-endpoint

🔗 Creating pull request...
✓ PR created: https://github.com/user/repo/pull/42

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ Workflow complete!

Next steps:
  1. Review PR: https://github.com/user/repo/pull/42
  2. Request code review
  3. Address feedback if needed
```

---

## Validation Failures

### Example: Tests Failing
```
Claude:
🔍 Running pre-commit validation...

✓ Format check passed
❌ Test execution failed

Test failures:
  FAIL: TestItemService_CreateItem/negative_price
  FAIL: TestItemService_UpdateItem/missing_title

Coverage: 76.2% (below 80% target)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🚫 Cannot commit: Tests must pass

Recommendations:
  1. Fix failing tests in internal/service/item_service_test.go
  2. Add missing tests to reach 80% coverage
  3. Run: go test -v ./internal/service to see details

Run /test-and-verify for detailed analysis.
```

### Example: Not on Feature Branch
```
Claude:
🔍 Running pre-commit validation...

✓ Format check passed
✓ Tests passed
✓ Vet passed

❌ Git status check failed

Current branch: main

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🚫 Cannot commit: Must use feature branch

Recommendations:
  1. Create feature branch: git checkout -b feat/your-feature
  2. Move changes to new branch
  3. Run /commit-push-pr again

Never commit directly to main/master.
```

---

## Tool Coordination

### Tools Used
- **Bash**: Git operations, Go toolchain execution
- **Read**: Analyze git diff for commit message generation
- **Grep**: Parse test output and coverage reports
- **Sequential**: Complex analysis for commit message generation

### Integration with Other Automation
- **pre-commit hook**: May run as part of git commit
- **go-test-verifier subagent**: Detailed test analysis if failures occur
- **post-tool-use hook**: Already formatted files before this runs

---

## Configuration

From `.claude/settings.json`:
```json
{
  "preCommitChecks": {
    "format": true,
    "tests": true,
    "vet": true,
    "minCoverage": 80
  },
  "commitMessage": {
    "format": "conventional",
    "includeIssueNumber": true
  }
}
```

---

## Boundaries

### Will Do ✓
- Run comprehensive pre-commit validation
- Generate conventional commit messages from analysis
- Push changes to remote
- Create pull requests (if gh CLI available)
- **Block commits if validation fails**

### Will Not Do ✗
- Commit without passing tests (unless --force)
- Skip validation without explicit flag
- Modify code to fix test failures
- Force push to main/master
- Squash or rebase commits
- Amend existing commits without confirmation

---

## Troubleshooting

### "gh CLI not found"
**Solution**: Install GitHub CLI
```bash
# macOS
brew install gh

# Linux
sudo apt install gh

# Authenticate
gh auth login
```

### "Coverage below threshold"
**Solution**: Add more tests
```bash
# Find uncovered code
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run detailed analysis
/test-and-verify
```

### "Vet warnings found"
**Solution**: Fix reported issues
```bash
# See detailed warnings
go vet ./...

# Common fixes:
# - Remove unused variables
# - Fix printf format strings
# - Address shadow variables
```

---

## Best Practices

1. **Always run validation**: Don't use --skip-tests unless truly necessary
2. **Feature branches**: Never commit directly to main/master
3. **Atomic commits**: One logical change per commit
4. **Descriptive messages**: Generated messages should be clear and actionable
5. **Test first**: Ensure tests pass before committing

## Success Criteria

- ✅ All validation checks pass
- ✅ Conventional commit message generated
- ✅ Changes pushed to remote
- ✅ PR created (if gh available)
- ✅ No manual intervention required for happy path
