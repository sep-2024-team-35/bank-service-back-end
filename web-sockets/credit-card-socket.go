package web_sockets

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sep-2024-team-35/bank-servce-back-end/dto"
	"github.com/sep-2024-team-35/bank-servce-back-end/services"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // ili proveri origin
	},
}

type CreditCardWSHandler struct {
	paymentService services.PaymentService
	accountService services.AccountService
	conn           *websocket.Conn
}

func NewCreditCardWSHandler(paymentSvc services.PaymentService, accountSvc services.AccountService) *CreditCardWSHandler {
	return &CreditCardWSHandler{
		paymentService: paymentSvc,
		accountService: accountSvc,
	}
}

func (h *CreditCardWSHandler) HandleWS(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade:", err)
		return
	}
	h.conn = ws
	defer ws.Close()

	// Welcome message
	h.conn.WriteMessage(websocket.TextMessage, []byte("Welcome to the Credit Card WebSocket!"))

	for {
		_, msg, err := h.conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		var cardDetails dto.CardDetailsDTO
		if err := json.Unmarshal(msg, &cardDetails); err != nil {
			fmt.Println("Invalid JSON:", string(msg))
			continue
		}

		// Prosledi servisu za izvršenje plaćanja
		err = h.paymentService.ExecutePayment(cardDetails)
		if err != nil {
			fmt.Println("Payment execution error:", err)
		}
	}
}

// Poziv sličan openCreditCardForm
func (h *CreditCardWSHandler) OpenCreditCardForm(paymentID string, amount float64) error {
	msg := fmt.Sprintf("%s,%.2f", paymentID, amount)
	return h.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}
