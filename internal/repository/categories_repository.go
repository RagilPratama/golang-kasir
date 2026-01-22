package repository

import (
	"errors"
	"kasir-api/internal/models"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

type CategoriesRepository interface {
	GetAll() []models.Category
	GetByID(id int) (*models.Category, error)
	Create(c models.Category) models.Category
	Update(id int, c models.Category) (*models.Category, error)
	Delete(id int) error
}

type MemoryCategoryRepository struct {
	categories []models.Category
}

func NewMemoryCategoryRepository() CategoriesRepository {
	return &MemoryCategoryRepository{
		categories: []models.Category{
			{
				ID:          1,
				Name:        "Makanan",
				Description: "Kategori makanan",
			},
			{
				ID:          2,
				Name:        "Minuman",
				Description: "Kategori minuman",
			},
		},
	}
}

func (r *MemoryCategoryRepository) GetAll() []models.Category {
	return r.categories
}

func (r *MemoryCategoryRepository) GetByID(id int) (*models.Category, error) {
	for i := range r.categories {
		if r.categories[i].ID == id {
			return &r.categories[i], nil
		}
	}
	return nil, ErrCategoryNotFound
}

func (r *MemoryCategoryRepository) Create(c models.Category) models.Category {
	c.ID = len(r.categories) + 1
	r.categories = append(r.categories, c)
	return c
}

func (r *MemoryCategoryRepository) Update(id int, updateData models.Category) (*models.Category, error) {
	for i := range r.categories {
		if r.categories[i].ID == id {
			r.categories[i].Name = updateData.Name
			r.categories[i].Description = updateData.Description
			return &r.categories[i], nil
		}
	}
	return nil, ErrCategoryNotFound
}

func (r *MemoryCategoryRepository) Delete(id int) error {
	for i := range r.categories {
		if r.categories[i].ID == id {
			r.categories = append(r.categories[:i], r.categories[i+1:]...)
			return nil
		}
	}
	return ErrCategoryNotFound
}
