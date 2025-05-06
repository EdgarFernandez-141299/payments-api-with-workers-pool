package refund

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/gommon/uid"

	orderEntities "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/order/entities"

	paymentEntities "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_order/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/refund/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	testOrderID, _   = uid.NewUniqueID(uid.WithID("order-id"))
	testPaymentID, _ = uid.NewUniqueID(uid.WithID("payment-id"))
	testRefundID, _  = uid.NewUniqueID(uid.WithID("refund-id"))
	testRefundID2, _ = uid.NewUniqueID(uid.WithID("refund-id-2"))
)

// setupTestDB sets up an in-memory database for testing
func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate all required tables
	err = db.AutoMigrate(
		&orderEntities.OrderEntity{},
		&paymentEntities.PaymentOrderEntity{},
		&entities.RefundEntity{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestRefundReadRepository_GetRefundsByReferenceOrderID(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}

	repo := NewRefundReadRepository(db)

	t.Run("Refunds exist for order", func(t *testing.T) {
		ctx := context.Background()

		// Create test order
		order := orderEntities.OrderEntity{
			ID:               testOrderID.String(),
			ReferenceOrderID: "reference-order-1",
			EnterpriseID:     "enterprise-1",
		}
		result := db.Create(&order)
		assert.NoError(t, result.Error)

		// Create test payment
		payment := paymentEntities.PaymentOrderEntity{
			ID:             testPaymentID,
			OrderID:        testOrderID.String(),
			PaymentOrderID: "payment-1",
			EnterpriseID:   "enterprise-1",
		}
		result = db.Create(&payment)
		assert.NoError(t, result.Error)

		// Create test refunds
		refund1 := entities.RefundEntity{
			ID:        testRefundID,
			PaymentID: testPaymentID.String(),
			Amount:    decimal.NewFromFloat(50),
			Reason:    "partial refund",
		}
		result = db.Create(&refund1)
		assert.NoError(t, result.Error)

		refund2 := entities.RefundEntity{
			ID:        testRefundID2,
			PaymentID: testPaymentID.String(),
			Amount:    decimal.NewFromFloat(25),
			Reason:    "additional refund",
		}
		result = db.Create(&refund2)
		assert.NoError(t, result.Error)

		// Test repository method
		refunds, err := repo.GetRefundsByReferenceOrderID(
			ctx,
			"reference-order-1",
			testPaymentID.String(),
			"enterprise-1",
		)

		assert.NoError(t, err)
		assert.Len(t, refunds, 2)
		assert.Equal(t, testRefundID.String(), refunds[0].ID.String())
		assert.Equal(t, testRefundID2.String(), refunds[1].ID.String())
	})

	t.Run("No refunds exist for order", func(t *testing.T) {
		ctx := context.Background()

		// Generar un ID único para cada prueba
		uniqueOrderID, _ := uid.NewUniqueID(uid.WithID(fmt.Sprintf("order-id-%d", time.Now().UnixNano())))
		// Usar el ID único en la inserción
		result := db.Create(&orderEntities.OrderEntity{
			ID:               uniqueOrderID.String(),
			ReferenceOrderID: "reference-order-2",
			EnterpriseID:     "enterprise-1",
		})
		assert.NoError(t, result.Error)

		// Test repository method with non-existent refunds
		refunds, err := repo.GetRefundsByReferenceOrderID(
			ctx,
			"reference-order-2",
			"non-existent-payment",
			"enterprise-1",
		)

		assert.NoError(t, err)
		assert.Empty(t, refunds)
	})

	t.Run("Invalid enterprise ID", func(t *testing.T) {
		ctx := context.Background()

		// Test with wrong enterprise ID
		refunds, err := repo.GetRefundsByReferenceOrderID(
			ctx,
			"reference-order-1",
			testPaymentID.String(),
			"wrong-enterprise",
		)

		assert.NoError(t, err)
		assert.Empty(t, refunds)
	})
}
