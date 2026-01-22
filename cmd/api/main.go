package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/internal/handlers"
	"kasir-api/internal/repository"
	"net/http"
)

func main() {
	// Initialize Repository
	produkRepo := repository.NewMemoryProdukRepository()

	// Initialize Handler
	produkHandler := handlers.NewProdukHandler(produkRepo)

	// Setup Routes
	// Root handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{
			"message": "Welcome to Kasir API",
			"status":  "running",
		}
		json.NewEncoder(w).Encode(response)
	})

	// GET detail produk
	// PUT update produk
	http.HandleFunc("/api/produk/", produkHandler.HandleProdukDetail)

	// GET produk
	// POST produk
	http.HandleFunc("/api/produk", produkHandler.HandleProdukList)

	fmt.Println("server running di localhost:8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}
