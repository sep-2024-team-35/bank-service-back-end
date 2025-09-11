package repositories

import (
	"errors"
	"github.com/google/uuid"

	"github.com/sep-2024-team-35/bank-servce-back-end/models"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Save(transaction *models.Transaction) (*models.Transaction, error)
	FindByPaymentRequestID(paymentRequestID uuid.UUID) (*models.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Save(transaction *models.Transaction) (*models.Transaction, error) {
	if err := r.db.Save(transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *transactionRepository) FindByPaymentRequestID(paymentRequestID uuid.UUID) (*models.Transaction, error) {
	var tx models.Transaction
	err := r.db.
		Where("payment_request_id = ?", paymentRequestID).
		Order("id desc").
		First(&tx).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &tx, nil
}
