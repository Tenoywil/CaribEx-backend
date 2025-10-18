# CaribX Backend

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](Dockerfile)

> Blockchain Money Transfer & Marketplace for Jamaica & the Caribbean

A secure, high-performance backend enabling P2P transfers, marketplace listings, wallet management, and order processing — built with Domain-Driven Design (DDD), strong security, and two-level caching.

---

## 🚀 Features

- **🔐 Web3 Authentication**: Sign-In With Ethereum (SIWE) for wallet-based auth
- **💰 Wallet Management**: Multi-currency wallet with transaction ledger
- **🛍️ Marketplace**: Product listings, shopping cart, and order processing
- **⚡ High Performance**: Two-level caching (L1 in-memory + L2 Redis)
- **🔒 Security First**: Row-Level Security (RLS), PII redaction, rate limiting
- **📊 Observability**: Prometheus metrics, OpenTelemetry tracing, structured logging
- **🎯 Domain-Driven Design**: Clean architecture with clear separation of concerns
- **🐳 Docker Ready**: Full containerization with Docker Compose support

---

## 📋 Table of Contents

- [Quick Start](#-quick-start)
- [Project Structure](#-project-structure)
- [API Documentation](#-api-documentation)
- [Architecture](#-architecture)
- [Development](#-development)
- [Testing](#-testing)
- [Deployment](#-deployment)
- [Contributing](#-contributing)
- [License](#-license)

---

## 🏁 Quick Start

### Prerequisites

- Go 1.23 or higher
- Docker & Docker Compose
- PostgreSQL 15+ (or use Docker)
- Redis 7+ (or use Docker)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/Tenoywil/CaribEx-backend.git
   cd CaribEx-backend
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start services with Docker Compose**
   ```bash
   make docker-up
   ```

   This will start:
   - PostgreSQL (localhost:5432)
   - Redis (localhost:6379)
   - API Server (localhost:8080)

4. **Run migrations**
   ```bash
   make migrate-up
   ```

5. **Access the API**
   ```bash
   curl http://localhost:8080/healthz
   ```

### Alternative: Run Locally

```bash
# Install dependencies
make deps

# Build the application
make build

# Run the server
make run-dev
```

---

## 📁 Project Structure

```
CaribEx-backend/
├── cmd/
│   └── api-server/          # Application entry point
│       └── main.go
├── internal/
│   ├── domain/              # Domain models & business logic
│   │   ├── user/           # User aggregate
│   │   ├── wallet/         # Wallet aggregate
│   │   ├── product/        # Product aggregate
│   │   ├── cart/           # Cart aggregate
│   │   └── order/          # Order aggregate
│   ├── repository/         # Data access layer
│   │   └── postgres/       # PostgreSQL repositories
│   ├── usecase/            # Application business logic
│   ├── controller/         # HTTP request handlers
│   └── routes/             # Route definitions
├── pkg/
│   ├── cache/              # L1 & L2 cache implementations
│   ├── config/             # Configuration management
│   ├── logger/             # Structured logging
│   ├── middleware/         # HTTP middleware
│   └── monitoring/         # Metrics & tracing
├── migrations/             # Database migrations
├── tests/                  # Integration & E2E tests
├── docs/                   # Documentation
│   ├── API.md             # API documentation
│   └── ARCHITECTURE.md    # Architecture details
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── README.md
```

---

## 📚 API Documentation

### Base URL

```
http://localhost:8080/v1
```

### Core Endpoints

#### Authentication
- `GET /v1/auth/nonce` - Generate SIWE nonce
- `POST /v1/auth/siwe` - Authenticate with signature
- `GET /v1/auth/me` - Get current user

#### Wallet
- `GET /v1/wallet` - Get wallet balance
- `POST /v1/wallet/send` - Send funds
- `GET /v1/wallet/transactions` - Transaction history

#### Products
- `GET /v1/products` - List products
- `GET /v1/products/:id` - Get product details
- `POST /v1/products` - Create product (seller only)
- `PUT /v1/products/:id` - Update product
- `DELETE /v1/products/:id` - Delete product

#### Cart & Orders
- `GET /v1/cart` - Get current cart
- `POST /v1/cart/items` - Add item to cart
- `PUT /v1/cart/items/:id` - Update cart item
- `DELETE /v1/cart/items/:id` - Remove cart item
- `POST /v1/orders` - Checkout and create order
- `GET /v1/orders` - Get user orders

**Full API documentation**: [docs/API.md](docs/API.md)

---

## 🏗️ Architecture

CaribX follows **Domain-Driven Design (DDD)** principles with clean architecture:

### Layers

1. **Domain Layer** (`internal/domain/`)
   - Pure business logic
   - Domain models and aggregates
   - Repository interfaces

2. **Repository Layer** (`internal/repository/`)
   - Data access implementations
   - PostgreSQL with pgx
   - Transaction management

3. **Use Case Layer** (`internal/usecase/`)
   - Application business logic
   - Orchestrates domain operations
   - Cross-cutting concerns

4. **Controller Layer** (`internal/controller/`)
   - HTTP request handlers
   - Request/response mapping
   - Error handling

5. **Infrastructure Layer** (`pkg/`)
   - Caching (L1 + L2)
   - Configuration
   - Logging & monitoring
   - Middleware

### Key Technologies

- **Language**: Go 1.23+
- **Database**: PostgreSQL 15 with pgx/pgxpool
- **Cache**: Redis 7 + in-memory (ristretto)
- **Logging**: zerolog (structured JSON logs)
- **Metrics**: Prometheus
- **Tracing**: OpenTelemetry
- **Auth**: SIWE (Sign-In With Ethereum)
- **Security**: Row-Level Security (RLS), rate limiting, circuit breakers

**Detailed architecture**: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)

---

## 💻 Development

### Prerequisites

Install development tools:

```bash
# golangci-lint for linting
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# golang-migrate for database migrations
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# air for hot reload (optional)
go install github.com/cosmtrek/air@latest
```

### Available Make Commands

```bash
make help           # Show all available commands
make build          # Build the API server
make run-dev        # Run in development mode
make test           # Run tests
make test-coverage  # Run tests with coverage report
make lint           # Run linters
make fmt            # Format code
make clean          # Clean build artifacts
make docker-build   # Build Docker image
make docker-up      # Start Docker Compose services
make docker-down    # Stop Docker Compose services
make migrate-up     # Run database migrations up
make migrate-down   # Run database migrations down
make watch          # Run with hot reload (requires air)
```

### Development Workflow

1. **Create a feature branch**
   ```bash
   git checkout -b feature/my-feature
   ```

2. **Make changes and test**
   ```bash
   make test
   make lint
   ```

3. **Run locally**
   ```bash
   make run-dev
   ```

4. **Commit and push**
   ```bash
   git add .
   git commit -m "feat: add new feature"
   git push origin feature/my-feature
   ```

### Database Migrations

Create a new migration:

```bash
migrate create -ext sql -dir migrations -seq migration_name
```

This creates two files:
- `migrations/NNNNNN_migration_name.up.sql`
- `migrations/NNNNNN_migration_name.down.sql`

Apply migrations:

```bash
make migrate-up
```

Rollback migrations:

```bash
make migrate-down
```

---

## 🧪 Testing

### Run All Tests

```bash
make test
```

### Run Tests with Coverage

```bash
make test-coverage
```

This generates `coverage.html` for viewing in a browser.

### Test Structure

```
tests/
├── unit/           # Unit tests for domain logic
├── integration/    # Integration tests with database
└── e2e/           # End-to-end tests
```

### Writing Tests

**Unit Test Example** (`internal/domain/wallet/wallet_test.go`):

```go
func TestWallet_Debit(t *testing.T) {
    wallet := &Wallet{Balance: 100.0}
    err := wallet.Debit(50.0)
    assert.NoError(t, err)
    assert.Equal(t, 50.0, wallet.Balance)
}
```

**Integration Test Example** (`tests/integration/user_test.go`):

```go
func TestUserRepository_Create(t *testing.T) {
    // Use testcontainers for ephemeral database
    // Create user, verify in database
}
```

---

## 🚀 Deployment

### Docker

**Build image**:

```bash
docker build -t caribx-backend:latest .
```

**Run container**:

```bash
docker run -p 8080:8080 --env-file .env caribx-backend:latest
```

### Docker Compose

**Production-like setup**:

```bash
docker-compose up -d
```

### Kubernetes

*(Coming soon)*

Example deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: caribx-backend
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: api
        image: caribx-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          value: postgres-service
        - name: REDIS_HOST
          value: redis-service
```

### Environment Variables

Required environment variables (see `.env.example`):

- `PORT`: Server port (default: 8080)
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`: Database config
- `REDIS_HOST`, `REDIS_PORT`: Redis config
- `SESSION_SECRET`: Session encryption key
- `JWT_SECRET`: JWT signing key

**Security Note**: Never commit `.env` files. Use secrets management in production.

---

## 🔒 Security

### Authentication

- **SIWE (Sign-In With Ethereum)**: Wallet-based authentication
- **Session cookies**: HTTP-only, Secure, SameSite
- **JWT tokens**: For machine-to-machine auth

### Authorization

- **RBAC**: Role-based access control (customer, seller, admin)
- **Resource ownership**: Users can only access their own resources

### Data Protection

- **Row-Level Security (RLS)**: PostgreSQL policies for tenant isolation
- **PII redaction**: Sensitive data scrubbed from logs
- **Input validation**: go-playground/validator
- **Rate limiting**: Per-IP and per-user token buckets
- **Circuit breakers**: Protect against cascading failures

### Best Practices

- All passwords/secrets via environment variables
- HTTPS only in production
- Regular security audits
- Dependency vulnerability scanning
- Principle of least privilege

---

## 📊 Monitoring

### Health Checks

- `GET /healthz` - Liveness probe
- `GET /readyz` - Readiness probe (checks DB and Redis)

### Metrics

Prometheus metrics available at `/metrics`:

- Request counts and durations
- Cache hit/miss rates
- Database connection pool stats
- Error rates
- Custom business metrics

### Logging

Structured JSON logs with zerolog:

```json
{
  "level": "info",
  "time": "2025-10-18T12:00:00Z",
  "request_id": "uuid",
  "method": "GET",
  "path": "/v1/products",
  "status": 200,
  "duration": 15.5,
  "message": "request completed"
}
```

### Tracing

OpenTelemetry tracing for distributed request tracking:

- HTTP handlers
- Database queries
- Cache operations
- External API calls

---

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Commit Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, etc.)
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 👥 Team

- **Tenoy Williams** - Lead Developer - [@Tenoywil](https://github.com/Tenoywil)

---

## 🙏 Acknowledgments

- [Go Community](https://golang.org/community)
- [Domain-Driven Design](https://domainlanguage.com/ddd/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Ethereum Foundation](https://ethereum.org/) for SIWE

---

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/Tenoywil/CaribEx-backend/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Tenoywil/CaribEx-backend/discussions)
- **Email**: support@caribx.io

---

## 🗺️ Roadmap

- [x] Core API implementation
- [x] SIWE authentication
- [x] Wallet management
- [x] Marketplace features
- [ ] GraphQL API
- [ ] WebSocket support for real-time updates
- [ ] Multi-blockchain support
- [ ] Mobile SDK
- [ ] Advanced analytics
- [ ] Machine learning recommendations

---

**Built with ❤️ for the Caribbean community**