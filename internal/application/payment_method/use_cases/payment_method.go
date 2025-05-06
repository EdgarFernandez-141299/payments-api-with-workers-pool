package usecases

import (
	"context"
	"fmt"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/payment_method/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/payment_method/dto/response"
	repositories "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_method"
)

type PaymentMethodUseCasesIF interface {
	Create(
		ctx context.Context,
		request request.PaymentMethodRequest,
		enterpriseId string,
	) (response.PaymentMethodResponse, error)
}

type PaymentMethodUseCases struct {
	repository repositories.PaymentMethodRepositoryIF
}

func NewPaymentMethodUseCases(repository repositories.PaymentMethodRepositoryIF) PaymentMethodUseCasesIF {
	return &PaymentMethodUseCases{
		repository: repository,
	}
}

func (p *PaymentMethodUseCases) Create(
	ctx context.Context,
	request request.PaymentMethodRequest,
	enterpriseId string,
) (response.PaymentMethodResponse, error) {
	entity := entities.NewPaymentMethodEntity(
		request.Name,
		request.Code,
		request.Description,
		enterpriseId,
	)

	if err := p.repository.Create(ctx, entity); err != nil {
		return response.PaymentMethodResponse{}, fmt.Errorf("error creating payment method")
	}

	return response.NewPaymentMethodResponse(
		entity.ID.String(),
		entity.Name,
		entity.Code,
		entity.Description,
	), nil
}
