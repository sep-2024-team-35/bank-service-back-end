package services

import (
	"github.com/google/uuid"
	"github.com/sep-2024-team-35/bank-servce-back-end/dto"
	"github.com/sep-2024-team-35/bank-servce-back-end/models"
	"github.com/sep-2024-team-35/bank-servce-back-end/repositories"
)

type AccountService interface {
	GetMerchantAccount(req *models.PaymentRequest) (*models.Account, error)
	RegisterNewMerchant(registrationDto *dto.MerchantRegistrationDto) (*models.Account, error)
}

type accountService struct {
	accountRepo repositories.AccountRepository
}

func NewAccountService(repo repositories.AccountRepository) AccountService {
	return &accountService{
		accountRepo: repo,
	}
}

func (s *accountService) GetMerchantAccount(req *models.PaymentRequest) (*models.Account, error) {
	return s.accountRepo.FindByMerchantIDAndPassword(req.MerchantID, req.MerchantPassword)
}

func (s *accountService) RegisterNewMerchant(registrationDto *dto.MerchantRegistrationDto) (*models.Account, error) {
	account := &models.Account{
		ID:               uuid.New(),
		CardHolderName:   registrationDto.AccountHolderName,
		MerchantID:       registrationDto.MerchantID,
		MerchantPassword: registrationDto.MerchantPassword,
		MerchantAccount:  true,
		Balance:          models.NewZeroBalance(),
	}

	return s.accountRepo.Save(account)
}
