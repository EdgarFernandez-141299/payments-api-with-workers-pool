package aggregate

import (
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/events"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

func WhenPaymentPartialRefunded(o *Order, ev events.OrderPaymentPartialRefundedEvent) {
	for i, payment := range o.OrderPayments {
		if payment.ID == ev.PaymentID {
			o.OrderPayments[i].TotalRefunded = payment.TotalRefunded.Add(ev.Amount)

			if o.OrderPayments[i].TotalRefunded.Equals(payment.Total) {
				o.OrderPayments[i].Status = enums.PaymentRefunded
			} else {
				o.OrderPayments[i].Status = enums.PartiallyRefunded
			}
		}
	}

	if o.HasPaymentRefundable() {
		o.Status = value_objects.OrderStatusPartiallyRefunded()
	} else {
		o.Status = value_objects.OrderStatusRefunded()
	}
}
