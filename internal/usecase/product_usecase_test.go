package usecase

import (
	"errors"
	"testing"

	"github.com/Tenoywil/CaribEx-backend/internal/domain/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductRepository is a mock implementation of product.Repository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(p *product.Product) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockProductRepository) GetByID(id string) (*product.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Product), args.Error(1)
}

func (m *MockProductRepository) GetByIDWithCategory(id string) (*product.ProductWithCategory, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.ProductWithCategory), args.Error(1)
}

func (m *MockProductRepository) List(filters map[string]interface{}, page, pageSize int) ([]*product.Product, int, error) {
	args := m.Called(filters, page, pageSize)
	return args.Get(0).([]*product.Product), args.Int(1), args.Error(2)
}

func (m *MockProductRepository) ListWithCategory(filters map[string]interface{}, page, pageSize int, sortBy, sortOrder string) ([]*product.ProductWithCategory, int, error) {
	args := m.Called(filters, page, pageSize, sortBy, sortOrder)
	return args.Get(0).([]*product.ProductWithCategory), args.Int(1), args.Error(2)
}

func (m *MockProductRepository) Update(p *product.Product) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockProductRepository) UpdateQuantity(id string, quantity int) error {
	args := m.Called(id, quantity)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductRepository) GetCategories() ([]*product.Category, error) {
	args := m.Called()
	return args.Get(0).([]*product.Category), args.Error(1)
}

func TestUpdateProductQuantity_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := NewProductUseCase(mockRepo)

	productID := "test-product-id"
	newQuantity := 50

	// Mock product exists
	mockRepo.On("GetByID", productID).Return(&product.Product{
		ID:       productID,
		Quantity: 100,
	}, nil)

	// Mock successful update
	mockRepo.On("UpdateQuantity", productID, newQuantity).Return(nil)

	err := uc.UpdateProductQuantity(productID, newQuantity)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProductQuantity_NegativeQuantity(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := NewProductUseCase(mockRepo)

	productID := "test-product-id"
	negativeQuantity := -5

	err := uc.UpdateProductQuantity(productID, negativeQuantity)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quantity cannot be negative")
	mockRepo.AssertNotCalled(t, "GetByID")
	mockRepo.AssertNotCalled(t, "UpdateQuantity")
}

func TestUpdateProductQuantity_ProductNotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := NewProductUseCase(mockRepo)

	productID := "non-existent-product"
	newQuantity := 50

	// Mock product not found
	mockRepo.On("GetByID", productID).Return(nil, errors.New("product not found"))

	err := uc.UpdateProductQuantity(productID, newQuantity)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product not found")
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "UpdateQuantity")
}

func TestUpdateProductQuantity_ZeroQuantity(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := NewProductUseCase(mockRepo)

	productID := "test-product-id"
	zeroQuantity := 0

	// Mock product exists
	mockRepo.On("GetByID", productID).Return(&product.Product{
		ID:       productID,
		Quantity: 10,
	}, nil)

	// Mock successful update
	mockRepo.On("UpdateQuantity", productID, zeroQuantity).Return(nil)

	err := uc.UpdateProductQuantity(productID, zeroQuantity)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProductQuantity_RepositoryError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := NewProductUseCase(mockRepo)

	productID := "test-product-id"
	newQuantity := 50

	// Mock product exists
	mockRepo.On("GetByID", productID).Return(&product.Product{
		ID:       productID,
		Quantity: 100,
	}, nil)

	// Mock repository error during update
	mockRepo.On("UpdateQuantity", productID, newQuantity).Return(errors.New("database error"))

	err := uc.UpdateProductQuantity(productID, newQuantity)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockRepo.AssertExpectations(t)
}
