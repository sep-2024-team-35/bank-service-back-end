package dto

import "github.com/shopspring/decimal"

type ExternalTransactionRequestDTO struct {
	AcquirerOrderID      string          `json:"acquirerOrderId"`
	AcquirerTimestamp    string          `json:"acquirerTimestamp"`
	PrimaryAccountNumber string          `json:"primaryAccountNumber"`
	Amount               decimal.Decimal `json:"amount"`
	Currency             string          `json:"currency"`
}
