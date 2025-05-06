package errors

import "gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"

const refundCanNotAppliedCodeError = "REFUND_CAN_NOT_APPLIED_ERROR"

func NewRefundCanNotAppliedDueToPaymentStatusError() error {
	return domain.WrapBusinessError(nil, refundCanNotAppliedCodeError, "Payment status not processed", map[string]interface{}{})
}

func NewRefundAmountIsGreaterThanTotalRefundableError() error {
	return domain.WrapBusinessError(nil, refundCanNotAppliedCodeError, "Refund amount is greater than total refundable", map[string]interface{}{})
}
