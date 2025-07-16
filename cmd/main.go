// @title Bank Service API
// @version 1.0
// @description Microservice simulating bank accounts and cards
// @host localhost:8443
// @BasePath /api
// @schemes https

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
	env := os.Getenv("APP_ENV")
	if env == "" {
		log.Fatal("❌ APP_ENV environment variable is not set. Please set it to 'acquirer' or 'issuer'.")
	}

	envFile := ".env." + strings.ToLower(env)

	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("❌ Failed to load env file %s: %v", envFile, err)
	}

	log.Printf("✅ Loaded environment file: %s", envFile)

	config.InitDB()
	config.RunMigrations()
	config.RunSeeder()

	r := routes.SetupRouter()
	if err := r.RunTLS(":8443", "cert/cert.pem", "cert/key.pem"); err != nil {
		log.Fatalf("❌ Failed to start HTTPS server: %v", err)
	}
}
