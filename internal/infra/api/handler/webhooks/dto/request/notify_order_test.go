package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWebhookOrderDTO(t *testing.T) {
	t.Run("should create a valid WebhookOrderDTO", func(t *testing.T) {
		// Arrange
		expectedOrder := Order{
			Token:            "test-token",
			MerchantID:       "merchant-123",
			PaymentMethod:    "credit_card",
			PaymentStatus:    "approved",
			Currency:         "MXN",
			TaxAmount:        100,
			ItemsTotalAmount: 1000,
			SubTotal:         900,
			TotalAmount:      1000,
			OrderID:          "order-123",
			TransactionID:    "trans-123",
			Metadata: map[string]interface{}{
				"test_key": "test_value",
			},
			Payment: Payment{
				Data: PaymentData{
					Metadata: map[string]interface{}{
						"payment_test": "payment_value",
					},
					FromCard: CardInfo{
						CardBrand:  "visa",
						FirstSix:   "411111",
						LastFour:   "1111",
						CardHolder: "John Doe",
					},
					Amount: MoneyAmount{
						Amount:   1000,
						Currency: "MXN",
					},
					UpdatedAt:  time.Now().Format(time.RFC3339),
					MethodType: "card",
					CreatedAt:  time.Now().Format(time.RFC3339),
					Merchant: PaymentMerchant{
						StoreCode: "store-123",
						ID:        "merchant-123",
					},
					ID:        "payment-123",
					Processor: "stripe",
					Customer: PaymentCustomer{
						Email: "test@example.com",
						ID:    "customer-123",
					},
					Status:                "approved",
					Reason:                "success",
					ExternalTransactionID: "ext-123",
					Installments:          1,
					AuthenticationMethod:  "3ds",
					ManualStatus:          "approved",
					AuthorizationCode:     "auth-123",
				},
			},
			Status:     "completed",
			UserID:     "user-123",
			CashChange: 0,
		}

		// Act
		dto := WebhookOrderDTO{
			Order: expectedOrder,
		}

		// Assert
		assert.Equal(t, expectedOrder, dto.Order)
		assert.Equal(t, expectedOrder.Token, dto.Order.Token)
		assert.Equal(t, expectedOrder.MerchantID, dto.Order.MerchantID)
		assert.Equal(t, expectedOrder.PaymentMethod, dto.Order.PaymentMethod)
		assert.Equal(t, expectedOrder.PaymentStatus, dto.Order.PaymentStatus)
		assert.Equal(t, expectedOrder.Currency, dto.Order.Currency)
		assert.Equal(t, expectedOrder.TaxAmount, dto.Order.TaxAmount)
		assert.Equal(t, expectedOrder.ItemsTotalAmount, dto.Order.ItemsTotalAmount)
		assert.Equal(t, expectedOrder.SubTotal, dto.Order.SubTotal)
		assert.Equal(t, expectedOrder.TotalAmount, dto.Order.TotalAmount)
		assert.Equal(t, expectedOrder.OrderID, dto.Order.OrderID)
		assert.Equal(t, expectedOrder.TransactionID, dto.Order.TransactionID)
		assert.Equal(t, expectedOrder.Metadata, dto.Order.Metadata)
		assert.Equal(t, expectedOrder.Status, dto.Order.Status)
		assert.Equal(t, expectedOrder.UserID, dto.Order.UserID)
		assert.Equal(t, expectedOrder.CashChange, dto.Order.CashChange)

		// Payment assertions
		assert.Equal(t, expectedOrder.Payment.Data.Metadata, dto.Order.Payment.Data.Metadata)
		assert.Equal(t, expectedOrder.Payment.Data.FromCard, dto.Order.Payment.Data.FromCard)
		assert.Equal(t, expectedOrder.Payment.Data.Amount, dto.Order.Payment.Data.Amount)
		assert.Equal(t, expectedOrder.Payment.Data.UpdatedAt, dto.Order.Payment.Data.UpdatedAt)
		assert.Equal(t, expectedOrder.Payment.Data.MethodType, dto.Order.Payment.Data.MethodType)
		assert.Equal(t, expectedOrder.Payment.Data.CreatedAt, dto.Order.Payment.Data.CreatedAt)
		assert.Equal(t, expectedOrder.Payment.Data.Merchant, dto.Order.Payment.Data.Merchant)
		assert.Equal(t, expectedOrder.Payment.Data.ID, dto.Order.Payment.Data.ID)
		assert.Equal(t, expectedOrder.Payment.Data.Processor, dto.Order.Payment.Data.Processor)
		assert.Equal(t, expectedOrder.Payment.Data.Customer, dto.Order.Payment.Data.Customer)
		assert.Equal(t, expectedOrder.Payment.Data.Status, dto.Order.Payment.Data.Status)
		assert.Equal(t, expectedOrder.Payment.Data.Reason, dto.Order.Payment.Data.Reason)
		assert.Equal(t, expectedOrder.Payment.Data.ExternalTransactionID, dto.Order.Payment.Data.ExternalTransactionID)
		assert.Equal(t, expectedOrder.Payment.Data.Installments, dto.Order.Payment.Data.Installments)
		assert.Equal(t, expectedOrder.Payment.Data.AuthenticationMethod, dto.Order.Payment.Data.AuthenticationMethod)
		assert.Equal(t, expectedOrder.Payment.Data.ManualStatus, dto.Order.Payment.Data.ManualStatus)
		assert.Equal(t, expectedOrder.Payment.Data.AuthorizationCode, dto.Order.Payment.Data.AuthorizationCode)
	})

	t.Run("should handle empty metadata", func(t *testing.T) {
		// Arrange
		order := Order{
			Metadata: map[string]interface{}{},
			Payment: Payment{
				Data: PaymentData{
					Metadata: map[string]interface{}{},
				},
			},
		}

		// Act
		dto := WebhookOrderDTO{
			Order: order,
		}

		// Assert
		assert.Empty(t, dto.Order.Metadata)
		assert.Empty(t, dto.Order.Payment.Data.Metadata)
	})

	t.Run("should handle zero amounts", func(t *testing.T) {
		// Arrange
		order := Order{
			TaxAmount:        0,
			ItemsTotalAmount: 0,
			SubTotal:         0,
			TotalAmount:      0,
			CashChange:       0,
			Payment: Payment{
				Data: PaymentData{
					Amount: MoneyAmount{
						Amount:   0,
						Currency: "MXN",
					},
				},
			},
		}

		// Act
		dto := WebhookOrderDTO{
			Order: order,
		}

		// Assert
		assert.Zero(t, dto.Order.TaxAmount)
		assert.Zero(t, dto.Order.ItemsTotalAmount)
		assert.Zero(t, dto.Order.SubTotal)
		assert.Zero(t, dto.Order.TotalAmount)
		assert.Zero(t, dto.Order.CashChange)
		assert.Zero(t, dto.Order.Payment.Data.Amount.Amount)
	})
}
