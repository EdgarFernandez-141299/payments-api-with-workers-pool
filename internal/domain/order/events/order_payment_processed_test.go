package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
)

func TestFromProcessedOrderCommand(t *testing.T) {
	// Arrange
	cmd := command.CreatePaymentOrderProcessedCommand{
		OrderID:           "order123",
		PaymentID:         "payment123",
		AuthorizationCode: "auth123",
		OrderStatusString: "processed",
	}

	// Act
	event := FromProcessedOrderCommand(cmd)

	// Assert
	assert.NotNil(t, event)
	assert.Equal(t, cmd.OrderID, event.OrderID)
	assert.Equal(t, cmd.PaymentID, event.PaymentID)
	assert.Equal(t, cmd.AuthorizationCode, event.AuthorizationCode)
	assert.Equal(t, cmd.OrderStatusString, event.OrderStatusString)
}
