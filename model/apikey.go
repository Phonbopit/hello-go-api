package model

type APIKey struct {
	ID           string `json:"id"`
	Key          string `json:"key"`            // Plaintext on create, hash when from DB
	Name         string `json:"name"`
	CreatedAt    int64  `json:"created_at"`
	LastUsedAt   *int64 `json:"last_used_at"`   // Pointer for SQL NULL
	RequestCount int64  `json:"request_count"`
}
