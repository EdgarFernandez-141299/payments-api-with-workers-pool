package repository

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/card/entities"
)

type CardWriteRepositoryIF interface {
	CreateCard(ctx context.Context, entity *entities.CardEntity) error
	DeleteCard(ctx context.Context, cardId string) error
}
