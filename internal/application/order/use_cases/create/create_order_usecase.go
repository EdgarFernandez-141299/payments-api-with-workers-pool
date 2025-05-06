package create

import (
	"context"
	"errors"
	"gitlab.com/clubhub.ai1/go-libraries/eventsourcing"
	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	errors2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/errors"

	log "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/observability/adapters"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/event_store"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/aggregate"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/response"
)

type CreateOrderUseCaseIF interface {
	Create(ctx context.Context, cmd command.CreateOrderCommand) (response.OrderResponseDTO, error)
}

type CreateOrderUseCaseImpl struct {
	repository event_store.OrderEventRepository
	log        log.Logger
}

func NewCreateOrderUseCase(
	logger log.Logger,
	repository event_store.OrderEventRepository,
) CreateOrderUseCaseIF {
	return &CreateOrderUseCaseImpl{
		repository: repository,
		log:        logger,
	}
}

func (c *CreateOrderUseCaseImpl) Create(oldCtx context.Context, cmd command.CreateOrderCommand) (response.OrderResponseDTO, error) {
	return decorators.TraceDecorator[response.OrderResponseDTO](
		oldCtx,
		"CreateOrderUseCaseImpl.Create",
		func(ctx context.Context, span decorators.Span) (response.OrderResponseDTO, error) {
			o := new(aggregate.Order)
			if err := c.repository.Get(ctx, cmd.ReferenceID, o); err != nil {
				if errors.Is(err, eventsourcing.ErrAggregateAlreadyExists) {
					return response.OrderResponseDTO{}, errors2.NewOrderAlreadyExistError(cmd.ReferenceID)
				}

				if !errors.Is(err, eventsourcing.ErrAggregateNotFound) {
					return response.OrderResponseDTO{}, err
				}
			}

			order, err := aggregate.Create(cmd)
			if err != nil {
				return response.OrderResponseDTO{}, err
			}

			err = c.repository.Create(ctx, order)

			if err != nil {
				return response.OrderResponseDTO{}, err
			}

			return response.OrderResponseDTO{
				ReferenceOrderID: order.ID,
			}, nil
		})
}
