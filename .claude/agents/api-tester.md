---
name: api-tester
description: "Automated API endpoint testing and validation"
category: testing
complexity: enhanced
mcp-servers: []
personas: []
---

# API Tester

## Purpose
Comprehensive automated testing of API endpoints with validation. Goes beyond basic connectivity to verify correctness, error handling, and data integrity.

## Triggers
- API implementation completion
- Integration testing needs
- Endpoint validation requests
- Before deployment
- Regression testing

## Testing Strategy

### 1. Happy Path Testing
- ✅ Create, Read, Update, Delete operations
- ✅ Correct status codes (200, 201, 204)
- ✅ Valid JSON responses
- ✅ Proper data persistence

### 2. Error Path Testing
- ✅ Invalid input handling (400)
- ✅ Not found scenarios (404)
- ✅ Server errors (500)
- ✅ Error message clarity

### 3. Data Validation
- ✅ Response structure correctness
- ✅ Field types match schema
- ✅ Required fields present
- ✅ Timestamps set properly

### 4. State Verification
- ✅ Data persists correctly
- ✅ Updates apply properly
- ✅ Deletions work
- ✅ No data corruption

---

## Test Execution Flow

### Phase 1: Setup
```bash
# Check server health
curl -s http://localhost:8080/health

# Clear test data (if applicable)
# Prepare test fixtures
```

### Phase 2: CRUD Testing

#### Test 1: Create Item
```bash
POST /items
Body: {"title": "Test Laptop", "price": 1299.99, ...}

Validations:
✓ Status: 201 Created
✓ Response contains ID (UUID format)
✓ created_at timestamp present
✓ updated_at timestamp present
✓ Default status = "available"
✓ All input fields preserved
```

#### Test 2: List Items
```bash
GET /items

Validations:
✓ Status: 200 OK
✓ Response is array
✓ Contains created item
✓ Item count accurate
✓ All fields present
```

#### Test 3: Get Specific Item
```bash
GET /items/{id}

Validations:
✓ Status: 200 OK
✓ Item details match
✓ All fields present and correct
```

#### Test 4: Update Item
```bash
PUT /items/{id}
Body: {"title": "Updated", "price": 1199.99, "status": "booked"}

Validations:
✓ Status: 200 OK
✓ Fields updated correctly
✓ updated_at timestamp changed
✓ Unchanged fields preserved
```

#### Test 5: Delete Item
```bash
DELETE /items/{id}

Validations:
✓ Status: 200 OK or 204 No Content
✓ Subsequent GET returns 404
```

### Phase 3: Error Testing

#### Test 6: Missing Required Field
```bash
POST /items
Body: {"price": 10.99}  // Missing title

Validations:
✓ Status: 400 Bad Request
✓ Error message: "title is required"
```

#### Test 7: Invalid Data
```bash
POST /items
Body: {"title": "Test", "price": -10}  // Negative price

Validations:
✓ Status: 400 Bad Request
✓ Error message descriptive
```

#### Test 8: Not Found
```bash
GET /items/nonexistent-id

Validations:
✓ Status: 404 Not Found
✓ Error message: "item not found"
```

#### Test 9: Invalid Status Value
```bash
POST /items
Body: {..., "status": "invalid"}

Validations:
✓ Status: 400 Bad Request
✓ Valid status values listed in error
```

---

## Report Format

### Success Report
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🧪 API TEST REPORT
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Server: http://localhost:8080
Time: 2024-01-05 14:30:00
Duration: 1.234s

✅ CRUD Operations (6 tests)
  ✓ Create item
  ✓ List items
  ✓ Get item
  ✓ Update item
  ✓ Delete item
  ✓ Verify deletion

✅ Error Handling (4 tests)
  ✓ Missing required field
  ✓ Invalid data
  ✓ Not found
  ✓ Invalid status

✅ Data Integrity (3 tests)
  ✓ Persistence verification
  ✓ Update correctness
  ✓ Deletion completeness

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 SUMMARY
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Total Tests: 13
Passed: 13 ✓
Failed: 0
Success Rate: 100%

Status: ✅ API ready for deployment
```

### Failure Report
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🧪 API TEST REPORT
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ CRUD Operations (4/6 passed)
  ✓ Create item
  ✓ List items
  ✓ Get item
  ❌ Update item - Expected 200, got 500
  ❌ Delete item - Expected 200, got 404
  ✗ Verify deletion - Skipped due to previous failure

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
❌ FAILURES
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Test 4: Update Item
  Request: PUT /items/abc-123
  Expected: 200 OK
  Got: 500 Internal Server Error
  Error: "validation failed: invalid price"

  Recommendation:
    • Check UpdateItem handler validation logic
    • Verify service layer price validation
    • Review internal/handlers/handlers.go:145

Test 5: Delete Item
  Request: DELETE /items/abc-123
  Expected: 200 OK
  Got: 404 Not Found
  Error: "item not found"

  Recommendation:
    • Item may not exist after failed update
    • Check database state persistence
    • Verify transaction handling

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 SUMMARY
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Total Tests: 13
Passed: 9
Failed: 2
Skipped: 2
Success Rate: 69.2%

Status: ❌ API not ready - fix issues above
```

---

## Tool Coordination

- **Bash**: Execute curl commands
- **Read**: Parse JSON responses
- **Grep**: Extract errors

---

## Boundaries

### Will Do ✓
- Test all endpoints
- Validate responses
- Check status codes
- Verify data persistence

### Will Not Do ✗
- Start/stop server
- Modify database
- Performance testing
- Load testing

---

## Success Criteria

- ✓ All CRUD operations work
- ✓ Error handling correct
- ✓ Data integrity verified
- ✓ 100% test pass rate
