package postgres

import (
	"context"
	"fmt"

	"github.com/Tenoywil/CaribEx-backend/internal/domain/order"
	"github.com/jackc/pgx/v5/pgxpool"
)

type orderRepository struct {
	db *pgxpool.Pool
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *pgxpool.Pool) order.Repository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(o *order.Order) error {
	query := `
		INSERT INTO orders (id, user_id, cart_id, status, total, payment_ref, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(context.Background(), query,
		o.ID, o.UserID, o.CartID, o.Status, o.Total, o.PaymentRef, o.CreatedAt, o.UpdatedAt)
	return err
}

func (r *orderRepository) GetByID(id string) (*order.Order, error) {
	query := `
		SELECT id, user_id, cart_id, status, total, payment_ref, created_at, updated_at
		FROM orders WHERE id = $1
	`
	var o order.Order
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&o.ID, &o.UserID, &o.CartID, &o.Status, &o.Total, &o.PaymentRef, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by id: %w", err)
	}
	return &o, nil
}

func (r *orderRepository) GetByUserID(userID string, page, pageSize int) ([]*order.Order, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM orders WHERE user_id = $1`
	err := r.db.QueryRow(context.Background(), countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	// Get orders
	query := `
		SELECT id, user_id, cart_id, status, total, payment_ref, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(context.Background(), query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	var orders []*order.Order
	for rows.Next() {
		var o order.Order
		err := rows.Scan(&o.ID, &o.UserID, &o.CartID, &o.Status, &o.Total, &o.PaymentRef, &o.CreatedAt, &o.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, &o)
	}

	return orders, total, nil
}

func (r *orderRepository) GetItems(orderID string) ([]*order.OrderItem, error) {
	query := `
		SELECT id, order_id, product_id, quantity, price
		FROM order_items WHERE order_id = $1
	`
	rows, err := r.db.Query(context.Background(), query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query order items: %w", err)
	}
	defer rows.Close()

	var items []*order.OrderItem
	for rows.Next() {
		var item order.OrderItem
		err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		items = append(items, &item)
	}

	return items, nil
}

func (r *orderRepository) UpdateStatus(orderID string, status order.OrderStatus) error {
	query := `
		UPDATE orders 
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.Exec(context.Background(), query, status, orderID)
	return err
}
