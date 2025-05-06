package entities

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestOrderEntityBuilder(t *testing.T) {
	builder := NewOrderEntityBuilder()

	id := "order123"
	userID := "user123"
	referenceOrderID := "ref123"
	totalAmount := decimal.NewFromFloat(100.50)
	countryCode := "US"
	currencyCode := "USD"
	status := "pending"
	enterpriseID := "enterprise123"
	metadata := map[string]interface{}{
		"key": "value",
	}
	allowCapture := true

	// Test SetMetadata with string
	entity := builder.
		SetID(id).
		SetUserID(userID).
		SetReferenceOrderID(referenceOrderID).
		SetTotalAmount(totalAmount).
		SetCountryCode(countryCode).
		SetCurrencyCode(currencyCode).
		SetStatus(status).
		SetEnterpriseID(enterpriseID).
		SetMetadata(metadata).
		SetAllowCapture(allowCapture).
		Build()

	assert.Equal(t, id, entity.ID)
	assert.Equal(t, userID, entity.UserID)
	assert.Equal(t, referenceOrderID, entity.ReferenceOrderID)
	assert.True(t, totalAmount.Equal(entity.TotalAmount))
	assert.Equal(t, countryCode, entity.CountryCode)
	assert.Equal(t, currencyCode, entity.CurrencyCode)
	assert.Equal(t, status, entity.Status)
	assert.Equal(t, enterpriseID, entity.EnterpriseID)
	assert.Equal(t, allowCapture, entity.AllowCapture)
	assert.Nil(t, entity.DeletedAt)

	var metadataJSON map[string]string
	err := json.Unmarshal(entity.Metadata, &metadataJSON)
	assert.NoError(t, err)
	assert.Equal(t, "value", metadataJSON["key"])

	// Test SetMetadataFromMap
	metadataMap := map[string]interface{}{
		"key2": "value2",
		"num":  123,
	}

	entity = NewOrderEntityBuilder().
		SetID(id).
		SetMetadataFromMap(metadataMap).
		Build()

	var metadataJSON2 map[string]interface{}
	err = json.Unmarshal(entity.Metadata, &metadataJSON2)
	assert.NoError(t, err)
	assert.Equal(t, "value2", metadataJSON2["key2"])
	assert.Equal(t, float64(123), metadataJSON2["num"])
}

func TestOrderEntityMethods(t *testing.T) {
	t.Run("IsEmpty", func(t *testing.T) {
		entity := OrderEntity{}
		assert.True(t, entity.IsEmpty())

		entity.ID = "123"
		assert.False(t, entity.IsEmpty())
	})

	t.Run("TableName", func(t *testing.T) {
		entity := OrderEntity{}
		assert.Equal(t, "order", entity.TableName())
	})

	t.Run("SetStatus", func(t *testing.T) {
		entity := &OrderEntity{}
		entity.SetStatus("completed")
		assert.Equal(t, "completed", entity.Status)
	})
}
