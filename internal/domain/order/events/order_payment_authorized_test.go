package events

import (
	"testing"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"

	"github.com/stretchr/testify/assert"
)

func TestFromAuthorizedOrderCommand(t *testing.T) {
	// Datos de prueba
	orderID := "test-order-123"
	paymentID := "test-payment-456"
	authCode := "AUTH123"
	orderStatus := "AUTHORIZED"
	cardData := command.CardData{
		CardBrand: "VISA",
		CardLast4: "1234",
		CardType:  "CREDIT",
	}

	// Crear comando de prueba
	cmd := command.CreatePaymentOrderAuthorizedCommand{
		OrderID:           orderID,
		PaymentID:         paymentID,
		AuthorizationCode: authCode,
		OrderStatusString: orderStatus,
		PaymentCard:       cardData,
	}

	// Convertir comando a evento
	event := FromAuthorizedOrderCommand(cmd)

	// Verificar que los datos se copiaron correctamente
	assert.Equal(t, orderID, event.OrderID)
	assert.Equal(t, paymentID, event.PaymentID)
	assert.Equal(t, authCode, event.AuthorizationCode)
	assert.Equal(t, orderStatus, event.OrderStatusString)
	assert.Equal(t, cardData, event.PaymentCard)
}
