package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"hello-go-api/model"
	"hello-go-api/store"
)

// APIKeyHandler handles API key management endpoints.
type APIKeyHandler struct {
	Store store.APIKeyStore
}

// CreateKeyRequest is the request body for creating an API key.
type CreateKeyRequest struct {
	Name string `json:"name"`
}

// CreateKey handles POST /admin/keys - creates a new API key.
func (h *APIKeyHandler) CreateKey(w http.ResponseWriter, r *http.Request) {
	var req CreateKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, `{"error":"name is required"}`, http.StatusBadRequest)
		return
	}

	// Generate random API key
	key, err := generateAPIKey()
	if err != nil {
		http.Error(w, `{"error":"failed to generate key"}`, http.StatusInternalServerError)
		return
	}

	// Generate random ID
	id, err := generateID()
	if err != nil {
		http.Error(w, `{"error":"failed to generate ID"}`, http.StatusInternalServerError)
		return
	}

	apiKey := model.APIKey{
		ID:           id,
		Key:          key,
		Name:         req.Name,
		CreatedAt:    time.Now().Unix(),
		LastUsedAt:   nil,
		RequestCount: 0,
	}

	if err := h.Store.Create(apiKey); err != nil {
		http.Error(w, `{"error":"failed to create key"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(apiKey)
}

// ListKeys handles GET /admin/keys - lists all API keys with stats.
func (h *APIKeyHandler) ListKeys(w http.ResponseWriter, r *http.Request) {
	keys, err := h.Store.List()
	if err != nil {
		http.Error(w, `{"error":"failed to fetch keys"}`, http.StatusInternalServerError)
		return
	}

	// Mask keys for security (hashes are returned from DB, show generic mask)
	for i := range keys {
		keys[i].Key = "sk_••••••••"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keys)
}

// generateAPIKey generates a random API key in format: sk_<32 hex chars>
func generateAPIKey() (string, error) {
	bytes := make([]byte, 16) // 16 bytes = 32 hex chars
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "sk_" + hex.EncodeToString(bytes), nil
}

// generateID generates a random ID for the API key.
func generateID() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// maskAPIKey masks the middle part of an API key for security.
// Example: sk_1a2b3c4d5e6f7890abcdef1234567890 → sk_...567890
func maskAPIKey(key string) string {
	if len(key) < 10 {
		return key
	}
	// Show prefix "sk_" and last 6 chars
	prefix := "sk_"
	if strings.HasPrefix(key, prefix) {
		return prefix + "..." + key[len(key)-6:]
	}
	return key
}

// RegisterRoutes registers API key management routes.
func (h *APIKeyHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /admin/keys", h.CreateKey)
	mux.HandleFunc("GET /admin/keys", h.ListKeys)
}
