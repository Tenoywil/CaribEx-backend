package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Tenoywil/CaribEx-backend/internal/controller"
	"github.com/Tenoywil/CaribEx-backend/internal/repository/postgres"
	"github.com/Tenoywil/CaribEx-backend/internal/repository/redis"
	"github.com/Tenoywil/CaribEx-backend/internal/routes"
	"github.com/Tenoywil/CaribEx-backend/internal/usecase"
	"github.com/Tenoywil/CaribEx-backend/pkg/config"
	"github.com/Tenoywil/CaribEx-backend/pkg/logger"
	"github.com/Tenoywil/CaribEx-backend/pkg/middleware"
	"github.com/Tenoywil/CaribEx-backend/pkg/storage"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	redisclient "github.com/redis/go-redis/v9"
)

func main() {
	fmt.Println("CaribEX Backend API Server")
	fmt.Println("Version: 0.1.0")

	// Initialize logger
	appLogger := logger.New()
	appLogger.Info("Starting CaribEX Backend API Server")

	// Load configuration
	cfg := config.Load()

	// Initialize database connection pool
	dbURL := cfg.DBConnectionString

	dbConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		appLogger.Error(err, "Failed to parse database config")
		os.Exit(1)
	}

	dbConfig.MaxConns = int32(cfg.DBMaxConnections)

	db, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		appLogger.Error(err, "Failed to connect to database")
		os.Exit(1)
	}
	defer db.Close()

	appLogger.Info("Database connection established")

	// Test database connection
	if err := db.Ping(context.Background()); err != nil {
		appLogger.Error(err, "Failed to ping database")
		os.Exit(1)
	}

	// Initialize Redis client
	redisClient := redisclient.NewClient(&redisclient.Options{
		Addr:     cfg.RedisConnectionString,
		Password: cfg.RedisPassword,
		Username: "default",
	})
	defer redisClient.Close()

	appLogger.Info("Redis connection established")

	// Test Redis connection
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		appLogger.Error(err, "Failed to ping Redis")
		os.Exit(1)
	}

	// Initialize repositories
	sessionRepo := redis.NewSessionRepository(redisClient)
	userRepo := postgres.NewUserRepository(db)
	productRepo := postgres.NewProductRepository(db)
	walletRepo := postgres.NewWalletRepository(db)
	cartRepo := postgres.NewCartRepository(db)
	orderRepo := postgres.NewOrderRepository(db)

	// Initialize storage service
	storageService, err := storage.NewSupabaseStorage(storage.Config{
		URL:         cfg.SupabaseURL,
		Key:         cfg.SupabaseKey,
		Bucket:      cfg.SupabaseBucket,
		MaxFileSize: cfg.StorageMaxFileSize,
	})
	if err != nil {
		appLogger.Error(err, "Failed to initialize storage service")
		os.Exit(1)
	}
	appLogger.Info("Storage service initialized")

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo)
	authUseCase := usecase.NewAuthUseCase(sessionRepo, userUseCase, cfg.SIWEDomain)
	productUseCase := usecase.NewProductUseCase(productRepo)
	walletUseCase := usecase.NewWalletUseCase(walletRepo)
	cartUseCase := usecase.NewCartUseCase(cartRepo)
	orderUseCase := usecase.NewOrderUseCase(orderRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authUseCase)
	userController := controller.NewUserController(userUseCase)
	productController := controller.NewProductController(productUseCase, storageService)
	walletController := controller.NewWalletController(walletUseCase)
	cartController := controller.NewCartController(cartUseCase)
	orderController := controller.NewOrderController(orderUseCase)

	// Set Gin mode
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.Default()

	// Setup CORS
	router.Use(middleware.SetupCORS(cfg.AllowedOriginsSlice))

	// Setup routes
	routes.SetupRoutes(router, authController, authUseCase, userController, productController, walletController, cartController, orderController)

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	appLogger.Info(fmt.Sprintf("Server starting on %s", addr))

	// Parse timeouts
	readTimeout, _ := time.ParseDuration(cfg.ServerReadTimeout)
	writeTimeout, _ := time.ParseDuration(cfg.ServerWriteTimeout)

	// Graceful shutdown
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	shutdownTimeout, _ := time.ParseDuration(cfg.ServerShutdownTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error(err, "Server forced to shutdown")
	}

	appLogger.Info("Server exited")
}
