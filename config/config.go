package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get raw DB connection: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Ping failed: %v", err)
	}

	fmt.Println("✅ Successfully connected to the database.")
	DB = db
}
