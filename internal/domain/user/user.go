package user

import "time"

// Role represents user roles in the system
type Role string

const (
	RoleCustomer Role = "customer"
	RoleSeller   Role = "seller"
	RoleAdmin    Role = "admin"
)

// User represents a user in the system
type User struct {
	ID            string    `json:"id"`
	Username      string    `json:"username"`
	WalletAddress string    `json:"wallet_address"`
	Role          Role      `json:"role"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Repository defines the interface for user data operations
type Repository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByWalletAddress(address string) (*User, error)
	Update(user *User) error
	Delete(id string) error
}
