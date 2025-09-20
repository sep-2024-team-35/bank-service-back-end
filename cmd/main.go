package main

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sep-2024-team-35/bank-servce-back-end/config"
	"github.com/sep-2024-team-35/bank-servce-back-end/routes"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("🔥 Panic recovered: %v", r)
		}
	}()
	
	env := os.Getenv("APP_ENV")
	if env == "" {
		log.Fatal("❌ APP_ENV is not set.")
	}

	envFile := ".env." + strings.ToLower(env)
	if _, err := os.Stat(envFile); err == nil {
		_ = godotenv.Load(envFile)
	}

	config.LoadEncryptionKeyGlobally()

	config.InitDB()
	config.RunMigrations()
	config.RunSeeder()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := routes.SetupRouter().Run(":" + port); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}
