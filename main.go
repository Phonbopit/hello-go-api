package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"hello-go-api/handler"
	"hello-go-api/middleware"
	"hello-go-api/store"
)

func main() {
	// Open database connection
	db, err := sql.Open("sqlite3", "products.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Initialize stores
	productStore, err := store.NewSQLiteStore(db)
	if err != nil {
		log.Fatal("Failed to initialize product store:", err)
	}

	apiKeyStore, err := store.NewSQLiteAPIKeyStore(db)
	if err != nil {
		log.Fatal("Failed to initialize API key store:", err)
	}

	// Initialize handlers
	productHandler := &handler.ProductHandler{Store: productStore}
	apiKeyHandler := &handler.APIKeyHandler{Store: apiKeyStore}

	// Create router and middleware
	mux := http.NewServeMux()
	authMiddleware := middleware.RequireAPIKey(apiKeyStore)

	// Register routes
	apiKeyHandler.RegisterRoutes(mux)
	productHandler.RegisterRoutes(mux, authMiddleware)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
