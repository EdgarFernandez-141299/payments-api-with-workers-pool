package integration_events

type paidIntegrationEventOrderStatus string

const (
	paymentProcessed               paidIntegrationEventOrderStatus = "processed"
	paymentOrderProcessedEventType                                 = "order-payment-processed"
)

type OrderPaymentProcessedIntegrationEvent struct {
	baseOrderIntegrationEvent
	AuthorizationCode string `json:"authorization_code"`
}

func NewOrderPaymentProcessedIntegrationEvent(
	params IntegrationEventsParams,
	authorizationCode string,
	orderStatus string,
) OrderPaymentProcessedIntegrationEvent {
	return OrderPaymentProcessedIntegrationEvent{
		baseOrderIntegrationEvent: PaymentOrderIntegrationEventWithStatus(
			params, paymentOrderProcessedEventType, string(paymentProcessed), orderStatus,
		),
		AuthorizationCode: authorizationCode,
	}
}
