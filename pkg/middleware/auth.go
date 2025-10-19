package middleware

import (
	"net/http"

	"github.com/Tenoywil/CaribEx-backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// AuthMiddleware creates a middleware that validates session authentication
func AuthMiddleware(authUseCase *usecase.AuthUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get session ID from cookie
		sessionID, err := ctx.Cookie("session_id")
		if err != nil {
			log.Debug().Msg("no session cookie found")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			ctx.Abort()
			return
		}

		// Validate session
		session, err := authUseCase.ValidateSession(ctx.Request.Context(), sessionID)
		if err != nil {
			log.Debug().Err(err).Str("session_id", sessionID).Msg("invalid session")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired session"})
			ctx.Abort()
			return
		}

		// Set user info in context
		ctx.Set("user_id", session.UserID)
		ctx.Set("wallet_address", session.WalletAddress)
		ctx.Set("session_id", session.ID)

		log.Debug().
			Str("user_id", session.UserID).
			Str("wallet", session.WalletAddress).
			Msg("request authenticated")

		ctx.Next()
	}
}

// OptionalAuthMiddleware creates a middleware that optionally validates authentication
// If authenticated, it sets user context; otherwise, it continues without error
func OptionalAuthMiddleware(authUseCase *usecase.AuthUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get session ID from cookie
		sessionID, err := ctx.Cookie("session_id")
		if err != nil {
			// No session, continue without auth
			ctx.Next()
			return
		}

		// Validate session
		session, err := authUseCase.ValidateSession(ctx.Request.Context(), sessionID)
		if err != nil {
			// Invalid session, continue without auth
			ctx.Next()
			return
		}

		// Set user info in context
		ctx.Set("user_id", session.UserID)
		ctx.Set("wallet_address", session.WalletAddress)
		ctx.Set("session_id", session.ID)

		ctx.Next()
	}
}
