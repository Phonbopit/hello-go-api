package store

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"hello-go-api/model"
)

type APIKeyStore interface {
	Create(key model.APIKey) error
	GetByKey(key string) (model.APIKey, error)
	IncrementUsage(key string) error
	List() ([]model.APIKey, error)
}

type SQLiteAPIKeyStore struct {
	db *sql.DB
}

func NewSQLiteAPIKeyStore(db *sql.DB) (*SQLiteAPIKeyStore, error) {
	schema := `
	CREATE TABLE IF NOT EXISTS api_keys (
		id TEXT PRIMARY KEY,
		key_hash TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		created_at INTEGER NOT NULL,
		last_used_at INTEGER,
		request_count INTEGER DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_api_keys_hash ON api_keys(key_hash);
	`

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to create api_keys schema: %w", err)
	}

	return &SQLiteAPIKeyStore{db: db}, nil
}

func (s *SQLiteAPIKeyStore) Create(key model.APIKey) error {
	keyHash := hashAPIKey(key.Key)
	_, err := s.db.Exec(
		`INSERT INTO api_keys (id, key_hash, name, created_at, last_used_at, request_count)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		key.ID, keyHash, key.Name, key.CreatedAt, key.LastUsedAt, key.RequestCount,
	)
	return err
}

func (s *SQLiteAPIKeyStore) GetByKey(key string) (model.APIKey, error) {
	var k model.APIKey
	keyHash := hashAPIKey(key)

	err := s.db.QueryRow(
		`SELECT id, key_hash, name, created_at, last_used_at, request_count
		 FROM api_keys WHERE key_hash = ?`, keyHash,
	).Scan(&k.ID, &k.Key, &k.Name, &k.CreatedAt, &k.LastUsedAt, &k.RequestCount)

	if err == sql.ErrNoRows {
		return k, fmt.Errorf("API key not found")
	}
	return k, err
}

func (s *SQLiteAPIKeyStore) IncrementUsage(key string) error {
	now := time.Now().Unix()
	keyHash := hashAPIKey(key)
	_, err := s.db.Exec(
		`UPDATE api_keys
		 SET request_count = request_count + 1, last_used_at = ?
		 WHERE key_hash = ?`,
		now, keyHash,
	)
	return err
}

func (s *SQLiteAPIKeyStore) List() ([]model.APIKey, error) {
	rows, err := s.db.Query(
		`SELECT id, key_hash, name, created_at, last_used_at, request_count
		 FROM api_keys ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []model.APIKey
	for rows.Next() {
		var k model.APIKey
		if err := rows.Scan(&k.ID, &k.Key, &k.Name, &k.CreatedAt, &k.LastUsedAt, &k.RequestCount); err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}

	return keys, rows.Err()
}

func hashAPIKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}
