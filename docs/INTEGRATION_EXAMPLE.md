# Integration Example

This document shows how to integrate the SIWE authentication into your main.go file.

## Complete main.go Example

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Tenoywil/CaribEx-backend/internal/controller"
	"github.com/Tenoywil/CaribEx-backend/internal/repository/postgres"
	redisRepo "github.com/Tenoywil/CaribEx-backend/internal/repository/redis"
	"github.com/Tenoywil/CaribEx-backend/internal/routes"
	"github.com/Tenoywil/CaribEx-backend/internal/usecase"
	"github.com/Tenoywil/CaribEx-backend/pkg/config"
	"github.com/Tenoywil/CaribEx-backend/pkg/logger"
	redisClient "github.com/Tenoywil/CaribEx-backend/pkg/redis"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize logger
	logger.Init()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configuration")
	}

	// Initialize database
	dbPool, err := initDatabase(cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize database")
	}
	defer dbPool.Close()

	// Initialize Redis
	redis, err := redisClient.NewClient(cfg.Redis)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize Redis")
	}
	defer redis.Close()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(dbPool)
	productRepo := postgres.NewProductRepository(dbPool)
	walletRepo := postgres.NewWalletRepository(dbPool)
	cartRepo := postgres.NewCartRepository(dbPool)
	orderRepo := postgres.NewOrderRepository(dbPool)
	sessionRepo := redisRepo.NewSessionRepository(redis)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo)
	productUseCase := usecase.NewProductUseCase(productRepo)
	walletUseCase := usecase.NewWalletUseCase(walletRepo)
	cartUseCase := usecase.NewCartUseCase(cartRepo, productRepo)
	orderUseCase := usecase.NewOrderUseCase(orderRepo, cartRepo, walletRepo)
	authUseCase := usecase.NewAuthUseCase(sessionRepo, userUseCase, cfg.Auth.SIWEDomain)

	// Initialize controllers
	authController := controller.NewAuthController(authUseCase)
	userController := controller.NewUserController(userUseCase)
	productController := controller.NewProductController(productUseCase)
	walletController := controller.NewWalletController(walletUseCase)
	cartController := controller.NewCartController(cartUseCase)
	orderController := controller.NewOrderController(orderUseCase)

	// Initialize Gin router
	router := gin.Default()

	// Setup CORS
	router.Use(corsMiddleware(cfg.Server.AllowedOrigins))

	// Setup routes
	routes.SetupRoutes(
		router,
		authController,
		authUseCase,
		userController,
		productController,
		walletController,
		cartController,
		orderController,
	)

	// Start server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Graceful shutdown
	go func() {
		log.Info().Str("addr", srv.Addr).Msg("starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}

	log.Info().Msg("server exited")
}

func initDatabase(cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database,
	)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxConnections)
	poolConfig.MaxConnIdleTime = cfg.MaxIdleTime
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().
		Str("host", cfg.Host).
		Int("port", cfg.Port).
		Str("database", cfg.Database).
		Msg("connected to database")

	return pool, nil
}

func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Check if origin is allowed
		for _, allowed := range allowedOrigins {
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
	}
}
```

## Key Integration Points

### 1. Redis Client Initialization

```go
redis, err := redisClient.NewClient(cfg.Redis)
if err != nil {
	log.Fatal().Err(err).Msg("failed to initialize Redis")
}
defer redis.Close()
```

### 2. Session Repository

```go
sessionRepo := redisRepo.NewSessionRepository(redis)
```

### 3. Auth Use Case

```go
authUseCase := usecase.NewAuthUseCase(
	sessionRepo,
	userUseCase,
	cfg.Auth.SIWEDomain, // Important: Pass SIWE domain from config
)
```

### 4. Auth Controller

```go
authController := controller.NewAuthController(authUseCase)
```

### 5. Routes Setup

```go
routes.SetupRoutes(
	router,
	authController,  // Add auth controller
	authUseCase,     // Add auth use case for middleware
	userController,
	productController,
	walletController,
	cartController,
	orderController,
)
```

### 6. CORS Configuration

**Important**: For cookie-based authentication to work, you must:

1. Set `Access-Control-Allow-Credentials: true`
2. Set specific `Access-Control-Allow-Origin` (not `*`)
3. Frontend must include `credentials: 'include'` in requests

```go
c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
```

## Dependencies to Install

Run this command to install all dependencies:

```bash
go mod tidy
```

This will download:
- `github.com/spruceid/siwe-go` - SIWE verification
- `github.com/redis/go-redis/v9` - Redis client
- All other existing dependencies

## Startup Checklist

Before starting the server:

1. ✅ Redis is running
   ```bash
   docker-compose up redis -d
   ```

2. ✅ PostgreSQL is running
   ```bash
   docker-compose up postgres -d
   ```

3. ✅ Environment variables are set
   ```bash
   cp .env.example .env
   # Edit .env with your values
   ```

4. ✅ Database migrations are applied
   ```bash
   make migrate-up
   ```

5. ✅ Dependencies are installed
   ```bash
   go mod tidy
   ```

6. ✅ Start the server
   ```bash
   make run-dev
   ```

## Testing the Integration

### 1. Health Check

```bash
curl http://localhost:8080/healthz
```

Expected: `{"status":"ok"}`

### 2. Get Nonce

```bash
curl http://localhost:8080/v1/auth/nonce
```

Expected: `{"nonce":"...","expires_at":"..."}`

### 3. Full Authentication Flow

See `docs/SIWE_AUTH.md` for complete testing instructions.

## Common Issues

### Redis Connection Failed

**Error**: `failed to connect to Redis`

**Solution**: 
- Ensure Redis is running: `docker-compose up redis -d`
- Check `REDIS_HOST` and `REDIS_PORT` in `.env`

### Database Connection Failed

**Error**: `failed to ping database`

**Solution**:
- Ensure PostgreSQL is running: `docker-compose up postgres -d`
- Check database credentials in `.env`
- Run migrations: `make migrate-up`

### SIWE Domain Mismatch

**Error**: `domain mismatch`

**Solution**:
- Set `SIWE_DOMAIN` in `.env` to match frontend domain
- Example: `SIWE_DOMAIN=localhost:3000` (no protocol)

### CORS Issues

**Error**: Frontend can't make requests

**Solution**:
- Set `ALLOWED_ORIGIN` in `.env` to frontend URL
- Example: `ALLOWED_ORIGIN=http://localhost:3000`
- Ensure frontend includes `credentials: 'include'`
