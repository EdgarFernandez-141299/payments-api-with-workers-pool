package events

import (
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
)

type PaymentProcessingStarted struct {
	PaymentOrder entities.PaymentOrder
}

type OrderPaymentProcessingStartedEventBuilder struct {
	event *PaymentProcessingStarted
}

func NewPaymentProcessingStartedEvent() *OrderPaymentProcessingStartedEventBuilder {
	return &OrderPaymentProcessingStartedEventBuilder{event: &PaymentProcessingStarted{}}
}

func (b *OrderPaymentProcessingStartedEventBuilder) WithPayment(paymentOrder entities.PaymentOrder) *OrderPaymentProcessingStartedEventBuilder {
	b.event.PaymentOrder = paymentOrder
	return b
}

func (b *OrderPaymentProcessingStartedEventBuilder) WithCollectionAccount(collectionAccount entities.CollectionAccount) *OrderPaymentProcessingStartedEventBuilder {
	b.event.PaymentOrder.CollectionAccount = collectionAccount
	return b
}

func (b *OrderPaymentProcessingStartedEventBuilder) Build() *PaymentProcessingStarted {
	return b.event
}

func FromProcessOrderCommand(cmd command.CreatePaymentOrderCommand) *PaymentProcessingStarted {
	event := NewPaymentProcessingStartedEvent().
		WithPayment(cmd.Payment).
		WithCollectionAccount(cmd.CollectionAccount).
		Build()

	return event
}
