package account

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type Repository interface {
	Close()
	PutAccount(ctx context.Context, a Account) error
	GetAccountById(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type postgresRepository struct {
	db *sql.DB
}

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

func (r *postgresRepository) Ping() error {
	return r.db.Close()
}

func (r *postgresRepository) PutAccount(ctx context.Context, a Account) error {
	// Implement this method
	_, err := r.db.ExecContext(ctx, "INSERT INTO accounts(id, name) VALUES($1, $2)", a.ID, a.Name)
	return err
}

func (r *postgresRepository) GetAccountById(ctx context.Context, id string) (*Account, error) {
	// Implement this method
	row := r.db.QueryRowContext(ctx, "SELECT id, name From account WHERE id = $1", id)
	a:= &Account{}
	if err :=  row.Scan(&a.ID, &a.Name); err != nil{
		return nil, err
	}
	return nil, nil
}

func (r *postgresRepository) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	// Implement this method
	rows, err :=  r.db.QueryContext(
		ctx, 
		"SELECT id, name  FROM accounts ORDER BY id DESC OFFSET $1 LIMIT $2",
		skip,
		take,
	)

	if err != nil{
		return nil, err
	}
	defer rows.Close()
	accounts := []Account{}

	for rows.Next(){
		a:= &Account{}
		if err =  rows.Scan(&a.ID, &a.Name); err == nil{
			accounts = append(accounts, *a)
		}
	}

	if err = rows.Err(); err != nil{
		return nil, err
	}

	return accounts, nil
}
