package postgres

import (
	"context"
	"fmt"

	"github.com/Tenoywil/CaribEx-backend/internal/domain/wallet"
	"github.com/jackc/pgx/v5/pgxpool"
)

type walletRepository struct {
	db *pgxpool.Pool
}

// NewWalletRepository creates a new wallet repository
func NewWalletRepository(db *pgxpool.Pool) wallet.Repository {
	return &walletRepository{db: db}
}

func (r *walletRepository) GetByUserID(userID string) (*wallet.Wallet, error) {
	query := `
		SELECT id, user_id, balance, currency, updated_at
		FROM wallets WHERE user_id = $1
	`
	var w wallet.Wallet
	err := r.db.QueryRow(context.Background(), query, userID).Scan(
		&w.ID, &w.UserID, &w.Balance, &w.Currency, &w.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet by user id: %w", err)
	}
	return &w, nil
}

func (r *walletRepository) CreateTransaction(tx *wallet.Transaction) error {
	query := `
		INSERT INTO transactions (id, wallet_id, type, amount, reference, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(context.Background(), query,
		tx.ID, tx.WalletID, tx.Type, tx.Amount, tx.Reference, tx.Status, tx.CreatedAt)
	return err
}

func (r *walletRepository) GetTransactions(walletID string, page, pageSize int) ([]*wallet.Transaction, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM transactions WHERE wallet_id = $1`
	err := r.db.QueryRow(context.Background(), countQuery, walletID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	// Get transactions
	query := `
		SELECT id, wallet_id, type, amount, reference, status, created_at
		FROM transactions
		WHERE wallet_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(context.Background(), query, walletID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*wallet.Transaction
	for rows.Next() {
		var tx wallet.Transaction
		err := rows.Scan(&tx.ID, &tx.WalletID, &tx.Type, &tx.Amount, &tx.Reference, &tx.Status, &tx.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, &tx)
	}

	return transactions, total, nil
}

func (r *walletRepository) UpdateBalance(walletID string, amount float64) error {
	query := `
		UPDATE wallets 
		SET balance = balance + $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.Exec(context.Background(), query, amount, walletID)
	return err
}
