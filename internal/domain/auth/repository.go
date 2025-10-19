package auth

import "context"

// SessionRepository defines the interface for session storage
type SessionRepository interface {
	// SaveSession stores a session
	SaveSession(ctx context.Context, session *Session) error
	
	// GetSession retrieves a session by ID
	GetSession(ctx context.Context, sessionID string) (*Session, error)
	
	// DeleteSession removes a session
	DeleteSession(ctx context.Context, sessionID string) error
	
	// SaveNonce stores a nonce
	SaveNonce(ctx context.Context, nonce *Nonce) error
	
	// GetNonce retrieves a nonce by value
	GetNonce(ctx context.Context, nonceValue string) (*Nonce, error)
	
	// DeleteNonce removes a nonce
	DeleteNonce(ctx context.Context, nonceValue string) error
}
