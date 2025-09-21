package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sep-2024-team-35/bank-servce-back-end/dto"
	"github.com/sep-2024-team-35/bank-servce-back-end/models"
	"github.com/sep-2024-team-35/bank-servce-back-end/repositories"
	"github.com/sep-2024-team-35/bank-servce-back-end/utils"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"time"
)

type PaymentService interface {
	CreateRequest(dto dto.PaymentRequestDTO) (*models.PaymentRequest, error)
	Pay(cardDetailsDTO dto.CardDetailsDTO, paymentRequestID string) (*dto.PSPResponseDTO, error)
	ExternalPay(externalTransaction dto.ExternalTransactionRequestDTO) (*dto.PCCResponseDTO, error)
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

	request, err := s.mapToPaymentRequest(requestDto)
	if err != nil {
		log.Printf("[ERROR] Failed to map PaymentRequestDTO to PaymentRequest: %v", err)
		return nil, fmt.Errorf("failed to map payment request: %w", err)
	}

	var savedRequest *models.PaymentRequest

	err = s.paymentRepository.DB().Transaction(func(tx *gorm.DB) error {
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

func (s *paymentService) Pay(
	cardDetailsDTO dto.CardDetailsDTO,
	paymentRequestID string,
) (*dto.PSPResponseDTO, error) {
	issuerAccount, err := s.accountService.FindAccountByPAN(cardDetailsDTO.PrimaryAccountNumber)

	paymentRequest, err := s.paymentRepository.GetByID(paymentRequestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment request: %w", err)
	}

	transaction, err := s.transactionRepository.FindByPaymentRequestID(paymentRequest.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find transaction: %w", err)
	}

	if transaction == nil {
		return nil, fmt.Errorf("transaction not found")
	}

	var pspResponse *dto.PSPResponseDTO
	if issuerAccount != nil {
		if issuerAccount.CCV != cardDetailsDTO.SecurityCode {
			return nil, fmt.Errorf("invalid security code (CCV)")
		}

		expTime, err := time.Parse("01/06", cardDetailsDTO.ExpirationDate)
		if err != nil {
			return nil, fmt.Errorf("invalid expiration date format, expected MM/YY")
		}

		lastDay := time.Date(expTime.Year(), expTime.Month()+1, 0, 23, 59, 59, 0, time.UTC)
		if time.Now().After(lastDay) {
			return nil, fmt.Errorf("card expired")
		}

		pspResponse, err = s.processIntrabankTransaction(issuerAccount, paymentRequest, transaction)
		if err != nil {
			return nil, fmt.Errorf("intrabank processing failed: %w", err)
		}
	} else {
		pspResponse, err = s.processInterbankTransaction(cardDetailsDTO, transaction, paymentRequest)
		if err != nil {
			return nil, fmt.Errorf("interbank processing failed: %w", err)
		}
	}

	switch transaction.Status {
	case "SUCCESS":
		pspResponse.RedirectURL = paymentRequest.SuccessURL
	case "FAILED":
		pspResponse.RedirectURL = paymentRequest.FailedURL
	case "ERROR":
		pspResponse.RedirectURL = paymentRequest.ErrorURL
	default:
		pspResponse.RedirectURL = paymentRequest.ErrorURL
	}

	return pspResponse, nil
}

func (s *paymentService) processIntrabankTransaction(
	issuerAccount *models.Account,
	paymentRequest *models.PaymentRequest,
	transaction *models.Transaction,
) (*dto.PSPResponseDTO, error) {
	acquirerAccount, err := s.accountService.GetMerchantAccount(
		paymentRequest.MerchantID,
		paymentRequest.MerchantPassword,
	)

	if err != nil {
		transaction.Status = "FAILED"
		if _, updErr := s.transactionRepository.Update(transaction); updErr != nil {
			log.Printf("[ERROR] Failed to update transaction after merchant lookup failure: %v", updErr)
		}
		return nil, err
	}

	if !issuerAccount.Balance.GreaterThanOrEqual(transaction.Amount) {
		transaction.Status = "FAILED"
		if _, updErr := s.transactionRepository.Update(transaction); updErr != nil {
			log.Printf("[ERROR] Failed to update transaction after insufficient funds: %v", updErr)
		}
		return nil, errors.New("insufficient funds on issuer account")
	}

	err = s.accountService.DB().Transaction(func(tx *gorm.DB) error {
		issuerAccount.Balance = issuerAccount.Balance.Sub(transaction.Amount)
		acquirerAccount.Balance = acquirerAccount.Balance.Add(transaction.Amount)

		transaction.Status = "COMPLETED"
		transaction.IssuerOrderID = generateRandomOrderID()
		transaction.IssuerTimestamp = time.Now().Format(time.RFC3339)

		if _, e := s.accountService.UpdateTransactional(tx, issuerAccount); e != nil {
			return fmt.Errorf("failed to update issuer account: %w", e)
		}
		if _, e := s.accountService.UpdateTransactional(tx, acquirerAccount); e != nil {
			return fmt.Errorf("failed to update acquirer account: %w", e)
		}
		if _, e := s.transactionRepository.UpdateTransactional(tx, transaction); e != nil {
			return fmt.Errorf("failed to update transaction: %w", e)
		}

		return nil
	})

	if err != nil {
		transaction.Status = "ERROR"
		if _, updErr := s.transactionRepository.Update(transaction); updErr != nil {
			log.Printf("[ERROR] Failed to update transaction after DB error: %v", updErr)
		}
		return nil, err
	}

	pspResponse := &dto.PSPResponseDTO{
		Status:            transaction.Status,
		RedirectURL:       "", // TODO add redirect url from request
		AcquirerOrderID:   transaction.AcquirerOrderID,
		AcquirerTimeStamp: transaction.AcquirerTimestamp,
		PaymentID:         fmt.Sprintf("%d", transaction.ID),
		MerchantOrderID:   transaction.MerchantOrderID,
	}

	return pspResponse, nil
}

func (s *paymentService) processInterbankTransaction(
	cardDetails dto.CardDetailsDTO,
	transaction *models.Transaction,
	paymentRequest *models.PaymentRequest,
) (*dto.PSPResponseDTO, error) {
	if cardDetails.PrimaryAccountNumber == "" {
		transaction.Status = "ERROR"
		if _, e := s.transactionRepository.Update(transaction); e != nil {
			maskedPAN := maskPAN(cardDetails.PrimaryAccountNumber)
			return nil, fmt.Errorf("failed to update transaction after empty PAN: %s", maskedPAN)
		}
		return nil, fmt.Errorf("primary account number is empty")
	}

	request := dto.ExternalTransactionRequestDTO{
		ID:                   transaction.ID,
		AcquirerOrderID:      transaction.AcquirerOrderID,
		AcquirerTimestamp:    transaction.AcquirerTimestamp,
		Amount:               transaction.Amount,
		MerchantOrderID:      transaction.MerchantOrderID,
		MerchantTimestamp:    transaction.MerchantTimestamp,
		Currency:             transaction.Currency,
		PrimaryAccountNumber: cardDetails.PrimaryAccountNumber,
		CardHolderName:       cardDetails.CardHolderName,
		ExpirationDate:       cardDetails.ExpirationDate,
		SecurityCode:         cardDetails.SecurityCode,
		PaymentRequestID:     cardDetails.PaymentRequestID,
	}

	payload, err := json.Marshal(request)
	if err != nil {
		transaction.Status = "ERROR"
		if _, e := s.transactionRepository.Update(transaction); e != nil {
			return nil, fmt.Errorf("failed to update transaction after marshal error: %w", e)
		}
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post("https://pcc-provider-api.azurewebsites.net/api/transactions", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		transaction.Status = "ERROR"
		if _, e := s.transactionRepository.Update(transaction); e != nil {
			return nil, fmt.Errorf("failed to update transaction after PCC post error: %w", e)
		}
		return nil, fmt.Errorf("failed to contact PCC: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("[ERROR] Failed to close response body: %v", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		transaction.Status = "FAILED"
		if _, e := s.transactionRepository.Update(transaction); e != nil {
			return nil, fmt.Errorf("failed to update transaction after PCC status error: %w", e)
		}
		return nil, fmt.Errorf("PCC responded with status %d", resp.StatusCode)
	}

	var pccResp dto.PCCResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&pccResp); err != nil {
		transaction.Status = "ERROR"
		if _, e := s.transactionRepository.Update(transaction); e != nil {
			return nil, fmt.Errorf("failed to update transaction after decode error: %w", e)
		}
		return nil, fmt.Errorf("failed to decode PCC response: %w", err)
	}

	if pccResp.Status == "" || pccResp.AcquirerOrderID == "" {
		transaction.Status = "ERROR"
		if _, e := s.transactionRepository.Update(transaction); e != nil {
			return nil, fmt.Errorf("failed to update transaction after invalid PCC response: %w", e)
		}
		return nil, fmt.Errorf("invalid PCC response: %+v", pccResp)
	}

	transaction.IssuerOrderID = pccResp.IssuerOrderID
	transaction.IssuerTimestamp = pccResp.IssuerTimestamp
	transaction.Status = pccResp.Status

	if pccResp.Status == "SUCCESS" {
		acquirerAccount, err := s.accountService.GetMerchantAccount(paymentRequest.MerchantID, paymentRequest.MerchantPassword)
		if err != nil {
			transaction.Status = "FAILED"
			if _, e := s.transactionRepository.Update(transaction); e != nil {
				return nil, fmt.Errorf("failed to update transaction after getting merchant account: %w", e)
			}
			return nil, fmt.Errorf("failed to get merchant account: %w", err)
		}

		acquirerAccount.Balance = acquirerAccount.Balance.Add(transaction.Amount)

		err = s.accountService.DB().Transaction(func(tx *gorm.DB) error {
			if _, e := s.accountService.UpdateTransactional(tx, acquirerAccount); e != nil {
				return fmt.Errorf("failed to update acquirer account: %w", e)
			}
			if _, e := s.transactionRepository.UpdateTransactional(tx, transaction); e != nil {
				return fmt.Errorf("failed to update transaction: %w", e)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("transactional update failed: %w", err)
		}
	} else {
		if _, e := s.transactionRepository.Update(transaction); e != nil {
			return nil, fmt.Errorf("failed to update transaction for non-success PCC: %w", e)
		}
	}

	pspResponse := &dto.PSPResponseDTO{
		Status:            transaction.Status,
		RedirectURL:       "",
		AcquirerOrderID:   transaction.AcquirerOrderID,
		AcquirerTimeStamp: transaction.AcquirerTimestamp,
		PaymentID:         fmt.Sprintf("%d", transaction.ID),
		MerchantOrderID:   transaction.MerchantOrderID,
	}

	return pspResponse, nil
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

func (s *paymentService) mapToPaymentRequest(dto dto.PaymentRequestDTO) (*models.PaymentRequest, error) {
	merchantTime, err := utils.ParseMerchantTimestamp(dto.MerchantTimestamp)
	if err != nil {
		return nil, fmt.Errorf("invalid merchant timestamp: %w", err)
	}

	return &models.PaymentRequest{
		ID:                uuid.New(),
		MerchantID:        dto.MerchantID,
		MerchantPassword:  dto.MerchantPassword,
		Amount:            dto.Amount,
		MerchantOrderID:   dto.MerchantOrderId,
		MerchantTimestamp: merchantTime,
		SuccessURL:        dto.SuccessUrl,
		FailedURL:         dto.FailedUrl,
		ErrorURL:          dto.ErrorUrl,
	}, nil
}

func (s *paymentService) buildTransaction(request *models.PaymentRequest) *models.Transaction {
	return &models.Transaction{
		ID:                uuid.New(),
		Amount:            request.Amount,
		Currency:          "RSD",
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

func (s *paymentService) ExternalPay(dto dto.ExternalTransactionRequestDTO) (*dto.PCCResponseDTO, error) {
	var transaction *models.Transaction

	issuerAccount, err := s.accountService.FindAccountByPAN(dto.PrimaryAccountNumber)
	if err != nil {
		maskedPAN := maskPAN(dto.PrimaryAccountNumber)
		log.Printf("[ERROR] Issuer account not found for PAN=%s", maskedPAN)
		transaction = buildExternalTransaction(dto, "FAILED")
		_, _ = s.transactionRepository.Save(transaction)
		return buildPCCResponse(transaction), fmt.Errorf("issuer account not found")
	}

	if issuerAccount.CCV != dto.SecurityCode {
		return nil, fmt.Errorf("invalid security code (CCV)")
	}

	expTime, err := time.Parse("01/06", dto.ExpirationDate)
	if err != nil {
		return nil, fmt.Errorf("invalid expiration date format, expected MM/YY")
	}

	lastDay := time.Date(expTime.Year(), expTime.Month()+1, 0, 23, 59, 59, 0, time.UTC)
	if time.Now().After(lastDay) {
		return nil, fmt.Errorf("card expired")
	}

	if !issuerAccount.Balance.GreaterThanOrEqual(dto.Amount) {
		log.Printf("[ERROR] Insufficient funds: Balance=%.2f, Required=%.2f",
			issuerAccount.Balance.InexactFloat64(), dto.Amount.InexactFloat64())
		transaction = buildExternalTransaction(dto, "FAILED")
		_, _ = s.transactionRepository.Save(transaction)
		return buildPCCResponse(transaction), fmt.Errorf("insufficient funds")
	}

	transaction = buildExternalTransaction(dto, "SUCCESS")
	transaction.IssuerOrderID = generateRandomOrderID()
	transaction.IssuerTimestamp = time.Now().Format(time.RFC3339)

	err = s.accountService.DB().Transaction(func(tx *gorm.DB) error {
		issuerAccount.Balance = issuerAccount.Balance.Sub(transaction.Amount)

		if _, err := s.accountService.UpdateTransactional(tx, issuerAccount); err != nil {
			return err
		}
		if _, err := s.transactionRepository.SaveTransactional(tx, transaction); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Printf("[ERROR] DB transaction failed: %v", err)
		transaction.Status = "ERROR"
		_, _ = s.transactionRepository.Update(transaction)
		return buildPCCResponse(transaction), err
	}

	log.Printf("[INFO] External transaction SUCCESS: OrderID=%s", transaction.IssuerOrderID)
	return buildPCCResponse(transaction), nil
}

func buildExternalTransaction(dto dto.ExternalTransactionRequestDTO, status string) *models.Transaction {
	return &models.Transaction{
		ID:                dto.ID,
		Amount:            dto.Amount,
		Currency:          dto.Currency,
		MerchantOrderID:   dto.MerchantOrderID,
		MerchantTimestamp: dto.MerchantTimestamp,
		AcquirerOrderID:   dto.AcquirerOrderID,
		AcquirerTimestamp: dto.AcquirerTimestamp,
		Status:            status,
	}
}

func buildPCCResponse(tx *models.Transaction) *dto.PCCResponseDTO {
	if tx.Status != "SUCCESS" {
		return &dto.PCCResponseDTO{
			Status:            tx.Status,
			AcquirerOrderID:   tx.AcquirerOrderID,
			AcquirerTimestamp: tx.AcquirerTimestamp,
			IssuerOrderID:     "",
			IssuerTimestamp:   "",
		}
	}

	return &dto.PCCResponseDTO{
		Status:            tx.Status,
		AcquirerOrderID:   tx.AcquirerOrderID,
		AcquirerTimestamp: tx.AcquirerTimestamp,
		IssuerOrderID:     tx.IssuerOrderID,
		IssuerTimestamp:   tx.IssuerTimestamp,
	}
}

func maskPAN(pan string) string {
	if len(pan) <= 4 {
		return "****"
	}
	return strings.Repeat("*", len(pan)-4) + pan[len(pan)-4:]
}
