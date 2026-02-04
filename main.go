package main

import (
	"fmt"
	"log"
	"net/http"

	"hello-go-api/handler"
	"hello-go-api/store"
)

func main() {
	db, err := store.NewSQLiteStore("products.db")

	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	defer db.Close() // Close DB when main() exits

	productHandler := &handler.ProductHandler{Store: db}

	mux := http.NewServeMux()
	productHandler.RegisterRoutes(mux)

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("Database: products.db")
	http.ListenAndServe(":8080", mux)
}
