package service

import (
	"errors"
	"kasir-api/internal/models"
	"kasir-api/internal/repository"
)

type TransactionService interface {
	Checkout(items []models.CheckoutItem) (*models.Transaction, error)
}

type transactionService struct {
	repo        repository.TransactionRepository
	productRepo repository.ProductRepository
}

func NewTransactionService(repo repository.TransactionRepository, productRepo repository.ProductRepository) TransactionService {
	return &transactionService{repo: repo, productRepo: productRepo}
}

func (s *transactionService) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	var totalAmount int
	var details []models.TransactionDetail

	for _, item := range items {
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			return nil, errors.New("product not found with ID " + string(rune(item.ProductID)))
		}

		if product.Stock < item.Quantity {
			return nil, errors.New("insufficient stock for product " + product.Name)
		}

		subtotal := product.Price * item.Quantity
		totalAmount += subtotal

		details = append(details, models.TransactionDetail{
			ProductID:   product.ID,
			ProductName: product.Name,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	transaction := &models.Transaction{
		TotalAmount: totalAmount,
		Details:     details,
	}

	err := s.repo.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
