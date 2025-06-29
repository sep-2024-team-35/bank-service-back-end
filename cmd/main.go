// @title           Bank Service API
// @version         1.0
// @description     Microservice simulating bank accounts and cards
// @host            localhost:8080
// @BasePath        /api
package main

import (
	"github.com/joho/godotenv"
	"github.com/sep-2024-team-35/bank-servce-back-end/config"
	"github.com/sep-2024-team-35/bank-servce-back-end/routes"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, relying on system environment variables")
	}

	config.InitDB()

	r := routes.SetupRouter()
	if err := r.RunTLS(":8443", "cert/cert.pem", "cert/key.pem"); err != nil {
		log.Fatalf("❌ Failed to start HTTPS server: %v", err)
	}

}
