package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sep-2024-team-35/bank-servce-back-end/dto"
	"github.com/sep-2024-team-35/bank-servce-back-end/services"
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default: // linux
		err = exec.Command("xdg-open", url).Start()
	}
	if err != nil {
		log.Printf("❌ Failed to open browser: %v", err)
	}
}

type PaymentHandler struct {
	paymentService     services.PaymentService
	transactionService *services.TransactionService
}

func NewPaymentHandler(paymentService services.PaymentService, transactionService *services.TransactionService) *PaymentHandler {
	return &PaymentHandler{
		paymentService:     paymentService,
		transactionService: transactionService,
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
		log.Printf("[ERROR] Invalid payment request payload: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	log.Printf("[INFO] Received payment request: MerchantID=%s, OrderID=%s, Amount=%s",
		paymentRequest.MerchantId, paymentRequest.MerchantOrderId, paymentRequest.Amount.String())

	savedRequest, err := h.paymentService.CreateRequest(paymentRequest)
	if err != nil {
		log.Printf("[ERROR] Failed to persist payment request for MerchantID=%s: %v", paymentRequest.MerchantId, err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create payment request"})
		return
	}
	log.Printf("[INFO] Payment request saved: ID=%s, Amount=%s", savedRequest.ID.String(), savedRequest.Amount.String())

	// Open browser with frontend form
	redirectURL := fmt.Sprintf("https://ebanksep-fe.azurewebsites.net/card?paymentID=%s", savedRequest.ID.String())

	openBrowser(redirectURL)
	log.Printf("[INFO] Browser opened with URL: %s", redirectURL)

	// Add transaction
	_, err = h.transactionService.AddTransactionAcquirer(savedRequest)
	if err != nil {
		log.Printf("[ERROR] Failed to create transaction for payment request ID=%s: %v", savedRequest.ID.String(), err)
		c.JSON(http.StatusNotAcceptable, nil)
		return
	}
	log.Printf("[INFO] Transaction successfully created for payment request ID=%s (Credit Card)", savedRequest.ID.String())

	c.JSON(http.StatusCreated, map[string]string{"paymentRequestID": savedRequest.ID.String()})
}
