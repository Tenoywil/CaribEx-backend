package usecase

import (
	"errors"
	"time"

	"github.com/Tenoywil/CaribEx-backend/internal/domain/wallet"
	"github.com/google/uuid"
)

// WalletUseCase handles wallet business logic
type WalletUseCase struct {
	walletRepo wallet.Repository
}

// NewWalletUseCase creates a new wallet use case
func NewWalletUseCase(walletRepo wallet.Repository) *WalletUseCase {
	return &WalletUseCase{walletRepo: walletRepo}
}

// GetWalletByUserID retrieves a wallet by user ID
func (uc *WalletUseCase) GetWalletByUserID(userID string) (*wallet.Wallet, error) {
	return uc.walletRepo.GetByUserID(userID)
}

// SendFunds sends funds from a wallet
func (uc *WalletUseCase) SendFunds(walletID string, amount float64, reference string) (*wallet.Transaction, error) {
	// Get wallet to check balance
	w, err := uc.walletRepo.GetByUserID(walletID)
	if err != nil {
		return nil, err
	}

	if w.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	// Create debit transaction
	tx := &wallet.Transaction{
		ID:        uuid.New().String(),
		WalletID:  w.ID,
		Type:      wallet.TransactionTypeDebit,
		Amount:    amount,
		Reference: reference,
		Status:    wallet.TransactionStatusPending,
		CreatedAt: time.Now(),
	}

	// Create transaction record
	err = uc.walletRepo.CreateTransaction(tx)
	if err != nil {
		return nil, err
	}

	// Update wallet balance
	err = uc.walletRepo.UpdateBalance(w.ID, -amount)
	if err != nil {
		return nil, err
	}

	// Update transaction status
	tx.Status = wallet.TransactionStatusSuccess

	return tx, nil
}

// ReceiveFunds receives funds to a wallet
func (uc *WalletUseCase) ReceiveFunds(walletID string, amount float64, reference string) (*wallet.Transaction, error) {
	w, err := uc.walletRepo.GetByUserID(walletID)
	if err != nil {
		return nil, err
	}

	// Create credit transaction
	tx := &wallet.Transaction{
		ID:        uuid.New().String(),
		WalletID:  w.ID,
		Type:      wallet.TransactionTypeCredit,
		Amount:    amount,
		Reference: reference,
		Status:    wallet.TransactionStatusPending,
		CreatedAt: time.Now(),
	}

	// Create transaction record
	err = uc.walletRepo.CreateTransaction(tx)
	if err != nil {
		return nil, err
	}

	// Update wallet balance
	err = uc.walletRepo.UpdateBalance(w.ID, amount)
	if err != nil {
		return nil, err
	}

	// Update transaction status
	tx.Status = wallet.TransactionStatusSuccess

	return tx, nil
}

// GetTransactions retrieves transaction history
func (uc *WalletUseCase) GetTransactions(walletID string, page, pageSize int) ([]*wallet.Transaction, int, error) {
	return uc.walletRepo.GetTransactions(walletID, page, pageSize)
}
