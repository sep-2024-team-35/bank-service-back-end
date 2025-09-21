package dto

import (
	"github.com/shopspring/decimal"
)

type PaymentRequestDTO struct {
	MerchantID        string          `json:"merchantId" binding:"required"`
	MerchantPassword  string          `json:"merchantPassword" binding:"required"`
	Amount            decimal.Decimal `json:"amount" binding:"required"`
	MerchantOrderId   string          `json:"merchantOrderId" binding:"required"`
	MerchantTimestamp string          `json:"merchantOrderTimeStamp" binding:"required"`
	SuccessUrl        string          `json:"successUrl" binding:"required,url"`
	FailedUrl         string          `json:"failedUrl" binding:"required,url"`
	ErrorUrl          string          `json:"errorUrl" binding:"required,url"`
}
