package services

import (
	"errors"
	"github.com/google/uuid"
	"github.com/sep-2024-team-35/bank-servce-back-end/dto"
	"github.com/sep-2024-team-35/bank-servce-back-end/models"
	"github.com/sep-2024-team-35/bank-servce-back-end/repositories"
	"gorm.io/gorm"
	"log"
	"time"
)

type PaymentService interface {
	CreateRequest(dto dto.PaymentRequestDto) (*models.PaymentRequest, error)
}

type paymentService struct {
	paymentRepository     repositories.PaymentRepository
	transactionRepository repositories.TransactionRepository
	accountService        AccountService
}

func NewPaymentService(repo repositories.PaymentRepository, accSvc AccountService, transactionRepository repositories.TransactionRepository) PaymentService {
	return &paymentService{
		paymentRepository:     repo,
		accountService:        accSvc,
		transactionRepository: transactionRepository,
	}
}

func (s *paymentService) CreateRequest(requestDto dto.PaymentRequestDto) (*models.PaymentRequest, error) {
	exists, err := s.accountService.isMerchantAccountExisting(requestDto.MerchantID, requestDto.MerchantPassword)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("merchant account does not exist")
	}

	request := &models.PaymentRequest{
		ID:                uuid.New(),
		MerchantID:        requestDto.MerchantID,
		MerchantPassword:  requestDto.MerchantPassword,
		Amount:            requestDto.Amount,
		MerchantOrderID:   requestDto.MerchantOrderId,
		MerchantTimestamp: requestDto.MerchantTimestamp,
		SuccessURL:        requestDto.SuccessUrl,
		FailedURL:         requestDto.FailedUrl,
		ErrorURL:          requestDto.ErrorUrl,
	}

	var savedRequest *models.PaymentRequest

	err = s.paymentRepository.DB().Transaction(func(tx *gorm.DB) error {
		var err error

		savedRequest, err = s.paymentRepository.SaveRequestTransactional(tx, request)
		if err != nil {
			return err
		}

		newTransaction := &models.Transaction{
			Amount:            savedRequest.Amount.InexactFloat64(),
			MerchantOrderID:   savedRequest.MerchantOrderID,
			MerchantTimestamp: savedRequest.MerchantTimestamp.Format(time.RFC3339),
			AcquirerOrderID:   generateRandomOrderID(),
			AcquirerTimestamp: time.Now().Format(time.RFC3339),
			IssuerOrderID:     "", // TODO: implement
			IssuerTimestamp:   "", // TODO: implement
			Status:            "CREATED",
			PaymentRequestID:  request.ID,
		}

		_, err = s.transactionRepository.SaveTransactional(tx, newTransaction)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("[ERROR] Failed to create payment request and transaction for MerchantID=%s: %v", requestDto.MerchantID, err)
		return nil, err
	}

	log.Printf("[INFO] Payment request and transaction successfully created: ID=%s", savedRequest.ID.String())
	return savedRequest, nil
}
