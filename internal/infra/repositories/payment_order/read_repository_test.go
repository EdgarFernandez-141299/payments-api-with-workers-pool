package repositories

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/gommon/uid"
	orderEntities "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/order/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_order/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	idOrder, _ = uid.NewUniqueID(uid.WithID("id"))
)

// setupTestDB sets up an in-memory database for testing
func setupTestReadDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Migrar ambas tablas: PaymentOrderEntity y OrderEntity
	err = db.AutoMigrate(&entities.PaymentOrderEntity{}, &orderEntities.OrderEntity{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestPaymentOrderRepository_GetPaymentOrderByID(t *testing.T) {
	db, err := setupTestReadDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}

	repo := NewPaymentOrderReadRepository(db)

	t.Run("Payment order exists", func(t *testing.T) {
		ctx := context.Background()

		order := orderEntities.OrderEntity{
			ID:               idOrder.String(),
			ReferenceOrderID: "referenceID",
			EnterpriseID:     "enterpriseID",
		}
		result := db.Create(&order)
		assert.NoError(t, result.Error)

		payment := entities.PaymentOrderEntity{
			ID:             idOrder,
			OrderID:        idOrder.String(),
			PaymentOrderID: "referenceID",
			EnterpriseID:   "enterpriseID",
			Reference:      "referenceID",
		}
		result = db.Create(&payment)
		assert.NoError(t, result.Error)

		// Verificar que los registros se crearon correctamente
		var checkOrder orderEntities.OrderEntity
		result = db.First(&checkOrder, "id = ?", idOrder)
		assert.NoError(t, result.Error)

		var checkPayment entities.PaymentOrderEntity
		result = db.First(&checkPayment, "id = ?", idOrder)
		assert.NoError(t, result.Error)

		found, err := repo.GetPaymentOrderByReference(
			ctx,
			"referenceID",
			"referenceID",
			"enterpriseID",
		)

		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, payment.ID, found.ID)
	})
}
