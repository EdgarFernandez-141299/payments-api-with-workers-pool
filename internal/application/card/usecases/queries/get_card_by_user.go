package queries

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/adapters/repository"
	errorsBusiness "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/card/errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
)

type GetCardByUserUsecaseIF interface {
	GetCardByIDAndUserID(ctx context.Context, userID, cardID, enterpriseID string) (entities.Card, error)
}

type GetCardByUserUsecase struct {
	repository repository.CardReadRepositoryIF
}

func NewGetCardByUserUsecase(
	repository repository.CardReadRepositoryIF,
) GetCardByUserUsecaseIF {
	return &GetCardByUserUsecase{
		repository: repository,
	}
}

func (p *GetCardByUserUsecase) GetCardByIDAndUserID(
	ctx context.Context, userID, cardID, enterpriseID string,
) (entities.Card, error) {
	card, err := p.repository.GetCardByUserID(ctx, userID, cardID, enterpriseID)

	if err != nil {
		return entities.Card{}, errorsBusiness.NewCardIsNotFromMember(userID, err)
	}

	return entities.Card{
		ID:                card.ID.String(),
		ExternalCardID:    card.ExternalCardID,
		UserID:            card.UserID,
		CardHolder:        card.CardHolder,
		Alias:             card.Alias,
		Bin:               card.Bin,
		LastFour:          card.LastFour,
		Brand:             card.Brand,
		ExpirationDate:    card.ExpirationDate.String(),
		CardType:          card.CardType,
		Status:            card.Status,
		IsDefault:         card.IsDefault,
		IsRecurrent:       card.IsRecurrent,
		RetryAttempts:     card.RetryAttempts,
		EnterpriseID:      card.EnterpriseID,
		CardFailureReason: card.CardFailureReason,
		CardFailureCode:   card.CardFailureCode,
	}, nil
}
