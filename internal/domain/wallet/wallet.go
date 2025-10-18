package wallet

import "time"

// Currency represents supported currencies
type Currency string

const (
	CurrencyJAM  Currency = "JAM"
	CurrencyUSD  Currency = "USD"
	CurrencyUSDC Currency = "USDC"
)

// Wallet represents a user's wallet
type Wallet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Balance   float64   `json:"balance"`
	Currency  Currency  `json:"currency"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeCredit TransactionType = "credit"
	TransactionTypeDebit  TransactionType = "debit"
)

// TransactionStatus represents transaction status
type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "pending"
	TransactionStatusSuccess TransactionStatus = "success"
	TransactionStatusFailed  TransactionStatus = "failed"
)

// Transaction represents a wallet transaction
type Transaction struct {
	ID        string            `json:"id"`
	WalletID  string            `json:"wallet_id"`
	Type      TransactionType   `json:"type"`
	Amount    float64           `json:"amount"`
	Reference string            `json:"reference"`
	Status    TransactionStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
}

// Repository defines the interface for wallet data operations
type Repository interface {
	GetByUserID(userID string) (*Wallet, error)
	CreateTransaction(tx *Transaction) error
	GetTransactions(walletID string, page, pageSize int) ([]*Transaction, int, error)
	UpdateBalance(walletID string, amount float64) error
}
