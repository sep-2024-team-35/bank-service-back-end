package dto

type MerchantRegistrationDto struct {
	AccountHolderName string `json:"accountHolderName"`
	MerchantID        string `json:"merchantId"`
	MerchantPassword  string `json:"merchantPassword"`
}
