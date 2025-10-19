package postgres

import (
	"context"
	"fmt"

	"github.com/Tenoywil/CaribEx-backend/internal/domain/product"
	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepository struct {
	db *pgxpool.Pool
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *pgxpool.Pool) product.Repository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(p *product.Product) error {
	query := `
		INSERT INTO products (id, seller_id, title, description, price, quantity, images, category_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Exec(context.Background(), query,
		p.ID, p.SellerID, p.Title, p.Description, p.Price, p.Quantity, p.Images, p.CategoryID, p.IsActive, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *productRepository) GetByID(id string) (*product.Product, error) {
	query := `
		SELECT id, seller_id, title, description, price, quantity, images, category_id, is_active, created_at, updated_at
		FROM products WHERE id = $1
	`
	var p product.Product
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&p.ID, &p.SellerID, &p.Title, &p.Description, &p.Price, &p.Quantity, &p.Images, &p.CategoryID, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}
	return &p, nil
}

func (r *productRepository) GetByIDWithCategory(id string) (*product.ProductWithCategory, error) {
	query := `
		SELECT p.id, p.seller_id, p.title, p.description, p.price, p.quantity, p.images, 
		       p.category_id, p.is_active, p.created_at, p.updated_at,
		       c.id, c.name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
	`
	var p product.ProductWithCategory
	var categoryID, categoryName *string
	
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&p.ID, &p.SellerID, &p.Title, &p.Description, &p.Price, &p.Quantity, &p.Images, 
		&p.CategoryID, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
		&categoryID, &categoryName)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}
	
	// Set category if it exists
	if categoryID != nil && categoryName != nil {
		p.Category = &product.Category{
			ID:   *categoryID,
			Name: *categoryName,
		}
	}
	
	return &p, nil
}

func (r *productRepository) List(filters map[string]interface{}, page, pageSize int) ([]*product.Product, int, error) {
	offset := (page - 1) * pageSize

	// Build query with filters
	whereClause := "WHERE is_active = true"
	args := []interface{}{}
	argCount := 1

	if categoryID, ok := filters["category_id"]; ok {
		whereClause += fmt.Sprintf(" AND category_id = $%d", argCount)
		args = append(args, categoryID)
		argCount++
	}

	if search, ok := filters["search"]; ok {
		whereClause += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", argCount, argCount)
		searchPattern := fmt.Sprintf("%%%s%%", search)
		args = append(args, searchPattern)
		argCount++
	}

	// Get total count
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products %s", whereClause)
	err := r.db.QueryRow(context.Background(), countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Get products
	query := fmt.Sprintf(`
		SELECT id, seller_id, title, description, price, quantity, images, category_id, is_active, created_at, updated_at
		FROM products
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argCount, argCount+1)
	args = append(args, pageSize, offset)

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	var products []*product.Product
	for rows.Next() {
		var p product.Product
		err := rows.Scan(&p.ID, &p.SellerID, &p.Title, &p.Description, &p.Price, &p.Quantity, &p.Images, &p.CategoryID, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, &p)
	}

	return products, total, nil
}

func (r *productRepository) ListWithCategory(filters map[string]interface{}, page, pageSize int, sortBy, sortOrder string) ([]*product.ProductWithCategory, int, error) {
	offset := (page - 1) * pageSize

	// Build query with filters
	whereClause := "WHERE p.is_active = true"
	args := []interface{}{}
	argCount := 1

	if categoryID, ok := filters["category_id"]; ok {
		whereClause += fmt.Sprintf(" AND p.category_id = $%d", argCount)
		args = append(args, categoryID)
		argCount++
	}

	if search, ok := filters["search"]; ok {
		whereClause += fmt.Sprintf(" AND (p.title ILIKE $%d OR p.description ILIKE $%d)", argCount, argCount)
		searchPattern := fmt.Sprintf("%%%s%%", search)
		args = append(args, searchPattern)
		argCount++
	}

	// Get total count
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products p %s", whereClause)
	err := r.db.QueryRow(context.Background(), countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Build ORDER BY clause
	orderByClause := "ORDER BY p.created_at DESC"
	validSortFields := map[string]string{
		"created_at": "p.created_at",
		"updated_at": "p.updated_at",
		"price":      "p.price",
		"title":      "p.title",
	}
	
	if sortBy != "" {
		if dbField, ok := validSortFields[sortBy]; ok {
			order := "DESC"
			if sortOrder == "asc" || sortOrder == "ASC" {
				order = "ASC"
			}
			orderByClause = fmt.Sprintf("ORDER BY %s %s", dbField, order)
		}
	}

	// Get products with category
	query := fmt.Sprintf(`
		SELECT p.id, p.seller_id, p.title, p.description, p.price, p.quantity, p.images, 
		       p.category_id, p.is_active, p.created_at, p.updated_at,
		       c.id, c.name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		%s
		%s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderByClause, argCount, argCount+1)
	args = append(args, pageSize, offset)

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	var products []*product.ProductWithCategory
	for rows.Next() {
		var p product.ProductWithCategory
		var categoryID, categoryName *string
		
		err := rows.Scan(
			&p.ID, &p.SellerID, &p.Title, &p.Description, &p.Price, &p.Quantity, &p.Images, 
			&p.CategoryID, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
			&categoryID, &categoryName)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}
		
		// Set category if it exists
		if categoryID != nil && categoryName != nil {
			p.Category = &product.Category{
				ID:   *categoryID,
				Name: *categoryName,
			}
		}
		
		products = append(products, &p)
	}

	return products, total, nil
}

func (r *productRepository) Update(p *product.Product) error {
	query := `
		UPDATE products 
		SET title = $1, description = $2, price = $3, quantity = $4, images = $5, category_id = $6, is_active = $7, updated_at = $8
		WHERE id = $9
	`
	_, err := r.db.Exec(context.Background(), query,
		p.Title, p.Description, p.Price, p.Quantity, p.Images, p.CategoryID, p.IsActive, p.UpdatedAt, p.ID)
	return err
}

func (r *productRepository) UpdateQuantity(id string, quantity int) error {
	query := `
		UPDATE products 
		SET quantity = $1, updated_at = $2
		WHERE id = $3
	`
	_, err := r.db.Exec(context.Background(), query, quantity, "now()", id)
	return err
}

func (r *productRepository) Delete(id string) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}

func (r *productRepository) GetCategories() ([]*product.Category, error) {
	query := `SELECT id, name FROM categories ORDER BY name`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []*product.Category
	for rows.Next() {
		var c product.Category
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, &c)
	}

	return categories, nil
}
