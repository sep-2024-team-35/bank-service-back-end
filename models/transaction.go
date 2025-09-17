package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID                int             `gorm:"primaryKey;autoIncrement" json:"id"` // TODO: promeni u UUID
	Amount            decimal.Decimal `gorm:"column:amount" json:"amount"`
	MerchantOrderID   string          `gorm:"column:merchant_order_id" json:"merchantOrderId"`
	MerchantTimestamp string          `gorm:"column:merchant_timestamp" json:"merchantTimestamp"`
	AcquirerOrderID   string          `gorm:"column:acquirer_order_id" json:"acquirerOrderId"`
	AcquirerTimestamp string          `gorm:"column:acquirer_timestamp" json:"acquirerTimestamp"`
	IssuerOrderID     string          `gorm:"column:issuer_order_id" json:"issuerOrderId"`
	IssuerTimestamp   string          `gorm:"column:issuer_timestamp" json:"issuerTimestamp"`
	Status            string          `gorm:"column:status" json:"status"`
	// TODO: currency field
	PaymentRequestID uuid.UUID      `gorm:"type:uuid;column:payment_request_id;not null" json:"paymentRequestId"`
	PaymentRequest   PaymentRequest `gorm:"foreignKey:PaymentRequestID" json:"-"`
}
