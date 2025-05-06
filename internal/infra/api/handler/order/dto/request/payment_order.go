package dto

import (
	"errors"

	vo "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
	paymentmethodsDto "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/request/payment_methods"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/constants"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
)

type CreatePaymentOrderRequestDTO struct {
	OrderID      string `json:"reference_order_id" validate:"required"` // reference to the order
	UserID       string `json:"user_id" validate:"required"`
	UserType     string `json:"user_type" validate:"required"`
	CurrencyCode string `json:"currency_code" validate:"required"`
	CountryCode  string `json:"country_code" validate:"required"`
	PaymentOrderRequestDTO
}

type PaymentOrderRequestDTO struct {
	PaymentOrderID   string           `json:"payment_order_id" validate:"required"`
	AssociatedOrigin string           `json:"associated_origin" validate:"required"`
	PaymentMethod    PaymentMethodDTO `json:"payment_method" validate:"required"`
	Amount           decimal.Decimal  `json:"amount" validate:"required"`
}

type PaymentMethodDTO struct {
	Type       string                         `json:"type" validate:"required"` // CCData, TERMINAL, etc
	CreditCard paymentmethodsDto.CreditCard   `json:"credit_card"`
	TokenCard  paymentmethodsDto.TokenCardDTO `json:"token_card"`
}

func (c *CreatePaymentOrderRequestDTO) Validate() error {
	err := validator.New().Struct(c)

	if err != nil {
		return err
	}

	if c.Amount.Sign() <= 0 {
		return errors.New("amount must be greater than 0")
	}

	userType := vo.NewUserTypeFromString(c.UserType)
	if !userType.IsValid() {
		return errors.New("unsupported user type")
	}

	_, err = vo.NewFromAssociatedOriginString(c.AssociatedOrigin)
	if err != nil {
		return err
	}

	if c.PaymentMethod.Type == constants.PaymentMethodCreditCard {
		if c.PaymentMethod.CreditCard.ID == "" || c.PaymentMethod.CreditCard.CVV == "" {
			return errors.New("credit card id and cvv are required")
		}
	}

	if c.PaymentMethod.Type == constants.PaymentMethodTokenCard {
		if c.PaymentMethod.TokenCard.Token == "" || c.PaymentMethod.TokenCard.CVV == "" {
			return errors.New("token card token and cvv are required")
		}

		if c.PaymentMethod.TokenCard.SaveCard && c.PaymentMethod.TokenCard.Card.Alias == "" {
			return errors.New("to save card, alias is required")
		}
	}

	return nil
}

func (c *CreatePaymentOrderRequestDTO) Command() (command.CreatePaymentOrderCommand, error) {
	userType := vo.NewUserTypeFromString(c.UserType)

	if !userType.IsValid() {
		return command.CreatePaymentOrderCommand{}, errors.New("unsupported user type")
	}

	user := entities.NewUser(userType, c.UserID)

	associatedOrigin, err := vo.NewFromAssociatedOriginString(c.AssociatedOrigin)
	if err != nil {
		return command.CreatePaymentOrderCommand{}, err
	}

	currencyCode, err := vo.NewCurrencyCode(c.CurrencyCode)
	if err != nil {
		return command.CreatePaymentOrderCommand{}, err
	}

	country, err := vo.NewCountryWithCode(c.CountryCode)
	if err != nil {
		return command.CreatePaymentOrderCommand{}, err
	}

	currencyAmount, err := vo.NewCurrencyAmount(currencyCode, c.Amount)
	if err != nil {
		return command.CreatePaymentOrderCommand{}, err
	}

	builder := command.NewCreatePaymentOrderCommandBuilder().
		WithReferenceOrderID(c.OrderID).
		WithPaymentOrderID(c.PaymentOrderID).
		WithAssociatedOrigin(associatedOrigin).
		WithCurrencyCode(currencyCode).
		WithCountryCode(country.Code).
		WithUser(user)

	switch c.PaymentMethod.Type {
	case constants.PaymentMethodCreditCard:
		if c.PaymentMethod.CreditCard.ID == "" || c.PaymentMethod.CreditCard.CVV == "" {
			return command.CreatePaymentOrderCommand{}, errors.New("credit card id and cvv are required")
		}
		paymentMethod := vo.NewCCPaymentMethod(c.PaymentMethod.CreditCard.ID, c.PaymentMethod.CreditCard.CVV)
		payment := entities.NewPaymentOrder(c.PaymentOrderID, associatedOrigin, currencyAmount, paymentMethod)

		return builder.WithPayment(payment).Build(), nil
	case constants.PaymentMethodTokenCard:
		if c.PaymentMethod.TokenCard.Token == "" || c.PaymentMethod.TokenCard.CVV == "" {
			return command.CreatePaymentOrderCommand{}, errors.New("token card token and cvv are required")
		}

		paymentMethod := vo.NewTokenCardPaymentMethod(c.PaymentMethod.TokenCard.Token,
			c.PaymentMethod.TokenCard.CVV,
			c.PaymentMethod.TokenCard.Card.Brand,
			c.PaymentMethod.TokenCard.Card.Last4,
			c.PaymentMethod.TokenCard.Card.Exp,
			c.PaymentMethod.TokenCard.Card.CardType,
			c.PaymentMethod.TokenCard.Card.Alias,
			c.PaymentMethod.TokenCard.SaveCard,
		)

		payment := entities.NewPaymentOrder(c.PaymentOrderID, associatedOrigin, currencyAmount, paymentMethod)

		return builder.WithPayment(payment).Build(), nil
	case constants.TerminalPaymentMethod:
		paymentMethod := vo.NewDevicePaymentMethod()
		payment := entities.NewPaymentOrder(c.PaymentOrderID, associatedOrigin, currencyAmount, paymentMethod)

		return builder.WithPayment(payment).Build(), nil
	default:
		return command.CreatePaymentOrderCommand{}, errors.New("unsupported payment method")
	}
}
