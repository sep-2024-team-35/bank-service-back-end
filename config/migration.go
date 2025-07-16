package config

import (
	"github.com/sep-2024-team-35/bank-servce-back-end/models"
	"log"
)

func RunMigrations() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get raw DB connection: %v", err)
	}

	_, err = sqlDB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err != nil {
		log.Fatalf("❌ Failed to create uuid-ossp extension: %v", err)
	}

	_ = DB.Migrator().DropTable(
		&models.Transaction{},
		&models.PaymentRequest{},
		&models.BankIdentifierNumber{},
		&models.Account{},
	)

	err = DB.AutoMigrate(
		&models.Account{},
		&models.BankIdentifierNumber{},
		&models.PaymentRequest{},
		&models.Transaction{},
	)
	if err != nil {
		log.Fatalf("❌ Failed to migrate database: %v", err)
	}
	log.Println("✅ Database migrated successfully.")
}
