package dto

type MerchantRegistrationDTO struct {
	AccountHolderName string `json:"accountHolderName"`
	MerchantID        string `json:"merchantId"`
	MerchantPassword  string `json:"merchantPassword"`
}
