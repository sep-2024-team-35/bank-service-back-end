package services

import (
	"github.com/google/uuid"
	"log"
	"math/rand"
	"time"

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

func (s *TransactionService) AddTransactionAcquirer(paymentRequest *models.PaymentRequest) (*models.Transaction, error) {
	tx := &models.Transaction{
		ID:                uuid.New(),
		Amount:            paymentRequest.Amount,
		MerchantOrderID:   paymentRequest.MerchantOrderID,
		MerchantTimestamp: paymentRequest.MerchantTimestamp.Format(time.RFC3339),
		AcquirerOrderID:   generateRandomOrderID(),
		AcquirerTimestamp: time.Now().Format(time.RFC3339),
		IssuerOrderID:     "",
		IssuerTimestamp:   "",
		Status:            "CREATED",
		PaymentRequestID:  paymentRequest.ID,
	}

	savedTx, err := s.repo.Save(tx)
	if err != nil {
		return nil, err
	}

	log.Printf("Created new acquirer transaction: %s", savedTx.ID)
	return savedTx, nil
}

func generateRandomOrderID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	id := make([]byte, 16)
	rand.Seed(time.Now().UnixNano())

	for i := range id {
		idx := rand.Intn(len(charset))
		id[i] = charset[idx]
	}
	return string(id)
}
