package handlers

import (
	"encoding/json"
	"kasir-api/internal/models"
	"kasir-api/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) HandleCategoryList(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.GetAll(w, r)
	} else if r.Method == "POST" {
		h.Create(w, r)
	}
}

func (h *CategoryHandler) HandleCategoryDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.GetDetailCategories(w, r)
	} else if r.Method == "PUT" {
		h.Update(w, r)
	} else if r.Method == "DELETE" {
		h.Delete(w, r)
	}
}

// GetAll godoc
// @Summary Get all categories
// @Description Get list of all categories
// @Tags category
// @Accept json
// @Produce json
// @Success 200 {array} models.Category
// @Router /category [get]
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	category := h.service.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// Create godoc
// @Summary Create new category
// @Description Create a new category
// @Tags category
// @Accept json
// @Produce json
// @Param category body models.Category true "Category Data"
// @Success 200 {object} models.Category
// @Failure 400 {string} string "invalid request"
// @Router /category [post]
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	category = h.service.Create(category)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// GetDetailCategories godoc
// @Summary Get detail category
// @Description Get detail of a category by ID
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.Category
// @Failure 400 {string} string "invalid id"
// @Failure 404 {string} string "category not found"
// @Router /category/{id} [get]
func (h *CategoryHandler) GetDetailCategories(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	category, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// Update godoc
// @Summary Update category
// @Description Update an existing category
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body models.Category true "Category Data"
// @Success 200 {object} models.Category
// @Failure 400 {string} string "invalid request/id"
// @Failure 404 {string} string "category not found"
// @Router /category/{id} [put]
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated, err := h.service.Update(id, category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// Delete godoc
// @Summary Delete category
// @Description Delete a category by ID
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 204 {object} nil
// @Failure 400 {string} string "invalid id"
// @Failure 404 {string} string "category not found"
// @Router /category/{id} [delete]
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
