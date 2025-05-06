package card

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/gommon/uid"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/card/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/card/entities"
	"gorm.io/gorm"
)

func TestCardRepository_CreateCard(t *testing.T) {
	t.Run("should create a card", func(t *testing.T) {
		db, err := setupTestDB()
		if err != nil {
			t.Fatalf("failed to setup test db: %v", err)
		}
		repo := NewCardWriteRepository(db)

		ctx := context.Background()
		card := entities.CardEntity{
			ID:          uid.GenerateID(),
			UserID:      "test-user-id",
			Status:      enums.Active.String(),
			IsRecurrent: false,
		}

		err = repo.CreateCard(ctx, &card)

		var result entities.CardEntity
		db.First(&result, "id = ?", card.ID)
		assert.Equal(t, card.ID, result.ID)
		assert.Equal(t, enums.Active.String(), result.Status)
		assert.False(t, result.IsRecurrent)
	})

	t.Run("should create a card when card be default", func(t *testing.T) {
		db, err := setupTestDB()
		if err != nil {
			t.Fatalf("failed to setup test db: %v", err)
		}
		repo := NewCardWriteRepository(db)

		ctx := context.Background()
		card := entities.CardEntity{
			ID:          uid.GenerateID(),
			UserID:      "test-user-id",
			Status:      enums.Default.String(),
			IsRecurrent: false,
		}

		err = repo.CreateCard(ctx, &card)

		var result entities.CardEntity
		db.First(&result, "id = ?", card.ID)
		assert.Equal(t, card.ID, result.ID)
		assert.Equal(t, enums.Default.String(), result.Status)
		assert.False(t, result.IsRecurrent)
	})

	t.Run("should create a card when card is recurrent", func(t *testing.T) {
		db, err := setupTestDB()
		if err != nil {
			t.Fatalf("failed to setup test db: %v", err)
		}
		repo := NewCardWriteRepository(db)

		ctx := context.Background()
		card := entities.CardEntity{
			ID:          uid.GenerateID(),
			UserID:      "test-user-id",
			Status:      enums.Default.String(),
			IsRecurrent: true,
		}

		err = repo.CreateCard(ctx, &card)

		var result entities.CardEntity
		db.First(&result, "id = ?", card.ID)
		assert.Equal(t, card.ID, result.ID)
		assert.Equal(t, enums.Default.String(), result.Status)
		assert.True(t, result.IsRecurrent)
		assert.NoError(t, err)
	})
}

func TestCardRepository_DeleteCard(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	repo := NewCardWriteRepository(db)

	ctx := context.Background()
	card := entities.CardEntity{
		ID:     uid.GenerateID(),
		UserID: "test-user-id",
	}

	db.Create(&card)
	err = repo.DeleteCard(ctx, card.ID.String())
	assert.NoError(t, err)

	var result entities.CardEntity
	err = db.First(&result, "external_card_id = ?", card.ID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
