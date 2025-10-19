package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Tenoywil/CaribEx-backend/internal/domain/auth"
	"github.com/Tenoywil/CaribEx-backend/internal/domain/user"
	"github.com/rs/zerolog/log"
	"github.com/spruceid/siwe-go"
)

// AuthUseCase handles authentication business logic
type AuthUseCase struct {
	sessionRepo auth.SessionRepository
	userUseCase *UserUseCase
	domain      string
}

// NewAuthUseCase creates a new auth use case
func NewAuthUseCase(
	sessionRepo auth.SessionRepository,
	userUseCase *UserUseCase,
	domain string,
) *AuthUseCase {
	return &AuthUseCase{
		sessionRepo: sessionRepo,
		userUseCase: userUseCase,
		domain:      domain,
	}
}

// GenerateNonce creates a new nonce for SIWE authentication
func (uc *AuthUseCase) GenerateNonce(ctx context.Context) (*auth.Nonce, error) {
	nonce := auth.NewNonce()
	
	if err := uc.sessionRepo.SaveNonce(ctx, nonce); err != nil {
		log.Error().Err(err).Msg("failed to save nonce")
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	log.Info().Str("nonce", nonce.Value).Msg("nonce generated")
	return nonce, nil
}

// VerifySIWE verifies a SIWE message and signature
func (uc *AuthUseCase) VerifySIWE(
	ctx context.Context,
	message, signature string,
) (*auth.Session, *user.User, error) {
	// Log the received message for debugging
	log.Debug().Str("message", message).Str("signature", signature).Msg("received SIWE authentication request")
	
	// Parse the SIWE message
	siweMessage, err := siwe.ParseMessage(message)
	if err != nil {
		log.Error().Err(err).Str("message", message).Msg("failed to parse SIWE message")
		return nil, nil, fmt.Errorf("invalid SIWE message: %w", err)
	}

	// Verify the domain matches
	if siweMessage.GetDomain() != uc.domain {
		return nil, nil, fmt.Errorf("domain mismatch: expected %s, got %s", uc.domain, siweMessage.GetDomain())
	}

	// Verify the nonce exists and is valid
	nonce, err := uc.sessionRepo.GetNonce(ctx, siweMessage.GetNonce())
	if err != nil {
		log.Error().Err(err).Str("nonce", siweMessage.GetNonce()).Msg("nonce not found or expired")
		return nil, nil, fmt.Errorf("invalid or expired nonce")
	}

	// Verify the signature
	_, err = siweMessage.VerifyEIP191(signature)
	if err != nil {
		log.Error().Err(err).Msg("signature verification failed")
		return nil, nil, fmt.Errorf("invalid signature: %w", err)
	}

	// Get the wallet address from the message
	// Note: VerifyEIP191 already validates that the signature matches the address
	walletAddress := strings.ToLower(siweMessage.GetAddress().Hex())

	// Delete the used nonce
	if err := uc.sessionRepo.DeleteNonce(ctx, nonce.Value); err != nil {
		log.Warn().Err(err).Msg("failed to delete used nonce")
	}

	// Get or create user
	u, err := uc.userUseCase.GetUserByWalletAddress(walletAddress)
	if err != nil {
		// User doesn't exist, create new user with customer role
		log.Info().Str("wallet", walletAddress).Msg("creating new user")
		u, err = uc.userUseCase.CreateUser(
			fmt.Sprintf("user_%s", walletAddress[:8]),
			walletAddress,
			user.RoleCustomer,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to create user")
			return nil, nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Create session
	session := auth.NewSession(u.ID, walletAddress, 24*time.Hour)
	if err := uc.sessionRepo.SaveSession(ctx, session); err != nil {
		log.Error().Err(err).Msg("failed to save session")
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	log.Info().
		Str("user_id", u.ID).
		Str("wallet", walletAddress).
		Str("session_id", session.ID).
		Msg("user authenticated via SIWE")

	return session, u, nil
}

// ValidateSession checks if a session is valid
func (uc *AuthUseCase) ValidateSession(ctx context.Context, sessionID string) (*auth.Session, error) {
	session, err := uc.sessionRepo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session: %w", err)
	}

	if session.IsExpired() {
		uc.sessionRepo.DeleteSession(ctx, sessionID)
		return nil, fmt.Errorf("session expired")
	}

	return session, nil
}

// Logout invalidates a session
func (uc *AuthUseCase) Logout(ctx context.Context, sessionID string) error {
	if err := uc.sessionRepo.DeleteSession(ctx, sessionID); err != nil {
		log.Error().Err(err).Str("session_id", sessionID).Msg("failed to delete session")
		return fmt.Errorf("failed to logout: %w", err)
	}

	log.Info().Str("session_id", sessionID).Msg("user logged out")
	return nil
}
