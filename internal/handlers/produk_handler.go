package handlers

import (
	"encoding/json"
	"kasir-api/internal/models"
	"kasir-api/internal/repository"
	"net/http"
	"strconv"
	"strings"
)

type ProdukHandler struct {
	repo repository.ProdukRepository
}

func NewProdukHandler(repo repository.ProdukRepository) *ProdukHandler {
	return &ProdukHandler{repo: repo}
}

func (h *ProdukHandler) HandleProdukList(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.getProduk(w, r)
	} else if r.Method == "POST" {
		h.createProduk(w, r)
	}
}

func (h *ProdukHandler) HandleProdukDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.getDetailProduk(w, r)
	} else if r.Method == "PUT" {
		h.updateProduk(w, r)
	}
}

func (h *ProdukHandler) getProduk(w http.ResponseWriter, r *http.Request) {
	produk := h.repo.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produk)
}

func (h *ProdukHandler) createProduk(w http.ResponseWriter, r *http.Request) {
	var produkBaru models.Produk
	err := json.NewDecoder(r.Body).Decode(&produkBaru)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	created := h.repo.Create(produkBaru)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(created)
}

func (h *ProdukHandler) getDetailProduk(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	p, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *ProdukHandler) updateProduk(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var updateData models.Produk
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	updated, err := h.repo.Update(id, updateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}
