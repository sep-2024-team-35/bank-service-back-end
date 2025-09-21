package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Account struct {
	ID                      uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MerchantAccount         bool      `gorm:"column:merchant_account;not null" json:"merchantAccount"`
	MerchantPassword        string    `gorm:"column:merchant_password" json:"merchantPassword,omitempty"`
	MerchantID              string    `gorm:"column:merchant_id;unique" json:"merchantId,omitempty"`
	PrimaryAccountNumber    string    `gorm:"column:primary_account_number" json:"primaryAccountNumber,omitempty"`
	Number                  string    `gorm:"column:number" json:"number,omitempty"`
	BankIdentifierCode      string    `gorm:"column:bank_identifier_code" json:"bankIdentifierCode,omitempty"`
	CardHolderName          string    `gorm:"column:card_holder_name;not null" json:"cardHolderName"`
	EncryptedBalance        string    `gorm:"column:encrypted_balance" json:"-"`
	EncryptedExpirationDate string    `gorm:"column:encrypted_expiration_date" json:"-"`
	PANHash                 string    `gorm:"column:pan_hash;index" json:"-"`
	CCV                     string    `gorm:"column:ccv;" json:"ccv,omitempty"`

	Balance        decimal.Decimal `gorm:"-" json:"balance,omitempty"`
	ExpirationDate time.Time       `gorm:"-" json:"expirationDate,omitempty"`
}

func NewZeroBalance() decimal.Decimal {
	return decimal.NewFromInt(0)
}
