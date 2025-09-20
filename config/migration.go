package config

import (
	"log"

	"github.com/sep-2024-team-35/bank-servce-back-end/models"
)

func RunMigrations() {
	_, err := DB.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get raw DB connection: %v", err)
	}
	
	_ = DB.Migrator().DropTable(
		&models.Transaction{},
		&models.PaymentRequest{},
		&models.Account{},
	)

	err = DB.AutoMigrate(
		&models.Account{},
		&models.PaymentRequest{},
		&models.Transaction{},
	)
	if err != nil {
		log.Fatalf("❌ Failed to migrate database: %v", err)
	}

	log.Println("✅ Database migrated successfully.")
}
