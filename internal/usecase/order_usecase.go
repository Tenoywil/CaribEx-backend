package usecase

import (
	"time"

	"github.com/Tenoywil/CaribEx-backend/internal/domain/order"
	"github.com/google/uuid"
)

// OrderUseCase handles order business logic
type OrderUseCase struct {
	orderRepo order.Repository
}

// NewOrderUseCase creates a new order use case
func NewOrderUseCase(orderRepo order.Repository) *OrderUseCase {
	return &OrderUseCase{orderRepo: orderRepo}
}

// CreateOrder creates a new order
func (uc *OrderUseCase) CreateOrder(userID, cartID string, total float64, paymentRef string) (*order.Order, error) {
	o := &order.Order{
		ID:         uuid.New().String(),
		UserID:     userID,
		CartID:     cartID,
		Status:     order.OrderStatusPending,
		Total:      total,
		PaymentRef: paymentRef,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := uc.orderRepo.Create(o)
	if err != nil {
		return nil, err
	}

	return o, nil
}

// GetOrderByID retrieves an order by ID
func (uc *OrderUseCase) GetOrderByID(id string) (*order.Order, error) {
	return uc.orderRepo.GetByID(id)
}

// GetOrdersByUserID retrieves all orders for a user
func (uc *OrderUseCase) GetOrdersByUserID(userID string, page, pageSize int) ([]*order.Order, int, error) {
	return uc.orderRepo.GetByUserID(userID, page, pageSize)
}

// GetOrderItems retrieves all items in an order
func (uc *OrderUseCase) GetOrderItems(orderID string) ([]*order.OrderItem, error) {
	return uc.orderRepo.GetItems(orderID)
}

// UpdateOrderStatus updates the status of an order
func (uc *OrderUseCase) UpdateOrderStatus(orderID string, status order.OrderStatus) error {
	return uc.orderRepo.UpdateStatus(orderID, status)
}
