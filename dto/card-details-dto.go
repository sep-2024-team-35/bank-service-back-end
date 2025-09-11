package dto

type CardDetailsDTO struct {
	PrimaryAccountNumber string `json:"primaryAccountNumber,omitempty"`
	CardHolderName       string `json:"cardHolderName,omitempty"`
	ExpirationDate       string `json:"expirationDate,omitempty"`
	SecurityCode         string `json:"securityCode,omitempty"`
	PaymentRequestID     string `json:"paymentRequestId,omitempty"`
}
