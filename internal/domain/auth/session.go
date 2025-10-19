package auth

import (
	"time"

	"github.com/google/uuid"
)

// Session represents an authenticated user session
type Session struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	WalletAddress string    `json:"wallet_address"`
	Nonce         string    `json:"nonce"`
	ExpiresAt     time.Time `json:"expires_at"`
	CreatedAt     time.Time `json:"created_at"`
}

// NewSession creates a new session
func NewSession(userID, walletAddress string, duration time.Duration) *Session {
	now := time.Now()
	return &Session{
		ID:            uuid.New().String(),
		UserID:        userID,
		WalletAddress: walletAddress,
		ExpiresAt:     now.Add(duration),
		CreatedAt:     now,
	}
}

// IsExpired checks if the session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// Nonce represents a SIWE nonce
type Nonce struct {
	Value     string    `json:"nonce"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// NewNonce creates a new nonce with 10 minute expiration
func NewNonce() *Nonce {
	now := time.Now()
	return &Nonce{
		Value:     uuid.New().String(),
		ExpiresAt: now.Add(10 * time.Minute),
		CreatedAt: now,
	}
}

// IsExpired checks if the nonce has expired
func (n *Nonce) IsExpired() bool {
	return time.Now().After(n.ExpiresAt)
}
