package services

import "os"

type BankService struct {
	bic string
}

func NewBankService() *BankService {
	return &BankService{
		bic: os.Getenv("BIC"),
	}
}

func (s *BankService) GetBIC() string {
	return s.bic
}
