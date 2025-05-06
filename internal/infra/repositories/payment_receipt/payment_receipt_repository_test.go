package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/common/testutils"
	domainEntities "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/payment_receipt/command"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_receipt/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

const testReceiptUrl = "https://example.com/receipt.pdf"

func setupTestDB(t *testing.T) (*gorm.DB, *testutils.PostgresContainer, error) {
	container, err := testutils.SetupPostgresContainer(t)
	if err != nil {
		return nil, nil, err
	}

	db, err := gorm.Open(postgres.Open(container.URI), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	err = db.AutoMigrate(&entities.PaymentReceiptDTO{})
	if err != nil {
		return nil, nil, err
	}

	return db, container, nil
}

func TestGetByPaymentID(t *testing.T) {
	db, container, err := setupTestDB(t)
	require.NoError(t, err)

	defer func() {
		ctx := context.Background()
		container.Terminate(ctx)
	}()

	repo := NewPaymentReceiptRepository(db)

	t.Run("Successfully retrieves a payment receipt by payment ID", func(t *testing.T) {
		ctx := context.Background()

		currencyCode, err := value_objects.NewCurrencyCode("USD")
		require.NoError(t, err)

		amount, err := value_objects.NewCurrencyAmount(currencyCode, decimal.NewFromFloat(100.50))
		require.NoError(t, err)

		country, err := value_objects.NewCountryWithCode("USA")
		require.NoError(t, err)

		now := time.Now()

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

		cmd := command.CreatePaymentReceiptCommand{
			UserID:           "test-user-id",
			EnterpriseID:     "test-enterprise-id",
			Email:            "test@example.com",
			ReferenceOrderID: "test-order-id",
			PaymentID:        "test-payment-id-get",
			PaymentStatus:    "completed",
			PaymentAmount:    amount,
			PaymentCountry:   country,
			PaymentMethod:    paymentMethod,
			PaymentDate:      now.Format(time.RFC3339),
			CreatedAt:        now,
			UpdatedAt:        now,
		}

		entity := domainEntities.NewPaymentReceiptEntity(cmd).WithReceiptURL(testReceiptUrl)
		err = repo.CreatePaymentReceipt(ctx, entity)
		require.NoError(t, err)

		receipt, err := repo.GetByPaymentID(ctx, cmd.PaymentID)
		require.NoError(t, err)

		assert.Equal(t, cmd.UserID, receipt.UserID)
		assert.Equal(t, cmd.EnterpriseID, receipt.EnterpriseID)
		assert.Equal(t, cmd.Email, receipt.Email)
		assert.Equal(t, cmd.ReferenceOrderID, receipt.ReferenceOrderID)
		assert.Equal(t, cmd.PaymentID, receipt.PaymentID)
		assert.Equal(t, cmd.PaymentStatus, receipt.PaymentStatus)
		assert.Equal(t, cmd.PaymentAmount.Value, receipt.PaymentAmount)
		assert.Equal(t, cmd.PaymentCountry.Code, receipt.PaymentCountryCode)
		assert.Equal(t, cmd.PaymentAmount.Code.Code, receipt.PaymentCurrencyCode)
		assert.Equal(t, string(cmd.PaymentMethod.Type), string(receipt.PaymentMethod))
		assert.Equal(t, cmd.PaymentDate, receipt.PaymentDate)
		assert.Equal(t, testReceiptUrl, receipt.FileURL)
	})

	t.Run("Returns error when payment receipt does not exist", func(t *testing.T) {
		ctx := context.Background()

		_, err := repo.GetByPaymentID(ctx, "non-existent-payment-id")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "record not found")
	})
}

func TestCreatePaymentReceipt(t *testing.T) {
	db, container, err := setupTestDB(t)
	require.NoError(t, err)

	defer func() {
		ctx := context.Background()
		container.Terminate(ctx)
	}()

	repo := NewPaymentReceiptRepository(db)

	t.Run("Successfully creates a payment receipt", func(t *testing.T) {
		ctx := context.Background()

		currencyCode, err := value_objects.NewCurrencyCode("USD")
		require.NoError(t, err)

		amount, err := value_objects.NewCurrencyAmount(currencyCode, decimal.NewFromFloat(100.50))
		require.NoError(t, err)

		country, err := value_objects.NewCountryWithCode("USA")
		require.NoError(t, err)

		cardInfo := value_objects.CardInfo{
			CardID: "card-456",
			CVV:    "456",
		}

		paymentMethod := value_objects.PaymentMethod{
			Type: enums.CCMethod,
			CCData: value_objects.PaymentMethodData[value_objects.CardInfo]{
				Data: cardInfo,
			},
		}

		now := time.Now()
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

		entity := domainEntities.NewPaymentReceiptEntity(cmd).WithReceiptURL(testReceiptUrl)
		err = repo.CreatePaymentReceipt(ctx, entity)
		require.NoError(t, err)

		var count int64
		err = db.Model(&entities.PaymentReceiptDTO{}).Where("payment_id = ?", cmd.PaymentID).Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)

		receipt, err := repo.GetByPaymentID(ctx, cmd.PaymentID)
		require.NoError(t, err)

		assert.Equal(t, cmd.UserID, receipt.UserID)
		assert.Equal(t, cmd.EnterpriseID, receipt.EnterpriseID)
		assert.Equal(t, cmd.Email, receipt.Email)
		assert.Equal(t, cmd.ReferenceOrderID, receipt.ReferenceOrderID)
		assert.Equal(t, cmd.PaymentID, receipt.PaymentID)
		assert.Equal(t, cmd.PaymentStatus, receipt.PaymentStatus)
		assert.Equal(t, cmd.PaymentAmount.Value, receipt.PaymentAmount)
		assert.Equal(t, cmd.PaymentCountry.Code, receipt.PaymentCountryCode)
		assert.Equal(t, cmd.PaymentAmount.Code.Code, receipt.PaymentCurrencyCode)
		assert.Equal(t, string(cmd.PaymentMethod.Type), string(receipt.PaymentMethod))
		assert.Equal(t, cmd.PaymentDate, receipt.PaymentDate)
		assert.Equal(t, testReceiptUrl, receipt.FileURL)
	})

	t.Run("Returns error when database operation fails", func(t *testing.T) {
		ctx := context.Background()

		invalidDB, err := gorm.Open(postgres.Open(container.URI), &gorm.Config{})
		require.NoError(t, err)

		err = invalidDB.Exec("DROP TABLE IF EXISTS payment_receipt").Error
		require.NoError(t, err)

		err = invalidDB.Exec("CREATE TABLE payment_receipt (id VARCHAR(14) PRIMARY KEY)").Error
		require.NoError(t, err)

		mockRepo := NewPaymentReceiptRepository(invalidDB)

		currencyCode, err := value_objects.NewCurrencyCode("USD")
		require.NoError(t, err)

		amount, err := value_objects.NewCurrencyAmount(currencyCode, decimal.NewFromFloat(100.50))
		require.NoError(t, err)

		country, err := value_objects.NewCountryWithCode("USA")
		require.NoError(t, err)

		cardInfo := value_objects.CardInfo{
			CardID: "card-789",
			CVV:    "789",
		}

		paymentMethod := value_objects.PaymentMethod{
			Type: enums.CCMethod,
			CCData: value_objects.PaymentMethodData[value_objects.CardInfo]{
				Data: cardInfo,
			},
		}

		now := time.Now()
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

		entity := domainEntities.NewPaymentReceiptEntity(cmd)
		err = mockRepo.CreatePaymentReceipt(ctx, entity)
		require.Error(t, err)
	})
}
