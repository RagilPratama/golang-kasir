package handlers

import (
	"encoding/json"
	"kasir-api/internal/models"
	"kasir-api/internal/repository"
	"net/http"
	"strconv"
	"strings"
)

type CategoryHandler struct {
	repo repository.CategoriesRepository
}

func NewCategoryHandler(repo repository.CategoriesRepository) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

func (h *CategoryHandler) HandleCategoryList(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.getAll(w, r)
	} else if r.Method == "POST" {
		h.create(w, r)
	}
}

func (h *CategoryHandler) HandleCategoryDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.getDetailCategories(w, r)
	} else if r.Method == "PUT" {
		h.update(w, r)
	} else if r.Method == "DELETE" {
		h.delete(w, r)
	}
}

func (h *CategoryHandler) getAll(w http.ResponseWriter, r *http.Request) {
	category := h.repo.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) create(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	category = h.repo.Create(category)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) getDetailCategories(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	category, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) update(w http.ResponseWriter, r *http.Request) {
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

	updated, err := h.repo.Update(id, category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func (h *CategoryHandler) delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.repo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
