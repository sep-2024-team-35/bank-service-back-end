package repositories

import (
	"github.com/sep-2024-team-35/bank-servce-back-end/models"
	"gorm.io/gorm"
	"time"
)

type AccountRepository interface {
	FindByMerchantIDAndPassword(merchantID, merchantPassword string) (*models.Account, error)
	FindByPANAndSecurityCodeAndHolderNameAndExpirationDate(pan, code, holder string, expDate time.Time) (*models.Account, error)
	Save(account *models.Account) (*models.Account, error)
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) FindByMerchantIDAndPassword(merchantID, merchantPassword string) (*models.Account, error) {
	var account models.Account
	err := r.db.Where("merchant_id = ? AND merchant_password = ?", merchantID, merchantPassword).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) FindByPANAndSecurityCodeAndHolderNameAndExpirationDate(pan, code, holder string, expDate time.Time) (*models.Account, error) {
	var account models.Account
	err := r.db.Where("pan = ? AND security_code = ? AND card_holder_name = ? AND expiration_date = ?", pan, code, holder, expDate).
		First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) Save(account *models.Account) (*models.Account, error) {
	if err := r.db.Create(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}
