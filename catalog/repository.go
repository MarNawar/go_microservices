package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	// "github.com/MarNawar/microservices/catalog/pb"
	"github.com/lib/pq"
	// _ "github.com/lib/pq" // PostgreSQL driver
)

var (
	ErrNotFound = errors.New("entity not found")
)

type Repository interface {
	Close()
	PutProduct(ctx context.Context, p Product) error
	GetProductById(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type postgresRepository struct {
	db *sql.DB
}

// type productDocument struct {
// 	Name        string `json:"name"`
// 	Description string `json:"description"`
// 	Price       string `json:"price"`
// }

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() {
	if r.db != nil {
		r.db.Close()
	}
}

func (r *postgresRepository) PutProduct(ctx context.Context, p Product) error {
	// Implement this method
	_, err := r.db.ExecContext(ctx, "INSERT INTO products(id, name, description, price) VALUES($1, $2, $3, $4)", p.ID, p.Name, p.Description, p.Price)
	return err
}

func (r *postgresRepository) GetProductById(ctx context.Context, id string) (*Product, error) {
	// Implement this method
	row := r.db.QueryRowContext(ctx, "SELECT id, name, description, price From account WHERE id = $1", id)
	p := &Product{}
	if err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *postgresRepository) ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT id, name, description, price FROM accounts ORDER BY id DESC OFFSET $1 LIMIT $2",
		skip,
		take,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	products := []Product{}

	for rows.Next() {
		p := &Product{}
		if err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price); err == nil {
			products = append(products, *p)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *postgresRepository) ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	if len(ids) == 0 {
		return []Product{}, nil
	}

	query := "SELECT id, name, description, price FROM products WHERE id = ANY($1)"
	rows, err := r.db.QueryContext(ctx, query, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("failed to query products by IDs: %w", err)
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price); err != nil {
			return nil, fmt.Errorf("failed to scan product row: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", err)
	}

	return products, nil
}

func (r *postgresRepository) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	searchQuery := "%" + query + "%"

	rows, err := r.db.QueryContext(
		ctx,
		"SELECT id, name, description, price FROM accounts WHERE name ILike $3 ORDER BY id DESC OFFSET $1 LIMIT $2",
		skip,
		take,
		searchQuery,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse rows into []Product
	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil

}
