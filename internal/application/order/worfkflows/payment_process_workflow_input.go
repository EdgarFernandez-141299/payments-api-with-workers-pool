package worfkflows

import "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"

type PaymentWorkflowInput struct {
	OrderID         string
	PaymentCommands []command.CreatePaymentOrderCommand
}

func NewPaymentWorkflowInput(
	orderID string,
	paymentCommands []command.CreatePaymentOrderCommand,
) PaymentWorkflowInput {
	return PaymentWorkflowInput{
		OrderID:         orderID,
		PaymentCommands: paymentCommands,
	}
}
