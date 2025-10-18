package postgres

import (
	"context"
	"fmt"

	"github.com/Tenoywil/CaribEx-backend/internal/domain/cart"
	"github.com/jackc/pgx/v5/pgxpool"
)

type cartRepository struct {
	db *pgxpool.Pool
}

// NewCartRepository creates a new cart repository
func NewCartRepository(db *pgxpool.Pool) cart.Repository {
	return &cartRepository{db: db}
}

func (r *cartRepository) GetByUserID(userID string) (*cart.Cart, error) {
	query := `
		SELECT id, user_id, status, total, created_at, updated_at
		FROM carts WHERE user_id = $1 AND status = 'active'
		ORDER BY created_at DESC LIMIT 1
	`
	var c cart.Cart
	err := r.db.QueryRow(context.Background(), query, userID).Scan(
		&c.ID, &c.UserID, &c.Status, &c.Total, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart by user id: %w", err)
	}
	return &c, nil
}

func (r *cartRepository) GetItems(cartID string) ([]*cart.CartItem, error) {
	query := `
		SELECT id, cart_id, product_id, quantity, price, created_at, updated_at
		FROM cart_items WHERE cart_id = $1
	`
	rows, err := r.db.Query(context.Background(), query, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to query cart items: %w", err)
	}
	defer rows.Close()

	var items []*cart.CartItem
	for rows.Next() {
		var item cart.CartItem
		err := rows.Scan(&item.ID, &item.CartID, &item.ProductID, &item.Quantity, &item.Price, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cart item: %w", err)
		}
		items = append(items, &item)
	}

	return items, nil
}

func (r *cartRepository) AddItem(item *cart.CartItem) error {
	query := `
		INSERT INTO cart_items (id, cart_id, product_id, quantity, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (cart_id, product_id) 
		DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity, updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.Exec(context.Background(), query,
		item.ID, item.CartID, item.ProductID, item.Quantity, item.Price, item.CreatedAt, item.UpdatedAt)
	return err
}

func (r *cartRepository) UpdateItem(item *cart.CartItem) error {
	query := `
		UPDATE cart_items 
		SET quantity = $1, price = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.Exec(context.Background(), query,
		item.Quantity, item.Price, item.UpdatedAt, item.ID)
	return err
}

func (r *cartRepository) RemoveItem(itemID string) error {
	query := `DELETE FROM cart_items WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, itemID)
	return err
}

func (r *cartRepository) UpdateTotal(cartID string, total float64) error {
	query := `
		UPDATE carts 
		SET total = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.Exec(context.Background(), query, total, cartID)
	return err
}

func (r *cartRepository) SetStatus(cartID string, status cart.CartStatus) error {
	query := `
		UPDATE carts 
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.Exec(context.Background(), query, status, cartID)
	return err
}
