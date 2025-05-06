package aggregate

import (
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/events"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

func WhenPaymentRefunded(o *Order, ev events.OrderPaymentTotalRefundedEvent) {
	for i := range o.OrderPayments {
		if o.OrderPayments[i].ID == ev.PaymentID {
			o.OrderPayments[i].TotalRefunded = o.OrderPayments[i].Total
			o.OrderPayments[i].Status = enums.PaymentRefunded
		}
	}

	if o.HasPaymentRefundable() {
		o.Status = value_objects.OrderStatusPartiallyRefunded()
	} else {
		o.Status = value_objects.OrderStatusRefunded()
	}
}
