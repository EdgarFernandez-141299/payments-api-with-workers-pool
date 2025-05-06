package events

import (
	"testing"
	"time"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"

	"github.com/stretchr/testify/assert"
)

func TestFromReleasedOrderCommand(t *testing.T) {
	// Datos de prueba
	orderID := "test-order-123"
	paymentID := "test-payment-456"
	reason := "Payment released due to cancellation"

	// Crear comando de prueba
	cmd := command.PaymentOrderReleasedCommand{
		OrderID:   orderID,
		PaymentID: paymentID,
		Reason:    reason,
	}

	// Convertir comando a evento
	event := FromReleasedOrderCommand(cmd)

	// Verificar que los datos se copiaron correctamente
	assert.Equal(t, orderID, event.OrderID)
	assert.Equal(t, paymentID, event.PaymentID)
	assert.Equal(t, reason, event.Reason)

	// Verificar que ReleasedAt se estableció con la hora actual
	assert.NotZero(t, event.ReleasedAt)
	assert.True(t, time.Since(event.ReleasedAt) < time.Second, "ReleasedAt should be set to current time")
}
