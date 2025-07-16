package models

type BankIdentifierNumber struct {
	ID string `gorm:"primaryKey;column:id" json:"id"`
}
