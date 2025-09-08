package services

import (
	"github.com/google/uuid"
	"log"

	"github.com/sep-2024-team-35/bank-servce-back-end/models"
	"github.com/sep-2024-team-35/bank-servce-back-end/repositories"
)

type TransactionService struct {
	repo repositories.TransactionRepository
}

func NewTransactionService(repo repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) GetTransactionByPaymentRequestID(paymentRequestID uuid.UUID) (*models.Transaction, error) {
	return s.repo.FindByPaymentRequestID(paymentRequestID)
}

func (s *TransactionService) SuccessTransaction(tx *models.Transaction) (*models.Transaction, error) {
	tx.Status = "SUCCESS"
	savedTx, err := s.repo.Save(tx)
	if err != nil {
		return nil, err
	}
	log.Printf("Transaction %d marked as SUCCESS", tx.ID)
	return savedTx, nil
}

func (s *TransactionService) FailTransaction(tx *models.Transaction) (*models.Transaction, error) {
	tx.Status = "FAILED"
	savedTx, err := s.repo.Save(tx)
	if err != nil {
		return nil, err
	}
	log.Printf("Transaction %d marked as FAILED", tx.ID)
	return savedTx, nil
}

func (s *TransactionService) ErrorTransaction(tx *models.Transaction) (*models.Transaction, error) {
	tx.Status = "ERROR"
	savedTx, err := s.repo.Save(tx)
	if err != nil {
		return nil, err
	}
	log.Printf("Transaction %d marked as ERROR", tx.ID)
	return savedTx, nil
}
