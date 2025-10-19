package usecase

import (
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

// ListProducts retrieves a list of products with filters
func (uc *ProductUseCase) ListProducts(filters map[string]interface{}, page, pageSize int) ([]*product.Product, int, error) {
	return uc.productRepo.List(filters, page, pageSize)
}

// UpdateProduct updates product information
func (uc *ProductUseCase) UpdateProduct(p *product.Product) error {
	p.UpdatedAt = time.Now()
	return uc.productRepo.Update(p)
}

// DeleteProduct deletes a product
func (uc *ProductUseCase) DeleteProduct(id string) error {
	return uc.productRepo.Delete(id)
}

// GetCategories retrieves all product categories
func (uc *ProductUseCase) GetCategories() ([]*product.Category, error) {
	return uc.productRepo.GetCategories()
}
