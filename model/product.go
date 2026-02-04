package model

// Product represents a product in the system.
// The `json` tags control JSON serialization (field name mapping).
type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
