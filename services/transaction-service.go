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

// AddTransactionAcquirer kreira novi transaction za acquirer
func (s *TransactionService) AddTransactionAcquirer(paymentRequest *models.PaymentRequest) (*models.Transaction, error) {
	tx := &models.Transaction{
		Amount:            paymentRequest.Amount.InexactFloat64(),                // decimal.Decimal → float64
		MerchantOrderID:   paymentRequest.MerchantOrderID,                        // string → string (ok)
		MerchantTimestamp: paymentRequest.MerchantTimestamp.Format(time.RFC3339), // time.Time → string
		AcquirerOrderID:   generateRandomOrderID(),
		AcquirerTimestamp: time.Now().Format(time.RFC3339), // time.Time → string
		IssuerOrderID:     "",                              // prazno jer još nije poznato
		IssuerTimestamp:   "",                              // prazno string polje
		Status:            "CREATED",
		PaymentRequestID:  paymentRequest.ID, // sada možemo povezati transakciju sa request-om
	}

	savedTx, err := s.repo.Save(tx)
	if err != nil {
		return nil, err
	}

	log.Printf("Created new acquirer transaction: %s", savedTx.ID)
	return savedTx, nil
}

// Helper funkcija za generisanje random order ID
func generateRandomOrderID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	id := make([]byte, 16)
	rand.Seed(time.Now().UnixNano()) // seed random generator

	for i := range id {
		idx := rand.Intn(len(charset)) // nasumičan index u charset
		id[i] = charset[idx]
	}
	return string(id)
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
