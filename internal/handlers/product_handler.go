package handlers

import (
	"encoding/json"
	"kasir-api/internal/models"
	"kasir-api/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) HandleProductList(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.GetProduct(w, r)
	} else if r.Method == "POST" {
		h.CreateProduct(w, r)
	}
}

func (h *ProductHandler) HandleProductDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.GetDetailProduct(w, r)
	} else if r.Method == "PUT" {
		h.UpdateProduct(w, r)
	} else if r.Method == "DELETE" {
		h.DeleteProduct(w, r)
	}
}

// GetProduct godoc
// @Summary Get all product
// @Description Get list of all product
// @Tags product
// @Accept json
// @Produce json
// @Success 200 {array} models.Product
// @Router /product [get]
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	product := h.service.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// CreateProduct godoc
// @Summary Create new product
// @Description Create a new product
// @Tags product
// @Accept json
// @Produce json
// @Param product body models.Product true "Product Data"
// @Success 201 {object} models.Product
// @Failure 400 {string} string "invalid request"
// @Router /product [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var productBaru models.Product
	err := json.NewDecoder(r.Body).Decode(&productBaru)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	created := h.service.Create(productBaru)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(created)
}

// GetDetailProduct godoc
// @Summary Get detail product
// @Description Get detail of a product by ID
// @Tags product
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Product
// @Failure 400 {string} string "invalid id"
// @Failure 404 {string} string "product not found"
// @Router /product/{id} [get]
func (h *ProductHandler) GetDetailProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	p, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// UpdateProduct godoc
// @Summary Update product
// @Description Update an existing product
// @Tags product
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.Product true "Product Data"
// @Success 200 {object} models.Product
// @Failure 400 {string} string "invalid request/id"
// @Failure 404 {string} string "product not found"
// @Router /product/{id} [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var updateData models.Product
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	updated, err := h.service.Update(id, updateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Delete a product by ID
// @Tags product
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "invalid id"
// @Failure 404 {string} string "product not found"
// @Router /product/{id} [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Product Berhasil Dihapus"})
}
