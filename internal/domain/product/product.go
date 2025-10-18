package product

import "time"

// Product represents a marketplace product listing
type Product struct {
	ID          string    `json:"id"`
	SellerID    string    `json:"seller_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	Images      []string  `json:"images"`
	CategoryID  string    `json:"category_id"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Category represents a product category
type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Repository defines the interface for product data operations
type Repository interface {
	Create(product *Product) error
	GetByID(id string) (*Product, error)
	List(filters map[string]interface{}, page, pageSize int) ([]*Product, int, error)
	Update(product *Product) error
	Delete(id string) error
	GetCategories() ([]*Category, error)
}
