package controller

import (
	"net/http"
	"strconv"

	"github.com/Tenoywil/CaribEx-backend/internal/usecase"
	"github.com/gin-gonic/gin"
)

// WalletController handles HTTP requests for wallets
type WalletController struct {
	walletUseCase *usecase.WalletUseCase
}

// NewWalletController creates a new wallet controller
func NewWalletController(walletUseCase *usecase.WalletUseCase) *WalletController {
	return &WalletController{walletUseCase: walletUseCase}
}

// GetWallet handles GET /wallet
func (c *WalletController) GetWallet(ctx *gin.Context) {
	// TODO: Get user ID from authenticated user context
	userID := ctx.GetString("user_id")

	w, err := c.walletUseCase.GetWalletByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
		return
	}

	ctx.JSON(http.StatusOK, w)
}

// SendFundsRequest represents the request body for sending funds
type SendFundsRequest struct {
	Amount    float64 `json:"amount" binding:"required"`
	Reference string  `json:"reference"`
}

// SendFunds handles POST /wallet/send
func (c *WalletController) SendFunds(ctx *gin.Context) {
	var req SendFundsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Get user ID from authenticated user context
	userID := ctx.GetString("user_id")

	tx, err := c.walletUseCase.SendFunds(userID, req.Amount, req.Reference)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tx)
}

// ReceiveFundsRequest represents the request body for receiving funds
type ReceiveFundsRequest struct {
	Amount    float64 `json:"amount" binding:"required"`
	Reference string  `json:"reference"`
}

// ReceiveFunds handles POST /wallet/receive
func (c *WalletController) ReceiveFunds(ctx *gin.Context) {
	var req ReceiveFundsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Get user ID from authenticated user context
	userID := ctx.GetString("user_id")

	tx, err := c.walletUseCase.ReceiveFunds(userID, req.Amount, req.Reference)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tx)
}

// GetTransactions handles GET /wallet/transactions
func (c *WalletController) GetTransactions(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	
	// Ensure page and pageSize are within valid ranges
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// TODO: Get wallet ID from authenticated user context
	walletID := ctx.GetString("wallet_id")

	transactions, total, err := c.walletUseCase.GetTransactions(walletID, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"total":        total,
		"page":         page,
		"page_size":    pageSize,
		"total_pages":  (total + pageSize - 1) / pageSize,
	})
}
