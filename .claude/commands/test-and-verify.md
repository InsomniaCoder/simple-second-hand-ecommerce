---
name: test-and-verify
description: "Run Go tests with coverage analysis and quality verification"
category: testing
complexity: enhanced
mcp-servers: [sequential]
personas: []
---

# /test-and-verify - Comprehensive Test Execution

## Purpose
Execute tests with race detection, coverage analysis, and quality verification. Activates `go-test-verifier` subagent for detailed analysis and recommendations.

## Triggers
- Before committing changes
- After implementing features
- Coverage analysis requests
- Test quality verification needs

## Usage
```bash
/test-and-verify [--verbose] [--coverage-html] [--quick]
```

### Flags
- `--verbose`: Show detailed test output
- `--coverage-html`: Generate HTML coverage report
- `--quick`: Run short tests only (skip slow integration tests)

---

## Execution Steps

### Step 1: Test Execution with Race Detection
```bash
go test -v -race -cover ./...
```

**What it does**:
- Runs all tests with verbose output
- Detects race conditions with `-race` flag
- Calculates code coverage with `-cover`
- Tests all packages recursively

**Expected output**:
```
=== RUN   TestItemRepository_Create
=== RUN   TestItemRepository_Create/valid_item
=== RUN   TestItemRepository_Create/duplicate_id
--- PASS: TestItemRepository_Create (0.00s)
    --- PASS: TestItemRepository_Create/valid_item (0.00s)
    --- PASS: TestItemRepository_Create/duplicate_id (0.00s)
PASS
coverage: 92.5% of statements
ok      github.com/yourusername/go-crud-api-reference/internal/repository    0.123s
```

---

### Step 2: Coverage Report Generation
```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

**Purpose**: Detailed per-function coverage analysis

**Output format**:
```
github.com/.../repository/memory.go:15:     NewMemoryRepository      100.0%
github.com/.../repository/memory.go:22:     Create                   95.5%
github.com/.../repository/memory.go:35:     GetByID                  100.0%
github.com/.../repository/memory.go:47:     List                     100.0%
github.com/.../repository/memory.go:57:     Update                   87.5%
github.com/.../repository/memory.go:72:     Delete                   100.0%
total:                                      (statements)             92.5%
```

---

### Step 3: Static Analysis
```bash
go vet ./...
```

**Purpose**: Identify suspicious constructs
**Checks for**:
- Unreachable code
- Printf format errors
- Suspicious assignments
- Unused results
- Shadow variables

---

### Step 4: Activate go-test-verifier Subagent

The `go-test-verifier` subagent provides:
- Coverage gap analysis
- Test quality assessment
- Specific recommendations
- Risk-based prioritization

See [go-test-verifier.md](../ agents/go-test-verifier.md) for details.

---

## Example Output

### Success Case: All Tests Pass

```
User: /test-and-verify

Claude:
🧪 Running comprehensive test suite...

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 TEST EXECUTION
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Running: go test -v -race -cover ./...

✓ internal/models         : 3 tests   | 100.0% coverage
✓ internal/repository     : 18 tests  | 95.5% coverage
✓ internal/service        : 24 tests  | 89.2% coverage
✓ internal/handlers       : 32 tests  | 87.8% coverage

Total: 77 tests passed
Time: 2.341s
Race conditions: None detected

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📈 COVERAGE ANALYSIS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Overall Coverage: 88.7% ✓ (target: 80%)

Package Breakdown:
  internal/models      : 100.0% ✓✓
  internal/repository  :  95.5% ✓✓
  internal/service     :  89.2% ✓
  internal/handlers    :  87.8% ✓

Critical Functions:
  CreateItem           :  95.0% ✓
  UpdateItem           :  92.5% ✓
  DeleteItem           : 100.0% ✓
  ListItems            :  90.0% ✓

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🔬 STATIC ANALYSIS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Running: go vet ./...
✓ No issues found

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ QUALITY ASSESSMENT
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Test Quality Metrics:
  Table-Driven Tests : 62/77 (80.5%) ✓
  Test Naming       : Conventional ✓
  Race-Free         : Yes ✓
  No Skipped Tests  : Yes ✓

Status: Ready for commit ✓

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
💡 RECOMMENDATIONS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Optional improvements:
  • Consider adding edge case tests for UpdateItem (price boundaries)
  • Add integration test for full CRUD workflow
  • Document complex test scenarios in comments

All quality gates passed. Ready for /commit-push-pr
```

---

### Failure Case: Tests Failing

```
User: /test-and-verify

Claude:
🧪 Running comprehensive test suite...

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 TEST EXECUTION
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Running: go test -v -race -cover ./...

✓ internal/models         : 3 tests   | 100.0% coverage
✓ internal/repository     : 18 tests  | 95.5% coverage
❌ internal/service        : 3 failed, 21 passed
✓ internal/handlers       : 32 tests  | 87.8% coverage

Total: 74 passed, 3 failed
Time: 1.892s

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
❌ TEST FAILURES
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

FAIL: TestItemService_CreateItem/negative_price
  Location: internal/service/item_service_test.go:67
  Error: expected error for negative price, got nil

FAIL: TestItemService_UpdateItem/empty_title
  Location: internal/service/item_service_test.go:145
  Error: validation should reject empty title after update

FAIL: TestItemService_DeleteItem
  Location: internal/service/item_service_test.go:203
  Error: panic: runtime error: invalid memory address

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📈 COVERAGE ANALYSIS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Overall Coverage: 76.2% ⚠️ (below 80% target)

Missing Coverage:
  service.go:145     Error handler not tested
  service.go:167     Edge case: nil item validation
  handlers.go:89     Panic recovery path not tested

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🔬 STATIC ANALYSIS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Running: go vet ./...
⚠️  1 issue found

internal/service/item_service.go:92:
  suspicious use of nil in comparison

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
❌ QUALITY ASSESSMENT
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Issues:
  ❌ Tests failing (3 failures)
  ❌ Coverage below threshold (76.2% < 80%)
  ⚠️  Vet warnings present

Status: Cannot commit until issues resolved

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🔧 RECOMMENDED ACTIONS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

High Priority:
  1. Fix failing tests in item_service_test.go
     - Add validation for negative prices (line 67)
     - Add validation for empty titles (line 145)
     - Fix nil pointer in DeleteItem (line 203)

  2. Fix vet warning in item_service.go:92
     - Use proper nil check pattern

  3. Add missing coverage:
     - Test error handling paths in service.go:145
     - Test nil item validation in service.go:167
     - Test panic recovery in handlers.go:89

Next Steps:
  1. Run failing tests individually for details:
     go test -v ./internal/service -run TestItemService_CreateItem/negative_price

  2. Fix issues and rerun:
     /test-and-verify

  3. When all pass, commit:
     /commit-push-pr
```

---

### Coverage Gap Analysis

When coverage is below threshold, detailed gap analysis:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🔍 COVERAGE GAP ANALYSIS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Uncovered Code Paths:

1. Error Handling (HIGH RISK)
   File: internal/service/item_service.go
   Lines: 145-152
   Issue: Error path when repository fails not tested

   Suggested Test:
   ```go
   {
       name: "repository error on create",
       item: validItem,
       setup: func(m *mockRepository) {
           m.CreateFunc = func(ctx context.Context, item *models.Item) error {
               return errors.New("db error")
           }
       },
       wantErr: true,
   },
   ```

2. Validation Edge Case (MEDIUM RISK)
   File: internal/service/item_service.go
   Lines: 167-169
   Issue: Nil item validation not tested

   Suggested Test:
   ```go
   {
       name: "nil item input",
       item: nil,
       wantErr: true,
       errType: service.ErrInvalidInput,
   },
   ```

3. Recovery Path (LOW RISK)
   File: internal/handlers/middleware.go
   Lines: 89-95
   Issue: Panic recovery middleware not tested
   Note: Tested implicitly in integration tests

Coverage Improvement Plan:
  1. Add tests for items 1 & 2 above (+4.2% coverage)
  2. Add boundary tests for price validation (+1.8% coverage)
  3. Add concurrent access tests (+2.1% coverage)
  → Total projected coverage: 84.3% ✓
```

---

## Tool Coordination

### Tools Used
- **Bash**: Execute go test, go vet, coverage commands
- **Read**: Parse test output and coverage files
- **Grep**: Extract specific test failures and coverage gaps
- **Sequential MCP**: Complex analysis for recommendations

### Integration with Subagents
**go-test-verifier**: Automatically activated for:
- Detailed coverage analysis
- Test quality assessment
- Specific recommendations
- Risk-based prioritization

---

## Quality Metrics Tracked

### Test Execution
- ✅ Total tests passed/failed
- ✅ Execution time
- ✅ Race condition detection
- ✅ Package-level results

### Coverage
- ✅ Overall coverage percentage
- ✅ Per-package breakdown
- ✅ Per-function coverage
- ✅ Critical function coverage

### Test Quality
- ✅ Table-driven test percentage
- ✅ Test naming conventions
- ✅ Skipped tests
- ✅ Test isolation

### Static Analysis
- ✅ Vet warnings
- ✅ Suspicious constructs
- ✅ Common mistakes

---

## Configuration

From `.claude/settings.json`:
```json
{
  "testing": {
    "minCoverage": 80,
    "targetCoverage": 90,
    "testCommand": "go test -v -race -cover ./...",
    "coverageCommand": "go test -coverprofile=coverage.out ./..."
  }
}
```

---

## Boundaries

### Will Do ✓
- Execute comprehensive test suite
- Generate coverage reports
- Run static analysis
- Provide detailed recommendations
- Activate subagent for deep analysis
- **Block commits on critical issues**

### Will Not Do ✗
- Generate test code (use separate command)
- Modify existing tests
- Skip checks without explicit flag
- Lower quality standards for convenience
- Fix tests automatically

---

## Quick Reference

### Common Commands
```bash
# Standard execution
/test-and-verify

# Verbose output
/test-and-verify --verbose

# Generate HTML coverage
/test-and-verify --coverage-html

# Quick tests only (skip slow)
/test-and-verify --quick

# Test specific package
go test -v ./internal/service

# Test specific function
go test -v ./internal/service -run TestItemService_CreateItem

# Coverage for specific package
go test -cover ./internal/service
```

### Coverage Thresholds
- **80%**: Minimum (enforced)
- **90%**: Target
- **95%+**: Critical business logic

---

## Success Criteria

- ✅ All tests pass (no failures)
- ✅ No race conditions detected
- ✅ Coverage ≥80% (enforced)
- ✅ No vet warnings
- ✅ Test quality metrics met
- ✅ Ready for commit
