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
		h.GetProduk(w, r)
	} else if r.Method == "POST" {
		h.CreateProduk(w, r)
	}
}

func (h *ProdukHandler) HandleProdukDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.GetDetailProduk(w, r)
	} else if r.Method == "PUT" {
		h.UpdateProduk(w, r)
	} else if r.Method == "DELETE" {
		h.DeleteProduk(w, r)
	}
}

// GetProduk godoc
// @Summary Get all produk
// @Description Get list of all produk
// @Tags produk
// @Accept json
// @Produce json
// @Success 200 {array} models.Produk
// @Router /produk [get]
func (h *ProdukHandler) GetProduk(w http.ResponseWriter, r *http.Request) {
	produk := h.repo.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produk)
}

// CreateProduk godoc
// @Summary Create new produk
// @Description Create a new produk
// @Tags produk
// @Accept json
// @Produce json
// @Param produk body models.Produk true "Produk Data"
// @Success 201 {object} models.Produk
// @Failure 400 {string} string "invalid request"
// @Router /produk [post]
func (h *ProdukHandler) CreateProduk(w http.ResponseWriter, r *http.Request) {
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

// GetDetailProduk godoc
// @Summary Get detail produk
// @Description Get detail of a produk by ID
// @Tags produk
// @Accept json
// @Produce json
// @Param id path int true "Produk ID"
// @Success 200 {object} models.Produk
// @Failure 400 {string} string "invalid id"
// @Failure 404 {string} string "produk not found"
// @Router /produk/{id} [get]
func (h *ProdukHandler) GetDetailProduk(w http.ResponseWriter, r *http.Request) {
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

// UpdateProduk godoc
// @Summary Update produk
// @Description Update an existing produk
// @Tags produk
// @Accept json
// @Produce json
// @Param id path int true "Produk ID"
// @Param produk body models.Produk true "Produk Data"
// @Success 200 {object} models.Produk
// @Failure 400 {string} string "invalid request/id"
// @Failure 404 {string} string "produk not found"
// @Router /produk/{id} [put]
func (h *ProdukHandler) UpdateProduk(w http.ResponseWriter, r *http.Request) {
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

// DeleteProduk godoc
// @Summary Delete produk
// @Description Delete a produk by ID
// @Tags produk
// @Accept json
// @Produce json
// @Param id path int true "Produk ID"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "invalid id"
// @Failure 404 {string} string "produk not found"
// @Router /produk/{id} [delete]
func (h *ProdukHandler) DeleteProduk(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.repo.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Produk Berhasil Dihapus"})
}
