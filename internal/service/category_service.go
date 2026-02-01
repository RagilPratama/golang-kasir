package service

import (
	"kasir-api/internal/models"
	"kasir-api/internal/repository"
)

type CategoryService interface {
	GetAll() []models.Category
	GetByID(id int) (*models.Category, error)
	Create(c models.Category) models.Category
	Update(id int, c models.Category) (*models.Category, error)
	Delete(id int) error
}

type categoryService struct {
	repo repository.CategoriesRepository
}

func NewCategoryService(repo repository.CategoriesRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetAll() []models.Category {
	return s.repo.GetAll()
}

func (s *categoryService) GetByID(id int) (*models.Category, error) {
	return s.repo.GetByID(id)
}

func (s *categoryService) Create(c models.Category) models.Category {
	return s.repo.Create(c)
}

func (s *categoryService) Update(id int, c models.Category) (*models.Category, error) {
	return s.repo.Update(id, c)
}

func (s *categoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
