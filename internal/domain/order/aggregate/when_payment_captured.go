package aggregate

import (
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/events"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

func WhenPaymentCaptured(o *Order, ev events.OrderPaymentCapturedEvent) {
	o.Status = value_objects.OrderStatusProcessed()
	for i := range o.OrderPayments {
		if o.OrderPayments[i].ID == ev.PaymentID {
			o.OrderPayments[i].Status = enums.PaymentProcessed
			break
		}
	}
}
