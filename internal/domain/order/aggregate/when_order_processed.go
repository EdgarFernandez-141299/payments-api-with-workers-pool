package aggregate

import (
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/events"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

func WhenOrderPaymentProcessed(o *Order, ev events.OrderPaymentProcessedEvent) {
	o.Status = value_objects.OrderStatusProcessed()
	for i := range o.OrderPayments {
		if o.OrderPayments[i].ID == ev.PaymentID {
			o.OrderPayments[i].Status = enums.PaymentProcessed
			o.OrderPayments[i].AuthorizationCode = ev.AuthorizationCode
			o.OrderPayments[i].PaymentCard = entities.CardData{
				CardBrand: ev.PaymentCard.CardBrand,
				CardLast4: ev.PaymentCard.CardLast4,
				CardType:  ev.PaymentCard.CardType,
			}
			break
		}
	}
}
