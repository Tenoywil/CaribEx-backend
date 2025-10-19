package controller

import (
	"net/http"
	"strconv"

	"github.com/Tenoywil/CaribEx-backend/internal/usecase"
	"github.com/gin-gonic/gin"
)

// OrderController handles HTTP requests for orders
type OrderController struct {
	orderUseCase *usecase.OrderUseCase
}

// NewOrderController creates a new order controller
func NewOrderController(orderUseCase *usecase.OrderUseCase) *OrderController {
	return &OrderController{orderUseCase: orderUseCase}
}

// CreateOrderRequest represents the request body for creating an order
type CreateOrderRequest struct {
	CartID     string  `json:"cart_id" binding:"required"`
	PaymentRef string  `json:"payment_ref"`
	Total      float64 `json:"total" binding:"required"`
}

// CreateOrder handles POST /orders
func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var req CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Get user ID from authenticated user context
	userID := ctx.GetString("user_id")

	order, err := c.orderUseCase.CreateOrder(userID, req.CartID, req.Total, req.PaymentRef)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, order)
}

// GetOrder handles GET /orders/:id
func (c *OrderController) GetOrder(ctx *gin.Context) {
	id := ctx.Param("id")

	order, err := c.orderUseCase.GetOrderByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	items, err := c.orderUseCase.GetOrderItems(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"order": order,
		"items": items,
	})
}

// ListOrders handles GET /orders
func (c *OrderController) ListOrders(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	// Ensure page and pageSize are within valid ranges
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// TODO: Get user ID from authenticated user context
	userID := ctx.GetString("user_id")

	orders, total, err := c.orderUseCase.GetOrdersByUserID(userID, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"orders":      orders,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + pageSize - 1) / pageSize,
	})
}
