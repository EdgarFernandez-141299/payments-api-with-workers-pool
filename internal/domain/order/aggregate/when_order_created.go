package aggregate

import (
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/events"
	vo "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

// WhenOrderCreated applies the OrderCreated to the Order aggregate.
// If the provided Order is nil, a new instance is created and the event data is applied.
//
// Parameters:
//   - o: *Order - The existing Order instance or nil if creating a new one.
//   - ev: events.OrderCreated - The event containing order creation data.
func WhenOrderCreated(o *Order, ev events.OrderCreated) {
	if o == nil {
		o = &Order{}
	}

	o.Metadata = ev.Metadata
	o.OrderPayments = make([]entities.PaymentOrder, 0)
	o.ID = ev.ID
	o.TotalAmount = ev.TotalAmount
	o.Currency = ev.TotalAmount.Code
	o.PhoneNumber = ev.PhoneNumber
	o.CountryCode = ev.CountryCode
	o.User = ev.User
	o.CreatedAt = ev.CreatedAt
	o.Status = vo.OrderStatusProcessing()
	o.EnterpriseID = ev.EnterpriseID
	o.Email = ev.Email
	o.WebhookUrl = ev.WebhookUrl
	o.AllowCapture = ev.AllowCapture
}
