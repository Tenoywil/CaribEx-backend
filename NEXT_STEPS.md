# Next Steps to Complete SIWE Integration

## ✅ What's Already Done

All SIWE authentication code has been implemented:
- Auth domain models (Session, Nonce)
- Redis session repository
- SIWE verification use case
- Auth controller with endpoints
- Auth middleware for protected routes
- Routes configured with authentication
- Configuration updated
- Documentation created

## 🔧 Required Actions

### 1. Install Dependencies

Run this command in your terminal:

```bash
go mod tidy
```

This will download:
- `github.com/spruceid/siwe-go@v0.2.1`
- `github.com/redis/go-redis/v9@v9.7.0`

### 2. Update Your main.go

Your `cmd/api-server/main.go` needs to be updated to include the auth components.

**Add these imports:**
```go
import (
    // ... existing imports ...
    redisRepo "github.com/Tenoywil/CaribEx-backend/internal/repository/redis"
    redisClient "github.com/Tenoywil/CaribEx-backend/pkg/redis"
)
```

**Initialize Redis client (after database initialization):**
```go
// Initialize Redis
redis, err := redisClient.NewClient(cfg.Redis)
if err != nil {
    log.Fatal().Err(err).Msg("failed to initialize Redis")
}
defer redis.Close()
```

**Add session repository:**
```go
// Initialize repositories
userRepo := postgres.NewUserRepository(dbPool)
// ... other repos ...
sessionRepo := redisRepo.NewSessionRepository(redis)  // Add this
```

**Add auth use case:**
```go
// Initialize use cases
userUseCase := usecase.NewUserUseCase(userRepo)
// ... other use cases ...
authUseCase := usecase.NewAuthUseCase(sessionRepo, userUseCase, cfg.Auth.SIWEDomain)  // Add this
```

**Add auth controller:**
```go
// Initialize controllers
authController := controller.NewAuthController(authUseCase)  // Add this
userController := controller.NewUserController(userUseCase)
// ... other controllers ...
```

**Update routes.SetupRoutes call:**
```go
// Setup routes
routes.SetupRoutes(
    router,
    authController,  // Add this
    authUseCase,     // Add this
    userController,
    productController,
    walletController,
    cartController,
    orderController,
)
```

**Add CORS middleware (before routes setup):**
```go
// Setup CORS
router.Use(func(c *gin.Context) {
    origin := c.Request.Header.Get("Origin")
    
    // Check if origin is allowed
    for _, allowed := range cfg.Server.AllowedOrigins {
        if origin == allowed {
            c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
            break
        }
    }

    c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

    if c.Request.Method == "OPTIONS" {
        c.AbortWithStatus(204)
        return
    }

    c.Next()
})
```

See `docs/INTEGRATION_EXAMPLE.md` for a complete main.go example.

### 3. Configure Environment

Ensure your `.env` file has:

```bash
# SIWE Configuration
SIWE_DOMAIN=localhost:3000

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# CORS
ALLOWED_ORIGIN=http://localhost:3000
```

### 4. Start Redis

```bash
docker-compose up redis -d
```

Or if you don't have docker-compose configured for Redis, add to `docker-compose.yml`:

```yaml
services:
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  redis_data:
```

### 5. Test the Setup

```bash
# Start the server
go run cmd/api-server/main.go

# In another terminal, test nonce generation
curl http://localhost:8080/v1/auth/nonce
```

Expected response:
```json
{
  "nonce": "uuid-string",
  "expires_at": "timestamp"
}
```

## 📁 Files Created/Modified

### New Files Created:
- `internal/domain/auth/session.go`
- `internal/domain/auth/repository.go`
- `internal/repository/redis/session_repository.go`
- `internal/usecase/auth_usecase.go`
- `internal/controller/auth_controller.go`
- `pkg/middleware/auth.go`
- `pkg/redis/client.go`
- `docs/SIWE_AUTH.md`
- `docs/INTEGRATION_EXAMPLE.md`
- `SIWE_SETUP_COMPLETE.md`

### Modified Files:
- `go.mod` - Added SIWE and Redis dependencies
- `internal/routes/routes.go` - Added auth routes and middleware
- `pkg/config/config.go` - Added SIWEDomain field
- `.env.example` - Added SIWE_DOMAIN

## 🧪 Testing Checklist

- [ ] `go mod tidy` runs successfully
- [ ] Redis is running
- [ ] Server starts without errors
- [ ] `GET /v1/auth/nonce` returns a nonce
- [ ] `GET /healthz` returns OK
- [ ] Protected routes return 401 without auth
- [ ] Frontend can authenticate with wagmi

## 📚 Documentation

- **Setup Guide**: `SIWE_SETUP_COMPLETE.md`
- **Auth Documentation**: `docs/SIWE_AUTH.md`
- **Integration Example**: `docs/INTEGRATION_EXAMPLE.md`
- **API Reference**: `docs/API.md`

## 🎯 Frontend Integration

Once backend is running, integrate with wagmi:

```typescript
// Install dependencies
npm install wagmi viem siwe

// Use the auth flow
const { nonce } = await fetch('http://localhost:8080/v1/auth/nonce')
  .then(r => r.json());

// Sign with wagmi
const signature = await signMessageAsync({ message });

// Authenticate
await fetch('http://localhost:8080/v1/auth/siwe', {
  method: 'POST',
  credentials: 'include',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ message, signature }),
});
```

See `docs/SIWE_AUTH.md` for complete frontend integration guide.

## ⚠️ Important Notes

1. **CORS**: Must set `Access-Control-Allow-Credentials: true` and specific origin (not `*`)
2. **Cookies**: Frontend must use `credentials: 'include'` in all fetch requests
3. **Domain**: `SIWE_DOMAIN` should NOT include protocol (use `localhost:3000` not `http://localhost:3000`)
4. **Redis**: Required for session storage - must be running before starting server
5. **Production**: Update cookie settings to `secure: true` when using HTTPS

## 🚀 Ready to Go!

The codebase is now fully configured for wagmi SIWE authentication. Just follow the steps above to complete the integration.
