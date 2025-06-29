package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sep-2024-team-35/bank-servce-back-end/handlers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/sep-2024-team-35/bank-servce-back-end/docs"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		api.GET("/test", handlers.TestHandler)
	}

	return r
}
