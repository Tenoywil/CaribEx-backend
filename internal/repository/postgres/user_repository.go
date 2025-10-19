package postgres

import (
	"context"
	"fmt"

	"github.com/Tenoywil/CaribEx-backend/internal/domain/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *pgxpool.Pool) user.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(u *user.User) error {
	query := `
		INSERT INTO users (id, username, wallet_address, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(context.Background(), query,
		u.ID, u.Username, u.WalletAddress, u.Role, u.CreatedAt, u.UpdatedAt)
	return err
}

func (r *userRepository) GetByID(id string) (*user.User, error) {
	query := `
		SELECT id, username, wallet_address, role, created_at, updated_at
		FROM users WHERE id = $1
	`
	var u user.User
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&u.ID, &u.Username, &u.WalletAddress, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &u, nil
}

func (r *userRepository) GetByWalletAddress(address string) (*user.User, error) {
	query := `
		SELECT id, username, wallet_address, role, created_at, updated_at
		FROM users WHERE wallet_address = $1
	`
	var u user.User
	err := r.db.QueryRow(context.Background(), query, address).Scan(
		&u.ID, &u.Username, &u.WalletAddress, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by wallet address: %w", err)
	}
	return &u, nil
}

func (r *userRepository) Update(u *user.User) error {
	query := `
		UPDATE users 
		SET username = $1, wallet_address = $2, role = $3, updated_at = $4
		WHERE id = $5
	`
	_, err := r.db.Exec(context.Background(), query,
		u.Username, u.WalletAddress, u.Role, u.UpdatedAt, u.ID)
	return err
}

func (r *userRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}
