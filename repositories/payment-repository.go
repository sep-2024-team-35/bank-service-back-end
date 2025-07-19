package repositories

import (
	"github.com/sep-2024-team-35/bank-servce-back-end/models"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	SaveRequest(request *models.PaymentRequest) (*models.PaymentRequest, error)
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) SaveRequest(request *models.PaymentRequest) (*models.PaymentRequest, error) {
	if err := r.db.Create(request).Error; err != nil {
		return nil, err
	}
	return request, nil
}
