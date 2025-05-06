package usecases

import (
	"context"
	"errors"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/payment_concept/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/payment_concept/dto/response"
	repository "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_concept"
)

type PaymentConceptUsecaseIF interface {
	Create(
		ctx context.Context,
		payment request.PaymentConceptRequest,
		enterpriseId string,
	) (response.PaymentConceptResponse, error)
}

type PaymentConceptUsecaseImpl struct {
	repository repository.PaymentConceptRepositoryIF
}

func NewPaymentConceptUsecase(repository repository.PaymentConceptRepositoryIF) PaymentConceptUsecaseIF {
	return &PaymentConceptUsecaseImpl{
		repository: repository,
	}
}

func (p *PaymentConceptUsecaseImpl) Create(
	ctx context.Context,
	payment request.PaymentConceptRequest,
	enterpriseId string,
) (response.PaymentConceptResponse, error) {
	entity := entities.NewPaymentConceptEntity(payment, enterpriseId)
	err := p.repository.Create(ctx, entity)

	if err != nil {
		return response.PaymentConceptResponse{}, errors.New("error creating payment concept")
	}

	return response.NewPaymentConceptResponse(
		entity.ID.String(),
		entity.Name,
		entity.Code,
		entity.Description,
		entity.CreatedAt.String(),
	), nil
}
