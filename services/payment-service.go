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
	"gorm.io/gorm"
	"log"
	"net/http"
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

func (s *paymentService) Pay(
	cardDetailsDTO dto.CardDetailsDTO,
	paymentRequestID string,
) (*dto.PSPResponseDTO, error) {

	log.Printf("[DEBUG] Request PAN: %s", cardDetailsDTO.PrimaryAccountNumber)

	// 1. Pronađi issuer account po PAN
	issuerAccount, err := s.accountService.FindAccountByPAN(cardDetailsDTO.PrimaryAccountNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to find issuer account: %w", err)
	}

	// 2. Pronađi payment request
	paymentRequest, err := s.paymentRepository.GetByID(paymentRequestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment request: %w", err)
	}

	// 3. Pronađi postojeću transakciju
	transaction, err := s.transactionRepository.FindByPaymentRequestID(paymentRequest.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find transaction: %w", err)
	}

	if transaction == nil {
		return nil, fmt.Errorf("transaction not found")
	}

	// 4. Obradi transakciju
	var pspResponse *dto.PSPResponseDTO
	if issuerAccount != nil {
		// intrabank
		pspResponse, err = s.processIntrabankTransaction(issuerAccount, paymentRequest, transaction)
		if err != nil {
			return nil, fmt.Errorf("intrabank processing failed: %w", err)
		}
	} else {
		// interbank
		pspResponse, err = s.processInterbankTransaction(cardDetailsDTO, transaction, paymentRequest)
		if err != nil {
			return nil, fmt.Errorf("interbank processing failed: %w", err)
		}
	}

	// 5. Vraćamo PSPResponseDTO
	return pspResponse, nil
}

// Same banks
func (s *paymentService) processIntrabankTransaction(
	issuerAccount *models.Account,
	paymentRequest *models.PaymentRequest,
	transaction *models.Transaction,
) (*dto.PSPResponseDTO, error) {

	// 1. Business errors up front

	// 1a. Merchant lookup
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

	// 1b. Insufficient funds
	if !issuerAccount.Balance.GreaterThanOrEqual(transaction.Amount) {
		transaction.Status = "FAILED"
		if _, updErr := s.transactionRepository.Update(transaction); updErr != nil {
			log.Printf("[ERROR] Failed to update transaction after insufficient funds: %v", updErr)
		}
		return nil, errors.New("insufficient funds on issuer account")
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

	// 3. If something went wrong in the DB transaction, mark it ERROR
	if err != nil {
		transaction.Status = "ERROR"
		if _, updErr := s.transactionRepository.Update(transaction); updErr != nil {
			log.Printf("[ERROR] Failed to update transaction after DB error: %v", updErr)
		}
		return nil, err
	}

	// 4. Success → kreiraj PSPResponseDTO
	pspResponse := &dto.PSPResponseDTO{
		Status:            transaction.Status,
		RedirectURL:       "", // po potrebi dodaj redirect URL
		AcquirerOrderID:   transaction.AcquirerOrderID,
		AcquirerTimeStamp: transaction.AcquirerTimestamp,
		PaymentID:         fmt.Sprintf("%d", transaction.ID), // ili UUID ako koristiš
	}

	return pspResponse, nil
}

// Different banks
func (s *paymentService) processInterbankTransaction(
	cardDetails dto.CardDetailsDTO,
	transaction *models.Transaction,
	paymentRequest *models.PaymentRequest,
) (*dto.PSPResponseDTO, error) {

	log.Printf("[DEBUG] Starting processInterbankTransaction for transaction ID: %d", transaction.ID)

	if cardDetails.PrimaryAccountNumber == "" {
		transaction.Status = "ERROR"
		if _, e := s.transactionRepository.Update(transaction); e != nil {
			return nil, fmt.Errorf("failed to update transaction after empty PAN: %w", e)
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
	resp, err := client.Post("http://pcc-container:8081/api/transactions", "application/json", bytes.NewBuffer(payload))
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

	// setuj polja transakcije
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
		// FAILED ili ERROR
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
	log.Printf("[DEBUG] ExternalPay request: PAN=%s, Amount=%.2f", dto.PrimaryAccountNumber, dto.Amount.InexactFloat64())

	var transaction *models.Transaction

	// 1. Pronađi issuer nalog
	issuerAccount, err := s.accountService.FindAccountByPAN(dto.PrimaryAccountNumber)
	if err != nil {
		log.Printf("[ERROR] Issuer account not found for PAN=%s", dto.PrimaryAccountNumber)
		transaction = buildExternalTransaction(dto, "FAILED")
		_, _ = s.transactionRepository.Save(transaction)
		return buildPCCResponse(transaction), fmt.Errorf("issuer account not found")
	}

	// 2. Proveri stanje
	if !issuerAccount.Balance.GreaterThanOrEqual(dto.Amount) {
		log.Printf("[ERROR] Insufficient funds: Balance=%.2f, Required=%.2f",
			issuerAccount.Balance.InexactFloat64(), dto.Amount.InexactFloat64())
		transaction = buildExternalTransaction(dto, "FAILED")
		_, _ = s.transactionRepository.Save(transaction)
		return buildPCCResponse(transaction), fmt.Errorf("insufficient funds")
	}

	// 3. SUCCESS → atomicno update stanja i transakcije
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
		//PaymentRequestID:  dto.PaymentRequestID,
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
