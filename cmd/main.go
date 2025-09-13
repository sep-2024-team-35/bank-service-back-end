package main

import (
	"github.com/joho/godotenv"
	"github.com/sep-2024-team-35/bank-servce-back-end/config"
	"github.com/sep-2024-team-35/bank-servce-back-end/routes"
	"log"
	"os"
	"strings"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("🔥 Panic recovered in main: %v", r)
		}
	}()

	log.Println("🚀 Starting Bank Service backend...")

	env := os.Getenv("APP_ENV")
	if env == "" {
		log.Fatal("❌ APP_ENV environment variable is not set. Please set it to 'acquirer' or 'issuer'.")
	}

	log.Printf("🌍 APP_ENV = %s", env)

	envFile := ".env." + strings.ToLower(env)
	if _, err := os.Stat(envFile); err == nil {
		if err := godotenv.Load(envFile); err != nil {
			log.Fatalf("❌ Failed to load env file %s: %v", envFile, err)
		}
		log.Printf("✅ Loaded local environment file: %s", envFile)
	} else {
		log.Printf("☁️ No local env file found (%s). Using cloud environment variables.", envFile)
	}

	// Initialize database
	log.Println("🔧 Initializing database...")
	config.InitDB()
	log.Println("✅ Database initialized.")

	log.Println("📦 Running migrations...")
	config.RunMigrations()
	log.Println("✅ Migrations completed.")

	log.Println("🌱 Seeding database...")
	config.RunSeeder()
	log.Println("✅ Database seeded.")

	// Use dynamic PORT from environment for Azure
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback za lokalni dev
	}
	log.Printf("🛡️ Starting HTTP server on port %s...", port)

	r := routes.SetupRouter()
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start HTTP server: %v", err)
	}
}
