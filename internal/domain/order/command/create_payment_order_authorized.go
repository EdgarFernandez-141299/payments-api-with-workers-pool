package command

type CreatePaymentOrderAuthorizedCommand struct {
	OrderID           string
	PaymentID         string
	AuthorizationCode string
	OrderStatusString string
	PaymentCard       CardData
}

func NewCreatePaymentOrderAuthorizedCommand(
	orderID string,
	paymentID string,
	authorizationCode string,
	orderStatusString string,
	paymentCard CardData,
) CreatePaymentOrderAuthorizedCommand {
	return CreatePaymentOrderAuthorizedCommand{
		OrderID:           orderID,
		PaymentID:         paymentID,
		AuthorizationCode: authorizationCode,
		OrderStatusString: orderStatusString,
		PaymentCard:       paymentCard,
	}
}
