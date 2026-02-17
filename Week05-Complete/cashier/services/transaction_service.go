package services

import (
	"cashier/models"
	"cashier/repositories"
	"time"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(userID int, items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
	return s.repo.CreateTransaction(userID, items, useLock)
}

func (s *TransactionService) GetReport(startDate, endDate time.Time) (*models.Report, error) {
	return s.repo.GetReport(startDate, endDate)
}