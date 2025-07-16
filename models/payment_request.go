package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type PaymentRequest struct {
	ID                uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	MerchantID        string          `gorm:"column:merchant_id;not null" json:"merchantId"`
	MerchantPassword  string          `gorm:"column:merchant_password;not null" json:"merchantPassword"`
	Amount            decimal.Decimal `gorm:"type:numeric(15,2);column:amount;not null" json:"amount"`
	MerchantOrderID   string          `gorm:"column:merchant_order_id;not null" json:"merchantOrderId"`
	MerchantTimestamp time.Time       `gorm:"column:merchant_timestamp;not null" json:"merchantTimestamp"`
	SuccessURL        string          `gorm:"column:success_url;not null" json:"successUrl"`
	FailedURL         string          `gorm:"column:failed_url;not null" json:"failedUrl"`
	ErrorURL          string          `gorm:"column:error_url;not null" json:"errorUrl"`
}
