package events

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
)

func TestFromCapturedOrderCommand(t *testing.T) {
	// Arrange
	orderID := "test-order-id"
	paymentID := "test-payment-id"
	amount := decimal.NewFromFloat(100.50)
	cmd := command.PaymentOrderCapturedCommand{
		OrderID:   orderID,
		PaymentID: paymentID,
		Amount:    amount,
	}

	// Act
	event := FromCapturedOrderCommand(cmd)

	// Assert
	assert.NotNil(t, event)
	assert.Equal(t, orderID, event.OrderID)
	assert.Equal(t, paymentID, event.PaymentID)
	assert.Equal(t, amount, event.Amount)
	assert.NotZero(t, event.CapturedAt)
	assert.True(t, time.Now().After(event.CapturedAt) || time.Now().Equal(event.CapturedAt))
}
