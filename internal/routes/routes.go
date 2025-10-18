package routes

import (
	"github.com/Tenoywil/CaribEx-backend/internal/controller"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	router *gin.Engine,
	userController *controller.UserController,
	productController *controller.ProductController,
	walletController *controller.WalletController,
	cartController *controller.CartController,
	orderController *controller.OrderController,
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
		// User routes
		users := v1.Group("/users")
		{
			users.POST("", userController.CreateUser)
			users.GET("/:id", userController.GetUser)
			users.GET("/wallet/:address", userController.GetUserByWallet)
			users.PUT("/:id", userController.UpdateUser)
			users.DELETE("/:id", userController.DeleteUser)
		}

		// Product routes
		products := v1.Group("/products")
		{
			products.POST("", productController.CreateProduct)
			products.GET("", productController.ListProducts)
			products.GET("/:id", productController.GetProduct)
			products.PUT("/:id", productController.UpdateProduct)
			products.DELETE("/:id", productController.DeleteProduct)
		}

		// Category routes
		v1.GET("/categories", productController.GetCategories)

		// Wallet routes
		wallet := v1.Group("/wallet")
		{
			wallet.GET("", walletController.GetWallet)
			wallet.POST("/send", walletController.SendFunds)
			wallet.POST("/receive", walletController.ReceiveFunds)
			wallet.GET("/transactions", walletController.GetTransactions)
		}

		// Cart routes
		cart := v1.Group("/cart")
		{
			cart.GET("", cartController.GetCart)
			cart.POST("/items", cartController.AddItem)
			cart.PUT("/items/:id", cartController.UpdateItem)
			cart.DELETE("/items/:id", cartController.RemoveItem)
		}

		// Order routes
		orders := v1.Group("/orders")
		{
			orders.POST("", orderController.CreateOrder)
			orders.GET("", orderController.ListOrders)
			orders.GET("/:id", orderController.GetOrder)
		}
	}
}
