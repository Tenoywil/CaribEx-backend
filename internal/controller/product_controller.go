package controller

import (
	"net/http"
	"strconv"

	"github.com/Tenoywil/CaribEx-backend/internal/usecase"
	"github.com/Tenoywil/CaribEx-backend/pkg/storage"
	"github.com/gin-gonic/gin"
)

// ProductController handles HTTP requests for products
type ProductController struct {
	productUseCase *usecase.ProductUseCase
	storageService storage.Service
}

// NewProductController creates a new product controller
func NewProductController(productUseCase *usecase.ProductUseCase, storageService storage.Service) *ProductController {
	return &ProductController{
		productUseCase: productUseCase,
		storageService: storageService,
	}
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

	p, err := c.productUseCase.GetProductByIDWithCategory(id)
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

	filters := make(map[string]interface{})
	if categoryID := ctx.Query("category_id"); categoryID != "" {
		filters["category_id"] = categoryID
	}
	if search := ctx.Query("search"); search != "" {
		filters["search"] = search
	}
	
	// Get sort parameters
	sortBy := ctx.DefaultQuery("sort_by", "created_at")
	sortOrder := ctx.DefaultQuery("sort_order", "desc")

	products, total, err := c.productUseCase.ListProductsWithCategory(filters, page, pageSize, sortBy, sortOrder)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"products":    products,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
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

// UpdateProductQuantityRequest represents the request body for updating product quantity
type UpdateProductQuantityRequest struct {
	Quantity int `json:"quantity" binding:"required,min=0"`
}

// UpdateProductQuantity handles PATCH /products/:id/quantity
func (c *ProductController) UpdateProductQuantity(ctx *gin.Context) {
	id := ctx.Param("id")

	var req UpdateProductQuantityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.productUseCase.UpdateProductQuantity(id, req.Quantity)
	if err != nil {
		if err.Error() == "product not found" || err.Error() == "failed to get product by id: no rows in result set" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get updated product to return
	p, err := c.productUseCase.GetProductByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve updated product"})
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

// UploadImageRequest represents a single image upload response
type UploadImageResponse struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
}

// UploadImage handles POST /products/upload-image for standalone image uploads
func (c *ProductController) UploadImage(ctx *gin.Context) {
	// Parse multipart form
	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}
	defer file.Close()

	// Upload to storage
	url, err := c.storageService.UploadFile(ctx.Request.Context(), file, header, "products")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, UploadImageResponse{
		URL:      url,
		Filename: header.Filename,
	})
}

// CreateProductMultipart handles POST /products with multipart/form-data
func (c *ProductController) CreateProductMultipart(ctx *gin.Context) {
	// Get seller ID from authenticated user context
	sellerID := ctx.GetString("user_id")

	// Parse form data
	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil { // 10MB max
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form data"})
		return
	}

	// Extract fields
	title := ctx.PostForm("title")
	description := ctx.PostForm("description")
	priceStr := ctx.PostForm("price")
	quantityStr := ctx.PostForm("quantity")
	categoryID := ctx.PostForm("category_id")

	// Validate required fields
	if title == "" || priceStr == "" || quantityStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "title, price, and quantity are required"})
		return
	}

	// Parse numeric fields
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid price format"})
		return
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid quantity format"})
		return
	}

	// Process uploaded images
	var imageURLs []string
	form := ctx.Request.MultipartForm
	files := form.File["images"]

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open uploaded file"})
			return
		}

		url, err := c.storageService.UploadFile(ctx.Request.Context(), file, fileHeader, "products")
		file.Close()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		imageURLs = append(imageURLs, url)
	}

	// Create product
	p, err := c.productUseCase.CreateProduct(sellerID, title, description, price, quantity, imageURLs, categoryID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, p)
}
