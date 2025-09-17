package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sep-2024-team-35/bank-servce-back-end/config"
	"github.com/sep-2024-team-35/bank-servce-back-end/handlers"
	"github.com/sep-2024-team-35/bank-servce-back-end/repositories"
	"github.com/sep-2024-team-35/bank-servce-back-end/services"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"

	_ "github.com/sep-2024-team-35/bank-servce-back-end/docs"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		// TODO: correct origins and methods
		AllowOrigins:     []string{"https://localhost:8443", "http://localhost:8080", "http://localhost:3000", "https://ebanksep-fe.azurewebsites.net"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// === Repositories ===
	accountRepo := repositories.NewAccountRepository(config.DB)
	paymentRepo := repositories.NewPaymentRepository(config.DB)
	transactionRepo := repositories.NewTransactionRepository(config.DB)

	// === Services ===
	accountService := services.NewAccountService(accountRepo)
	transactionService := services.NewTransactionService(transactionRepo)
	paymentService := services.NewPaymentService(paymentRepo, accountService, transactionRepo)

	// === Handlers ===
	accountHandler := handlers.NewAccountHandler(accountService)
	paymentHandler := handlers.NewPaymentHandler(paymentService, transactionService)

	// === API rute ===
	api := r.Group("/api")
	{
		// Account routes
		api.POST("/account/register", accountHandler.RegisterNewMerchant)

		// Payment routes
		api.POST("/payment/create/request", paymentHandler.CreateRequest)
		api.PATCH("/payment/:paymentID/pay", paymentHandler.Pay)

	}

	return r
}
