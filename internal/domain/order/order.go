package order

import "time"

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order represents a customer order
type Order struct {
	ID         string      `json:"id"`
	UserID     string      `json:"user_id"`
	CartID     string      `json:"cart_id"`
	Status     OrderStatus `json:"status"`
	Total      float64     `json:"total"`
	PaymentRef string      `json:"payment_ref"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID        string  `json:"id"`
	OrderID   string  `json:"order_id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// Repository defines the interface for order data operations
type Repository interface {
	Create(order *Order) error
	GetByID(id string) (*Order, error)
	GetByUserID(userID string, page, pageSize int) ([]*Order, int, error)
	GetItems(orderID string) ([]*OrderItem, error)
	UpdateStatus(orderID string, status OrderStatus) error
}
