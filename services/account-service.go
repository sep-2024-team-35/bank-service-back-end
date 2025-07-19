package services

import (
	"errors"
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

func (s *accountService) RegisterNewMerchant(regDto *dto.MerchantRegistrationDto) (*models.Account, error) {
	if err := validateMerchantRegistrationInput(regDto); err != nil {
		return nil, err
	}

	exists, err := s.isMerchantAccountExisting(regDto.MerchantID, regDto.MerchantPassword)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("merchant account already exists")
	}

	account := createMerchantAccountFromDto(regDto)

	return s.accountRepo.Save(account)
}

func validateMerchantRegistrationInput(regDto *dto.MerchantRegistrationDto) error {
	if regDto.AccountHolderName == "" || regDto.MerchantID == "" || regDto.MerchantPassword == "" {
		return errors.New("all fields are required")
	}
	return nil
}

func (s *accountService) isMerchantAccountExisting(merchantID, merchantPassword string) (bool, error) {
	account, err := s.accountRepo.FindByMerchantIDAndPassword(merchantID, merchantPassword)
	if err != nil {
		return false, err
	}
	return account != nil, nil
}

func createMerchantAccountFromDto(regDto *dto.MerchantRegistrationDto) *models.Account {
	return &models.Account{
		ID:               uuid.New(),
		CardHolderName:   regDto.AccountHolderName,
		MerchantID:       regDto.MerchantID,
		MerchantPassword: regDto.MerchantPassword,
		MerchantAccount:  true,
		Balance:          models.NewZeroBalance(),
	}
}
