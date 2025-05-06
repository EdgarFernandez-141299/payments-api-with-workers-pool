package entities

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
)

func TestPaymentOrderEntity_TableName(t *testing.T) {
	entity := PaymentOrderEntity{}
	assert.Equal(t, "payment", entity.TableName())
}

func TestPaymentOrderEntityBuilder(t *testing.T) {
	builder := NewPaymentOrderEntity()

	orderID := "order123"
	associatedOrigin := "web"
	currencyCode := "USD"
	countryCode := "US"
	cardID := "card123"
	cardDetail := `{"number":"4111111111111111","expiry":"12/25"}`
	paymentMethod := "credit_card"
	collectionAccountID := "account123"
	metadata := `{"key":"value"}`
	status := "pending"
	totalAmount := decimal.NewFromFloat(100.50)
	reference := "ref123"
	failureReason := "declined"
	failureCode := "insufficient_funds"
	enterpriseID := "enterprise123"
	ipAddress := "192.168.1.1"
	deviceFingerprint := "fingerprint123"
	paymentOrderID := "payment_order_123"
	paymentFlow := enums.Autocapture
	authorizedAt := time.Now().Add(-1 * time.Hour)
	capturedAt := time.Now().Add(-30 * time.Minute)
	releasedAt := time.Now().Add(-15 * time.Minute)
	authorizationCode := "AUTH123456"

	// Test all setter methods
	entity := builder.
		SetOrderID(orderID).
		SetAssociatedOrigin(associatedOrigin).
		SetCurrencyCode(currencyCode).
		SetCountryCode(countryCode).
		SetCardID(cardID).
		SetCardDetail(cardDetail).
		SetPaymentMethod(paymentMethod).
		SetCollectionAccountID(collectionAccountID).
		SetMetadata(metadata).
		SetStatus(status).
		SetTotalAmount(totalAmount).
		SetReference(reference).
		SetFailureReason(failureReason).
		SetFailureCode(failureCode).
		SetEnterpriseID(enterpriseID).
		SetIPAddress(ipAddress).
		SetDeviceFingerprint(deviceFingerprint).
		SetTransactionDate().
		SetPaymentOrderID(paymentOrderID).
		SetPaymentFlow(paymentFlow).
		SetAuthorizedAt(authorizedAt).
		SetCapturedAt(capturedAt).
		SetReleasedAt(releasedAt).
		SetAuthorizationCode(authorizationCode).
		Build()

	// Verify all fields were set correctly
	assert.NotNil(t, entity.ID)
	assert.Equal(t, orderID, entity.OrderID)
	assert.Equal(t, associatedOrigin, entity.AssociatedOrigin)
	assert.Equal(t, currencyCode, entity.CurrencyCode)
	assert.Equal(t, countryCode, entity.CountryCode)
	assert.Equal(t, cardID, entity.CardID)
	assert.Equal(t, json.RawMessage(cardDetail), entity.CardDetail)
	assert.Equal(t, paymentMethod, entity.PaymentMethod)
	assert.Equal(t, collectionAccountID, entity.CollectionAccountID)
	assert.Equal(t, json.RawMessage(metadata), entity.Metadata)
	assert.Equal(t, status, entity.Status)
	assert.True(t, totalAmount.Equal(entity.TotalAmount))
	assert.Equal(t, reference, entity.Reference)
	assert.Equal(t, failureReason, entity.FailureReason)
	assert.Equal(t, failureCode, entity.FailureCode)
	assert.Equal(t, enterpriseID, entity.EnterpriseID)
	assert.Equal(t, ipAddress, entity.IPAddress)
	assert.Equal(t, deviceFingerprint, entity.DeviceFingerprint)
	assert.WithinDuration(t, time.Now(), entity.TransactionDate, time.Second)
	assert.Equal(t, paymentOrderID, entity.PaymentOrderID)
	assert.Equal(t, paymentFlow.String(), entity.PaymentFlow)
	assert.Equal(t, authorizedAt, entity.AuthorizedAt)
	assert.Equal(t, capturedAt, entity.CapturedAt)
	assert.Equal(t, releasedAt, entity.ReleasedAt)
	assert.Equal(t, authorizationCode, entity.AuthorizationCode)
}

func TestPaymentOrderEntityBuilder_OptionalFields(t *testing.T) {
	builder := NewPaymentOrderEntity()

	entity := builder.
		SetOrderID("order123").
		SetCurrencyCode("USD").
		SetTotalAmount(decimal.NewFromFloat(100.50)).
		Build()

	assert.NotNil(t, entity.ID)
	assert.Equal(t, "order123", entity.OrderID)
	assert.Equal(t, "USD", entity.CurrencyCode)
	assert.True(t, decimal.NewFromFloat(100.50).Equal(entity.TotalAmount))

	assert.Empty(t, entity.AssociatedOrigin)
	assert.Empty(t, entity.CountryCode)
	assert.Empty(t, entity.CardID)
	assert.Empty(t, entity.CardDetail)
	assert.Empty(t, entity.PaymentMethod)
	assert.Empty(t, entity.CollectionAccountID)
	assert.Empty(t, entity.Metadata)
	assert.Empty(t, entity.Status)
	assert.Empty(t, entity.Reference)
	assert.Empty(t, entity.FailureReason)
	assert.Empty(t, entity.FailureCode)
	assert.Empty(t, entity.EnterpriseID)
	assert.Empty(t, entity.IPAddress)
	assert.Empty(t, entity.DeviceFingerprint)
	assert.Empty(t, entity.PaymentOrderID)
	assert.Empty(t, entity.PaymentFlow)
	assert.True(t, entity.AuthorizedAt.IsZero())
	assert.True(t, entity.CapturedAt.IsZero())
	assert.True(t, entity.ReleasedAt.IsZero())
}

func TestPaymentOrderEntityBuilder_JSONFields(t *testing.T) {
	builder := NewPaymentOrderEntity()

	invalidJSON := "invalid json"
	entity := builder.
		SetCardDetail(invalidJSON).
		SetMetadata(invalidJSON).
		Build()

	assert.Equal(t, json.RawMessage(invalidJSON), entity.CardDetail)
	assert.Equal(t, json.RawMessage(invalidJSON), entity.Metadata)
}

func TestPaymentOrderEntityBuilder_ZeroValues(t *testing.T) {
	builder := NewPaymentOrderEntity()

	entity := builder.
		SetOrderID("").
		SetCurrencyCode("").
		SetTotalAmount(decimal.Zero).
		Build()

	assert.NotNil(t, entity.ID)
	assert.Empty(t, entity.OrderID)
	assert.Empty(t, entity.CurrencyCode)
	assert.True(t, decimal.Zero.Equal(entity.TotalAmount))
}

func TestPaymentOrderEntityBuilder_TableName(t *testing.T) {
	entity := PaymentOrderEntity{}
	assert.Equal(t, "payment", entity.TableName())
}

func TestPaymentOrderEntityBuilder_MultipleBuilds(t *testing.T) {
	builder := NewPaymentOrderEntity()

	entity1 := builder.
		SetOrderID("order1").
		SetCurrencyCode("USD").
		SetTotalAmount(decimal.NewFromFloat(100.50)).
		Build()

	entity2 := builder.
		SetOrderID("order2").
		SetCurrencyCode("EUR").
		SetTotalAmount(decimal.NewFromFloat(200.75)).
		Build()

	assert.NotEqual(t, entity1.ID, entity2.ID)
	assert.Equal(t, "order1", entity1.OrderID)
	assert.Equal(t, "order2", entity2.OrderID)
	assert.Equal(t, "USD", entity1.CurrencyCode)
	assert.Equal(t, "EUR", entity2.CurrencyCode)
	assert.True(t, decimal.NewFromFloat(100.50).Equal(entity1.TotalAmount))
	assert.True(t, decimal.NewFromFloat(200.75).Equal(entity2.TotalAmount))
}
