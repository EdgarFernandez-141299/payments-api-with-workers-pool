package dto

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	vo "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
)

type CreateOrderRequestDTO struct {
	ReferenceOrderID string                 `json:"reference_order_id" validate:"required"`
	CurrencyCode     string                 `json:"currency_code" validate:"required"`
	CountryCode      string                 `json:"country_code" validate:"required"`
	UserType         string                 `json:"user_type" validate:"required"`    // member, employee, etc
	UserID           string                 `json:"user_id" validate:"required"`      // user id
	TotalAmount      decimal.Decimal        `json:"total_amount" validate:"required"` // cambiar a decimal.Decimal
	PhoneNumber      string                 `json:"phone_number"`
	Email            string                 `json:"email" validate:"required"`
	BillingAddress   Address                `json:"billing_address" validate:"required"`
	Metadata         map[string]interface{} `json:"metadata"`
	WebhookUrl       string                 `json:"webhook_url" validate:"required"`
	ClientOrigin     string                 `json:"client_origin"`
	AllowCapture     bool                   `json:"allow_capture"`
}

type Address struct {
	ZipCode     string `json:"zip_code" validate:"required"`
	Street      string `json:"street" validate:"required"`
	CountryCode string `json:"country" validate:"required"`
	City        string `json:"city" validate:"required"`
}

func (c *CreateOrderRequestDTO) Validate() error {
	err := validator.New().Struct(c)

	if err != nil {
		return err
	}

	if c.TotalAmount.Sign() <= 0 {
		return errors.New("total amount must be greater than 0")
	}

	return nil
}

func (c *CreateOrderRequestDTO) ToCommand(enterpriseID string) (command.CreateOrderCommand, error) {
	country, err := vo.NewCountryWithCode(c.CountryCode)

	if err != nil {
		return command.CreateOrderCommand{}, err
	}

	currencyCode, err := vo.NewCurrencyCode(c.CurrencyCode)

	if err != nil {
		return command.CreateOrderCommand{}, err
	}

	currencyAmount, err := vo.NewCurrencyAmount(currencyCode, c.TotalAmount)

	if err != nil {
		return command.CreateOrderCommand{}, err
	}

	userType := vo.NewUserTypeFromString(c.UserType)

	user := entities.NewUser(userType, c.UserID)

	webhookUrl := vo.NewWebhookUrl(c.WebhookUrl)

	return command.NewCreateOrderCommandBuilder().
		WithReferenceID(c.ReferenceOrderID).
		WithTotalAmount(currencyAmount).
		WithPhoneNumber(c.PhoneNumber).
		WithEnterpriseID(enterpriseID).
		WithEmail(c.Email).
		WithUser(user).
		WithCountryCode(country).
		WithBillingAddress(
			vo.NewAddress(
				c.BillingAddress.ZipCode,
				c.BillingAddress.Street,
				c.BillingAddress.City,
				country,
			),
		).
		WithCurrencyCode(currencyCode).
		WithCountryCode(country).
		WithMetadata(c.Metadata).
		WithWebhookUrl(webhookUrl).
		WithAllowCapture(c.AllowCapture).
		Build(), nil
}
