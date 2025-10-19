package controller

import (
	"net/http"
	"strconv"

	"github.com/Tenoywil/CaribEx-backend/internal/usecase"
	"github.com/gin-gonic/gin"
)

// BlockchainController handles HTTP requests for blockchain operations
type BlockchainController struct {
	blockchainUseCase *usecase.BlockchainUseCase
}

// NewBlockchainController creates a new blockchain controller
func NewBlockchainController(blockchainUseCase *usecase.BlockchainUseCase) *BlockchainController {
	return &BlockchainController{blockchainUseCase: blockchainUseCase}
}

// VerifyTransactionRequest represents the request body for transaction verification
type VerifyTransactionRequest struct {
	TxHash  string `json:"txHash" binding:"required"`
	ChainID int64  `json:"chainId" binding:"required"`
}

// VerifyTransactionResponse represents the response for transaction verification
type VerifyTransactionResponse struct {
	Status    string `json:"status"`
	TxHash    string `json:"txHash"`
	Message   string `json:"message"`
	From      string `json:"from,omitempty"`
	To        string `json:"to,omitempty"`
	Value     string `json:"value,omitempty"`
	ChainID   int64  `json:"chainId,omitempty"`
	IsPending bool   `json:"isPending,omitempty"`
}

// VerifyTransaction handles POST /v1/wallet/verify-transaction
func (c *BlockchainController) VerifyTransaction(ctx *gin.Context) {
	var req VerifyTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Get user ID from authenticated user context
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Verify and log the transaction
	tx, err := c.blockchainUseCase.VerifyAndLogTransaction(userID, req.TxHash, req.ChainID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"txHash":  req.TxHash,
			"error":   err.Error(),
			"message": "Transaction verification failed",
		})
		return
	}

	response := VerifyTransactionResponse{
		Status:  "verified",
		TxHash:  req.TxHash,
		Message: "Transaction successfully verified",
		From:    tx.From,
		To:      tx.To,
		ChainID: tx.ChainID,
	}

	ctx.JSON(http.StatusOK, response)
}

// GetTransactionStatus handles GET /v1/wallet/transaction-status
func (c *BlockchainController) GetTransactionStatus(ctx *gin.Context) {
	txHash := ctx.Query("txHash")
	if txHash == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "txHash query parameter is required"})
		return
	}

	chainID := int64(1) // Default to Ethereum mainnet
	if chainIDStr := ctx.Query("chainId"); chainIDStr != "" {
		if parsed, err := strconv.ParseInt(chainIDStr, 10, 64); err == nil {
			chainID = parsed
		}
	}

	// Get verification details without logging
	verification, err := c.blockchainUseCase.GetTransactionVerification(txHash, chainID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"txHash": txHash,
			"error":  err.Error(),
		})
		return
	}

	response := VerifyTransactionResponse{
		Status:    "success",
		TxHash:    verification.TxHash,
		Message:   "Transaction status retrieved",
		From:      verification.From,
		To:        verification.To,
		Value:     verification.Value,
		ChainID:   verification.ChainID,
		IsPending: verification.IsPending,
	}

	if verification.IsPending {
		response.Message = "Transaction is pending"
	} else if !verification.Verified {
		response.Status = "failed"
		response.Message = "Transaction failed on-chain"
	}

	ctx.JSON(http.StatusOK, response)
}
