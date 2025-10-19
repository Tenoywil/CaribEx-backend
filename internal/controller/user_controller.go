package controller

import (
	"net/http"

	"github.com/Tenoywil/CaribEx-backend/internal/domain/user"
	"github.com/Tenoywil/CaribEx-backend/internal/usecase"
	"github.com/gin-gonic/gin"
)

// UserController handles HTTP requests for users
type UserController struct {
	userUseCase *usecase.UserUseCase
}

// NewUserController creates a new user controller
func NewUserController(userUseCase *usecase.UserUseCase) *UserController {
	return &UserController{userUseCase: userUseCase}
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Username      string    `json:"username" binding:"required"`
	WalletAddress string    `json:"wallet_address" binding:"required"`
	Role          user.Role `json:"role" binding:"required"`
}

// CreateUser handles POST /users
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := c.userUseCase.CreateUser(req.Username, req.WalletAddress, req.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, u)
}

// GetUser handles GET /users/:id
func (c *UserController) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")

	u, err := c.userUseCase.GetUserByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, u)
}

// GetUserByWallet handles GET /users/wallet/:address
func (c *UserController) GetUserByWallet(ctx *gin.Context) {
	address := ctx.Param("address")

	u, err := c.userUseCase.GetUserByWalletAddress(address)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, u)
}

// UpdateUser handles PUT /users/:id
func (c *UserController) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")

	var u user.User
	if err := ctx.ShouldBindJSON(&u); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u.ID = id
	err := c.userUseCase.UpdateUser(&u)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, u)
}

// DeleteUser handles DELETE /users/:id
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.userUseCase.DeleteUser(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
