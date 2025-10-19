# CaribEX Backend Architecture

## Overview

CaribEX backend is built using **Domain-Driven Design (DDD)** principles with a clean architecture approach. The system is designed for high performance, security, and scalability, featuring a two-level caching strategy and comprehensive observability.

---

## Architecture Layers

### 1. Domain Layer (`internal/domain/`)

Contains pure business logic and domain models. No dependencies on infrastructure.

**Responsibilities**:
- Define domain entities and value objects
- Define repository interfaces
- Contain business rules and validations
- Aggregate roots (User, Wallet, Product, Cart, Order)

**Key Packages**:
- `user/`: User identity and roles
- `wallet/`: Wallet and transaction management
- `product/`: Product catalog and categories
- `cart/`: Shopping cart operations
- `order/`: Order processing and fulfillment

### 2. Repository Layer (`internal/repository/`)

Implements data access patterns for each domain.

**Responsibilities**:
- Implement domain repository interfaces
- Handle database connections and queries
- Transaction management
- Data mapping between domain models and database schemas

**Technologies**:
- PostgreSQL with pgx/pgxpool
- Connection pooling
- Prepared statements
- Row-Level Security (RLS) for multi-tenancy

### 3. Use Case Layer (`internal/usecase/`)

Contains application business logic and orchestrates domain operations.

**Responsibilities**:
- Coordinate multiple domain operations
- Implement business workflows
- Handle cross-cutting concerns
- Apply business policies

**Examples**:
- Checkout flow (cart → order → inventory decrement → wallet debit)
- Send funds (validate balance → create transaction → update ledger)
- Product listing (fetch from cache → query DB → populate cache)

### 4. Controller Layer (`internal/controller/`)

HTTP request handlers that adapt external requests to use case calls.

**Responsibilities**:
- Parse and validate HTTP requests
- Call appropriate use cases
- Format HTTP responses
- Handle errors and status codes

### 5. Routes Layer (`internal/routes/`)

Defines API routes and applies middleware chains.

**Responsibilities**:
- Route registration
- Middleware application
- Route grouping (public, authenticated, admin)

### 6. Infrastructure Layer (`pkg/`)

Cross-cutting concerns and infrastructure services.

**Packages**:
- `config/`: Configuration management
- `cache/`: L1 (in-memory) and L2 (Redis) caching
- `logger/`: Structured logging with zerolog
- `monitoring/`: Prometheus metrics and OpenTelemetry tracing
- `middleware/`: HTTP middleware (auth, rate limiting, CORS, etc.)

---

## Data Flow

### Read Flow (with caching)

```
Client Request
    ↓
Controller (validate request)
    ↓
Use Case (check L1 cache)
    ↓ (miss)
Use Case (check L2 Redis)
    ↓ (miss)
Repository (query PostgreSQL)
    ↓
Use Case (populate L2 & L1)
    ↓
Controller (format response)
    ↓
Client Response
```

### Write Flow

```
Client Request
    ↓
Controller (validate request)
    ↓
Use Case (business logic)
    ↓
Repository (transactional write)
    ↓
Use Case (invalidate L1 & L2 caches)
    ↓
Use Case (emit events/metrics)
    ↓
Controller (format response)
    ↓
Client Response
```

---

## Caching Strategy

### L1 Cache (In-Memory)

**Implementation**: ristretto or custom map with singleflight

**Characteristics**:
- Process-local cache
- Sub-millisecond latency
- LRU/LFU eviction
- Max size: 100MB (configurable)

**Use Cases**:
- Hot product listings
- User session data
- Frequently accessed categories

**TTL**: 1-5 minutes

### L2 Cache (Redis)

**Implementation**: go-redis/redis/v8

**Characteristics**:
- Shared across instances
- 1-2ms latency
- Persistence optional
- Key expiration

**Use Cases**:
- Product details (`product:<id>`)
- Product list pages (`products:page:<p>`)
- Wallet balances (`wallet:<user_id>`)
- Recommendation lists

**TTL**: 5-15 minutes (varies by data type)

### Cache Invalidation

**Write-through**: On mutation, update DB first, then delete cache keys.

**Strategies**:
- Delete specific keys (e.g., `product:<id>`)
- Delete pattern keys (e.g., `products:*`)
- Lazy repopulation on next read

**Stampede Protection**:
- Use `singleflight` to coalesce concurrent cache-miss requests
- Only one goroutine fetches from DB while others wait

---

## Security

### Authentication

**SIWE (Sign-In With Ethereum)**:
1. Client requests nonce from `/auth/nonce`
2. Client signs message with wallet
3. Server verifies signature
4. Server creates session and sets HTTP-only cookie

**Session Management**:
- HTTP-only, Secure, SameSite cookies
- Short TTL (24h default)
- Refresh tokens for extended sessions
- Redis-backed session store (optional)

**JWT Tokens**:
- For machine-to-machine (M2M) communication
- Short-lived (1h)
- Signed with HS256/RS256

### Authorization

**Role-Based Access Control (RBAC)**:
- Roles: `customer`, `seller`, `admin`
- Middleware checks user role before controller execution

**Resource Ownership**:
- Sellers can only modify their own products
- Users can only access their own wallets/orders

### Data Protection

**Row-Level Security (RLS)**:
- PostgreSQL policies enforce tenant isolation
- Set `app.current_tenant` via session variable

**PII Handling**:
- Hash sensitive IDs in logs
- Never log wallet addresses or balances
- Sanitize error messages

**Input Validation**:
- go-playground/validator for DTO validation
- Sanitize strings (trim, normalize)
- SQL injection protection via parameterized queries

### Rate Limiting

**Per-IP Limits**:
- Token bucket algorithm
- Redis-backed for cluster mode

**Per-User Limits**:
- Applied after authentication
- Wallet operations: stricter limits

**Circuit Breaker**:
- Wrap external service calls (price feeds, blockchain RPCs)
- Use sony/gobreaker

---

## Observability

### Logging

**Framework**: zerolog

**Structured Fields**:
- `request_id`: Unique per request
- `user_hash`: Hashed user ID
- `method`, `path`, `status`, `duration`
- `error`: Error details (sanitized)

**Log Levels**:
- DEBUG: Development only
- INFO: Normal operations
- WARN: Recoverable errors
- ERROR: Unexpected errors

**PII Redaction**:
- Automatic scrubbing of sensitive fields
- Configurable redaction patterns

### Metrics

**Framework**: Prometheus

**Metrics Types**:
- **Counters**: Request counts, error counts
- **Histograms**: Request duration, database query time
- **Gauges**: Cache size, connection pool size, goroutine count

**Custom Metrics**:
- `CaribEX_requests_total{method, path, status}`
- `CaribEX_request_duration_seconds{method, path}`
- `CaribEX_cache_hits_total{level}`
- `CaribEX_cache_misses_total{level}`
- `CaribEX_wallet_transactions_total{type, status}`

### Tracing

**Framework**: OpenTelemetry

**Instrumentation**:
- HTTP handlers (automatic)
- Database queries (pgx tracing)
- Cache operations
- External API calls

**Exporters**:
- Jaeger (development)
- Cloud provider (production)

### Health Checks

**Endpoints**:
- `/healthz`: Basic liveness check
- `/readyz`: Readiness check (DB + Redis connectivity)

**Kubernetes Integration**:
```yaml
livenessProbe:
  httpGet:
    path: /healthz
    port: 8080
readinessProbe:
  httpGet:
    path: /readyz
    port: 8080
```

---

## Database Schema

See `migrations/000001_init_schema.up.sql` for the complete schema.

**Key Tables**:
- `users`: User accounts and wallet addresses
- `wallets`: User balances per currency
- `transactions`: Immutable ledger of all wallet operations
- `products`: Marketplace listings
- `categories`: Product categories
- `carts` + `cart_items`: Shopping cart state
- `orders` + `order_items`: Order records
- `refresh_tokens`: Session refresh tokens

**Indexes**:
- Users: `wallet_address`
- Products: `seller_id`, `category_id`, `is_active`
- Transactions: `wallet_id`, `status`, `created_at` (DESC)
- Orders: `user_id`, `status`, `created_at` (DESC)

**Constraints**:
- Foreign keys with appropriate `ON DELETE` actions
- Check constraints for enums and positive amounts
- Unique constraints for wallet addresses and cart item pairs

---

## Testing Strategy

### Unit Tests

**Target**: Domain logic and use cases

**Tools**: Go standard `testing` package + `testify`

**Examples**:
- Wallet debit/credit validation
- Product quantity decrement
- Cart total calculation

### Integration Tests

**Target**: Repository layer + database

**Tools**: `testcontainers-go` or docker-compose

**Setup**:
- Spin up ephemeral Postgres container
- Run migrations
- Execute tests
- Teardown container

### Controller Tests

**Target**: HTTP handlers

**Tools**: `httptest` package

**Examples**:
- Request validation
- Authentication checks
- Response formatting

### E2E Tests

**Target**: Complete user flows

**Scenarios**:
- SIWE login → browse products → add to cart → checkout → verify wallet debit
- Seller creates product → customer purchases → order status updates

---

## Deployment

### Docker

**Build**:
```bash
docker build -t CaribEX-backend:latest .
```

**Run**:
```bash
docker run -p 8080:8080 --env-file .env CaribEX-backend:latest
```

### Docker Compose

**Development**:
```bash
docker-compose up -d
```

Starts:
- PostgreSQL (port 5432)
- Redis (port 6379)
- API Server (port 8080)

### Kubernetes

**Manifests**: (to be created)
- Deployment
- Service
- Ingress
- ConfigMap
- Secret

**Considerations**:
- Horizontal Pod Autoscaler (HPA)
- Persistent Volume Claims (PVCs) for databases
- Resource limits and requests

---

## Configuration Management

### Environment Variables

All configuration via environment variables (see `.env.example`).

**Categories**:
- Server (port, host, timeouts)
- Database (connection details, pool settings)
- Redis (connection details)
- Auth (secrets, token TTL)
- Cache (enable flags, TTL, sizes)

### Secrets Management

**Development**: `.env` file

**Production**:
- Kubernetes Secrets
- HashiCorp Vault
- Cloud provider secret managers (AWS Secrets Manager, GCP Secret Manager)

---

## Performance Considerations

### Database

- **Connection Pooling**: Max 25 connections (configurable)
- **Prepared Statements**: Reuse for common queries
- **Indexes**: Strategic indexes on foreign keys and filter columns
- **Query Optimization**: EXPLAIN ANALYZE for slow queries

### Caching

- **Hit Rate Target**: >80% for L1, >60% for L2
- **Eviction Policy**: LRU/LFU to keep hot data
- **TTL Tuning**: Balance freshness vs. hit rate

### Concurrency

- **Goroutines**: Bounded concurrency for external calls
- **Worker Pools**: For background tasks
- **Context Timeouts**: Cancel slow operations

### Resource Limits

- **Memory**: 512MB-2GB per instance (depends on cache size)
- **CPU**: 0.5-2 cores per instance
- **File Descriptors**: Increase limit for high concurrency

---

## Error Handling

### Error Types

- **Validation Errors**: 400 Bad Request
- **Authentication Errors**: 401 Unauthorized
- **Authorization Errors**: 403 Forbidden
- **Not Found Errors**: 404 Not Found
- **Rate Limit Errors**: 429 Too Many Requests
- **Internal Errors**: 500 Internal Server Error

### Error Response Format

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable message",
    "details": {}
  }
}
```

### Logging Errors

- **Client Errors (4xx)**: WARN level
- **Server Errors (5xx)**: ERROR level
- **Include**: `request_id`, `path`, `method`, stack trace (5xx only)

---

## Future Enhancements

- GraphQL API for flexible queries
- WebSocket support for real-time updates
- Event-driven architecture with message queues (NATS, Kafka)
- Multi-region deployment with data replication
- Machine learning for product recommendations
- Advanced fraud detection
- Support for multiple blockchain networks

---

## References

- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Redis Documentation](https://redis.io/documentation)
- [Prometheus Best Practices](https://prometheus.io/docs/practices/)
- [OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)
