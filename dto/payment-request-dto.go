package dto

import (
	"github.com/shopspring/decimal"
	"time"
)

type PaymentRequestDto struct {
	MerchantId        string          `json:"merchant_id" binding:"required"`
	MerchantPassword  string          `json:"merchant_password" binding:"required"`
	Amount            decimal.Decimal `json:"amount" binding:"required"`
	MerchantOrderId   string          `json:"merchant_order_id" binding:"required"`
	MerchantTimestamp time.Time       `json:"merchant_timestamp" binding:"required"`
	SuccessUrl        string          `json:"success_url" binding:"required,url"`
	FailedUrl         string          `json:"failed_url" binding:"required,url"`
	ErrorUrl          string          `json:"error_url" binding:"required,url"`
}
