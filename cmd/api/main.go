package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/internal/handlers"
	"kasir-api/internal/repository"
	"net/http"
)

func main() {
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

	// Initialize Repository
	produkRepo := repository.NewMemoryProdukRepository()
	categoryRepo := repository.NewMemoryCategoryRepository()
	// Initialize Handler
	produkHandler := handlers.NewProdukHandler(produkRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	// GET detail produk
	// PUT update produk
	http.HandleFunc("/api/produk/", produkHandler.HandleProdukDetail)
	// GET produk
	// POST produk
	http.HandleFunc("/api/produk", produkHandler.HandleProdukList)

	//Category
	http.HandleFunc("/api/category", categoryHandler.HandleCategoryList)
	http.HandleFunc("/api/category/", categoryHandler.HandleCategoryDetail)

	fmt.Println("server running di localhost:8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}
