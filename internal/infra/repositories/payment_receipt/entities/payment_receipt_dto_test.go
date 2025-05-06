package entities

import (
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/clubhub.ai1/gommon/uid"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/payment_receipt/command"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

func TestPaymentReceiptDTO_ToDomain(t *testing.T) {
	t.Run("Successfully converts DTO to domain entity with CC payment method", func(t *testing.T) {
		id := uid.GenerateID()
		now := time.Now()

		dto := PaymentReceiptDTO{
			ID:                  id.String(),
			PaymentReceiptID:    id.String(),
			UserID:              "test-user-id",
			EnterpriseID:        "test-enterprise-id",
			Email:               "test@example.com",
			ReferenceOrderID:    "test-order-id",
			PaymentID:           "test-payment-id",
			PaymentStatus:       "completed",
			PaymentAmount:       decimal.NewFromFloat(100.50),
			PaymentCountryCode:  "USA",
			PaymentCurrencyCode: "USD",
			PaymentMethod:       enums.CCMethod.String(),
			PaymentDate:         now.Format(time.RFC3339),
			FileURL:             "https://example.com/receipt.pdf",
			CreatedAt:           now,
			UpdatedAt:           now,
		}

		entity := dto.ToDomain()

		assert.Equal(t, dto.ID, entity.ID)
		assert.Equal(t, dto.PaymentReceiptID, entity.PaymentReceiptID)
		assert.Equal(t, dto.UserID, entity.UserID)
		assert.Equal(t, dto.EnterpriseID, entity.EnterpriseID)
		assert.Equal(t, dto.Email, entity.Email)
		assert.Equal(t, dto.ReferenceOrderID, entity.ReferenceOrderID)
		assert.Equal(t, dto.PaymentID, entity.PaymentID)
		assert.Equal(t, dto.PaymentStatus, entity.PaymentStatus)
		assert.Equal(t, dto.PaymentAmount, entity.PaymentAmount)
		assert.Equal(t, dto.PaymentCountryCode, entity.PaymentCountryCode)
		assert.Equal(t, dto.PaymentCurrencyCode, entity.PaymentCurrencyCode)
		assert.Equal(t, enums.CCMethod, entity.PaymentMethod)
		assert.Equal(t, dto.PaymentDate, entity.PaymentDate)
		assert.Equal(t, dto.FileURL, entity.FileURL)
		assert.Equal(t, dto.CreatedAt, entity.CreatedAt)
		assert.Equal(t, dto.UpdatedAt, entity.UpdatedAt)
	})

	t.Run("Successfully converts DTO to domain entity with device payment method", func(t *testing.T) {
		id := uid.GenerateID()
		now := time.Now()

		dto := PaymentReceiptDTO{
			ID:                  id.String(),
			PaymentReceiptID:    "RCPT-" + id.String(),
			UserID:              "test-user-id",
			EnterpriseID:        "test-enterprise-id",
			Email:               "test@example.com",
			ReferenceOrderID:    "test-order-id",
			PaymentID:           "test-payment-id",
			PaymentStatus:       "completed",
			PaymentAmount:       decimal.NewFromFloat(100.50),
			PaymentCountryCode:  "USA",
			PaymentCurrencyCode: "USD",
			PaymentMethod:       enums.PaymentDevice.String(),
			PaymentDate:         now.Format(time.RFC3339),
			FileURL:             "https://example.com/receipt.pdf",
			CreatedAt:           now,
			UpdatedAt:           now,
		}

		entity := dto.ToDomain()

		assert.Equal(t, enums.PaymentDevice, entity.PaymentMethod)
	})
}

func TestPaymentReceiptDTO_ToCurrencyAmount(t *testing.T) {
	t.Run("Successfully converts DTO to CurrencyAmount", func(t *testing.T) {
		dto := PaymentReceiptDTO{
			PaymentAmount:       decimal.NewFromFloat(100.50),
			PaymentCurrencyCode: "USD",
		}

		currencyAmount, err := dto.ToCurrencyAmount()

		require.NoError(t, err)
		assert.Equal(t, dto.PaymentAmount, currencyAmount.Value)
		assert.Equal(t, dto.PaymentCurrencyCode, currencyAmount.Code.Code)
	})

	t.Run("Returns error for invalid currency code", func(t *testing.T) {
		dto := PaymentReceiptDTO{
			PaymentAmount:       decimal.NewFromFloat(100.50),
			PaymentCurrencyCode: "INVALID",
		}

		_, err := dto.ToCurrencyAmount()

		require.Error(t, err)
	})
}

func TestPaymentReceiptDTO_ToCountry(t *testing.T) {
	t.Run("Successfully converts DTO to Country", func(t *testing.T) {
		dto := PaymentReceiptDTO{
			PaymentCountryCode: "USA",
		}

		country, err := dto.ToCountry()

		require.NoError(t, err)
		assert.Equal(t, dto.PaymentCountryCode, country.Code)
	})

	t.Run("Returns error for invalid country code", func(t *testing.T) {
		dto := PaymentReceiptDTO{
			PaymentCountryCode: "INVALID",
		}

		_, err := dto.ToCountry()

		require.Error(t, err)
	})
}

func TestFromCommand(t *testing.T) {
	t.Run("Successfully creates DTO from command with CC payment method", func(t *testing.T) {
		now := time.Now()

		currencyCode, err := value_objects.NewCurrencyCode("USD")
		require.NoError(t, err)

		amount, err := value_objects.NewCurrencyAmount(currencyCode, decimal.NewFromFloat(100.50))
		require.NoError(t, err)

		country, err := value_objects.NewCountryWithCode("USA")
		require.NoError(t, err)

		cardInfo := value_objects.CardInfo{
			CardID: "card-123",
			CVV:    "123",
		}

		paymentMethod := value_objects.PaymentMethod{
			Type: enums.CCMethod,
			CCData: value_objects.PaymentMethodData[value_objects.CardInfo]{
				Data: cardInfo,
			},
		}

		// Create command
		cmd := command.CreatePaymentReceiptCommand{
			UserID:           "test-user-id",
			EnterpriseID:     "test-enterprise-id",
			Email:            "test@example.com",
			ReferenceOrderID: "test-order-id",
			PaymentID:        "test-payment-id",
			PaymentStatus:    "completed",
			PaymentAmount:    amount,
			PaymentCountry:   country,
			PaymentMethod:    paymentMethod,
			PaymentDate:      now.Format(time.RFC3339),
			CreatedAt:        now,
			UpdatedAt:        now,
		}

		dto := FromEntity(entities.NewPaymentReceiptEntity(cmd))

		assert.NotEmpty(t, dto.ID)
		assert.Equal(t, dto.ID, dto.PaymentReceiptID)
		assert.Equal(t, cmd.UserID, dto.UserID)
		assert.Equal(t, cmd.EnterpriseID, dto.EnterpriseID)
		assert.Equal(t, cmd.Email, dto.Email)
		assert.Equal(t, cmd.ReferenceOrderID, dto.ReferenceOrderID)
		assert.Equal(t, cmd.PaymentID, dto.PaymentID)
		assert.Equal(t, cmd.PaymentStatus, dto.PaymentStatus)
		assert.Equal(t, cmd.PaymentAmount.Value, dto.PaymentAmount)
		assert.Equal(t, cmd.PaymentCountry.Code, dto.PaymentCountryCode)
		assert.Equal(t, cmd.PaymentAmount.Code.Code, dto.PaymentCurrencyCode)
		assert.Equal(t, cmd.PaymentMethod.Type.String(), dto.PaymentMethod)
		assert.Equal(t, cmd.PaymentDate, dto.PaymentDate)
		assert.True(t, cmd.CreatedAt.Equal(dto.CreatedAt), "CreatedAt times should be equal")
		assert.True(t, cmd.UpdatedAt.Equal(dto.UpdatedAt), "UpdatedAt times should be equal")

		assert.Equal(t, cmd.PaymentMethod.Type.String(), dto.PaymentMethod)
	})

	t.Run("Successfully creates DTO from command with device payment method", func(t *testing.T) {
		now := time.Now()

		currencyCode, err := value_objects.NewCurrencyCode("USD")
		require.NoError(t, err)

		amount, err := value_objects.NewCurrencyAmount(currencyCode, decimal.NewFromFloat(100.50))
		require.NoError(t, err)

		country, err := value_objects.NewCountryWithCode("USA")
		require.NoError(t, err)

		paymentMethod := value_objects.PaymentMethod{
			Type: enums.PaymentDevice,
			DeviceData: value_objects.PaymentMethodData[value_objects.PaymentDevice]{
				Data: value_objects.PaymentDevice{},
			},
		}

		cmd := command.CreatePaymentReceiptCommand{
			UserID:           "test-user-id",
			EnterpriseID:     "test-enterprise-id",
			Email:            "test@example.com",
			ReferenceOrderID: "test-order-id",
			PaymentID:        "test-payment-id",
			PaymentStatus:    "completed",
			PaymentAmount:    amount,
			PaymentCountry:   country,
			PaymentMethod:    paymentMethod,
			PaymentDate:      now.Format(time.RFC3339),
			CreatedAt:        now,
			UpdatedAt:        now,
		}

		dto := FromEntity(entities.NewPaymentReceiptEntity(cmd))

		assert.Equal(t, cmd.PaymentMethod.Type.String(), dto.PaymentMethod)
	})
}
