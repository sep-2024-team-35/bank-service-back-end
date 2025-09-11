package services

import (
	"errors"
	"log"

	"github.com/sep-2024-team-35/bank-servce-back-end/dto"
	"github.com/sep-2024-team-35/bank-servce-back-end/repositories"
)

type ExecutePaymentService struct {
	paymentRepo        repositories.PaymentRepository
	accountRepo        repositories.AccountRepository
	bankService        *BankService
	transactionService *TransactionService
	// pspNotificationSvc PspNotificationService
}

func NewExecutePaymentService(
	paymentRepo repositories.PaymentRepository,
	accountRepo repositories.AccountRepository,
	bankSvc *BankService,
	transactionSvc *TransactionService,
	// pspNotificationSvc PspNotificationService,
) *ExecutePaymentService {
	return &ExecutePaymentService{
		paymentRepo:        paymentRepo,
		accountRepo:        accountRepo,
		bankService:        bankSvc,
		transactionService: transactionSvc,
		// pspNotificationSvc: pspNotificationSvc,
	}
}

// Glavna funkcija za izvršenje plaćanja
func (s *ExecutePaymentService) ExecutePayment(cardDetails dto.CardDetailsDTO) (bool, error) {
	// 1. Učitaj PaymentRequest
	paymentRequest, err := s.paymentRepo.GetByID(cardDetails.PaymentRequestID)
	if err != nil || paymentRequest == nil {
		return false, errors.New("payment request not found")
	}

	// 2. Učitaj merchant account
	merchantAccount, err := s.accountRepo.FindByMerchantID(paymentRequest.MerchantID)
	if err != nil || merchantAccount == nil {
		return false, errors.New("merchant account not found")
	}

	// 3. Učitaj ili kreiraj transaction
	transaction, err := s.transactionService.GetTransactionByPaymentRequestID(paymentRequest.ID)
	if err != nil {
		return false, err
	}
	if transaction != nil && transaction.Status == "SUCCESS" {
		// s.emitAlreadyPaidEvent(transaction)
		return false, nil
	}

	// 4. Validacija kartice (datum isteka)
	if !IsValidExpirationDate(cardDetails.ExpirationDate) {
		// s.emitFailedEvent(transaction)
		return false, nil
	}

	// 5. Provera PAN prefix-a
	if s.checkIfPanBelongsToAcquirerBank(cardDetails.PrimaryAccountNumber) {
		// Internal transaction
		issuerAccount, _ := s.accountRepo.FindByPAN(cardDetails.PrimaryAccountNumber)
		if issuerAccount == nil {
			// s.emitFailedEvent(transaction)
			return false, nil
		}

		if issuerAccount.Balance.GreaterThanOrEqual(paymentRequest.Amount) {
			issuerAccount.Balance = issuerAccount.Balance.Sub(paymentRequest.Amount)
			merchantAccount.Balance = merchantAccount.Balance.Add(paymentRequest.Amount)
			if _, err := s.accountRepo.Update(issuerAccount); err != nil {
				log.Printf("Failed to update issuer account: %v", err)
			}

			if _, err := s.accountRepo.Update(merchantAccount); err != nil {
				log.Printf("Failed to update merchant account: %v", err)
			}

			if _, err := s.transactionService.SuccessTransaction(transaction); err != nil {
				log.Printf("Failed to mark transaction as SUCCESS: %v", err)
			}
			// s.emitSuccessEvent(transaction)
			return true, nil
		} else {
			if _, err := s.transactionService.FailTransaction(transaction); err != nil {
				log.Printf("Failed to mark transaction as FAILED: %v", err)
			}
			// s.emitFailedEvent(transaction)
			return false, nil
		}
	} else {
		// External PCC call - za sada samo mock
		log.Println("Would call external PCC service here")
		if _, err := s.transactionService.ErrorTransaction(transaction); err != nil {
			log.Printf("Failed to mark transaction as ERROR: %v", err)
		}
		// s.emitErrorEvent(transaction)
		return false, nil
	}
}

// Helper funkcije
func (s *ExecutePaymentService) checkIfPanBelongsToAcquirerBank(pan string) bool {
	bic := s.bankService.GetBIC()
	return len(pan) >= 4 && pan[:4] == bic
}

/*
// Event emitter funkcije
func (s *ExecutePaymentService) emitAlreadyPaidEvent(tx *models.Transaction) {
	log.Printf("Transaction %d already paid\n", tx.ID)
	s.pspNotificationSvc.SendTransactionResult(tx)
}

func (s *ExecutePaymentService) emitSuccessEvent(tx *models.Transaction) {
	log.Printf("Transaction %d successful\n", tx.ID)
	s.pspNotificationSvc.SendTransactionResult(tx)
}

func (s *ExecutePaymentService) emitFailedEvent(tx *models.Transaction) {
	log.Printf("Transaction %d failed\n", tx.ID)
	s.pspNotificationSvc.SendTransactionResult(tx)
}

func (s *ExecutePaymentService) emitErrorEvent(tx *models.Transaction) {
	log.Printf("Transaction %d error\n", tx.ID)
	s.pspNotificationSvc.SendTransactionResult(tx)
}
*/

func IsValidExpirationDate(exp string) bool {
	if len(exp) != 5 || exp[2] != '/' {
		return false
	}
	return true
}
