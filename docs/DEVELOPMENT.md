# Development Guide

This guide provides detailed information for developers working on the CaribEX backend.

## Table of Contents

- [Getting Started](#getting-started)
- [Project Architecture](#project-architecture)
- [Development Workflow](#development-workflow)
- [Common Tasks](#common-tasks)
- [Debugging](#debugging)
- [Performance Optimization](#performance-optimization)
- [Security Guidelines](#security-guidelines)

## Getting Started

### Initial Setup

1. **Install Go 1.23 or higher**
   ```bash
   # macOS
   brew install go@1.23
   
   # Ubuntu/Debian
   wget https://go.dev/dl/go1.23.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.23.linux-amd64.tar.gz
   ```

2. **Install development tools**
   ```bash
   # golangci-lint
   curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
   
   # golang-migrate
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   
   # air (hot reload)
   go install github.com/cosmtrek/air@latest
   ```

3. **Clone and setup**
   ```bash
   git clone https://github.com/Tenoywil/CaribEx-backend.git
   cd CaribEx-backend
   cp .env.example .env
   make deps
   ```

4. **Start services**
   ```bash
   make docker-up
   make migrate-up
   make run-dev
   ```

### IDE Setup

#### VS Code

Recommended extensions:
- Go (golang.go)
- Docker (ms-azuretools.vscode-docker)
- PostgreSQL (ckolkman.vscode-postgres)

Settings (`.vscode/settings.json`):
```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "package",
  "editor.formatOnSave": true,
  "[go]": {
    "editor.defaultFormatter": "golang.go"
  }
}
```

#### GoLand / IntelliJ IDEA

- Enable Go modules integration
- Set Go version to 1.23+
- Configure golangci-lint integration

## Project Architecture

### Domain-Driven Design

CaribEX follows DDD principles:

1. **Domain Layer** - Core business logic
   - Entities (User, Product, Order)
   - Value Objects (Money, Address)
   - Aggregates (Wallet, Cart)
   - Domain Services

2. **Application Layer** - Use cases
   - Coordinates domain operations
   - Manages transactions
   - Implements workflows

3. **Infrastructure Layer** - Technical details
   - Database access
   - Caching
   - External APIs
   - Logging

4. **Interface Layer** - API
   - HTTP handlers
   - Request/response DTOs
   - Route registration

### Package Organization

```
internal/domain/user/
├── user.go          # Entity definition
├── repository.go    # Repository interface
└── errors.go        # Domain errors

internal/repository/postgres/
└── user_repository.go  # Repository implementation

internal/usecase/
└── user_usecase.go     # Application logic

internal/controller/
└── user_controller.go  # HTTP handlers
```

## Development Workflow

### Feature Development

1. **Create feature branch**
   ```bash
   git checkout -b feature/user-authentication
   ```

2. **Write domain logic first**
   - Define entities and value objects
   - Write repository interfaces
   - Implement business rules

3. **Implement infrastructure**
   - Create repository implementations
   - Add database queries
   - Implement caching

4. **Create use cases**
   - Coordinate domain operations
   - Handle transactions
   - Add validation

5. **Add HTTP layer**
   - Create controllers
   - Define routes
   - Add middleware

6. **Write tests**
   - Unit tests for domain logic
   - Integration tests for repositories
   - Controller tests for HTTP

### Testing Strategy

#### Unit Tests

Test domain logic in isolation:

```go
// internal/domain/wallet/wallet_test.go
func TestWallet_Credit(t *testing.T) {
    wallet := &Wallet{Balance: 100.0}
    err := wallet.Credit(50.0)
    
    assert.NoError(t, err)
    assert.Equal(t, 150.0, wallet.Balance)
}
```

#### Integration Tests

Test with real database:

```go
// tests/integration/wallet_repository_test.go
func TestWalletRepository_Create(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer db.Close()
    
    repo := postgres.NewWalletRepository(db)
    wallet := &wallet.Wallet{
        UserID:   "user-123",
        Balance:  100.0,
        Currency: "JAM",
    }
    
    err := repo.Create(wallet)
    assert.NoError(t, err)
    assert.NotEmpty(t, wallet.ID)
}
```

#### E2E Tests

Test complete workflows:

```go
// tests/e2e/checkout_test.go
func TestCheckoutFlow(t *testing.T) {
    // 1. Authenticate user
    // 2. Add product to cart
    // 3. Checkout
    // 4. Verify order created
    // 5. Verify wallet debited
}
```

## Common Tasks

### Adding a New Endpoint

1. **Define domain model** (`internal/domain/product/product.go`)
2. **Create repository interface** (`internal/domain/product/repository.go`)
3. **Implement repository** (`internal/repository/postgres/product_repository.go`)
4. **Create use case** (`internal/usecase/product_usecase.go`)
5. **Add controller** (`internal/controller/product_controller.go`)
6. **Register routes** (`internal/routes/routes.go`)
7. **Write tests**
8. **Update API documentation** (`docs/API.md`)

### Adding a Database Migration

```bash
# Create migration files
migrate create -ext sql -dir migrations -seq add_user_email

# Edit migrations/NNNNNN_add_user_email.up.sql
ALTER TABLE users ADD COLUMN email VARCHAR(255);

# Edit migrations/NNNNNN_add_user_email.down.sql
ALTER TABLE users DROP COLUMN email;

# Apply migration
make migrate-up
```

### Adding Middleware

```go
// pkg/middleware/auth.go
func RequireAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Check authentication
        if !isAuthenticated(r) {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

Register middleware:
```go
// internal/routes/routes.go
router.Use(middleware.RequireAuth)
```

### Adding Cache

```go
// internal/usecase/product_usecase.go
func (u *ProductUseCase) GetByID(id string) (*product.Product, error) {
    // Check L1 cache
    if cached, ok := u.l1Cache.Get(fmt.Sprintf("product:%s", id)); ok {
        return cached.(*product.Product), nil
    }
    
    // Check L2 cache (Redis)
    if cached, err := u.l2Cache.Get(fmt.Sprintf("product:%s", id)); err == nil {
        return cached.(*product.Product), nil
    }
    
    // Fetch from database
    p, err := u.repo.GetByID(id)
    if err != nil {
        return nil, err
    }
    
    // Populate caches
    u.l1Cache.Set(fmt.Sprintf("product:%s", id), p)
    u.l2Cache.Set(fmt.Sprintf("product:%s", id), p)
    
    return p, nil
}
```

## Debugging

### Logging

Add debug logs:

```go
logger.Debug().
    Str("user_id", userID).
    Str("request_id", requestID).
    Msg("processing request")
```

Enable debug logs:
```bash
export LOG_LEVEL=debug
make run-dev
```

### Database Queries

Log SQL queries:

```bash
export DB_LOG_LEVEL=debug
```

Or use pgx logging:
```go
config.ConnConfig.Logger = logger
config.ConnConfig.LogLevel = pgx.LogLevelDebug
```

### Profiling

CPU profiling:
```bash
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof
```

Memory profiling:
```bash
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

### Debugging with Delve

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug the application
dlv debug ./cmd/api-server

# Set breakpoint
(dlv) break main.main
(dlv) continue
```

## Performance Optimization

### Database Optimization

1. **Use indexes**
   ```sql
   CREATE INDEX idx_users_wallet_address ON users(wallet_address);
   ```

2. **Use prepared statements**
   ```go
   stmt, err := db.Prepare("SELECT * FROM users WHERE id = $1")
   ```

3. **Batch operations**
   ```go
   batch := &pgx.Batch{}
   for _, user := range users {
       batch.Queue("INSERT INTO users (...) VALUES ($1, $2)", user.ID, user.Name)
   }
   results := db.SendBatch(batch)
   ```

### Caching Strategy

1. **Cache hot data** (product lists, categories)
2. **Use short TTLs** for frequently changing data
3. **Implement cache warming** for critical data
4. **Use singleflight** to prevent stampedes

### Concurrency

1. **Use worker pools** for parallel processing
2. **Set timeouts** on all operations
3. **Limit goroutines** to prevent resource exhaustion

```go
// Worker pool example
workers := 10
jobs := make(chan Job, 100)
results := make(chan Result, 100)

for w := 0; w < workers; w++ {
    go worker(jobs, results)
}
```

## Security Guidelines

### Input Validation

Always validate user input:

```go
type CreateUserInput struct {
    Username      string `json:"username" validate:"required,min=3,max=50"`
    WalletAddress string `json:"wallet_address" validate:"required,eth_addr"`
}

func (c *UserController) Create(w http.ResponseWriter, r *http.Request) {
    var input CreateUserInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }
    
    if err := validate.Struct(input); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Process valid input
}
```

### SQL Injection Prevention

Always use parameterized queries:

```go
// Good
db.Query("SELECT * FROM users WHERE id = $1", userID)

// Bad (vulnerable to SQL injection)
db.Query("SELECT * FROM users WHERE id = '" + userID + "'")
```

### Authentication & Authorization

1. **Always check authentication** before processing requests
2. **Verify resource ownership** for user-specific resources
3. **Use RBAC** for feature access control

```go
func (c *ProductController) Update(w http.ResponseWriter, r *http.Request) {
    userID := getUserFromContext(r)
    productID := chi.URLParam(r, "id")
    
    product, err := c.useCase.GetByID(productID)
    if err != nil {
        http.Error(w, "Not found", http.StatusNotFound)
        return
    }
    
    // Verify ownership
    if product.SellerID != userID {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }
    
    // Process update
}
```

### Secrets Management

Never hardcode secrets:

```go
// Bad
dbPassword := "mysecretpassword"

// Good
dbPassword := os.Getenv("DB_PASSWORD")
if dbPassword == "" {
    log.Fatal("DB_PASSWORD environment variable is required")
}
```

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Domain-Driven Design](https://domainlanguage.com/ddd/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Redis Documentation](https://redis.io/documentation)

## Getting Help

- Check [existing issues](https://github.com/Tenoywil/CaribEx-backend/issues)
- Ask in [discussions](https://github.com/Tenoywil/CaribEx-backend/discussions)
- Contact maintainers

Happy coding! 🚀
