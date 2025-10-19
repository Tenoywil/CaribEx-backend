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
	"github.com/Tenoywil/CaribEx-backend/pkg/blockchain"
	"github.com/Tenoywil/CaribEx-backend/pkg/config"
	"github.com/Tenoywil/CaribEx-backend/pkg/logger"
	"github.com/Tenoywil/CaribEx-backend/pkg/middleware"
	"github.com/Tenoywil/CaribEx-backend/pkg/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

	// Initialize blockchain RPC client (optional - only if RPC_URL is configured)
	if cfg.RPCURL != "" {
		if err := blockchain.InitRPC(cfg.RPCURL); err != nil {
			appLogger.Error(err, "Failed to initialize blockchain RPC client")
			// Don't exit - blockchain features will be unavailable but app can still run
		} else {
			defer blockchain.Close()
			appLogger.Info("Blockchain RPC client initialized")
		}
	} else {
		appLogger.Info("Blockchain RPC URL not configured - blockchain features disabled")
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

	// Initialize S3-compatible uploader for Supabase Storage
	var s3Service *storage.S3Service
	if cfg.SupabaseS3AccessKeyID != "" && cfg.SupabaseS3SecretAccessKey != "" {
		appLogger.Info("Initializing S3-compatible storage for Supabase")

		sess, err := session.NewSession(&aws.Config{
			Region:           aws.String(cfg.SupabaseRegion),
			Credentials:      credentials.NewStaticCredentials(cfg.SupabaseS3AccessKeyID, cfg.SupabaseS3SecretAccessKey, ""),
			Endpoint:         aws.String(cfg.SupabaseStorageURL),
			S3ForcePathStyle: aws.Bool(true),
		})
		if err != nil {
			appLogger.Error(err, "Failed to create S3 session")
			os.Exit(1)
		}

		s3Uploader := s3manager.NewUploader(sess)
		s3Client := s3.New(sess)
		s3Service = storage.NewS3Service(s3Uploader, s3Client, cfg.SupabaseBucket)

		appLogger.Info("S3-compatible storage initialized successfully")
	} else {
		appLogger.Info("S3 credentials not configured - file uploads will use basic storage service")
	}

	// S3Service is now available for use in controllers
	_ = s3Service // TODO: Pass to controllers that need file upload functionality

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo)
	authUseCase := usecase.NewAuthUseCase(sessionRepo, userUseCase, cfg.SIWEDomain)
	productUseCase := usecase.NewProductUseCase(productRepo)
	walletUseCase := usecase.NewWalletUseCase(walletRepo)
	cartUseCase := usecase.NewCartUseCase(cartRepo)
	orderUseCase := usecase.NewOrderUseCase(orderRepo)
	blockchainUseCase := usecase.NewBlockchainUseCase(walletRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authUseCase)
	userController := controller.NewUserController(userUseCase)
	productController := controller.NewProductController(productUseCase, storageService)
	walletController := controller.NewWalletController(walletUseCase)
	cartController := controller.NewCartController(cartUseCase)
	orderController := controller.NewOrderController(orderUseCase)
	blockchainController := controller.NewBlockchainController(blockchainUseCase)

	// Set Gin mode
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.Default()

	// Setup CORS
	router.Use(middleware.SetupCORS(cfg.AllowedOriginsSlice))

	// Setup routes
	routes.SetupRoutes(router, authController, authUseCase, userController, productController, walletController, cartController, orderController, blockchainController)

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
