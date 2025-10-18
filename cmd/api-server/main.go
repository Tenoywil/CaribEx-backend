package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Tenoywil/CaribEx-backend/internal/controller"
	"github.com/Tenoywil/CaribEx-backend/internal/repository/postgres"
	"github.com/Tenoywil/CaribEx-backend/internal/routes"
	"github.com/Tenoywil/CaribEx-backend/internal/usecase"
	"github.com/Tenoywil/CaribEx-backend/pkg/config"
	"github.com/Tenoywil/CaribEx-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	fmt.Println("CaribEX Backend API Server")
	fmt.Println("Version: 0.1.0")

	// Initialize logger
	appLogger := logger.New()
	appLogger.Info("Starting CaribEX Backend API Server")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		appLogger.Error(err, "Failed to load configuration")
		os.Exit(1)
	}

	// Initialize database connection pool
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Database,
	)

	dbConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		appLogger.Error(err, "Failed to parse database config")
		os.Exit(1)
	}

	dbConfig.MaxConns = int32(cfg.Database.MaxConnections)

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

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	productRepo := postgres.NewProductRepository(db)
	walletRepo := postgres.NewWalletRepository(db)
	cartRepo := postgres.NewCartRepository(db)
	orderRepo := postgres.NewOrderRepository(db)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo)
	productUseCase := usecase.NewProductUseCase(productRepo)
	walletUseCase := usecase.NewWalletUseCase(walletRepo)
	cartUseCase := usecase.NewCartUseCase(cartRepo)
	orderUseCase := usecase.NewOrderUseCase(orderRepo)

	// Initialize controllers
	userController := controller.NewUserController(userUseCase)
	productController := controller.NewProductController(productUseCase)
	walletController := controller.NewWalletController(walletUseCase)
	cartController := controller.NewCartController(cartUseCase)
	orderController := controller.NewOrderController(orderUseCase)

	// Set Gin mode
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, userController, productController, walletController, cartController, orderController)

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	appLogger.Info(fmt.Sprintf("Server starting on %s", addr))

	// Graceful shutdown
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
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

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error(err, "Server forced to shutdown")
	}

	appLogger.Info("Server exited")
}
