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
// @Param request body dto.PaymentRequestDTO true "Payment request data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} dto.ErrorResponse
// @Router /payment/create/request [post]
func (h *PaymentHandler) CreateRequest(c *gin.Context) {
	var paymentRequest dto.PaymentRequestDTO

	if err := c.ShouldBindJSON(&paymentRequest); err != nil {
		log.Printf("[ERROR] Invalid payment request payload: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	log.Printf("[INFO] Received payment request: MerchantID=%s, OrderID=%s, Amount=%s",
		paymentRequest.MerchantID, paymentRequest.MerchantOrderId, paymentRequest.Amount.String())

	savedRequest, err := h.paymentService.CreateRequest(paymentRequest)
	if err != nil {
		log.Printf("[ERROR] Failed to persist payment request for MerchantID=%s: %v", paymentRequest.MerchantID, err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create payment request"})
		return
	}
	log.Printf("[INFO] Payment request saved: ID=%s, Amount=%s", savedRequest.ID.String(), savedRequest.Amount.String())

	redirectURL := fmt.Sprintf("https://ebanksep-fe.azurewebsites.net/card?paymentID=%s", savedRequest.ID.String())
	log.Printf("[INFO] Redirecting client to URL: %s", redirectURL)
	c.Redirect(http.StatusSeeOther, redirectURL)

	c.JSON(http.StatusCreated, map[string]string{"paymentRequestID": savedRequest.ID.String()})
}

// Pay godoc
// @Summary Complete a payment
// @Description Submits card details to complete a previously created payment request.
// Updates the transaction status from CREATED to PAID (or FAILED if validation fails).
// @Tags payments
// @Accept json
// @Produce json
// @Param paymentID path string true "ID of the payment request"
// @Param request body dto.CardDetailsDTO true "Card details for payment processing"
// @Success 200 {object} map[string]string "Transaction completed successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid input data"
// @Failure 404 {object} dto.ErrorResponse "Payment request not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /payment/{paymentID}/pay [patch]
func (h *PaymentHandler) Pay(c *gin.Context) {
	// 1. Uzmi paymentID iz path-a
	paymentID := c.Param("paymentID")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "paymentID is required"})
		return
	}

	// 2. Bind-uj JSON body u CardDetailsDTO
	var cardDetails dto.CardDetailsDTO
	if err := c.ShouldBindJSON(&cardDetails); err != nil {
		log.Printf("[ERROR] Invalid card details payload: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// 3. Pozovi servis
	transaction, err := h.paymentService.Pay(cardDetails, paymentID)
	if err != nil {
		log.Printf("[ERROR] Failed to process payment for PaymentID=%s: %v", paymentID, err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// 4. Vrati odgovor
	c.JSON(http.StatusOK, map[string]string{
		"status": transaction.Status,
	})
}
