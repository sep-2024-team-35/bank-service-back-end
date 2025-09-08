package services

import (
	"errors"
	"github.com/sep-2024-team-35/bank-servce-back-end/dto"
	"github.com/sep-2024-team-35/bank-servce-back-end/models"
	"github.com/sep-2024-team-35/bank-servce-back-end/repositories"
)

type PaymentService interface {
	CreateRequest(dto dto.PaymentRequestDto) (*models.PaymentRequest, error)
}

type paymentService struct {
	repo           repositories.PaymentRepository
	accountService AccountService
}

func NewPaymentService(repo repositories.PaymentRepository, accSvc AccountService) PaymentService {
	return &paymentService{
		repo:           repo,
		accountService: accSvc,
	}
}

func (s *paymentService) CreateRequest(dto dto.PaymentRequestDto) (*models.PaymentRequest, error) {
	exists, err := s.accountService.isMerchantAccountExisting(dto.MerchantId, dto.MerchantPassword)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("merchant account does not exist")
	}

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
