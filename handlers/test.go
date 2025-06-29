package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// TestHandler godoc
// @Summary      Test endpoint
// @Description  Returns Hello, world!
// @Tags         Test
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string
// @Router       /api/test [get]
func TestHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, world!"})
}
