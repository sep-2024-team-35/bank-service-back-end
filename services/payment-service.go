package services

import (
	"github.com/sep-2024-team-35/bank-servce-back-end/dto"
	"github.com/sep-2024-team-35/bank-servce-back-end/models"
	"github.com/sep-2024-team-35/bank-servce-back-end/repositories"
)

type PaymentService interface {
	CreateRequest(dto dto.PaymentRequestDto) (*models.PaymentRequest, error)
}

type paymentService struct {
	repo repositories.PaymentRepository
}

func NewPaymentService(repo repositories.PaymentRepository) PaymentService {
	return &paymentService{repo: repo}
}

func (s *paymentService) CreateRequest(dto dto.PaymentRequestDto) (*models.PaymentRequest, error) {
	request := &models.PaymentRequest{
		MerchantID:        dto.MerchantId,
		MerchantPassword:  dto.MerchantPassword,
		Amount:            dto.Amount,
		MerchantOrderID:   dto.MerchantOrderId,
		MerchantTimestamp: dto.MerchantTimestamp,
		SuccessURL:        dto.SuccessUrl,
		FailedURL:         dto.FailedUrl,
		ErrorURL:          dto.ErrorUrl,
	}

	return s.repo.SaveRequest(request)
}
