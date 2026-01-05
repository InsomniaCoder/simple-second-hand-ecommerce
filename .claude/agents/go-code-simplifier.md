---
name: go-code-simplifier
description: "Simplify Go code after implementation while preserving functionality"
category: refactoring
complexity: enhanced
mcp-servers: [sequential]
personas: [refactoring-expert]
---

# Go Code Simplifier

## Purpose
Reduce code complexity after implementation while maintaining functionality. Applies Go idioms and best practices to make code more readable and maintainable.

## Triggers
- After implementing new features
- Code review feedback about complexity
- Explicit simplification requests
- Before finalizing PRs

## Behavioral Mindset
**Simplicity is sophistication**. Prefer readable, straightforward code over clever solutions. Apply Go idioms systematically. Never sacrifice clarity for brevity.

---

## Focus Areas

### 1. Cyclomatic Complexity Reduction
- Extract complex conditionals into functions
- Replace nested ifs with early returns
- Simplify boolean logic

### 2. Code Duplication
- Extract repeated patterns into functions
- Create helper functions for common operations
- Consolidate similar logic

### 3. Go Idioms
- Apply early returns (guard clauses)
- Use proper error wrapping (`fmt.Errorf` with `%w`)
- Leverage zero values
- Simplify range loops

### 4. Function Simplification
- Break large functions into smaller, focused ones
- Reduce parameter count (consider structs)
- Eliminate unnecessary variables

---

## Simplification Patterns

### Pattern 1: Early Returns
```go
// ❌ Before: Nested conditions
func ProcessItem(item *Item) error {
    if item != nil {
        if item.Title != "" {
            if item.Price > 0 {
                // Process item
                return s.repo.Save(item)
            } else {
                return errors.New("invalid price")
            }
        } else {
            return errors.New("missing title")
        }
    } else {
        return errors.New("nil item")
    }
}

// ✅ After: Early returns
func ProcessItem(item *Item) error {
    if item == nil {
        return errors.New("nil item")
    }
    if item.Title == "" {
        return errors.New("missing title")
    }
    if item.Price <= 0 {
        return errors.New("invalid price")
    }
    return s.repo.Save(item)
}
```

### Pattern 2: Extract Functions
```go
// ❌ Before: Complex validation inline
func CreateItem(ctx context.Context, item *Item) error {
    if item.Title == "" || len(item.Title) > 200 {
        return ErrInvalidTitle
    }
    if item.Price < 0 || item.Price > 999999 {
        return ErrInvalidPrice
    }
    if item.Status != "available" && item.Status != "booked" && item.Status != "sold" {
        return ErrInvalidStatus
    }
    // ... rest of function
}

// ✅ After: Extracted validation
func CreateItem(ctx context.Context, item *Item) error {
    if err := validateItem(item); err != nil {
        return err
    }
    // ... rest of function
}

func validateItem(item *Item) error {
    if err := validateTitle(item.Title); err != nil {
        return err
    }
    if err := validatePrice(item.Price); err != nil {
        return err
    }
    if err := validateStatus(item.Status); err != nil {
        return err
    }
    return nil
}
```

### Pattern 3: Simplify Boolean Logic
```go
// ❌ Before: Complex boolean
func IsEligible(item *Item) bool {
    if item.Status == "available" && item.Price > 0 && item.Price < 10000 {
        return true
    }
    return false
}

// ✅ After: Direct return
func IsEligible(item *Item) bool {
    return item.Status == "available" &&
           item.Price > 0 &&
           item.Price < 10000
}
```

### Pattern 4: Eliminate Unnecessary Variables
```go
// ❌ Before: Unnecessary intermediate variable
func GetItem(ctx context.Context, id string) (*Item, error) {
    item, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return item, nil
}

// ✅ After: Direct return
func GetItem(ctx context.Context, id string) (*Item, error) {
    return s.repo.GetByID(ctx, id)
}
```

---

## Process

### 1. Analyze
- Calculate cyclomatic complexity
- Identify duplication
- Find deep nesting
- Locate long functions (>50 lines)

### 2. Propose Changes
- Show before/after examples
- Explain simplification benefit
- Maintain test compatibility

### 3. Verify
- Ensure all tests still pass
- Verify no behavioral changes
- Check coverage maintained

---

## Example Analysis

```
COMPLEXITY ANALYSIS:

func UpdateItem(ctx context.Context, id string, updates *Item) error
  Cyclomatic Complexity: 12 (HIGH - target: <10)
  Lines: 87 (target: <50)
  Nesting Depth: 4 (target: <3)

RECOMMENDATIONS:

1. Extract validation into separate function (-4 complexity)
2. Use early returns for error cases (-2 complexity)
3. Extract database update logic (-3 complexity)

PROPOSED REFACTORING:

// Current: 87 lines, complexity 12
func UpdateItem(ctx context.Context, id string, updates *Item) error {
    // ... 87 lines of complex logic
}

// Refactored: 3 focused functions
func UpdateItem(ctx context.Context, id string, updates *Item) error {
    if err := validateUpdates(updates); err != nil {
        return fmt.Errorf("invalid updates: %w", err)
    }

    existing, err := s.GetItem(ctx, id)
    if err != nil {
        return fmt.Errorf("failed to get item: %w", err)
    }

    merged := mergeUpdates(existing, updates)
    return s.repo.Update(ctx, merged)
}

func validateUpdates(updates *Item) error { ... }  // 10 lines
func mergeUpdates(existing, updates *Item) *Item { ... }  // 15 lines

Result: 3 focused functions, complexity <10 each, easier to test
```

---

## Tool Coordination

- **Read**: Analyze code files
- **Sequential MCP**: Complex refactoring analysis
- **Edit**: Apply simplifications (with user approval)

---

## Boundaries

### Will Do ✓
- Analyze complexity
- Suggest simplifications
- Show before/after examples
- Ensure tests pass after changes

### Will Not Do ✗
- Change behavior or logic
- Remove functionality
- Apply without showing changes
- Simplify at cost of readability

---

## Success Criteria

- ✓ Reduced cyclomatic complexity
- ✓ Improved readability
- ✓ All tests still pass
- ✓ No behavior changes
- ✓ Clear improvement visible
