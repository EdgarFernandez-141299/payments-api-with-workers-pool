package card

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/gommon/uid"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/card/entities"
	userEntities "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/user/entities"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&entities.CardEntity{}, &userEntities.UserEntity{})
	return db, nil
}

func TestCheckCardExistence(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}

	repo := NewCardReadRepository(db)

	t.Run("Card is duplicate", func(t *testing.T) {
		ctx := context.Background()
		userID := "user123"
		enterpriseID := "enterprise123"
		lastFour := "1234"

		card := entities.CardEntity{
			UserID:       userID,
			EnterpriseID: enterpriseID,
			LastFour:     lastFour,
		}
		db.Create(&card)

		isDuplicate, err := repo.CheckCardExistence(ctx, userID, enterpriseID, lastFour)
		assert.NoError(t, err)
		assert.True(t, isDuplicate)
	})

	t.Run("Card is not duplicate", func(t *testing.T) {
		ctx := context.Background()
		userID := "user123"
		enterpriseID := "enterprise123"
		lastFour := "5678"

		isDuplicate, err := repo.CheckCardExistence(ctx, userID, enterpriseID, lastFour)
		assert.NoError(t, err)
		assert.False(t, isDuplicate)
	})

	t.Run("Error querying database", func(t *testing.T) {
		ctx := context.Background()
		userID := "user123"
		enterpriseID := "enterprise123"
		lastFour := "1234"

		db.AddError(errors.New("query error"))

		isDuplicate, err := repo.CheckCardExistence(ctx, userID, enterpriseID, lastFour)
		assert.Error(t, err)
		assert.False(t, isDuplicate)
	})
}

func TestCardRepository_GetCardsByUserID(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	repo := NewCardReadRepository(db)

	ctx := context.Background()

	user := userEntities.UserEntity{
		ID:             "test-user-id",
		UserType:       "test-user-type",
		EnterpriseID:   "test-enterprise-id",
		ExternalUserID: "test-external-user-id",
	}

	card := entities.CardEntity{
		ID:           uid.GenerateID(),
		UserID:       "test-user-id",
		CardHolder:   "John Doe",
		Bin:          "123456",
		LastFour:     "7890",
		EnterpriseID: "test-enterprise-id",
	}

	db.Create(&user)
	db.Create(&card)
	cards, err := repo.GetCardsByUserID(ctx, card.UserID, card.EnterpriseID)
	assert.NoError(t, err)
	assert.Len(t, cards, 1)
	assert.Equal(t, card.ID, cards[0].ID)
}

func TestCardRepository_GetCardByUserID(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	repo := NewCardReadRepository(db)

	ctx := context.Background()

	user := userEntities.UserEntity{
		ID:             "test-user-id",
		UserType:       "test-user-type",
		EnterpriseID:   "test-enterprise-id",
		ExternalUserID: "test-external-user-id",
	}

	card := entities.CardEntity{
		ID:     uid.GenerateID(),
		UserID: "test-user-id",
	}

	db.Create(&user)
	db.Create(&card)
	result, err := repo.GetCardByUserID(ctx, card.UserID, card.ID.String(), card.EnterpriseID)
	assert.NoError(t, err)
	assert.Equal(t, card.ID, result.ID)
}

func TestCardRepository_GetCardsExpiringSoon(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}

	repo := NewCardReadRepository(db)

	ctx := context.Background()

	user := userEntities.UserEntity{
		ID:             "test-user-id",
		UserType:       "test-user-type",
		EnterpriseID:   "test-enterprise-id",
		ExternalUserID: "test-external-user-id",
	}

	card := entities.CardEntity{
		ID:             uid.GenerateID(),
		UserID:         "test-user-id",
		CardHolder:     "John Doe",
		Bin:            "123456",
		LastFour:       "7890",
		EnterpriseID:   "test-enterprise-id",
		ExpirationDate: time.Now().AddDate(0, 1, 0),
	}

	db.Create(&user)
	db.Create(&card)

	expiringSoonCards, err := repo.GetCardsExpiringSoon(ctx, card.ExpirationDate.Month(), card.ExpirationDate.Year())
	assert.NoError(t, err)
	assert.Len(t, expiringSoonCards, 1)
}

func TestCardRepository_GetCardAndUserEmailByUserID(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}

	repo := NewCardReadRepository(db)

	ctx := context.Background()

	user := userEntities.UserEntity{
		ID:             "test-user-id",
		UserType:       "test-user-type",
		EnterpriseID:   "test-enterprise-id",
		ExternalUserID: "test-external-user-id",
	}

	card := entities.CardEntity{
		ID:           uid.GenerateID(),
		UserID:       "test-user-id",
		CardHolder:   "John Doe",
		Bin:          "123456",
		LastFour:     "7890",
		EnterpriseID: "test-enterprise-id",
	}

	db.Create(&user)
	db.Create(&card)

	result, err := repo.GetCardAndUserEmailByUserID(ctx, card.UserID, card.ID.String(), card.EnterpriseID)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, result.Email)
}
