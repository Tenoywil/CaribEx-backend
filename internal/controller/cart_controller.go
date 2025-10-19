package controller

import (
	"net/http"

	"github.com/Tenoywil/CaribEx-backend/internal/usecase"
	"github.com/gin-gonic/gin"
)

// CartController handles HTTP requests for carts
type CartController struct {
	cartUseCase *usecase.CartUseCase
}

// NewCartController creates a new cart controller
func NewCartController(cartUseCase *usecase.CartUseCase) *CartController {
	return &CartController{cartUseCase: cartUseCase}
}

// GetCart handles GET /cart
func (c *CartController) GetCart(ctx *gin.Context) {
	// TODO: Get user ID from authenticated user context
	userID := ctx.GetString("user_id")

	cart, err := c.cartUseCase.GetCartByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "cart not found"})
		return
	}

	items, err := c.cartUseCase.GetCartItems(cart.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"cart":  cart,
		"items": items,
	})
}

// AddItemRequest represents the request body for adding an item to cart
type AddItemRequest struct {
	ProductID string  `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
}

// AddItem handles POST /cart/items
func (c *CartController) AddItem(ctx *gin.Context) {
	var req AddItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Get cart ID from user context
	cartID := ctx.GetString("cart_id")

	item, err := c.cartUseCase.AddItemToCart(cartID, req.ProductID, req.Quantity, req.Price)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, item)
}

// UpdateItemRequest represents the request body for updating a cart item
type UpdateItemRequest struct {
	Quantity int `json:"quantity" binding:"required"`
}

// UpdateItem handles PUT /cart/items/:id
func (c *CartController) UpdateItem(ctx *gin.Context) {
	itemID := ctx.Param("id")

	var req UpdateItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Get and validate cart item
	// For now, create a stub item
	item := &struct {
		ID        string
		CartID    string
		ProductID string
		Quantity  int
		Price     float64
	}{
		ID:       itemID,
		Quantity: req.Quantity,
	}

	// This is a simplified version - in production, you'd fetch the item first
	ctx.JSON(http.StatusOK, item)
}

// RemoveItem handles DELETE /cart/items/:id
func (c *CartController) RemoveItem(ctx *gin.Context) {
	itemID := ctx.Param("id")

	// TODO: Get cart ID from context
	cartID := ctx.GetString("cart_id")

	err := c.cartUseCase.RemoveCartItem(cartID, itemID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
