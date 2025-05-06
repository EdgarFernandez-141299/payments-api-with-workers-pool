package card

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/card/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/card/entities"
	"gorm.io/gorm"
)

type CardWriteRepository struct {
	db *gorm.DB
}

func NewCardWriteRepository(db *gorm.DB) repository.CardWriteRepositoryIF {
	return &CardWriteRepository{
		db: db,
	}
}

func (c *CardWriteRepository) CreateCard(ctx context.Context, entity *entities.CardEntity) error {
	if entity.Status == enums.Default.String() {
		if err := c.updateDefaultCardsToActive(ctx, entity.UserID); err != nil {
			return err
		}
	}

	if entity.IsRecurrent {
		if err := c.updateRecurrentCardsToFalse(ctx, entity.UserID); err != nil {
			return err
		}
	}
	return c.db.WithContext(ctx).Create(&entity).Error
}

func (c *CardWriteRepository) updateRecurrentCardsToFalse(ctx context.Context, userId string) error {
	return c.db.WithContext(ctx).
		Model(&entities.CardEntity{}).
		Where("user_id = ? AND is_recurrent = ?", userId, true).
		Update("is_recurrent", false).Error
}

func (c *CardWriteRepository) DeleteCard(ctx context.Context, cardId string) error {
	return c.db.WithContext(ctx).Unscoped().Where("id = ?", cardId).
		Delete(&entities.CardEntity{}).Error
}

func (c *CardWriteRepository) updateDefaultCardsToActive(ctx context.Context, userId string) error {
	return c.db.WithContext(ctx).
		Model(&entities.CardEntity{}).
		Where("user_id = ? AND status = ?", userId, enums.Default.String()).
		Update("status", enums.Active.String()).Error
}
