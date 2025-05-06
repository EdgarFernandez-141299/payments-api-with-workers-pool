package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaymentMethodConstants(t *testing.T) {
	t.Run("should have correct value for PaymentMethodCreditCard", func(t *testing.T) {
		assert.Equal(t, "CCData", PaymentMethodCreditCard)
	})

	t.Run("should have correct value for TerminalPaymentMethod", func(t *testing.T) {
		assert.Equal(t, "TERMINAL", TerminalPaymentMethod)
	})
}
