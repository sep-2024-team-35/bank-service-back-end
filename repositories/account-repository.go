package repositories

import (
	"errors"
	"github.com/sep-2024-team-35/bank-servce-back-end/models"
	"gorm.io/gorm"
)

type AccountRepository interface {
	Save(account *models.Account) (*models.Account, error)
	Update(account *models.Account) (*models.Account, error)
	FindByMerchantID(merchantID string) (*models.Account, error)
	FindByPAN(pan string) (*models.Account, error)
	FindByMerchantIDAndPassword(merchantID string, password string) (*models.Account, error)
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Save(account *models.Account) (*models.Account, error) {
	if err := r.db.Create(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

func (r *accountRepository) Update(account *models.Account) (*models.Account, error) {
	if err := r.db.Save(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

func (r *accountRepository) FindByMerchantID(merchantID string) (*models.Account, error) {
	var account models.Account
	err := r.db.First(&account, "merchant_id = ?", merchantID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) FindByPAN(pan string) (*models.Account, error) {
	var account models.Account
	err := r.db.First(&account, "pan = ?", pan).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) FindByMerchantIDAndPassword(merchantID string, password string) (*models.Account, error) {
	var account models.Account
	if err := r.db.Where("merchant_id = ? AND merchant_password = ?", merchantID, password).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}
