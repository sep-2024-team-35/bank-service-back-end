package repositories

import (
	"errors"
	"fmt"
	"github.com/sep-2024-team-35/bank-servce-back-end/models"
	"gorm.io/gorm"
	"log"
)

type AccountRepository interface {
	Save(account *models.Account) (*models.Account, error)
	Update(account *models.Account) (*models.Account, error)
	UpdateTransactional(tx *gorm.DB, account *models.Account) (*models.Account, error)
	FindByMerchantID(merchantID string) (*models.Account, error)
	FindByPAN(pan string) (*models.Account, error)
	FindByMerchantIDAndPassword(merchantID string, password string) (*models.Account, error)
	DB() *gorm.DB
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

func (r *accountRepository) UpdateTransactional(tx *gorm.DB, account *models.Account) (*models.Account, error) {
	if tx == nil {
		return nil, errors.New("transaction object cannot be nil")
	}

	if err := tx.Model(&models.Account{}).
		Where("id = ?", account.ID).
		Updates(account).Error; err != nil {
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
	log.Printf("[Repo][FindByMerchantIDAndPassword] Called with merchantID=%s", merchantID)

	if r.db == nil {
		log.Printf("[Repo][FindByMerchantIDAndPassword][ERROR] Database connection is nil!")
		return nil, fmt.Errorf("database connection is nil")
	}

	var account models.Account
	err := r.db.Where("merchant_id = ? AND merchant_password = ?", merchantID, password).First(&account).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[Repo][FindByMerchantIDAndPassword] No account found for merchantID=%s", merchantID)
		return nil, nil
	}
	if err != nil {
		log.Printf("[Repo][FindByMerchantIDAndPassword][ERROR] Query failed for merchantID=%s: %v", merchantID, err)
		return nil, err
	}

	log.Printf("[Repo][FindByMerchantIDAndPassword] Account found: ID=%s, MerchantID=%s, Balance=%.2f",
		account.ID.String(), account.MerchantID, account.Balance)

	return &account, nil
}

func (r *accountRepository) DB() *gorm.DB {
	return r.db
}
