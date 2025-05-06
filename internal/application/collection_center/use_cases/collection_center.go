package usecases

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/utils"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_center/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_center/dto/response"
	repositories "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/collection_center"
)

type CollectionCenterUsecaseIF interface {
	Create(ctx context.Context,
		request request.CollectionCenterRequest,
		enterpriseID string) (response.CollectionCenterResponse, error)
}

type CollectionCenterUsecase struct {
	repository repositories.CollectionCenterRepositoryIF
}

func NewCollectionCenterUsecase(repository repositories.CollectionCenterRepositoryIF) CollectionCenterUsecaseIF {
	return &CollectionCenterUsecase{
		repository: repository,
	}
}

func (c *CollectionCenterUsecase) Create(ctx context.Context,
	request request.CollectionCenterRequest,
	enterpriseID string,
) (response.CollectionCenterResponse, error) {
	request.AvailableCurrencies = utils.RemoveDuplicateCurrencies(request.AvailableCurrencies)

	collectionCenterEntity := entities.NewCollectionCenterEntity(
		request,
		enterpriseID,
	)

	err := utils.ValidateCurrencies(collectionCenterEntity.AvailableCurrencies)

	if err != nil {
		return response.CollectionCenterResponse{}, err
	}

	err = c.repository.Create(ctx, collectionCenterEntity)

	if err != nil {
		return response.CollectionCenterResponse{}, err
	}

	return response.NewCollectionCenterResponse(
		collectionCenterEntity,
	), nil
}
