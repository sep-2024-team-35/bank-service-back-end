package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sep-2024-team-35/bank-servce-back-end/dto"
	"github.com/sep-2024-team-35/bank-servce-back-end/services"
	_ "github.com/sep-2024-team-35/bank-servce-back-end/services"
	"net/http"
)

type PaymentHandler struct {
	paymentService services.PaymentService
}

func NewPaymentHandler(paymentService services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// CreateRequest godoc
// @Summary Create a new payment request
// @Description Accepts payment request data from acquirer
// @Tags payments
// @Accept json
// @Produce json
// @Param request body dto.PaymentRequestDto true "Payment request data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} dto.ErrorResponse
// @Router /payment/create/request [post]
func (h *PaymentHandler) CreateRequest(c *gin.Context) {
	var paymentRequest dto.PaymentRequestDto

	if err := c.ShouldBindJSON(&paymentRequest); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	savedRequest, err := h.paymentService.CreateRequest(paymentRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create payment request"})
		return
	}

	c.JSON(http.StatusCreated, savedRequest)
}
