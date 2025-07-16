package config

import (
	"log"
	"os"
)

func RunSeeder() {
	sqlFile := "scripts/seed_acquirer.sql"
	if os.Getenv("APP_ENV") == "issuer" {
		sqlFile = "scripts/seed_issuer.sql"
	}

	content, err := os.ReadFile(sqlFile)
	if err != nil {
		log.Printf("⚠️  Seed file not found: %v", err)
		return
	}

	db, err := DB.DB()
	if err != nil {
		log.Fatalf("❌ Could not get DB: %v", err)
	}

	_, err = db.Exec(string(content))
	if err != nil {
		log.Printf("⚠️  Failed to execute seed: %v", err)
	} else {
		log.Println("✅ Seed data inserted successfully.")
	}
}
