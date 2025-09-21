package dto

type PSPResponseDTO struct {
	Status            string `json:"status"`
	RedirectURL       string `json:"redirectURL"`
	AcquirerOrderID   string `json:"acquirerOrderId"`
	AcquirerTimeStamp string `json:"acquirerTimeStamp"`
	PaymentID         string `json:"paymentID"`
	MerchantOrderID   string `json:"merchantOrderId"`
}
