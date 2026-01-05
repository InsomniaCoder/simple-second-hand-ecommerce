# Go CRUD API Reference Project - Claude Code Conventions

## Project Overview
Reference implementation demonstrating Boris's Claude Code best practices through a clean Go API with comprehensive automation. This project serves as a learning template for teams adopting Claude Code workflows.

**Use Case**: Personal e-commerce API for selling items before moving countries
**Architecture**: Clean architecture with interface-driven design
**Testing**: >80% coverage enforced by automation

---

## Go Code Standards

### Clean Code Principles
- **Simplicity First**: Prefer simple, readable code over clever solutions
- **Interface-Driven**: Use interfaces for testability and flexibility
- **Error Handling**: Never ignore errors, always provide context
- **Context Propagation**: Pass `context.Context` for cancellation and timeouts
- **Early Returns**: Reduce nesting with early returns and guard clauses

### Code Style Examples

#### ✅ Good: Descriptive names, early returns, clear error handling
```go
func (s *ItemService) CreateItem(ctx context.Context, item *models.Item) error {
    if err := validateItem(item); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    item.ID = uuid.New().String()
    item.CreatedAt = time.Now()
    item.UpdatedAt = time.Now()

    if err := s.repo.Create(ctx, item); err != nil {
        return fmt.Errorf("failed to create item: %w", err)
    }

    return nil
}
```

#### ❌ Bad: Nested ifs, unclear names, swallowed errors
```go
func (s *ItemService) CreateItem(ctx context.Context, i *models.Item) error {
    if i.Title != "" {
        i.ID = generateID()
        if err := s.repo.Create(ctx, i); err != nil {
            return err  // No context!
        }
        return nil
    }
    return errors.New("invalid")  // Vague error!
}
```

### Interface Patterns

#### Define interfaces in consumer packages
```go
// In service package, not repository package
type ItemRepository interface {
    Create(ctx context.Context, item *models.Item) error
    GetByID(ctx context.Context, id string) (*models.Item, error)
    List(ctx context.Context) ([]*models.Item, error)
    Update(ctx context.Context, item *models.Item) error
    Delete(ctx context.Context, id string) error
}
```

#### Keep interfaces small and focused (Interface Segregation)
```go
type ItemReader interface {
    GetByID(ctx context.Context, id string) (*models.Item, error)
    List(ctx context.Context) ([]*models.Item, error)
}

type ItemWriter interface {
    Create(ctx context.Context, item *models.Item) error
    Update(ctx context.Context, item *models.Item) error
    Delete(ctx context.Context, id string) error
}
```

### Error Handling Patterns

#### Wrap errors with context using fmt.Errorf and %w
```go
if err := s.repo.GetByID(ctx, id); err != nil {
    return fmt.Errorf("failed to get item %s: %w", id, err)
}
```

#### Define sentinel errors for known cases
```go
var (
    ErrNotFound     = errors.New("item not found")
    ErrInvalidInput = errors.New("invalid input")
    ErrAlreadyExists = errors.New("item already exists")
)
```

#### Check errors with errors.Is and errors.As
```go
if errors.Is(err, repository.ErrNotFound) {
    return nil, service.ErrItemNotFound
}

var validationErr *ValidationError
if errors.As(err, &validationErr) {
    // Handle validation error specifically
}
```

### Concurrency Patterns

#### Thread-safe repository with RWMutex
```go
type memoryRepository struct {
    mu    sync.RWMutex
    items map[string]*models.Item
}

func (r *memoryRepository) GetByID(ctx context.Context, id string) (*models.Item, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    item, ok := r.items[id]
    if !ok {
        return nil, ErrNotFound
    }
    return item, nil
}

func (r *memoryRepository) Create(ctx context.Context, item *models.Item) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, exists := r.items[item.ID]; exists {
        return ErrAlreadyExists
    }
    r.items[item.ID] = item
    return nil
}
```

---

## Testing Requirements

### Coverage Standards
- **Minimum**: 80% code coverage (enforced by pre-commit hook)
- **Target**: 90%+ for business logic (service layer)
- **All public functions** must have tests
- **Critical paths** (CRUD operations) require >95% coverage

### Test Structure - Table-Driven Tests

```go
func TestItemService_CreateItem(t *testing.T) {
    tests := []struct {
        name    string
        item    *models.Item
        setup   func(*mockRepository)
        wantErr bool
        errType error
    }{
        {
            name: "valid item",
            item: &models.Item{
                Title:       "Test Item",
                Description: "Test Description",
                Price:       10.99,
            },
            setup: func(m *mockRepository) {
                m.CreateFunc = func(ctx context.Context, item *models.Item) error {
                    return nil
                }
            },
            wantErr: false,
        },
        {
            name: "missing title",
            item: &models.Item{
                Description: "Test Description",
                Price:       10.99,
            },
            wantErr: true,
            errType: service.ErrInvalidInput,
        },
        {
            name: "negative price",
            item: &models.Item{
                Title: "Test Item",
                Price: -10.99,
            },
            wantErr: true,
            errType: service.ErrInvalidInput,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := &mockRepository{}
            if tt.setup != nil {
                tt.setup(mockRepo)
            }

            svc := service.NewItemService(mockRepo)
            err := svc.CreateItem(context.Background(), tt.item)

            if tt.wantErr {
                if err == nil {
                    t.Error("expected error, got nil")
                }
                if tt.errType != nil && !errors.Is(err, tt.errType) {
                    t.Errorf("expected error type %v, got %v", tt.errType, err)
                }
            } else {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
                // Verify side effects
                if tt.item.ID == "" {
                    t.Error("expected ID to be set")
                }
            }
        })
    }
}
```

### Test Organization
- Use `pkg/testutil/helpers.go` for common test utilities
- Create fixtures for test data (valid items, invalid items)
- Mock interfaces, not concrete types
- Test one behavior per test case
- Use descriptive test names: `TestFunction_Scenario_ExpectedResult`

---

## Commit Standards

### Conventional Commits Format
```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types
- **feat**: New feature (e.g., "feat(api): add item creation endpoint")
- **fix**: Bug fix (e.g., "fix(repository): handle concurrent access safely")
- **refactor**: Code refactoring without behavior change
- **test**: Adding or updating tests
- **docs**: Documentation changes
- **chore**: Maintenance tasks (dependencies, tooling)
- **perf**: Performance improvements

### Commit Message Examples

#### Feature Addition
```
feat(api): add item update endpoint

Implements PUT /items/{id} endpoint with validation
and partial update support. Uses ItemService for
business logic and proper error handling.

Closes #123
```

#### Bug Fix
```
fix(repository): handle concurrent access safely

Added RWMutex locks to prevent race conditions in
in-memory repository implementation. Read operations
use RLock, write operations use Lock.

Fixes #456
```

#### Test Addition
```
test(service): increase coverage for edge cases

Added tests for:
- Nil input validation
- Invalid status values
- Concurrent create operations
- Boundary conditions for price

Coverage: 87% → 92%
```

---

## Claude Code Workflows

### Daily Development Flow

1. **Start Session**
   - Verify tests pass: `/test-and-verify`
   - Check git status: `git status`

2. **Implement Feature**
   - Write code (files auto-format on save via post-tool-use hook)
   - Write tests alongside implementation

3. **Validate**
   - Run `/format-all` to ensure formatting
   - Run `/test-and-verify` to check coverage
   - Fix any issues identified

4. **Commit**
   - Use `/commit-push-pr` for automated workflow
   - Pre-commit hook runs validation automatically
   - Commit message generated from changes

### Code Review Checklist

Before requesting review, verify:
- [ ] All tests pass (`go test ./...`)
- [ ] Coverage meets threshold (80%+)
- [ ] Code formatted (`go fmt ./...`)
- [ ] No vet warnings (`go vet ./...`)
- [ ] Conventional commit message
- [ ] Documentation updated (if public API changed)
- [ ] No TODOs in production code

### Slash Commands Reference

**`/commit-push-pr`** - Automated git workflow
- Runs pre-commit validation (format, test, vet)
- Verifies not on main/master branch
- Generates conventional commit message
- Commits, pushes, and creates PR

**`/test-and-verify`** - Test execution with analysis
- Runs `go test -v -race -cover ./...`
- Generates coverage report
- Verifies minimum coverage (80%)
- Runs `go vet` for static analysis
- Activates `go-test-verifier` subagent for detailed report

**`/format-all`** - Code formatting
- Runs `go fmt ./...`
- Runs `goimports -w .` (if available)
- Reports formatted files

**`/api-test`** - API endpoint testing
- Health check for running server
- Tests all CRUD endpoints
- Validates responses
- Generates test report

### Subagents Reference

**`go-code-simplifier`** - Code simplification
- **When**: After feature implementation, code reviews
- **Actions**: Reduce complexity, extract patterns, apply Go idioms
- **Process**: Analyze → Propose changes → Verify tests pass

**`go-test-verifier`** - Test quality verification
- **When**: Pre-commit, coverage analysis
- **Actions**: Coverage analysis, quality checks, recommendations
- **Process**: Run tests → Parse coverage → Analyze quality → Detailed report
- **MCP**: Uses Sequential for complex analysis

**`api-tester`** - API endpoint testing
- **When**: API implementation, integration testing
- **Actions**: Test HTTP methods, validate responses, check JSON
- **Process**: Start server → Execute tests → Validate → Report

---

## Project Structure Rules

### Package Organization
```
cmd/            - Application entry points
  api/          - Main API server

internal/       - Private application code
  models/       - Domain models (entities)
  repository/   - Data access layer (interfaces + implementations)
  service/      - Business logic layer
  handlers/     - HTTP layer (routes, middleware)

pkg/            - Public libraries (reusable across projects)
  testutil/     - Test utilities and helpers

docs/           - Documentation
scripts/        - Utility scripts
```

### Import Organization
1. Standard library packages
2. External dependencies
3. Internal packages

```go
import (
    // Standard library
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    // External dependencies
    "github.com/google/uuid"

    // Internal packages
    "github.com/yourusername/go-crud-api-reference/internal/models"
    "github.com/yourusername/go-crud-api-reference/internal/repository"
)
```

### File Naming Conventions
- Implementation files: `item_service.go`, `memory.go`
- Test files: `item_service_test.go`, `memory_test.go`
- Interface definitions: `service.go`, `repository.go`
- All lowercase with underscores for multi-word names

---

## Automation Features

### Post-Tool-Use Hook
Automatically formats Go files after Edit/Write operations:
- Runs `go fmt` on modified `.go` files
- Runs `goimports` if available
- No manual formatting needed
- Ensures consistent style across team

### Pre-Commit Hook
Validates code before allowing commits:
- Ensures code is formatted (`go fmt`)
- Runs test suite (`go test`)
- Runs static analysis (`go vet`)
- **Blocks commit if any check fails**
- Prevents broken code from entering repository

---

## MCP Integration Examples

### Sequential Thinking for Complex Logic
When implementing complex features, Sequential MCP provides structured analysis:

```
User: "How should I handle item status transitions?"
Claude: *Activates Sequential MCP*
- Thought 1: Define valid status transitions (available → booked → sold)
- Thought 2: Identify business rules (can't unbuy, can unbook)
- Thought 3: Design validation function
- Thought 4: Consider edge cases
- Result: State machine approach with validation
```

### Context7 for Go Best Practices
For framework-specific questions, Context7 fetches official patterns:

```
User: "What's the best way to structure HTTP middleware?"
Claude: *Activates Context7 MCP*
- Fetches official Go HTTP middleware patterns
- Provides idiomatic examples with http.Handler
- Shows chaining approach
- Demonstrates error handling
```

### Serena for Project Memory
For session persistence and project context:

```
User: /sc:load
Claude: *Activates Serena MCP*
- Loads project context from memory
- Recalls previous implementation decisions
- Restores session state
- Continues from last checkpoint
```

---

## Performance Considerations

### Optimization Patterns
- Use `sync.Pool` for frequently allocated objects (if profiling shows need)
- Prefer value receivers for small structs (<64 bytes)
- Use buffered channels for producer-consumer patterns
- **Profile first**: Use `pprof` before optimizing

### Memory Management
```go
// Good: Value receiver for small struct
func (i Item) IsAvailable() bool {
    return i.Status == StatusAvailable
}

// Good: Pointer receiver for modification
func (r *memoryRepository) Create(ctx context.Context, item *models.Item) error {
    // ...
}
```

---

## Security Guidelines

### Input Validation
- Validate all user input at the service layer
- Sanitize data before storage
- Use appropriate types (avoid `interface{}` for user input)

### Logging Best Practices
```go
// ✅ Good: Log without sensitive data
log.Printf("Created item: id=%s, title=%s", item.ID, item.Title)

// ❌ Bad: Never log sensitive data
log.Printf("Created item: %+v", item)  // Don't log entire struct
```

### Error Messages
```go
// ✅ Good: Generic error to user, detailed error in logs
if err := authenticate(user); err != nil {
    log.Printf("authentication failed for user %s: %v", user.ID, err)
    return errors.New("authentication failed")  // Generic to client
}
```

---

## Common Pitfalls to Avoid

### 1. Ignoring Errors
```go
// ❌ Bad
_ = someFunction()

// ✅ Good
if err := someFunction(); err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

### 2. Not Using Context
```go
// ❌ Bad
func (s *Service) GetItem(id string) (*Item, error)

// ✅ Good
func (s *Service) GetItem(ctx context.Context, id string) (*Item, error)
```

### 3. Swallowing Error Context
```go
// ❌ Bad
if err != nil {
    return err  // Lost context
}

// ✅ Good
if err != nil {
    return fmt.Errorf("failed to process item: %w", err)
}
```

### 4. Not Testing Error Paths
```go
// Always test both success and error cases
tests := []struct{
    name string
    // ...
    wantErr bool
}{
    {"success case", ..., false},
    {"error case", ..., true},  // Don't forget this!
}
```

---

## Team Collaboration

### Shared Responsibilities
- **Update CLAUDE.md**: When discovering patterns or anti-patterns
- **Review PRs**: Use code review checklist
- **Share learnings**: Document mistakes in CLAUDE.md
- **Maintain automation**: Update slash commands and subagents

### When to Update CLAUDE.md
- Discovered a new Go idiom or pattern
- Found a common mistake that should be avoided
- Added a new testing pattern
- Updated commit message conventions
- Changed project structure

### PR Review Focus
1. Does it follow CLAUDE.md conventions?
2. Are tests comprehensive (>80% coverage)?
3. Is error handling proper (wrapped with context)?
4. Are interfaces used appropriately?
5. Is the commit message conventional?

---

## Extension Guidelines

### Adding New Endpoints
1. Add model changes (if needed) in `internal/models/`
2. Update repository interface in `internal/repository/`
3. Implement business logic in `internal/service/`
4. Add HTTP handlers in `internal/handlers/`
5. Update routes in `internal/handlers/routes.go`
6. Write tests for all layers
7. Update `docs/API.md`

### Adding New Entity Types
1. Define model in `internal/models/`
2. Create repository interface
3. Implement in-memory repository
4. Create service interface
5. Implement service with validation
6. Add handlers for CRUD operations
7. Write comprehensive tests (>80% coverage)

---

## Success Metrics

### Code Quality
- ✅ Test coverage >80% (enforced)
- ✅ All tests passing
- ✅ No `go vet` warnings
- ✅ Code formatted consistently
- ✅ No TODOs in production code

### Automation
- ✅ Post-tool-use hook auto-formats files
- ✅ Pre-commit hook prevents bad commits
- ✅ `/commit-push-pr` successfully validates
- ✅ `/test-and-verify` provides clear reports

### Development Velocity
- ✅ Fast feedback loops (automated testing)
- ✅ Reduced manual work (hooks, slash commands)
- ✅ Consistent code style (automation)
- ✅ Clear conventions (CLAUDE.md)

---

## Resources

### Official Go Documentation
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Standard Library](https://pkg.go.dev/std)

### Project Documentation
- [API Documentation](../docs/API.md)
- [Automation Guide](../docs/AUTOMATION.md)
- [README](../README.md)

### Claude Code Resources
- Boris's blog post on Claude Code best practices
- `.claude/commands/` for slash command examples
- `.claude/agents/` for subagent examples
