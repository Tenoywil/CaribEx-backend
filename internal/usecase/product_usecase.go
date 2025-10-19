package usecase

import (
	"fmt"
	"time"

	"github.com/Tenoywil/CaribEx-backend/internal/domain/product"
	"github.com/google/uuid"
)

// ProductUseCase handles product business logic
type ProductUseCase struct {
	productRepo product.Repository
}

// NewProductUseCase creates a new product use case
func NewProductUseCase(productRepo product.Repository) *ProductUseCase {
	return &ProductUseCase{productRepo: productRepo}
}

// CreateProduct creates a new product
func (uc *ProductUseCase) CreateProduct(sellerID, title, description string, price float64, quantity int, images []string, categoryID string) (*product.Product, error) {
	p := &product.Product{
		ID:          uuid.New().String(),
		SellerID:    sellerID,
		Title:       title,
		Description: description,
		Price:       price,
		Quantity:    quantity,
		Images:      images,
		CategoryID:  categoryID,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := uc.productRepo.Create(p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// GetProductByID retrieves a product by ID
func (uc *ProductUseCase) GetProductByID(id string) (*product.Product, error) {
	return uc.productRepo.GetByID(id)
}

// GetProductByIDWithCategory retrieves a product by ID with category details
func (uc *ProductUseCase) GetProductByIDWithCategory(id string) (*product.ProductWithCategory, error) {
	return uc.productRepo.GetByIDWithCategory(id)
}

// ListProducts retrieves a list of products with filters
func (uc *ProductUseCase) ListProducts(filters map[string]interface{}, page, pageSize int) ([]*product.Product, int, error) {
	return uc.productRepo.List(filters, page, pageSize)
}

// ListProductsWithCategory retrieves a list of products with category details and sorting
func (uc *ProductUseCase) ListProductsWithCategory(filters map[string]interface{}, page, pageSize int, sortBy, sortOrder string) ([]*product.ProductWithCategory, int, error) {
	return uc.productRepo.ListWithCategory(filters, page, pageSize, sortBy, sortOrder)
}

// UpdateProduct updates product information
func (uc *ProductUseCase) UpdateProduct(p *product.Product) error {
	p.UpdatedAt = time.Now()
	return uc.productRepo.Update(p)
}

// UpdateProductQuantity updates only the quantity of a product
func (uc *ProductUseCase) UpdateProductQuantity(id string, quantity int) error {
	// Validate quantity is not negative
	if quantity < 0 {
		return fmt.Errorf("quantity cannot be negative")
	}
	
	// Verify product exists
	_, err := uc.productRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}
	
	return uc.productRepo.UpdateQuantity(id, quantity)
}

// DeleteProduct deletes a product
func (uc *ProductUseCase) DeleteProduct(id string) error {
	return uc.productRepo.Delete(id)
}

// GetCategories retrieves all product categories
func (uc *ProductUseCase) GetCategories() ([]*product.Category, error) {
	return uc.productRepo.GetCategories()
}
