package command

import (
	"errors"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

type CreatePaymentOrderCommand struct {
	ReferenceOrderID  string
	ID                string
	Payment           entities.PaymentOrder
	User              entities.User
	AssociatedOrigin  value_objects.AssociatedOrigin
	CurrencyCode      value_objects.CurrencyCode
	CollectionAccount entities.CollectionAccount
	AuthorizationCode string
	Metadata          string
	PaymentFlow       enums.PaymentFlowEnum
	CountryCode       string
}

func (c CreatePaymentOrderCommand) WithCollectionAccount(account entities.CollectionAccount) CreatePaymentOrderCommand {
	newCommand := c

	newCommand.CollectionAccount = account

	return newCommand
}

func (c CreatePaymentOrderCommand) WithAuthorizationCode(authorizationCode string) CreatePaymentOrderCommand {
	newCommand := c

	newCommand.AuthorizationCode = authorizationCode

	return newCommand
}

func (c CreatePaymentOrderCommand) WithPaymentFlow(paymentFlow enums.PaymentFlowEnum) CreatePaymentOrderCommand {
	newCommand := c

	newCommand.PaymentFlow = paymentFlow
	newCommand.Payment.PaymentFlow = paymentFlow.String()

	return newCommand
}

func (c CreatePaymentOrderCommand) WithCardData(cardData entities.Card) CreatePaymentOrderCommand {
	newCommand := c

	newCommand.Payment.PaymentCard = entities.CardData{
		CardBrand: cardData.Brand,
		CardLast4: cardData.LastFour,
		CardType:  cardData.CardType,
	}

	return newCommand
}

func (c CreatePaymentOrderCommand) Validate() error {
	if c.ReferenceOrderID == "" {
		return errors.New("ReferenceOrderID is required")
	}

	if err := c.Payment.Validate(); err != nil {
		return err
	}

	if err := c.User.Validate(); err != nil {
		return err
	}

	return nil
}

type CreatePaymentOrderCommandBuilder struct {
	id                string
	referenceOrderID  string
	paymentOrder      entities.PaymentOrder
	user              entities.User
	associatedOrigin  value_objects.AssociatedOrigin
	currencyCode      value_objects.CurrencyCode
	collectionAccount entities.CollectionAccount
	paymentFlow       enums.PaymentFlowEnum
	countryCode       string
}

func NewCreatePaymentOrderCommandBuilder() *CreatePaymentOrderCommandBuilder {
	return &CreatePaymentOrderCommandBuilder{}
}

func (b *CreatePaymentOrderCommandBuilder) WithReferenceOrderID(referenceOrderID string) *CreatePaymentOrderCommandBuilder {
	b.referenceOrderID = referenceOrderID
	return b
}

func (b *CreatePaymentOrderCommandBuilder) WithPayment(paymentOrder entities.PaymentOrder) *CreatePaymentOrderCommandBuilder {
	b.paymentOrder = paymentOrder
	return b
}

func (b *CreatePaymentOrderCommandBuilder) WithUser(user entities.User) *CreatePaymentOrderCommandBuilder {
	b.user = user
	return b
}

func (b *CreatePaymentOrderCommandBuilder) WithPaymentOrderID(id string) *CreatePaymentOrderCommandBuilder {
	b.id = id
	return b
}

func (b *CreatePaymentOrderCommandBuilder) WithAssociatedOrigin(associatedOrigin value_objects.AssociatedOrigin) *CreatePaymentOrderCommandBuilder {
	b.associatedOrigin = associatedOrigin
	return b
}

func (b *CreatePaymentOrderCommandBuilder) WithCurrencyCode(currencyCode value_objects.CurrencyCode) *CreatePaymentOrderCommandBuilder {
	b.currencyCode = currencyCode
	return b
}

func (b *CreatePaymentOrderCommandBuilder) WithCountryCode(countryCode string) *CreatePaymentOrderCommandBuilder {
	b.countryCode = countryCode
	return b
}

func (b *CreatePaymentOrderCommandBuilder) WithCollectionAccount(collectionAccount entities.CollectionAccount) *CreatePaymentOrderCommandBuilder {
	b.collectionAccount = collectionAccount
	return b
}

func (b *CreatePaymentOrderCommandBuilder) WithPaymentFlow(paymentFlow enums.PaymentFlowEnum) *CreatePaymentOrderCommandBuilder {
	b.paymentFlow = paymentFlow
	return b
}

func (b *CreatePaymentOrderCommandBuilder) Build() CreatePaymentOrderCommand {
	return CreatePaymentOrderCommand{
		ReferenceOrderID:  b.referenceOrderID,
		Payment:           b.paymentOrder,
		User:              b.user,
		ID:                b.id,
		AssociatedOrigin:  b.associatedOrigin,
		CurrencyCode:      b.currencyCode,
		CollectionAccount: b.collectionAccount,
		PaymentFlow:       b.paymentFlow,
		CountryCode:       b.countryCode,
	}
}
