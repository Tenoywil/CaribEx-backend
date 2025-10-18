package controller

import (
	"net/http"
	"strconv"

	"github.com/Tenoywil/CaribEx-backend/internal/usecase"
	"github.com/gin-gonic/gin"
)

// ProductController handles HTTP requests for products
type ProductController struct {
	productUseCase *usecase.ProductUseCase
}

// NewProductController creates a new product controller
func NewProductController(productUseCase *usecase.ProductUseCase) *ProductController {
	return &ProductController{productUseCase: productUseCase}
}

// CreateProductRequest represents the request body for creating a product
type CreateProductRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	Price       float64  `json:"price" binding:"required"`
	Quantity    int      `json:"quantity" binding:"required"`
	Images      []string `json:"images"`
	CategoryID  string   `json:"category_id"`
}

// CreateProduct handles POST /products
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var req CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Get seller ID from authenticated user context
	sellerID := ctx.GetString("user_id")

	p, err := c.productUseCase.CreateProduct(sellerID, req.Title, req.Description, req.Price, req.Quantity, req.Images, req.CategoryID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, p)
}

// GetProduct handles GET /products/:id
func (c *ProductController) GetProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	p, err := c.productUseCase.GetProductByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	ctx.JSON(http.StatusOK, p)
}

// ListProducts handles GET /products
func (c *ProductController) ListProducts(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	filters := make(map[string]interface{})
	if categoryID := ctx.Query("category_id"); categoryID != "" {
		filters["category_id"] = categoryID
	}
	if search := ctx.Query("search"); search != "" {
		filters["search"] = search
	}

	products, total, err := c.productUseCase.ListProducts(filters, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"products":   products,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_pages": (total + pageSize - 1) / pageSize,
	})
}

// UpdateProduct handles PUT /products/:id
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	// Get existing product
	p, err := c.productUseCase.GetProductByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	// Bind updates
	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p.ID = id
	err = c.productUseCase.UpdateProduct(p)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, p)
}

// DeleteProduct handles DELETE /products/:id
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.productUseCase.DeleteProduct(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetCategories handles GET /categories
func (c *ProductController) GetCategories(ctx *gin.Context) {
	categories, err := c.productUseCase.GetCategories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, categories)
}
