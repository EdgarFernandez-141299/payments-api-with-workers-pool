package command

import (
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

type CreateOrderCommand struct {
	ReferenceID    string
	TotalAmount    value_objects.CurrencyAmount
	Email          string
	PhoneNumber    string
	User           entities.User
	CountryCode    value_objects.Country
	CurrencyCode   value_objects.CurrencyCode
	BillingAddress value_objects.Address
	EnterpriseID   string
	Metadata       map[string]interface{}
	WebhookUrl     value_objects.WebhookUrl
	AllowCapture   bool
}

func (cmd *CreateOrderCommand) Validate() error {
	if cmd.ReferenceID == "" {
		return errors.NewOrderCreateValidationError(errors.ErrInvalidOrderID)
	}

	if err := cmd.TotalAmount.Validate(); err != nil {
		return err
	}

	if cmd.PhoneNumber == "" {
		return errors.NewOrderCreateValidationError(errors.ErrInvalidOrderPhoneNumber)

	}

	if err := cmd.User.Validate(); err != nil {
		return err
	}

	if cmd.EnterpriseID == "" {
		return errors.NewOrderCreateValidationError(errors.ErrInvalidOrderEnterpriseID)
	}

	if !cmd.WebhookUrl.IsEmpty() {
		if err := cmd.WebhookUrl.ValidateUrl(); err != nil {
			return err
		}
	}

	return nil
}

type CreateOrderBuilder struct {
	cmd CreateOrderCommand
}

func NewCreateOrderCommandBuilder() *CreateOrderBuilder {
	return &CreateOrderBuilder{}
}

func (b *CreateOrderBuilder) WithReferenceID(referenceID string) *CreateOrderBuilder {
	b.cmd.ReferenceID = referenceID
	return b
}

func (b *CreateOrderBuilder) WithTotalAmount(totalAmount value_objects.CurrencyAmount) *CreateOrderBuilder {
	b.cmd.TotalAmount = totalAmount
	return b
}

func (b *CreateOrderBuilder) WithPhoneNumber(phoneNumber string) *CreateOrderBuilder {
	b.cmd.PhoneNumber = phoneNumber
	return b
}

func (b *CreateOrderBuilder) WithUser(user entities.User) *CreateOrderBuilder {
	b.cmd.User = user
	return b
}

func (b *CreateOrderBuilder) WithCountryCode(countryCode value_objects.Country) *CreateOrderBuilder {
	b.cmd.CountryCode = countryCode
	return b
}

func (b *CreateOrderBuilder) WithBillingAddress(address value_objects.Address) *CreateOrderBuilder {
	b.cmd.BillingAddress = address
	return b
}

func (b *CreateOrderBuilder) WithCurrencyCode(currencyCode value_objects.CurrencyCode) *CreateOrderBuilder {
	b.cmd.CurrencyCode = currencyCode
	return b
}

func (b *CreateOrderBuilder) WithEnterpriseID(enterpriseID string) *CreateOrderBuilder {
	b.cmd.EnterpriseID = enterpriseID
	return b
}

func (b *CreateOrderBuilder) WithEmail(email string) *CreateOrderBuilder {
	b.cmd.Email = email
	return b
}

func (b *CreateOrderBuilder) WithMetadata(metadata map[string]interface{}) *CreateOrderBuilder {
	if metadata == nil {
		b.cmd.Metadata = make(map[string]interface{})
		return b
	} else {
		b.cmd.Metadata = metadata
	}

	return b
}

func (b *CreateOrderBuilder) WithWebhookUrl(webhookUrl value_objects.WebhookUrl) *CreateOrderBuilder {
	b.cmd.WebhookUrl = webhookUrl
	return b
}

func (b *CreateOrderBuilder) WithAllowCapture(allowCapture bool) *CreateOrderBuilder {
	b.cmd.AllowCapture = allowCapture
	return b
}

func (b *CreateOrderBuilder) Build() CreateOrderCommand {
	return b.cmd
}
