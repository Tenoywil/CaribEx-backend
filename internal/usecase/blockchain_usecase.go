package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/Tenoywil/CaribEx-backend/internal/domain/wallet"
	"github.com/Tenoywil/CaribEx-backend/pkg/blockchain"
	"github.com/google/uuid"
)

// BlockchainUseCase handles blockchain transaction verification business logic
type BlockchainUseCase struct {
	walletRepo wallet.Repository
}

// NewBlockchainUseCase creates a new blockchain use case
func NewBlockchainUseCase(walletRepo wallet.Repository) *BlockchainUseCase {
	return &BlockchainUseCase{walletRepo: walletRepo}
}

// VerifyAndLogTransaction verifies an on-chain transaction and logs it to the database
func (uc *BlockchainUseCase) VerifyAndLogTransaction(userID, txHash string, chainID int64) (*wallet.Transaction, error) {
	// Validate chain ID
	if !blockchain.ValidateChainID(chainID) {
		return nil, errors.New("unsupported chain ID")
	}

	// Verify the transaction on-chain
	verification, err := blockchain.VerifyTransaction(txHash, chainID)
	if err != nil {
		return nil, err
	}

	// Check if transaction is still pending
	if verification.IsPending {
		return nil, errors.New("transaction is still pending")
	}

	// Check if transaction was successful
	if !verification.Verified {
		return nil, errors.New("transaction failed on-chain")
	}

	// Get user's wallet
	w, err := uc.walletRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Convert wei value to float (simplified - in production, use proper decimal handling)
	valueEth, err := blockchain.FormatValue(verification.Value)
	if err != nil {
		// If conversion fails, use 0 as amount
		valueEth = "0"
	}

	// Determine transaction type (credit if to address matches wallet, debit if from matches)
	txType := wallet.TransactionTypeCredit
	// Note: This is a simplified approach. In production, you'd need to map
	// Ethereum addresses to user wallet addresses properly

	// Create transaction log with value in reference for now
	tx := &wallet.Transaction{
		ID:        uuid.New().String(),
		WalletID:  w.ID,
		Type:      txType,
		Amount:    0, // Note: Blockchain value not directly stored in amount field
		Reference: fmt.Sprintf("Blockchain verification: %s (Value: %s ETH)", txHash, valueEth),
		Status:    wallet.TransactionStatusSuccess,
		CreatedAt: time.Now(),
		TxHash:    verification.TxHash,
		ChainID:   verification.ChainID,
		From:      verification.From,
		To:        verification.To,
	}

	// Store transaction in database
	err = uc.walletRepo.CreateTransaction(tx)
	if err != nil {
		return nil, err
	}

	// Note: Consider whether to update wallet balance based on verified transaction
	// This depends on your business logic (e.g., only for deposits, not all transactions)

	return tx, nil
}

// GetTransactionVerification retrieves verification details for a transaction hash
func (uc *BlockchainUseCase) GetTransactionVerification(txHash string, chainID int64) (*blockchain.TransactionVerification, error) {
	// Validate chain ID
	if !blockchain.ValidateChainID(chainID) {
		return nil, errors.New("unsupported chain ID")
	}

	// Verify the transaction on-chain
	verification, err := blockchain.VerifyTransaction(txHash, chainID)
	if err != nil {
		return nil, err
	}

	return verification, nil
}
