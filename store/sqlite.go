package store

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // SQLite driver (the _ means: import for side effects only)

	"hello-go-api/model"
)

type ProductStore interface {
	List() ([]model.Product, error)
	Get(id string) (model.Product, error)
	Create(p model.Product) error
	Delete(id string) error
}

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS products (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		price REAL NOT NULL
	);`

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) List() ([]model.Product, error) {
	rows, err := s.db.Query("SELECT id, name, price FROM products")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []model.Product

	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (s *SQLiteStore) Get(id string) (model.Product, error) {
	var p model.Product

	err := s.db.QueryRow("SELECT id, name, price FROM products WHERE id = ?", id).
		Scan(&p.ID, &p.Name, &p.Price)

	if err == sql.ErrNoRows {
		return p, fmt.Errorf("product not found")
	}
	if err != nil {
		return p, err
	}

	return p, nil
}

func (s *SQLiteStore) Create(p model.Product) error {
	_, err := s.db.Exec(
		"INSERT INTO products (id, name, price) VALUES (?, ?, ?)",
		p.ID, p.Name, p.Price,
	)
	return err
}

func (s *SQLiteStore) Delete(id string) error {
	result, err := s.db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}
