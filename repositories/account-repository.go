package repositories

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/sep-2024-team-35/bank-servce-back-end/config"
	"github.com/sep-2024-team-35/bank-servce-back-end/crypto"
	"github.com/sep-2024-team-35/bank-servce-back-end/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"log"
	"time"
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

func encryptAccountFields(account *models.Account) error {
	log.Printf("[Encrypt] Starting encryption for Account ID=%s", account.ID.String())

	if account.PrimaryAccountNumber != "" {
		account.PANHash = crypto.HashPAN(account.PrimaryAccountNumber)
		if _, err := crypto.Encrypt(account.PrimaryAccountNumber, config.EncryptionKey); err != nil {
			return fmt.Errorf("failed to encrypt PAN: %w", err)
		}
	}

	if account.CardHolderName != "" {
		if _, err := crypto.Encrypt(account.CardHolderName, config.EncryptionKey); err != nil {
			return fmt.Errorf("failed to encrypt CardHolderName: %w", err)
		}
	}

	if account.Balance.GreaterThan(decimal.Zero) {
		bStr := account.Balance.String()
		encBal, err := crypto.Encrypt(bStr, config.EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt Balance: %w", err)
		}
		account.EncryptedBalance = base64.StdEncoding.EncodeToString(encBal)
	}

	if !account.ExpirationDate.IsZero() {
		dateStr := account.ExpirationDate.Format(time.RFC3339)
		encDate, err := crypto.Encrypt(dateStr, config.EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt ExpirationDate: %w", err)
		}
		account.EncryptedExpirationDate = base64.StdEncoding.EncodeToString(encDate)
	}

	log.Printf("[Encrypt] Finished encryption for Account ID=%s", account.ID.String())
	return nil
}

func decryptAccountFields(account *models.Account) error {
	log.Printf("[Decrypt] Starting decryption for Account ID=%s", account.ID.String())

	if account.EncryptedBalance != "" {
		decoded, err := base64.StdEncoding.DecodeString(account.EncryptedBalance)
		if err == nil {
			if decrypted, derr := crypto.Decrypt(decoded, config.EncryptionKey); derr == nil {
				account.Balance, _ = decimal.NewFromString(decrypted)
			}
		}
	}

	if account.EncryptedExpirationDate != "" {
		decoded, err := base64.StdEncoding.DecodeString(account.EncryptedExpirationDate)
		if err == nil {
			if decrypted, derr := crypto.Decrypt(decoded, config.EncryptionKey); derr == nil {
				t, _ := time.Parse(time.RFC3339, decrypted)
				account.ExpirationDate = t
			}
		}
	}

	log.Printf("[Decrypt] Finished decryption for Account ID=%s", account.ID.String())
	return nil
}

func (r *accountRepository) Save(account *models.Account) (*models.Account, error) {
	if account.MerchantPassword != "" {
		hashed, err := crypto.HashPassword(account.MerchantPassword)
		if err != nil {
			return nil, err
		}
		account.MerchantPassword = hashed
	}

	if err := encryptAccountFields(account); err != nil {
		return nil, err
	}

	if err := r.db.Create(account).Error; err != nil {
		log.Printf("[Save][ERROR] Failed to create account: %v", err)
		return nil, err
	}
	log.Printf("[Save] Account saved successfully: ID=%s", account.ID.String())
	return account, nil
}

func (r *accountRepository) Update(account *models.Account) (*models.Account, error) {
	if err := encryptAccountFields(account); err != nil {
		return nil, err
	}

	if err := r.db.Save(account).Error; err != nil {
		log.Printf("[Update][ERROR] Failed to update account ID=%s: %v", account.ID.String(), err)
		return nil, err
	}
	log.Printf("[Update] Account updated successfully: ID=%s", account.ID.String())
	return account, nil
}

func (r *accountRepository) UpdateTransactional(tx *gorm.DB, account *models.Account) (*models.Account, error) {
	if tx == nil {
		return nil, errors.New("transaction object cannot be nil")
	}

	if err := encryptAccountFields(account); err != nil {
		return nil, err
	}

	if err := tx.Model(&models.Account{}).
		Where("id = ?", account.ID).
		Updates(account).Error; err != nil {
		log.Printf("[UpdateTransactional][ERROR] Failed for Account ID=%s: %v", account.ID.String(), err)
		return nil, err
	}

	log.Printf("[UpdateTransactional] Account updated successfully: ID=%s", account.ID.String())
	return account, nil
}

func (r *accountRepository) FindByMerchantID(merchantID string) (*models.Account, error) {
	var account models.Account
	err := r.db.First(&account, "merchant_id = ?", merchantID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[FindByMerchantID] No account found for merchantID=%s", merchantID)
		return nil, nil
	}
	if err != nil {
		log.Printf("[FindByMerchantID][ERROR] Query failed: %v", err)
		return nil, err
	}

	_ = decryptAccountFields(&account)
	log.Printf("[FindByMerchantID] Account retrieved successfully: ID=%s", account.ID.String())
	return &account, nil
}

func (r *accountRepository) FindByPAN(pan string) (*models.Account, error) {
	var account models.Account
	//err := r.db.First(&account, "primary_account_number = ?", pan).Error
	hashed := crypto.HashPAN(pan)
	err := r.db.First(&account, "pan_hash = ?", hashed).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[FindByPAN] No account found")
		return nil, nil
	}
	if err != nil {
		log.Printf("[FindByPAN][ERROR] Query failed: %v", err)
		return nil, err
	}

	_ = decryptAccountFields(&account)
	log.Printf("[FindByPAN] Account retrieved successfully: ID=%s", account.ID.String())
	return &account, nil
}

func (r *accountRepository) FindByMerchantIDAndPassword(merchantID string, password string) (*models.Account, error) {
	log.Printf("[FindByMerchantIDAndPassword] Called for merchantID=%s", merchantID)

	var account models.Account
	err := r.db.First(&account, "merchant_id = ?", merchantID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[FindByMerchantIDAndPassword] No account found for merchantID=%s", merchantID)
		return nil, nil
	}
	if err != nil {
		log.Printf("[FindByMerchantIDAndPassword][ERROR] Query failed: %v", err)
		return nil, err
	}

	if !crypto.CheckPasswordHash(account.MerchantPassword, password) {
		log.Printf("[FindByMerchantIDAndPassword] Password mismatch for merchantID=%s", merchantID)
		return nil, nil
	}

	_ = decryptAccountFields(&account)
	log.Printf("[FindByMerchantIDAndPassword] Account retrieved successfully: ID=%s", account.ID.String())
	return &account, nil
}

func (r *accountRepository) DB() *gorm.DB {
	return r.db
}
