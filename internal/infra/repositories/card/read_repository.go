package card

import (
	"context"
	"time"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/card/projections"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/card/entities"
	"gorm.io/gorm"
)

type CardReadRepository struct {
	db *gorm.DB
}

func NewCardReadRepository(db *gorm.DB) repository.CardReadRepositoryIF {
	return &CardReadRepository{
		db: db,
	}
}

func (c *CardReadRepository) GetCardsByUserID(
	ctx context.Context, userID, enterpriseID string,
) (entities.CardEntities, error) {
	var cards []entities.CardEntity
	err := c.db.WithContext(ctx).
		Where("user_id = ? AND enterprise_id = ?", userID, enterpriseID).
		Find(&cards).Error

	return cards, err
}

func (c *CardReadRepository) GetCardByUserID(ctx context.Context,
	userID, cardID, enterpriseID string) (entities.CardEntity, error) {
	var card entities.CardEntity
	err := c.db.WithContext(ctx).
		Where("user_id = ? AND card.id = ? AND card.enterprise_id = ?", userID, cardID, enterpriseID).
		First(&card).Error

	return card, err
}

func (c *CardReadRepository) CheckCardExistence(ctx context.Context, userID, enterpriseID, lastFour string) (bool, error) {
	var count int64
	err := c.db.WithContext(ctx).
		Model(&entities.CardEntity{}).
		Where("user_id = ? AND last_four = ? AND enterprise_id = ?", userID, lastFour, enterpriseID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (c *CardReadRepository) GetCardsExpiringSoon(
	ctx context.Context,
	expirationMonth time.Month,
	expirationYear int,
) ([]projections.NotificationCardExpiringSoonProjection, error) {
	var notificationCardExpiringSoonProjections []projections.NotificationCardExpiringSoonProjection

	startDate := time.Date(expirationYear, expirationMonth, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(expirationYear, expirationMonth+1, 0, 0, 0, 0, 0, time.UTC)

	err := c.db.WithContext(ctx).
		Model(&entities.CardEntity{}).
		Select(`"user".id as user_id, card.last_four, card.expiration_date, "user".email, "user".enterprise_id`).
		Joins(`JOIN "user" ON "user".id = card.user_id`).
		Where("card.expiration_date BETWEEN ? AND ?", startDate, endDate).
		Find(&notificationCardExpiringSoonProjections).Error

	return notificationCardExpiringSoonProjections, err
}

func (c *CardReadRepository) GetCardAndUserEmailByUserID(
	ctx context.Context,
	userID, cardID, enterpriseID string,
) (*projections.CardUserEmailProjection, error) {
	var cardUserEmailProjection projections.CardUserEmailProjection
	err := c.db.WithContext(ctx).
		Model(&entities.CardEntity{}).
		Select(`card.id, card.external_card_id, card.user_id, card.last_four, "user".email`).
		Joins(`JOIN "user" ON "user".id = card.user_id`).
		Where("card.user_id = ? AND card.id = ? AND card.enterprise_id = ?", userID, cardID, enterpriseID).
		First(&cardUserEmailProjection).Error

	return &cardUserEmailProjection, err
}
