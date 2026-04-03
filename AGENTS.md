# AGENTS.md - Development Guidelines for URL Shorten API

## Project Overview
Go URL shortening API using Gin web framework, GORM ORM, and PostgreSQL.

## Build Commands

```bash
# Build the application
make build
go build -o main ./cmd/server

# Run the server
make run
./main

# Run with Docker
make docker-run           # Start services
make docker-run-build    # Build and start
make docker-stop         # Stop services
make docker-stop-v       # Stop with volumes
```

## Test Commands

No tests currently exist. To add tests:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for specific package
go test -v ./internal/handlers

# Run tests matching pattern
go test -v -run TestCreateURL ./internal/handlers

# Run with coverage
go test -v -cover ./...
```

## Lint & Code Quality

```bash
# Run go vet
go vet ./...

# Format code
go fmt ./...

# Tidy go.mod
go mod tidy
```

## Code Style Guidelines

### Import Organization
Group imports: standard library first, then external packages, separated by blank line.
```go
import (
    "context"
    "fmt"
    "net/http"
    "time"

    "github.com/sanjog-pariyar/url-shorten-api/internal/config"
    "github.com/sanjog-pariyar/url-shorten-api/internal/models"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)
```

### Naming Conventions
- **Packages**: lowercase, concise (e.g., `handlers`, `services`, `repository`)
- **Types**: PascalCase (e.g., `URLHandler`, `URLService`)
- **Functions/Variables**: camelCase (e.g., `createShortURL`, `shortCode`)
- **Constants**: PascalCase for exported, camelCase for unexported
- **Database fields**: snake_case via GORM tags (e.g., `short_code`, `user_id`)

### Error Handling
- Define typed errors as package-level variables in the service layer
- Use `errors.Is()` for error comparison in handlers
- Return errors early; avoid nested error handling
- Log fatal errors with `log.Fatalf()` in main; return errors from handlers/services
```go
var (
    ErrInvalidURL     = errors.New("invalid URL")
    ErrURLNotFound    = errors.New("URL not found")
    ErrURLExpired     = errors.New("URL has expired")
)
```

### Architecture
- **Layers**: Handlers -> Services -> Repository
- **Dependency Injection**: Pass interfaces to constructors
- **Interfaces**: Define in the layer that uses them (services use repository interfaces)
- **Gin Context**: Extract request data in handlers; pass only typed data to services

### HTTP Response Patterns
- Use `c.JSON()` for JSON responses with proper status codes
- Use `gin.H{}` for error response maps
- Validate request bodies with `c.ShouldBindJSON()`
- Return early on validation/auth errors

### Database Models
- Use GORM tags for column mapping (`gorm:"primaryKey"`, `gorm:"uniqueIndex"`)
- Use JSON tags for API serialization (`json:"shortCode"`)
- Pointer fields for nullable values (`*uint64`, `*time.Time`)
- Default values via GORM (`gorm:"default:0"`)

### Configuration
- Use environment variables via `godotenv` or `os.Getenv()`
- Provide sensible defaults in `config.Load()`
- Group related config into structs (e.g., `RateLimitConfig`)

### Middleware
- Return `gin.HandlerFunc` from middleware constructors
- Use `c.Next()` to pass control to next handler
- Use `c.AbortWithStatusJSON()` to reject requests early
- Set context values with `c.Set()`; retrieve with `c.Get()`

### Code Organization
```
cmd/server/        # Entry point
internal/
  config/          # Configuration loading
  database/       # DB connection
  handlers/       # HTTP handlers (Gin)
  middleware/     # HTTP middleware
  models/         # Data models and DTOs
  repository/     # Data access layer
  services/       # Business logic
```

### General Patterns
- Use struct receivers for methods (e.g., `func (h *URLHandler) CreateShortURL`)
- Prefer interfaces for dependencies (e.g., `URLService` interface)
- Close resources with `defer` where applicable
- No comments unless explaining non-obvious logic
