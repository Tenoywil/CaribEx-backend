package cart

import "time"

// CartStatus represents the status of a cart
type CartStatus string

const (
	CartStatusActive     CartStatus = "active"
	CartStatusCheckedOut CartStatus = "checked_out"
)

// Cart represents a shopping cart
type Cart struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Status    CartStatus `json:"status"`
	Total     float64    `json:"total"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// CartItem represents an item in a cart
type CartItem struct {
	ID        string    `json:"id"`
	CartID    string    `json:"cart_id"`
	ProductID string    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Repository defines the interface for cart data operations
type Repository interface {
	GetByUserID(userID string) (*Cart, error)
	GetItems(cartID string) ([]*CartItem, error)
	AddItem(item *CartItem) error
	UpdateItem(item *CartItem) error
	RemoveItem(itemID string) error
	UpdateTotal(cartID string, total float64) error
	SetStatus(cartID string, status CartStatus) error
}
