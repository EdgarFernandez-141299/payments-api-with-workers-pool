package aggregate

import (
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/events"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

func WhenOrderPaymentFailed(o *Order, ev events.OrderPaymentFailedEvent) {
	o.Status = value_objects.OrderStatusFailed()
	for i := range o.OrderPayments {
		if o.OrderPayments[i].ID == ev.PaymentID {
			o.OrderPayments[i].Status = enums.PaymentFailed
			o.OrderPayments[i].FailureReason = ev.PaymentReason
			break
		}
	}
}
