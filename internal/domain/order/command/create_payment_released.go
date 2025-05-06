package command

type PaymentOrderReleasedCommand struct {
	OrderID   string
	PaymentID string
	Reason    string
}

func NewPaymentOrderReleasedCommand(orderID, paymentID, reason string) PaymentOrderReleasedCommand {
	return PaymentOrderReleasedCommand{
		OrderID:   orderID,
		PaymentID: paymentID,
		Reason:    reason,
	}
}
