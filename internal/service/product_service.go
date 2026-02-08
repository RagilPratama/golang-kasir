package service

import (
	"kasir-api/internal/models"
	"kasir-api/internal/repository"
)

type ProductService interface {
	GetAll(name string) []models.Product
	GetByID(id int) (*models.Product, error)
	Create(p models.Product) models.Product
	Update(id int, p models.Product) (*models.Product, error)
	Delete(id int) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetAll(name string) []models.Product {
	return s.repo.GetAll(name)
}

func (s *productService) GetByID(id int) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *productService) Create(p models.Product) models.Product {
	return s.repo.Create(p)
}

func (s *productService) Update(id int, p models.Product) (*models.Product, error) {
	return s.repo.Update(id, p)
}

func (s *productService) Delete(id int) error {
	return s.repo.Delete(id)
}
