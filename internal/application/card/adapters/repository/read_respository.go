package repository

import (
	"context"
	"time"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/card/projections"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/card/entities"
)

type CardReadRepositoryIF interface {
	GetCardsByUserID(ctx context.Context, userID, enterpriseID string) (entities.CardEntities, error)
	GetCardByUserID(ctx context.Context, userID, cardID, enterpriseID string) (entities.CardEntity, error)
	CheckCardExistence(ctx context.Context, userID, enterpriseID, lastFour string) (bool, error)
	GetCardsExpiringSoon(
		ctx context.Context,
		expirationMonth time.Month,
		expirationYear int,
	) ([]projections.NotificationCardExpiringSoonProjection, error)
	GetCardAndUserEmailByUserID(
		ctx context.Context,
		userID,
		cardID,
		enterpriseID string,
	) (*projections.CardUserEmailProjection, error)
}
