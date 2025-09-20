//package repositories
//
//import (
//	"encoding/base64"
//	"errors"
//	"fmt"
//	"github.com/sep-2024-team-35/bank-servce-back-end/config"
//	"github.com/sep-2024-team-35/bank-servce-back-end/crypto"
//	"github.com/sep-2024-team-35/bank-servce-back-end/models"
//	"gorm.io/gorm"
//	"log"
//)
//
//type AccountRepository interface {
//	Save(account *models.Account) (*models.Account, error)
//	Update(account *models.Account) (*models.Account, error)
//	UpdateTransactional(tx *gorm.DB, account *models.Account) (*models.Account, error)
//	FindByMerchantID(merchantID string) (*models.Account, error)
//	FindByPAN(pan string) (*models.Account, error)
//	FindByMerchantIDAndPassword(merchantID string, password string) (*models.Account, error)
//	DB() *gorm.DB
//}
//
//type accountRepository struct {
//	db *gorm.DB
//}
//
//func NewAccountRepository(db *gorm.DB) AccountRepository {
//	return &accountRepository{db: db}
//}
//
////func (r *accountRepository) Save(account *models.Account) (*models.Account, error) {
////	if err := r.db.Create(account).Error; err != nil {
////		return nil, err
////	}
////	return account, nil
////}
//
//func (r *accountRepository) Save(account *models.Account) (*models.Account, error) {
//	// TODO encrypt all sensitive data
//	// Encrypt PAN
//	encryptedPAN, err := crypto.Encrypt(account.PrimaryAccountNumber, config.EncryptionKey)
//	if err != nil {
//		return nil, fmt.Errorf("failed to encrypt PAN: %w", err)
//	}
//	account.PrimaryAccountNumber = base64.StdEncoding.EncodeToString(encryptedPAN)
//
//	// Encrypt Balance (as string)
//	//encryptedBalance, err := crypto.Encrypt(account.Balance.String(), config.EncryptionKey)
//	//if err != nil {
//	//	return nil, fmt.Errorf("failed to encrypt balance: %w", err)
//	//}
//	//account.Balance = base64.StdEncoding.EncodeToString(encryptedBalance)
//
//	// Save to DB
//	if err := r.db.Create(account).Error; err != nil {
//		return nil, err
//	}
//	return account, nil
//}
//
//func (r *accountRepository) Update(account *models.Account) (*models.Account, error) {
//	if err := r.db.Save(account).Error; err != nil {
//		return nil, err
//	}
//	return account, nil
//}
//
//func (r *accountRepository) UpdateTransactional(tx *gorm.DB, account *models.Account) (*models.Account, error) {
//	if tx == nil {
//		return nil, errors.New("transaction object cannot be nil")
//	}
//
//	if err := tx.Model(&models.Account{}).
//		Where("id = ?", account.ID).
//		Updates(account).Error; err != nil {
//		return nil, err
//	}
//
//	return account, nil
//}
//
//func (r *accountRepository) FindByMerchantID(merchantID string) (*models.Account, error) {
//	var account models.Account
//	err := r.db.First(&account, "merchant_id = ?", merchantID).Error
//	if errors.Is(err, gorm.ErrRecordNotFound) {
//		return nil, nil
//	}
//	if err != nil {
//		return nil, err
//	}
//	return &account, nil
//}
//
//func (r *accountRepository) FindByPAN(pan string) (*models.Account, error) {
//	var account models.Account
//	err := r.db.First(&account, "primary_account_number = ?", pan).Error
//	if errors.Is(err, gorm.ErrRecordNotFound) {
//		return nil, nil
//	}
//	if err != nil {
//		return nil, err
//	}
//	return &account, nil
//}
//
//func (r *accountRepository) FindByMerchantIDAndPassword(merchantID string, password string) (*models.Account, error) {
//	log.Printf("[Repo][FindByMerchantIDAndPassword] Called with merchantID=%s", merchantID)
//
//	if r.db == nil {
//		log.Printf("[Repo][FindByMerchantIDAndPassword][ERROR] Database connection is nil!")
//		return nil, fmt.Errorf("database connection is nil")
//	}
//
//	var account models.Account
//	err := r.db.Where("merchant_id = ? AND merchant_password = ?", merchantID, password).First(&account).Error
//
//	if errors.Is(err, gorm.ErrRecordNotFound) {
//		log.Printf("[Repo][FindByMerchantIDAndPassword] No account found for merchantID=%s", merchantID)
//		return nil, nil
//	}
//	if err != nil {
//		log.Printf("[Repo][FindByMerchantIDAndPassword][ERROR] Query failed for merchantID=%s: %v", merchantID, err)
//		return nil, err
//	}
//
//	log.Printf("[Repo][FindByMerchantIDAndPassword] Account found: ID=%s, MerchantID=%s, Balance=%.2f",
//		account.ID.String(), account.MerchantID, account.Balance)
//
//	return &account, nil
//}
//
//func (r *accountRepository) DB() *gorm.DB {
//	return r.db
//}

package repositories

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/sep-2024-team-35/bank-servce-back-end/config"
	"github.com/sep-2024-team-35/bank-servce-back-end/crypto"
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

func encryptAccountFields(account *models.Account) error {
	log.Printf("[Encrypt] Starting encryption for Account ID=%s", account.ID.String())

	// PAN
	if account.PrimaryAccountNumber != "" {
		log.Printf("[Encrypt][PAN] Original value: %s", account.PrimaryAccountNumber)
		encryptedPAN, err := crypto.Encrypt(account.PrimaryAccountNumber, config.EncryptionKey)
		if err != nil {
			log.Printf("[Encrypt][PAN][ERROR] Failed: %v", err)
			return fmt.Errorf("failed to encrypt PAN: %w", err)
		}
		encodedPAN := base64.StdEncoding.EncodeToString(encryptedPAN)
		log.Printf("[Encrypt][PAN] Encrypted (base64): %s", encodedPAN)
		account.PrimaryAccountNumber = encodedPAN
	}

	// CardHolderName
	if account.CardHolderName != "" {
		log.Printf("[Encrypt][CardHolderName] Original value: %s", account.CardHolderName)
		encryptedName, err := crypto.Encrypt(account.CardHolderName, config.EncryptionKey)
		if err != nil {
			log.Printf("[Encrypt][CardHolderName][ERROR] Failed: %v", err)
			return fmt.Errorf("failed to encrypt CardHolderName: %w", err)
		}
		encodedName := base64.StdEncoding.EncodeToString(encryptedName)
		log.Printf("[Encrypt][CardHolderName] Encrypted (base64): %s", encodedName)
		account.CardHolderName = encodedName
	}

	log.Printf("[Encrypt] Finished encryption for Account ID=%s", account.ID.String())
	return nil
}

func decryptAccountFields(account *models.Account) error {
	log.Printf("[Decrypt] Starting decryption for Account ID=%s", account.ID.String())

	// PAN
	if account.PrimaryAccountNumber != "" {
		log.Printf("[Decrypt][PAN] Encrypted (base64): %s", account.PrimaryAccountNumber)
		decodedPAN, err := base64.StdEncoding.DecodeString(account.PrimaryAccountNumber)
		if err == nil {
			decryptedPAN, derr := crypto.Decrypt(decodedPAN, config.EncryptionKey)
			if derr == nil {
				log.Printf("[Decrypt][PAN] Decrypted value: %s", decryptedPAN)
				account.PrimaryAccountNumber = decryptedPAN
			} else {
				log.Printf("[Decrypt][PAN][ERROR] Failed: %v", derr)
			}
		} else {
			log.Printf("[Decrypt][PAN][ERROR] Base64 decode failed: %v", err)
		}
	}

	// CardHolderName
	if account.CardHolderName != "" {
		log.Printf("[Decrypt][CardHolderName] Encrypted (base64): %s", account.CardHolderName)
		decodedName, err := base64.StdEncoding.DecodeString(account.CardHolderName)
		if err == nil {
			decryptedName, derr := crypto.Decrypt(decodedName, config.EncryptionKey)
			if derr == nil {
				log.Printf("[Decrypt][CardHolderName] Decrypted value: %s", decryptedName)
				account.CardHolderName = decryptedName
			} else {
				log.Printf("[Decrypt][CardHolderName][ERROR] Failed: %v", derr)
			}
		} else {
			log.Printf("[Decrypt][CardHolderName][ERROR] Base64 decode failed: %v", err)
		}
	}

	log.Printf("[Decrypt] Finished decryption for Account ID=%s", account.ID.String())
	return nil
}

//
// Repository methods
//

func (r *accountRepository) Save(account *models.Account) (*models.Account, error) {
	if err := encryptAccountFields(account); err != nil {
		return nil, err
	}

	if err := r.db.Create(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

func (r *accountRepository) Update(account *models.Account) (*models.Account, error) {
	if err := encryptAccountFields(account); err != nil {
		return nil, err
	}

	if err := r.db.Save(account).Error; err != nil {
		return nil, err
	}
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

	_ = decryptAccountFields(&account)
	return &account, nil
}

func (r *accountRepository) FindByPAN(pan string) (*models.Account, error) {
	var account models.Account
	err := r.db.First(&account, "primary_account_number = ?", pan).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	_ = decryptAccountFields(&account)
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

	_ = decryptAccountFields(&account)

	log.Printf("[Repo][FindByMerchantIDAndPassword] Account found: ID=%s, MerchantID=%s, Balance=%.2f",
		account.ID.String(), account.MerchantID, account.Balance)

	return &account, nil
}

func (r *accountRepository) DB() *gorm.DB {
	return r.db
}
