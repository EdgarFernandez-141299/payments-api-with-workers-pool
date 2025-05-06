package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

func TestNewOrderPaymentNotFoundError(t *testing.T) {
	t.Run("should create error with correct message and code", func(t *testing.T) {
		orderID := "order123"
		paymentID := "payment456"
		err := NewOrderPaymentNotFoundError(orderID, paymentID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no payment found for order ID order123 and payment ID payment456")
		assert.True(t, domain.IsBusinessErrorCode(err, orderPaymentNotFoundCodeError))
	})
}
