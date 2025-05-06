package aggregate

import (
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/events"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

func WhenPaymentReleased(o *Order, ev events.OrderPaymentReleasedEvent) {
	o.Status = value_objects.OrderStatusCanceled()
	for i := range o.OrderPayments {
		if o.OrderPayments[i].ID == ev.PaymentID {
			o.OrderPayments[i].Status = enums.PaymentCanceled
			break
		}
	}
}
