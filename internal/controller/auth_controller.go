package controller

import (
	"net/http"

	"github.com/Tenoywil/CaribEx-backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// AuthController handles authentication HTTP requests
type AuthController struct {
	authUseCase *usecase.AuthUseCase
}

// NewAuthController creates a new auth controller
func NewAuthController(authUseCase *usecase.AuthUseCase) *AuthController {
	return &AuthController{authUseCase: authUseCase}
}

// NonceResponse represents the nonce response
type NonceResponse struct {
	Nonce     string `json:"nonce"`
	ExpiresAt string `json:"expires_at"`
}

// GetNonce handles GET /auth/nonce
func (c *AuthController) GetNonce(ctx *gin.Context) {
	nonce, err := c.authUseCase.GenerateNonce(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate nonce"})
		return
	}

	ctx.JSON(http.StatusOK, NonceResponse{
		Nonce:     nonce.Value,
		ExpiresAt: nonce.ExpiresAt.Format(http.TimeFormat),
	})
}

// SIWERequest represents the SIWE authentication request
type SIWERequest struct {
	Message   string `json:"message" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

// SIWEResponse represents the SIWE authentication response
type SIWEResponse struct {
	User struct {
		ID            string `json:"id"`
		Username      string `json:"username"`
		WalletAddress string `json:"wallet_address"`
		Role          string `json:"role"`
	} `json:"user"`
	SessionID string `json:"session_id"`
}

// AuthenticateSIWE handles POST /auth/siwe
func (c *AuthController) AuthenticateSIWE(ctx *gin.Context) {
	var req SIWERequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, user, err := c.authUseCase.VerifySIWE(
		ctx.Request.Context(),
		req.Message,
		req.Signature,
	)
	if err != nil {
		log.Error().Err(err).Msg("SIWE verification failed")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authentication failed: " + err.Error()})
		return
	}

	// Set session cookie
	ctx.SetCookie(
		"session_id",
		session.ID,
		int(session.ExpiresAt.Sub(session.CreatedAt).Seconds()),
		"/",
		"",
		false, // Set to true in production with HTTPS
		true,  // HTTP only
	)

	// Return response
	response := SIWEResponse{
		SessionID: session.ID,
	}
	response.User.ID = user.ID
	response.User.Username = user.Username
	response.User.WalletAddress = user.WalletAddress
	response.User.Role = string(user.Role)

	ctx.JSON(http.StatusOK, response)
}

// GetMe handles GET /auth/me
func (c *AuthController) GetMe(ctx *gin.Context) {
	// Get user from context (set by auth middleware)
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	walletAddress, _ := ctx.Get("wallet_address")

	ctx.JSON(http.StatusOK, gin.H{
		"user_id":        userID,
		"wallet_address": walletAddress,
	})
}

// Logout handles POST /auth/logout
func (c *AuthController) Logout(ctx *gin.Context) {
	sessionID, err := ctx.Cookie("session_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no session found"})
		return
	}

	if err := c.authUseCase.Logout(ctx.Request.Context(), sessionID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
		return
	}

	// Clear cookie
	ctx.SetCookie(
		"session_id",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
