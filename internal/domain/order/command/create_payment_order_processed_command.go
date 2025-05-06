package command

type CardData struct {
	CardBrand string
	CardLast4 string
	CardType  string
}

type CreatePaymentOrderProcessedCommand struct {
	OrderID           string
	PaymentID         string
	AuthorizationCode string
	OrderStatusString string
	PaymentCard       CardData
}

func NewCreatePaymentOrderProcessedCommand(
	orderID string,
	paymentID string,
	authorizationCode string,
	orderStatusString string,
	paymentCard CardData,
) CreatePaymentOrderProcessedCommand {
	return CreatePaymentOrderProcessedCommand{
		OrderID:           orderID,
		PaymentID:         paymentID,
		AuthorizationCode: authorizationCode,
		OrderStatusString: orderStatusString,
		PaymentCard:       paymentCard,
	}
}
