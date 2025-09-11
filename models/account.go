package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Account struct {
	ID                   uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	MerchantAccount      bool            `gorm:"column:merchant_account;not null" json:"merchantAccount"`
	MerchantPassword     string          `gorm:"column:merchant_password" json:"merchantPassword,omitempty"`
	MerchantID           string          `gorm:"column:merchant_id;unique" json:"merchantId,omitempty"`
	PrimaryAccountNumber string          `gorm:"column:pan" json:"primaryAccountNumber,omitempty"`
	Number               string          `gorm:"column:number" json:"number,omitempty"`
	BankIdentifierCode   string          `gorm:"column:bank_identifier_code" json:"bankIdentifierCode,omitempty"`
	ExpirationDate       time.Time       `gorm:"column:expiration_date" json:"expirationDate,omitempty"`
	CardHolderName       string          `gorm:"column:card_holder_name;not null" json:"cardHolderName"`
	SecurityCode         string          `gorm:"column:security_code" json:"securityCode,omitempty"`
	Balance              decimal.Decimal `gorm:"type:numeric(15,2);column:balance;not null" json:"balance"`
}

func NewZeroBalance() decimal.Decimal {
	return decimal.NewFromInt(0)
}
