package repositories

import (
	"errors"
	"github.com/sep-2024-team-35/bank-servce-back-end/models"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	SaveRequest(request *models.PaymentRequest) (*models.PaymentRequest, error)
	SaveRequestTransactional(tx *gorm.DB, request *models.PaymentRequest) (*models.PaymentRequest, error)
	GetByID(id string) (*models.PaymentRequest, error)
	DB() *gorm.DB
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

func (r *paymentRepository) SaveRequestTransactional(tx *gorm.DB, request *models.PaymentRequest) (*models.PaymentRequest, error) {
	if err := tx.Create(request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (r *paymentRepository) GetByID(id string) (*models.PaymentRequest, error) {
	var request models.PaymentRequest
	if err := r.db.First(&request, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &request, nil
}

func (r *paymentRepository) DB() *gorm.DB {
	return r.db
}
