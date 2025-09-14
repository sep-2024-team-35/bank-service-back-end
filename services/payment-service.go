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
	CreateRequest(dto dto.PaymentRequestDTO) (*models.PaymentRequest, error)
	Pay(cardDetailsDTO dto.CardDetailsDTO, paymentRequestID string) (*models.Transaction, error)
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

func (s *paymentService) CreateRequest(requestDto dto.PaymentRequestDTO) (*models.PaymentRequest, error) {
	if err := s.validateMerchant(requestDto.MerchantID, requestDto.MerchantPassword); err != nil {
		return nil, err
	}

	request := s.mapToPaymentRequest(requestDto)

	var savedRequest *models.PaymentRequest

	err := s.paymentRepository.DB().Transaction(func(tx *gorm.DB) error {
		var err error

		savedRequest, err = s.paymentRepository.SaveRequestTransactional(tx, request)
		if err != nil {
			return err
		}

		transaction := s.buildTransaction(savedRequest)

		if _, err = s.transactionRepository.SaveTransactional(tx, transaction); err != nil {
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

func (s *paymentService) Pay(cardDetailsDTO dto.CardDetailsDTO, paymentRequestID string) (*models.Transaction, error) {
	issuerAccount, err := s.accountService.FindAccountByPAN(cardDetailsDTO.PrimaryAccountNumber)
	if err != nil {
		return nil, err
	}

	paymentRequest, err := s.paymentRepository.GetByID(paymentRequestID)
	if err != nil {
		return nil, err
	}

	transaction, err := s.transactionRepository.FindByPaymentRequestID(paymentRequest.ID)
	if err != nil {
		return nil, err
	}

	if issuerAccount != nil {
		if err := s.processIntrabankTransaction(issuerAccount, paymentRequest, transaction); err != nil {
			return nil, err
		}
	} else {
		if err := s.processInterbankTransaction(transaction); err != nil {
			return nil, err
		}
	}

	// 5. Vraćamo transakciju (posle obrade)
	return transaction, nil
}

// Same banks
func (s *paymentService) processIntrabankTransaction(
	issuerAccount *models.Account,
	paymentRequest *models.PaymentRequest,
	transaction *models.Transaction,
) error {
	// 1. Business errors up front

	// 1a. Merchant lookup
	acquirerAccount, err := s.accountService.GetMerchantAccount(
		paymentRequest.MerchantID,
		paymentRequest.MerchantPassword,
	)
	if err != nil {
		transaction.Status = "FAILED"
		// non-transactional update to persist status immediately
		if _, updErr := s.transactionRepository.Update(transaction); updErr != nil {
			// log updErr if you need observability
		}
		return err
	}

	// 1b. Insufficient funds
	if !issuerAccount.Balance.GreaterThanOrEqual(transaction.Amount) {
		transaction.Status = "FAILED"
		if _, updErr := s.transactionRepository.Update(transaction); updErr != nil {
			// log updErr
		}
		return errors.New("insufficient funds on issuer account")
	}

	// 2. All business checks passed → do atomic update of accounts + transaction
	err = s.accountService.DB().Transaction(func(tx *gorm.DB) error {
		// debit/credit
		issuerAccount.Balance = issuerAccount.Balance.Sub(transaction.Amount)
		acquirerAccount.Balance = acquirerAccount.Balance.Add(transaction.Amount)

		// finalize transaction record
		transaction.Status = "COMPLETED"
		transaction.IssuerOrderID = generateRandomOrderID()
		transaction.IssuerTimestamp = time.Now().Format(time.RFC3339)

		// persist inside tx
		if _, e := s.accountService.UpdateTransactional(tx, issuerAccount); e != nil {
			return e
		}
		if _, e := s.accountService.UpdateTransactional(tx, acquirerAccount); e != nil {
			return e
		}
		if _, e := s.transactionRepository.UpdateTransactional(tx, transaction); e != nil {
			return e
		}

		return nil
	})

	// 3. If something went wrong in the DB transaction, mark it ERROR
	if err != nil {
		transaction.Status = "ERROR"
		if _, updErr := s.transactionRepository.Update(transaction); updErr != nil {
			// log updErr
		}
		return err
	}

	// 4. Success
	return nil
}

//func (s *paymentService) processIntrabankTransaction(issuerAccount *models.Account, paymentRequest *models.PaymentRequest, transaction *models.Transaction) error {
//	acquirerAccount, err := s.accountService.GetMerchantAccount(paymentRequest.MerchantID, paymentRequest.MerchantPassword)
//	if err != nil {
//		return err
//	}
//
//	if !issuerAccount.Balance.GreaterThanOrEqual(transaction.Amount) {
//		return errors.New("insufficient funds on issuer account")
//	}
//
//	issuerAccount.Balance = issuerAccount.Balance.Sub(transaction.Amount)
//
//	// 4. Uvećaj balans merchant account-a
//	acquirerAccount.Balance = acquirerAccount.Balance.Add(transaction.Amount)
//
//	// TODO: implement method
//}

// Different banks
func (s *paymentService) processInterbankTransaction(transaction *models.Transaction) error {
	// TODO: implement method

	panic("implement me")
}

func (s *paymentService) validateMerchant(merchantID, merchantPassword string) error {
	exists, err := s.accountService.isMerchantAccountExisting(merchantID, merchantPassword)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("merchant account does not exist")
	}
	return nil
}

func (s *paymentService) mapToPaymentRequest(dto dto.PaymentRequestDTO) *models.PaymentRequest {
	return &models.PaymentRequest{
		ID:                uuid.New(),
		MerchantID:        dto.MerchantID,
		MerchantPassword:  dto.MerchantPassword,
		Amount:            dto.Amount,
		MerchantOrderID:   dto.MerchantOrderId,
		MerchantTimestamp: dto.MerchantTimestamp,
		SuccessURL:        dto.SuccessUrl,
		FailedURL:         dto.FailedUrl,
		ErrorURL:          dto.ErrorUrl,
	}
}

func (s *paymentService) buildTransaction(request *models.PaymentRequest) *models.Transaction {
	return &models.Transaction{
		Amount:            request.Amount,
		MerchantOrderID:   request.MerchantOrderID,
		MerchantTimestamp: request.MerchantTimestamp.Format(time.RFC3339),
		AcquirerOrderID:   generateRandomOrderID(),
		AcquirerTimestamp: time.Now().Format(time.RFC3339),
		IssuerOrderID:     "", // TODO: implement
		IssuerTimestamp:   "", // TODO: implement
		Status:            "CREATED",
		PaymentRequestID:  request.ID,
	}
}
