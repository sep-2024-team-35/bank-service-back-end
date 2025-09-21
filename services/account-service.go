package services

import (
	"errors"
	"github.com/google/uuid"
	"github.com/sep-2024-team-35/bank-servce-back-end/dto"
	"github.com/sep-2024-team-35/bank-servce-back-end/models"
	"github.com/sep-2024-team-35/bank-servce-back-end/repositories"
	"gorm.io/gorm"
)

type AccountService interface {
	UpdateTransactional(tx *gorm.DB, account *models.Account) (*models.Account, error)
	GetMerchantAccount(ID string, password string) (*models.Account, error)
	RegisterNewMerchant(registrationDto *dto.MerchantRegistrationDTO) (*models.Account, error)
	isMerchantAccountExisting(merchantID, merchantPassword string) (bool, error)
	FindAccountByPAN(pan string) (*models.Account, error)
	DB() *gorm.DB
}

type accountService struct {
	accountRepo repositories.AccountRepository
}

func (s *accountService) UpdateTransactional(tx *gorm.DB, account *models.Account) (*models.Account, error) {
	if tx == nil {
		return nil, errors.New("transaction object cannot be nil")
	}
	return s.accountRepo.UpdateTransactional(tx, account)
}

func (s *accountService) FindAccountByPAN(pan string) (*models.Account, error) {
	return s.accountRepo.FindByPAN(pan)
}

func NewAccountService(repo repositories.AccountRepository) AccountService {
	return &accountService{
		accountRepo: repo,
	}
}

func (s *accountService) GetMerchantAccount(ID string, password string) (*models.Account, error) {
	return s.accountRepo.FindByMerchantIDAndPassword(ID, password)
}

func (s *accountService) RegisterNewMerchant(regDto *dto.MerchantRegistrationDTO) (*models.Account, error) {
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

func validateMerchantRegistrationInput(regDto *dto.MerchantRegistrationDTO) error {
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

func createMerchantAccountFromDto(regDto *dto.MerchantRegistrationDTO) *models.Account {
	return &models.Account{
		ID:               uuid.New(),
		CardHolderName:   regDto.AccountHolderName,
		MerchantID:       regDto.MerchantID,
		MerchantPassword: regDto.MerchantPassword,
		MerchantAccount:  true,
		Balance:          models.NewZeroBalance(),
	}
}

func (s *accountService) DB() *gorm.DB {
	return s.accountRepo.DB()
}
