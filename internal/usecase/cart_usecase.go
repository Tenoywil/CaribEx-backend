package usecase

import (
	"time"

	"github.com/Tenoywil/CaribEx-backend/internal/domain/cart"
	"github.com/google/uuid"
)

// CartUseCase handles cart business logic
type CartUseCase struct {
	cartRepo cart.Repository
}

// NewCartUseCase creates a new cart use case
func NewCartUseCase(cartRepo cart.Repository) *CartUseCase {
	return &CartUseCase{cartRepo: cartRepo}
}

// GetCartByUserID retrieves a cart by user ID
func (uc *CartUseCase) GetCartByUserID(userID string) (*cart.Cart, error) {
	return uc.cartRepo.GetByUserID(userID)
}

// GetCartItems retrieves all items in a cart
func (uc *CartUseCase) GetCartItems(cartID string) ([]*cart.CartItem, error) {
	return uc.cartRepo.GetItems(cartID)
}

// AddItemToCart adds an item to the cart
func (uc *CartUseCase) AddItemToCart(cartID, productID string, quantity int, price float64) (*cart.CartItem, error) {
	item := &cart.CartItem{
		ID:        uuid.New().String(),
		CartID:    cartID,
		ProductID: productID,
		Quantity:  quantity,
		Price:     price,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := uc.cartRepo.AddItem(item)
	if err != nil {
		return nil, err
	}

	// Update cart total
	items, err := uc.cartRepo.GetItems(cartID)
	if err != nil {
		return nil, err
	}

	total := 0.0
	for _, i := range items {
		total += i.Price * float64(i.Quantity)
	}

	err = uc.cartRepo.UpdateTotal(cartID, total)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// UpdateCartItem updates a cart item
func (uc *CartUseCase) UpdateCartItem(item *cart.CartItem) error {
	item.UpdatedAt = time.Now()
	err := uc.cartRepo.UpdateItem(item)
	if err != nil {
		return err
	}

	// Update cart total
	items, err := uc.cartRepo.GetItems(item.CartID)
	if err != nil {
		return err
	}

	total := 0.0
	for _, i := range items {
		total += i.Price * float64(i.Quantity)
	}

	return uc.cartRepo.UpdateTotal(item.CartID, total)
}

// RemoveCartItem removes an item from the cart
func (uc *CartUseCase) RemoveCartItem(cartID, itemID string) error {
	err := uc.cartRepo.RemoveItem(itemID)
	if err != nil {
		return err
	}

	// Update cart total
	items, err := uc.cartRepo.GetItems(cartID)
	if err != nil {
		return err
	}

	total := 0.0
	for _, i := range items {
		total += i.Price * float64(i.Quantity)
	}

	return uc.cartRepo.UpdateTotal(cartID, total)
}

// CheckoutCart converts cart to checked out status
func (uc *CartUseCase) CheckoutCart(cartID string) error {
	return uc.cartRepo.SetStatus(cartID, cart.CartStatusCheckedOut)
}
