package command

type CreatePaymentOrderFailCommand struct {
	OrderID           string
	PaymentID         string
	PaymentReason     string
	OrderStatusString string
}

func NewCreatePaymentOrderFailCommand(
	orderID string,
	paymentID string,
	paymentReason,
	orderStatusString string,
) CreatePaymentOrderFailCommand {
	return CreatePaymentOrderFailCommand{
		OrderID:           orderID,
		PaymentID:         paymentID,
		PaymentReason:     paymentReason,
		OrderStatusString: orderStatusString,
	}
}
