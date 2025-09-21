package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func InitDB() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("❌ DATABASE_URL is not set")
	}

	log.Printf("🔗 Connecting to database: %s", dbURL)

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	log.Println("✅ Connected to database")
	DB = db
}
