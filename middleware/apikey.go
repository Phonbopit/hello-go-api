package middleware

import (
	"context"
	"net/http"

	"hello-go-api/store"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

const apiKeyIDKey contextKey = "api_key_id"

// RequireAPIKey is middleware that validates API keys and tracks usage.
// It checks the X-API-Key header, validates it against the database,
// and increments usage stats.
func RequireAPIKey(keyStore store.APIKeyStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("X-API-Key")
			if key == "" {
				http.Error(w, `{"error":"missing API key"}`, http.StatusUnauthorized)
				return
			}

			apiKey, err := keyStore.GetByKey(key)
			if err != nil {
				http.Error(w, `{"error":"invalid API key"}`, http.StatusUnauthorized)
				return
			}

			// Track usage (increment count, update last_used_at)
			if err := keyStore.IncrementUsage(key); err != nil {
				// Log error but don't fail the request
			}

			ctx := context.WithValue(r.Context(), apiKeyIDKey, apiKey.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}


func GetAPIKeyID(r *http.Request) string {
	if id, ok := r.Context().Value(apiKeyIDKey).(string); ok {
		return id
	}
	return ""
}
