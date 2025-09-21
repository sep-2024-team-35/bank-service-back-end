package dto

import "github.com/google/uuid"

type CardDetailsDTO struct {
	PrimaryAccountNumber string    `json:"primaryAccountNumber,omitempty"`
	CardHolderName       string    `json:"cardHolderName,omitempty"`
	ExpirationDate       string    `json:"expirationDate,omitempty"`
	SecurityCode         string    `json:"securityCode,omitempty"`
	PaymentRequestID     uuid.UUID `json:"paymentRequestId,omitempty"`
}
