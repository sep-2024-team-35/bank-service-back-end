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
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}

}
