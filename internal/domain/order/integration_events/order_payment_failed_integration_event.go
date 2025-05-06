package integration_events

const (
	failedStatus                = "failed"
	paymentOrderFailedEventType = "order-payment-failed"
)

type OrderFailedPaidIntegrationEvent struct {
	baseOrderIntegrationEvent
	FailureReason string `json:"failure_reason"`
}

func NewOrderFailedIntegrationEvent(
	params IntegrationEventsParams,
	reason string,
	orderStatus string,
) OrderFailedPaidIntegrationEvent {
	return OrderFailedPaidIntegrationEvent{
		baseOrderIntegrationEvent: PaymentOrderIntegrationEventWithStatus(
			params, paymentOrderFailedEventType, failedStatus, orderStatus,
		),
		FailureReason: reason,
	}
}
