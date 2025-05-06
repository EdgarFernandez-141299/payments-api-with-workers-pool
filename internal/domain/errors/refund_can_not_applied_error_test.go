package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRefundCanNotAppliedDueToPaymentStatusError(t *testing.T) {
	err := NewRefundCanNotAppliedDueToPaymentStatusError()
	assert.Error(t, err)
	assert.Equal(t, "Business Error code: REFUND_CAN_NOT_APPLIED_ERROR, message: Payment status not processed", err.Error())
}

func TestNewRefundAmountIsGreaterThanTotalRefundableError(t *testing.T) {
	err := NewRefundAmountIsGreaterThanTotalRefundableError()
	assert.Error(t, err)
	assert.Equal(t, "Business Error code: REFUND_CAN_NOT_APPLIED_ERROR, message: Refund amount is greater than total refundable", err.Error())
}
