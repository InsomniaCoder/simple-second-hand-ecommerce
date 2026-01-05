# Go CRUD API Reference Project

A comprehensive reference implementation demonstrating **Boris's Claude Code best practices** through a production-quality Go API.

## Overview

This project serves as a learning template for teams adopting Claude Code workflows. It showcases:

- **Slash Commands**: Automated workflows for testing, formatting, and git operations
- **Subagents**: Specialized agents for code quality, test verification, and API testing
- **Hooks**: Automatic formatting and pre-commit validation
- **Clean Architecture**: Interface-driven design with proper separation of concerns
- **Comprehensive Testing**: >80% coverage with table-driven tests
- **Team Conventions**: Shared standards in `.claude/CLAUDE.md`

## Project Context

Built for a personal e-commerce use case (selling items before moving), this API manages items with:
- Title, description, price
- Image URL
- Status (available, booked, sold)

## Quick Start

### Prerequisites
- Go 1.21 or higher
- Git
- (Optional) `goimports` for import organization
- (Optional) `gh` CLI for PR creation

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/go-crud-api-reference.git
cd go-crud-api-reference

# Install dependencies
go mod download

# Run tests
make test

# Start the server
make run
```

### Development Workflow

```bash
# Format code
make fmt

# Run tests with coverage
make coverage

# Build binary
make build

# Run all checks (format, test, vet)
make check
```

## Claude Code Features

This project includes comprehensive Claude Code automation:

### Slash Commands
- `/commit-push-pr` - Automated git workflow with validation
- `/test-and-verify` - Run tests with coverage analysis
- `/format-all` - Format all Go files
- `/api-test` - Test all API endpoints

### Subagents
- `go-code-simplifier` - Simplify complex code
- `go-test-verifier` - Verify test coverage and quality
- `api-tester` - Automated endpoint testing

### Hooks
- **post-tool-use**: Auto-format Go files after edits
- **pre-commit**: Validate tests before commits

For detailed automation documentation, see [docs/AUTOMATION.md](docs/AUTOMATION.md).

## API Endpoints

```
POST   /items          Create a new item
GET    /items          List all items
GET    /items/{id}     Get a specific item
PUT    /items/{id}     Update an item
DELETE /items/{id}     Delete an item
```

For complete API documentation, see [docs/API.md](docs/API.md).

## Project Structure

```
├── .claude/              # Claude Code automation
│   ├── CLAUDE.md         # Project conventions
│   ├── settings.json     # Permissions & hooks
│   ├── agents/          # Specialized subagents
│   ├── commands/        # Slash commands
│   └── hooks/           # Automation hooks
├── cmd/api/             # Application entry point
├── internal/
│   ├── models/          # Domain models
│   ├── repository/      # Data access layer
│   ├── service/         # Business logic
│   └── handlers/        # HTTP handlers
├── pkg/testutil/        # Test utilities
├── docs/                # Documentation
└── scripts/             # Utility scripts
```

## Architecture

This project follows clean architecture principles:

1. **Models Layer**: Domain entities and business rules
2. **Repository Layer**: Data access abstraction (in-memory implementation)
3. **Service Layer**: Business logic and validation
4. **Handlers Layer**: HTTP transport and routing

All layers communicate through interfaces, enabling:
- Easy testing with mocks
- Flexibility to swap implementations
- Clear separation of concerns

## Testing

Tests are table-driven and comprehensive:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
./scripts/test-coverage.sh
```

**Coverage target**: 80% minimum (enforced by pre-commit hook)

## Learning Resources

- **[CLAUDE.md](.claude/CLAUDE.md)** - Project conventions and standards
- **[AUTOMATION.md](docs/AUTOMATION.md)** - Claude Code automation guide
- **[API.md](docs/API.md)** - API specification and examples

## Extension Ideas

This reference can be extended with:
- PostgreSQL database implementation
- JWT authentication
- Pagination and filtering
- Docker containerization
- CI/CD pipeline
- Integration tests with Playwright

## License

MIT License - Use this as a learning template for your projects

## Acknowledgments

Inspired by Boris's blog post on Claude Code best practices, demonstrating:
- Team-shared conventions (CLAUDE.md)
- Workflow automation (slash commands)
- Quality enforcement (hooks and subagents)
- Verification-focused development
