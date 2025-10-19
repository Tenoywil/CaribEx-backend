package usecase

import (
	"time"

	"github.com/Tenoywil/CaribEx-backend/internal/domain/user"
	"github.com/google/uuid"
)

// UserUseCase handles user business logic
type UserUseCase struct {
	userRepo user.Repository
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(userRepo user.Repository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

// CreateUser creates a new user
func (uc *UserUseCase) CreateUser(username, walletAddress string, role user.Role) (*user.User, error) {
	u := &user.User{
		ID:            uuid.New().String(),
		Username:      username,
		WalletAddress: walletAddress,
		Role:          role,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := uc.userRepo.Create(u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// GetUserByID retrieves a user by ID
func (uc *UserUseCase) GetUserByID(id string) (*user.User, error) {
	return uc.userRepo.GetByID(id)
}

// GetUserByWalletAddress retrieves a user by wallet address
func (uc *UserUseCase) GetUserByWalletAddress(address string) (*user.User, error) {
	return uc.userRepo.GetByWalletAddress(address)
}

// UpdateUser updates user information
func (uc *UserUseCase) UpdateUser(u *user.User) error {
	u.UpdatedAt = time.Now()
	return uc.userRepo.Update(u)
}

// DeleteUser deletes a user
func (uc *UserUseCase) DeleteUser(id string) error {
	return uc.userRepo.Delete(id)
}
