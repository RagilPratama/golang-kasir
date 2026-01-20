package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

var produk = []Produk{
	{ID: 1, Nama: "Indomie Bangladesh", Harga: 7500, Stok: 20},
	{ID: 2, Nama: "Teh Tarik", Harga: 3000, Stok: 30},
}

func getDetailProduk(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	for _, p := range produk {
		if p.ID == id {
			json.NewEncoder(w).Encode(p)
			return
		}
	}
	http.Error(w, "produk not found", http.StatusNotFound)
}

func updateProduk(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	for i, p := range produk {
		if p.ID == id {
			var updateProduk Produk
			err := json.NewDecoder(r.Body).Decode(&updateProduk)
			if err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			produk[i] = updateProduk
			json.NewEncoder(w).Encode(produk[i])
			return
		}
	}
	http.Error(w, "produk not found", http.StatusNotFound)
}

func deleteProduk(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	for i, p := range produk {
		if p.ID == id {
			produk = append(produk[:i], produk[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Produk Berhasil Dihapus"})
			return
		}
	}
	http.Error(w, "produk not found", http.StatusNotFound)
}

func main() {
	//GET detail produk
	//PUT update produk
	//DELETE delete produk
	http.HandleFunc("/api/produk/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getDetailProduk(w, r)
		} else if r.Method == "PUT" {
			updateProduk(w, r)
		} else if r.Method == "DELETE" {
			deleteProduk(w, r)
		}
	})

	// GET produk
	// POST produk
	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)
		} else if r.Method == "POST" {
			// baca dari request
			var produkBaru Produk
			err := json.NewDecoder(r.Body).Decode(&produkBaru)
			if err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			produkBaru.ID = len(produk) + 1
			produk = append(produk, produkBaru)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(produkBaru)
		}

	})
	fmt.Println("server running di localhost:8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}
