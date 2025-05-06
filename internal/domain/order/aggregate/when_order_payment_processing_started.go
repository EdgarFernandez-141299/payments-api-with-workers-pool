package aggregate

import (
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/events"
	vo "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

func WhenOrderPaymentProcessingStarted(o *Order, event events.PaymentProcessingStarted) {
	o.OrderPayments = append(o.OrderPayments, event.PaymentOrder)
	o.Status = vo.OrderStatusProcessing()
}
