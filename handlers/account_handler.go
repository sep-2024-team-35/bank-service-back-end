package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sep-2024-team-35/bank-servce-back-end/dto"
	_ "github.com/sep-2024-team-35/bank-servce-back-end/dto"
	"github.com/sep-2024-team-35/bank-servce-back-end/services"
	"net/http"
)

type AccountHandler struct {
	Service services.AccountService
}

func NewAccountHandler(service services.AccountService) *AccountHandler {
	return &AccountHandler{Service: service}
}

// RegisterNewMerchant godoc
// @Summary Register a new merchant
// @Description Accepts merchant registration data
// @Tags accounts
// @Accept json
// @Produce json
// @Param request body dto.MerchantRegistrationDto true "Merchant registration data" example({"accountHolderName":"Petar Petrovic","merchantId":"merchant-123","merchantPassword":"securePass123"})
// @Success 201 {object} models.Account
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /account/register [post]
func (h *AccountHandler) RegisterNewMerchant(c *gin.Context) {
	var merchantDto dto.MerchantRegistrationDto
	if err := c.ShouldBindJSON(&merchantDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	account, err := h.Service.RegisterNewMerchant(&merchantDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register merchant"})
		return
	}

	c.JSON(http.StatusCreated, account)
}
