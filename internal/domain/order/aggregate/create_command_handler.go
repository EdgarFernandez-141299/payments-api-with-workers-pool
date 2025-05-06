package aggregate

import (
	"time"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/events"
)

// Create initializes a new Order aggregate based on the provided CreateOrderCommand.
// It assigns the order a ReferenceID using the WithID method and then generates
// an appropriate event using the FromCreateOrderCommand function. The event is
// tracked using the TrackChange method.
//
// Parameters:
// - cmd: A CreateOrderCommand containing the necessary information to create an order.
//
// Returns:
// - *Order: A pointer to the newly created Order aggregate.
// - error: An error, if the order creation fails at any stage.
func Create(cmd command.CreateOrderCommand) (*Order, error) {
	o := new(Order)

	err := o.WithID(cmd.ReferenceID)

	if err != nil {
		return nil, err
	}

	event := events.NewOrderCreatedEventBuilder().
		SetID(cmd.ReferenceID).
		SetTotalAmount(cmd.TotalAmount).
		SetPhoneNumber(cmd.PhoneNumber).
		SetUser(cmd.User).
		SetCreatedAt(time.Now().UTC()).
		SetCountryCode(cmd.CountryCode).
		SetBillingAddress(cmd.BillingAddress).
		SetEnterpriseID(cmd.EnterpriseID).
		SetEmail(cmd.Email).
		SetMetadata(cmd.Metadata).
		SetWebhookUrl(cmd.WebhookUrl).
		SetAllowCapture(cmd.AllowCapture).
		Build()

	o.TrackChange(o, event)

	return o, nil
}
