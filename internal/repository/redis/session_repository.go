package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Tenoywil/CaribEx-backend/internal/domain/auth"
	"github.com/redis/go-redis/v9"
)

// SessionRepository implements auth.SessionRepository using Redis
type SessionRepository struct {
	client *redis.Client
}

// NewSessionRepository creates a new Redis session repository
func NewSessionRepository(client *redis.Client) *SessionRepository {
	return &SessionRepository{client: client}
}

// SaveSession stores a session in Redis
func (r *SessionRepository) SaveSession(ctx context.Context, session *auth.Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	key := fmt.Sprintf("session:%s", session.ID)
	ttl := time.Until(session.ExpiresAt)
	
	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

// GetSession retrieves a session from Redis
func (r *SessionRepository) GetSession(ctx context.Context, sessionID string) (*auth.Session, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	
	data, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var session auth.Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	// Check if expired
	if session.IsExpired() {
		r.DeleteSession(ctx, sessionID)
		return nil, fmt.Errorf("session expired")
	}

	return &session, nil
}

// DeleteSession removes a session from Redis
func (r *SessionRepository) DeleteSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

// SaveNonce stores a nonce in Redis
func (r *SessionRepository) SaveNonce(ctx context.Context, nonce *auth.Nonce) error {
	data, err := json.Marshal(nonce)
	if err != nil {
		return fmt.Errorf("failed to marshal nonce: %w", err)
	}

	key := fmt.Sprintf("nonce:%s", nonce.Value)
	ttl := time.Until(nonce.ExpiresAt)
	
	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to save nonce: %w", err)
	}

	return nil
}

// GetNonce retrieves a nonce from Redis
func (r *SessionRepository) GetNonce(ctx context.Context, nonceValue string) (*auth.Nonce, error) {
	key := fmt.Sprintf("nonce:%s", nonceValue)
	
	data, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, fmt.Errorf("nonce not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	var nonce auth.Nonce
	if err := json.Unmarshal(data, &nonce); err != nil {
		return nil, fmt.Errorf("failed to unmarshal nonce: %w", err)
	}

	// Check if expired
	if nonce.IsExpired() {
		r.DeleteNonce(ctx, nonceValue)
		return nil, fmt.Errorf("nonce expired")
	}

	return &nonce, nil
}

// DeleteNonce removes a nonce from Redis
func (r *SessionRepository) DeleteNonce(ctx context.Context, nonceValue string) error {
	key := fmt.Sprintf("nonce:%s", nonceValue)
	
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete nonce: %w", err)
	}

	return nil
}
