package routes

import (
	"github.com/Tenoywil/CaribEx-backend/internal/controller"
	"github.com/Tenoywil/CaribEx-backend/internal/usecase"
	"github.com/Tenoywil/CaribEx-backend/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	router *gin.Engine,
	authController *controller.AuthController,
	authUseCase *usecase.AuthUseCase,
	userController *controller.UserController,
	productController *controller.ProductController,
	walletController *controller.WalletController,
	cartController *controller.CartController,
	orderController *controller.OrderController,
	blockchainController *controller.BlockchainController,
) {
	// Health check
	router.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	router.GET("/readyz", func(ctx *gin.Context) {
		// TODO: Check database and redis connectivity
		ctx.JSON(200, gin.H{"status": "ready"})
	})

	// API v1 routes
	v1 := router.Group("/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.GET("/nonce", authController.GetNonce)
			auth.POST("/siwe", authController.AuthenticateSIWE)
			auth.GET("/me", middleware.AuthMiddleware(authUseCase), authController.GetMe)
			auth.POST("/logout", middleware.AuthMiddleware(authUseCase), authController.Logout)
		}

		// User routes (protected)
		users := v1.Group("/users", middleware.AuthMiddleware(authUseCase))
		{
			users.POST("", userController.CreateUser)
			users.GET("/:id", userController.GetUser)
			users.GET("/wallet/:address", userController.GetUserByWallet)
			users.PUT("/:id", userController.UpdateUser)
			users.DELETE("/:id", userController.DeleteUser)
		}

		// Product routes (public read, protected write)
		products := v1.Group("/products")
		{
			products.GET("", productController.ListProducts)
			products.GET("/:id", productController.GetProduct)
			
			// Protected product routes
			productsProtected := products.Group("", middleware.AuthMiddleware(authUseCase))
			{
				productsProtected.POST("", productController.CreateProduct)
				productsProtected.PUT("/:id", productController.UpdateProduct)
				productsProtected.DELETE("/:id", productController.DeleteProduct)
			}
		}

		// Category routes (public)
		v1.GET("/categories", productController.GetCategories)

		// Wallet routes (protected)
		wallet := v1.Group("/wallet", middleware.AuthMiddleware(authUseCase))
		{
			wallet.GET("", walletController.GetWallet)
			wallet.POST("/send", walletController.SendFunds)
			wallet.POST("/receive", walletController.ReceiveFunds)
			wallet.GET("/transactions", walletController.GetTransactions)
			wallet.POST("/verify-transaction", blockchainController.VerifyTransaction)
			wallet.GET("/transaction-status", blockchainController.GetTransactionStatus)
		}

		// Cart routes (protected)
		cart := v1.Group("/cart", middleware.AuthMiddleware(authUseCase))
		{
			cart.GET("", cartController.GetCart)
			cart.POST("/items", cartController.AddItem)
			cart.PUT("/items/:id", cartController.UpdateItem)
			cart.DELETE("/items/:id", cartController.RemoveItem)
		}

		// Order routes (protected)
		orders := v1.Group("/orders", middleware.AuthMiddleware(authUseCase))
		{
			orders.POST("", orderController.CreateOrder)
			orders.GET("", orderController.ListOrders)
			orders.GET("/:id", orderController.GetOrder)
		}
	}
}
