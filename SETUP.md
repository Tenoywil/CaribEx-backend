# Quick Setup Guide

This guide will help you get CaribEX backend up and running in less than 5 minutes.

## Prerequisites Checklist

- [ ] Go 1.23+ installed ([Download](https://golang.org/dl/))
- [ ] Docker installed ([Download](https://www.docker.com/get-started))
- [ ] Docker Compose installed (included with Docker Desktop)
- [ ] Git installed

## Quick Start (Docker)

The fastest way to get started is using Docker Compose:

### 1. Clone the Repository

```bash
git clone https://github.com/Tenoywil/CaribEx-backend.git
cd CaribEx-backend
```

### 2. Set Up Environment Variables

```bash
cp .env.example .env
# The default values in .env.example work for local development
```

### 3. Start All Services

```bash
make docker-up
```

This single command will:
- Start PostgreSQL database on port 5432
- Start Redis cache on port 6379
- Build and start the API server on port 8080
- Run database migrations automatically

### 4. Verify It's Working

```bash
curl http://localhost:8080/healthz
```

You should see a health check response.

### 5. Explore the Seeded Marketplace Data

The database comes pre-loaded with sample Caribbean marketplace products! 🏝️

```bash
# Browse all products
curl http://localhost:8080/v1/products

# Search for coffee
curl "http://localhost:8080/v1/products?search=coffee"

# Get all categories
curl http://localhost:8080/v1/categories
```

**Sample data includes:**
- 8 seller accounts with wallets
- 40+ Caribbean-themed products
- Products across all categories (Electronics, Fashion, Food & Beverages, etc.)
- Realistic Jamaican pricing and descriptions

See `migrations/README.md` for complete details on seeded data.

**That's it! Your API is now running with sample data.** 🎉

## Quick Start (Local Development)

If you prefer to run services locally without Docker:

### 1. Start Infrastructure Services

You'll need PostgreSQL and Redis running. You can use Docker for these:

```bash
# Start only database and cache
docker-compose up -d postgres redis
```

### 2. Install Dependencies

```bash
make deps
```

### 3. Run Database Migrations

```bash
make migrate-up
```

### 4. Start the API Server

```bash
make run-dev
```

The API will be available at `http://localhost:8080`.

## Common Commands

```bash
# View all available commands
make help

# Build the application
make build

# Run tests
make test

# Run linter
make lint

# Format code
make fmt

# Stop Docker services
make docker-down

# View logs
docker-compose logs -f api

# Access database
docker exec -it CaribEX-postgres psql -U postgres -d CaribEX

# Access Redis CLI
docker exec -it CaribEX-redis redis-cli
```

## Project Structure Quick Reference

```
CaribEx-backend/
├── cmd/api-server/          # Application entry point
├── internal/
│   ├── domain/             # Business logic (user, wallet, product, etc.)
│   ├── repository/         # Database access
│   ├── usecase/           # Application logic
│   ├── controller/        # HTTP handlers
│   └── routes/            # Route definitions
├── pkg/                    # Shared packages (config, logger, cache, etc.)
├── migrations/            # Database migrations
├── docs/                  # Documentation
└── tests/                 # Tests
```

## API Endpoints Quick Reference

Once running, you can test these endpoints:

### Health Check
```bash
curl http://localhost:8080/healthz
```

### Authentication (Coming Soon)
```bash
# Get nonce for SIWE
curl http://localhost:8080/v1/auth/nonce

# Authenticate with signature
curl -X POST http://localhost:8080/v1/auth/siwe \
  -H "Content-Type: application/json" \
  -d '{"message": "...", "signature": "0x...", "wallet_address": "0x..."}'
```

### Products (Coming Soon)
```bash
# List products
curl http://localhost:8080/v1/products

# Get product details
curl http://localhost:8080/v1/products/{id}
```

For full API documentation, see [docs/API.md](docs/API.md).

## Troubleshooting

### Port Already in Use

If you get "port already in use" errors:

```bash
# Check what's using port 8080
lsof -i :8080

# Or change the port in .env
echo "PORT=8081" >> .env
```

### Database Connection Issues

```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# View PostgreSQL logs
docker logs CaribEX-postgres

# Restart PostgreSQL
docker-compose restart postgres
```

### Redis Connection Issues

```bash
# Check if Redis is running
docker ps | grep redis

# Test Redis connection
docker exec CaribEX-redis redis-cli ping
# Should return: PONG
```

### Migration Errors

```bash
# Check migration status
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/CaribEX?sslmode=disable" version

# Force to a specific version (use with caution)
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/CaribEX?sslmode=disable" force 1
```

### Build Errors

```bash
# Clean and rebuild
make clean
make deps
make build
```

## Next Steps

1. **Read the Documentation**
   - [API Documentation](docs/API.md) - Learn about available endpoints
   - [Architecture Guide](docs/ARCHITECTURE.md) - Understand the system design
   - [Development Guide](docs/DEVELOPMENT.md) - Deep dive into development

2. **Explore the Code**
   - Start with `cmd/api-server/main.go`
   - Look at domain models in `internal/domain/`
   - Check out the database schema in `migrations/`

3. **Start Developing**
   - Read [CONTRIBUTING.md](CONTRIBUTING.md)
   - Pick an issue from GitHub Issues
   - Make your first contribution!

4. **Join the Community**
   - Star the repository ⭐
   - Follow for updates
   - Join discussions

## Getting Help

If you run into issues:

1. Check the [troubleshooting section](#troubleshooting) above
2. Search [existing issues](https://github.com/Tenoywil/CaribEx-backend/issues)
3. Ask in [GitHub Discussions](https://github.com/Tenoywil/CaribEx-backend/discussions)
4. Create a new issue with details about your problem

## Environment Variables Reference

Key environment variables (see `.env.example` for all options):

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | API server port | `8080` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_NAME` | Database name | `CaribEX` |
| `REDIS_HOST` | Redis host | `localhost` |
| `REDIS_PORT` | Redis port | `6379` |
| `SESSION_SECRET` | Session encryption key | (change in production) |
| `JWT_SECRET` | JWT signing key | (change in production) |

## Development Tools

Optional but recommended tools:

### Install golangci-lint (for linting)
```bash
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

### Install golang-migrate (for migrations)
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Install air (for hot reload)
```bash
go install github.com/cosmtrek/air@latest
# Then run: make watch
```

## Quick Reference Card

```bash
# Development workflow
make docker-up      # Start all services
make migrate-up     # Run migrations
make run-dev        # Start API server
make test           # Run tests
make lint           # Check code quality

# Docker commands
make docker-build   # Build Docker image
make docker-down    # Stop all services

# Database commands
make migrate-up     # Apply migrations
make migrate-down   # Rollback migrations

# Cleanup
make clean          # Remove build artifacts
make docker-down    # Stop and remove containers
```

---

**Happy coding! 🚀**

For more detailed information, check out the [full documentation](docs/).
