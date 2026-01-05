---
name: api-test
description: "Test all API endpoints with automated validation"
category: testing
complexity: enhanced
mcp-servers: []
personas: []
---

# /api-test - API Endpoint Testing

## Purpose
Test all CRUD endpoints with sample data and validate responses. Provides quick verification that the API is functioning correctly.

## Triggers
- After implementing/modifying API endpoints
- Integration testing needs
- API validation before deployment
- Health check and smoke testing

## Usage
```bash
/api-test [--url http://localhost:8080] [--verbose]
```

### Flags
- `--url`: API base URL (default: http://localhost:8080)
- `--verbose`: Show full request/response details

---

## Execution Steps

### Step 1: Health Check
```bash
curl -s http://localhost:8080/health
```

**Purpose**: Verify API server is running
**Expected**: 200 OK response

---

### Step 2: CRUD Operations Test

#### 2.1 Create Item (POST /items)
```bash
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Laptop",
    "description": "MacBook Pro 16-inch",
    "price": 1299.99,
    "image": "https://example.com/laptop.jpg"
  }'
```

**Expected**:
- Status: 201 Created
- Response: Item with generated ID and timestamps

#### 2.2 List Items (GET /items)
```bash
curl -s http://localhost:8080/items
```

**Expected**:
- Status: 200 OK
- Response: Array of items including created item

#### 2.3 Get Specific Item (GET /items/{id})
```bash
curl -s http://localhost:8080/items/{id}
```

**Expected**:
- Status: 200 OK
- Response: Item details

#### 2.4 Update Item (PUT /items/{id})
```bash
curl -X PUT http://localhost:8080/items/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Laptop - Updated",
    "price": 1199.99,
    "status": "booked"
  }'
```

**Expected**:
- Status: 200 OK
- Response: Updated item

#### 2.5 Delete Item (DELETE /items/{id})
```bash
curl -X DELETE http://localhost:8080/items/{id}
```

**Expected**:
- Status: 200 OK or 204 No Content

---

### Step 3: Error Cases Test

#### 3.1 Invalid Input (POST with missing required fields)
```bash
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{"description": "No title"}'
```

**Expected**: 400 Bad Request with error message

#### 3.2 Not Found (GET non-existent ID)
```bash
curl -s http://localhost:8080/items/nonexistent-id
```

**Expected**: 404 Not Found

#### 3.3 Invalid Status (POST with invalid status)
```bash
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test",
    "price": 10.99,
    "status": "invalid-status"
  }'
```

**Expected**: 400 Bad Request

---

## Example Output

### Success - All Tests Pass

```
User: /api-test

Claude:
🧪 Running API endpoint tests...

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🏥 HEALTH CHECK
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

GET /health
✓ Status: 200 OK
✓ Server is running

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ CRUD OPERATIONS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

1. CREATE ITEM
   POST /items
   Request: {
     "title": "Test Laptop",
     "description": "MacBook Pro 16-inch",
     "price": 1299.99,
     "image": "https://example.com/laptop.jpg"
   }
   ✓ Status: 201 Created
   ✓ Response includes ID: 3f2a8b9c-1d2e-4f5a-8b7c-9d0e1f2a3b4c
   ✓ Timestamps set: created_at, updated_at
   ✓ Default status: available

2. LIST ITEMS
   GET /items
   ✓ Status: 200 OK
   ✓ Returns array
   ✓ Contains created item
   ✓ Count: 1 item

3. GET ITEM
   GET /items/3f2a8b9c-1d2e-4f5a-8b7c-9d0e1f2a3b4c
   ✓ Status: 200 OK
   ✓ Item details match
   ✓ All fields present

4. UPDATE ITEM
   PUT /items/3f2a8b9c-1d2e-4f5a-8b7c-9d0e1f2a3b4c
   Request: {
     "title": "Test Laptop - Updated",
     "price": 1199.99,
     "status": "booked"
   }
   ✓ Status: 200 OK
   ✓ Title updated correctly
   ✓ Price updated correctly
   ✓ Status updated to booked
   ✓ updated_at timestamp changed

5. DELETE ITEM
   DELETE /items/3f2a8b9c-1d2e-4f5a-8b7c-9d0e1f2a3b4c
   ✓ Status: 200 OK
   ✓ Item removed

6. VERIFY DELETE
   GET /items/3f2a8b9c-1d2e-4f5a-8b7c-9d0e1f2a3b4c
   ✓ Status: 404 Not Found

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🛡️ ERROR HANDLING
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

1. INVALID INPUT (missing required field)
   POST /items {"description": "No title"}
   ✓ Status: 400 Bad Request
   ✓ Error message: "title is required"

2. NOT FOUND
   GET /items/nonexistent-id
   ✓ Status: 404 Not Found
   ✓ Error message: "item not found"

3. INVALID STATUS
   POST /items {..., "status": "invalid"}
   ✓ Status: 400 Bad Request
   ✓ Error message: "invalid status value"

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 TEST SUMMARY
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Total Tests: 18
Passed: 18 ✓
Failed: 0
Success Rate: 100%

Time: 1.234s

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ ALL TESTS PASSED

API is functioning correctly!
```

---

### Failure - Server Not Running

```
User: /api-test

Claude:
🧪 Running API endpoint tests...

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🏥 HEALTH CHECK
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

GET /health
❌ Failed to connect: connection refused

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🚫 API SERVER NOT RUNNING
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Cannot run tests: API server not responding

To start the server:
  1. Run: make run
     OR
  2. Run: go run cmd/api/main.go

Once server is running, try again:
  /api-test
```

---

### Failure - Test Failures

```
User: /api-test

Claude:
🧪 Running API endpoint tests...

✓ Health check passed
✓ Create item passed
✓ List items passed
✓ Get item passed
❌ Update item FAILED
❌ Delete item FAILED

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
❌ TEST FAILURES
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

4. UPDATE ITEM
   PUT /items/3f2a8b9c-1d2e-4f5a-8b7c-9d0e1f2a3b4c
   Expected: 200 OK
   Got: 500 Internal Server Error
   Error: "failed to update item: validation error"

5. DELETE ITEM
   DELETE /items/3f2a8b9c-1d2e-4f5a-8b7c-9d0e1f2a3b4c
   Expected: 200 OK
   Got: 404 Not Found
   Error: "item not found"

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 TEST SUMMARY
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Total Tests: 18
Passed: 13
Failed: 2
Success Rate: 72.2%

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🔧 DEBUGGING SUGGESTIONS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

1. Check server logs for error details
2. Verify UpdateItem implementation in handlers
3. Check validation logic in service layer
4. Run unit tests: /test-and-verify

Issues found in:
  • internal/handlers/handlers.go (UpdateItem)
  • internal/service/item_service.go (validation)
```

---

## Tool Coordination

### Tools Used
- **Bash**: Execute curl commands
- **Read**: Parse JSON responses
- **Grep**: Extract error messages

### Integration with Subagents
**api-tester**: Can be activated for:
- More comprehensive testing
- Performance testing
- Load testing scenarios

---

## Test Scenarios Covered

### Happy Path
✅ Create new item
✅ List all items
✅ Get specific item
✅ Update item
✅ Delete item
✅ Verify deletion

### Error Cases
✅ Missing required fields
✅ Invalid input data
✅ Not found errors
✅ Invalid status values

### Edge Cases
- Empty list
- Concurrent operations (manual)
- Special characters in input
- Large data payloads

---

## Boundaries

### Will Do ✓
- Test all CRUD endpoints
- Validate HTTP status codes
- Check JSON response format
- Test error handling
- Provide clear test reports

### Will Not Do ✗
- Load or performance testing (use separate tool)
- Start/stop the API server
- Modify database state beyond tests
- Test authentication (future feature)

---

## Quick Reference

```bash
# Standard execution
/api-test

# Custom URL
/api-test --url http://localhost:3000

# Verbose output
/api-test --verbose

# Manual testing
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{"title":"Test","price":10.99}'
```

---

## Success Criteria

- ✅ All CRUD operations work
- ✅ Proper status codes returned
- ✅ Error handling works correctly
- ✅ JSON responses valid
- ✅ 100% test pass rate
