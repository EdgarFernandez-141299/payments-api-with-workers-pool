package command

import (
	"testing"

	bizErrors "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/errors"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

func TestNewCreateOrderCommand(t *testing.T) {
	userID := "user123"
	user := entities.NewUser(value_objects.NewUserType(value_objects.Member), userID)
	usdCurrencyCode, _ := value_objects.NewCurrencyCode("USD")
	totalAmountAsFloat := decimal.NewFromFloat(100.5)
	enterpriseID := "enterprise123"
	countryCode, _ := value_objects.NewCountryWithCode("MX")

	totalAmount, _ := value_objects.NewCurrencyAmount(usdCurrencyCode, totalAmountAsFloat)
	invalidAmount, _ := value_objects.NewCurrencyAmount(usdCurrencyCode, decimal.NewFromFloat(-10))
	invalidUser := entities.NewUser(value_objects.NewUserType(value_objects.Member), "")
	billingAddress := value_objects.NewAddress("1234", "street", "MX", value_objects.Country{Code: "MX"})
	email := "email@gmail.com"
	metadata := map[string]interface{}{}
	webhookUrl := value_objects.NewWebhookUrl("https://webhook.example.com/callback")

	t.Run("create order command - success", func(t *testing.T) {
		metadata["key"] = "value"

		cmd := NewCreateOrderCommandBuilder().
			WithReferenceID("1234").
			WithTotalAmount(totalAmount).
			WithPhoneNumber("3222332").
			WithUser(user).
			WithCurrencyCode(usdCurrencyCode).
			WithCountryCode(countryCode).
			WithEnterpriseID(enterpriseID).
			WithEmail(email).
			WithMetadata(metadata).
			WithBillingAddress(billingAddress).
			WithWebhookUrl(webhookUrl).
			Build()

		assert.NoError(t, cmd.Validate())
		assert.Equal(t, webhookUrl, cmd.WebhookUrl)
	})

	t.Run("create order command - invalid total amount", func(t *testing.T) {
		cmd := NewCreateOrderCommandBuilder().
			WithReferenceID("1234").
			WithTotalAmount(invalidAmount).
			WithPhoneNumber("3222332").
			WithUser(user).
			WithCurrencyCode(usdCurrencyCode).
			WithCountryCode(countryCode).
			WithEnterpriseID(enterpriseID).
			WithWebhookUrl(webhookUrl).
			Build()

		assert.Error(t, cmd.Validate(), bizErrors.NewInvalidAmountError(decimal.NewFromFloat(-10)).Error())
	})

	t.Run("create order command - missing user", func(t *testing.T) {
		cmd := NewCreateOrderCommandBuilder().
			WithReferenceID("1234").
			WithTotalAmount(totalAmount).
			WithPhoneNumber("3222332").
			WithUser(invalidUser).
			WithCurrencyCode(usdCurrencyCode).
			WithCountryCode(countryCode).
			WithEnterpriseID(enterpriseID).
			WithWebhookUrl(webhookUrl).
			Build()

		assert.Error(t, cmd.Validate(), bizErrors.NewInvalidUserIDError("").Error())
	})

	t.Run("with error metadata", func(t *testing.T) {
		matadata := map[string]interface{}{}
		matadata["key"] = "'"

		cmd := NewCreateOrderCommandBuilder().
			WithReferenceID("1234").
			WithTotalAmount(totalAmount).
			WithPhoneNumber("3222332").
			WithUser(user).
			WithCurrencyCode(usdCurrencyCode).
			WithCountryCode(countryCode).
			WithEnterpriseID(enterpriseID).
			WithEmail(email).
			WithMetadata(metadata).
			WithWebhookUrl(webhookUrl).
			Build()

		assert.NoError(t, cmd.Validate())
		assert.Equal(t, webhookUrl, cmd.WebhookUrl)
	})

	t.Run("create order command - missing enterprise ID", func(t *testing.T) {
		cmd := NewCreateOrderCommandBuilder().
			WithReferenceID("1234").
			WithTotalAmount(totalAmount).
			WithPhoneNumber("3222332").
			WithUser(user).
			WithCurrencyCode(usdCurrencyCode).
			WithCountryCode(countryCode).
			WithEnterpriseID("").
			WithWebhookUrl(webhookUrl).
			Build()

		assert.Error(t, cmd.Validate(), bizErrors.NewOrderCreateValidationError(bizErrors.ErrInvalidOrderEnterpriseID).Error())
	})

	t.Run("create order command - missing reference ID", func(t *testing.T) {
		cmd := NewCreateOrderCommandBuilder().
			WithReferenceID("").
			WithTotalAmount(totalAmount).
			WithPhoneNumber("3222332").
			WithUser(user).
			WithCurrencyCode(usdCurrencyCode).
			WithCountryCode(countryCode).
			WithEnterpriseID(enterpriseID).
			WithWebhookUrl(webhookUrl).
			Build()

		assert.Error(t, cmd.Validate(), bizErrors.NewOrderCreateValidationError(bizErrors.ErrInvalidOrderID).Error())
	})

	t.Run("create order command - missing phone number", func(t *testing.T) {
		cmd := NewCreateOrderCommandBuilder().
			WithReferenceID("1234").
			WithTotalAmount(totalAmount).
			WithPhoneNumber("").
			WithUser(user).
			WithCurrencyCode(usdCurrencyCode).
			WithCountryCode(countryCode).
			WithEnterpriseID(enterpriseID).
			WithWebhookUrl(webhookUrl).
			Build()

		assert.Error(t, cmd.Validate(), bizErrors.NewOrderCreateValidationError(bizErrors.ErrInvalidOrderPhoneNumber).Error())
	})

	t.Run("create order command - with empty webhook URL", func(t *testing.T) {
		emptyWebhookUrl := value_objects.NewWebhookUrl("")
		cmd := NewCreateOrderCommandBuilder().
			WithReferenceID("1234").
			WithTotalAmount(totalAmount).
			WithPhoneNumber("3222332").
			WithUser(user).
			WithCurrencyCode(usdCurrencyCode).
			WithCountryCode(countryCode).
			WithEnterpriseID(enterpriseID).
			WithEmail(email).
			WithMetadata(metadata).
			WithBillingAddress(billingAddress).
			WithWebhookUrl(emptyWebhookUrl).
			Build()

		assert.NoError(t, cmd.Validate())
		assert.Equal(t, emptyWebhookUrl, cmd.WebhookUrl)
	})

	t.Run("create order command - with invalid webhook URL", func(t *testing.T) {
		invalidWebhookUrl := value_objects.NewWebhookUrl("http://invalid-url")
		cmd := NewCreateOrderCommandBuilder().
			WithReferenceID("1234").
			WithTotalAmount(totalAmount).
			WithPhoneNumber("3222332").
			WithUser(user).
			WithCurrencyCode(usdCurrencyCode).
			WithCountryCode(countryCode).
			WithEnterpriseID(enterpriseID).
			WithEmail(email).
			WithMetadata(metadata).
			WithBillingAddress(billingAddress).
			WithWebhookUrl(invalidWebhookUrl).
			Build()

		assert.NoError(t, cmd.Validate())
		assert.Equal(t, invalidWebhookUrl, cmd.WebhookUrl)
	})

	t.Run("create order command - with allow capture", func(t *testing.T) {
		cmd := NewCreateOrderCommandBuilder().
			WithReferenceID("1234").
			WithTotalAmount(totalAmount).
			WithPhoneNumber("3222332").
			WithUser(user).
			WithCurrencyCode(usdCurrencyCode).
			WithCountryCode(countryCode).
			WithEnterpriseID(enterpriseID).
			WithEmail(email).
			WithMetadata(metadata).
			WithBillingAddress(billingAddress).
			WithWebhookUrl(webhookUrl).
			WithAllowCapture(true).
			Build()

		assert.NoError(t, cmd.Validate())
		assert.True(t, cmd.AllowCapture)
	})

	t.Run("create order command - without allow capture", func(t *testing.T) {
		cmd := NewCreateOrderCommandBuilder().
			WithReferenceID("1234").
			WithTotalAmount(totalAmount).
			WithPhoneNumber("3222332").
			WithUser(user).
			WithCurrencyCode(usdCurrencyCode).
			WithCountryCode(countryCode).
			WithEnterpriseID(enterpriseID).
			WithEmail(email).
			WithMetadata(metadata).
			WithBillingAddress(billingAddress).
			WithWebhookUrl(webhookUrl).
			WithAllowCapture(false).
			Build()

		assert.NoError(t, cmd.Validate())
		assert.False(t, cmd.AllowCapture)
	})
}
