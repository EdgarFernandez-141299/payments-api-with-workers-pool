package create

import (
	"context"

	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	queriesCard "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/usecases/queries"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/collection_account/queries"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/event_store"
	adapter "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/aggregate"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/response"
)

type CreatePaymentOrderUseCaseIF interface {
	CreatePaymentOrder(
		ctx context.Context,
		cmd command.CreatePaymentOrderCommand,
	) (response.PaymentOrderResponseDTO, error)
}

type CreatePaymentOrderUseCaseImpl struct {
	repository             event_store.OrderEventRepository
	collectionAccountQuery queries.GetCollectionAccountByRouteUsecaseIF
	cardQuery              queriesCard.GetCardByUserUsecaseIF
	paymentOrderAdapter    adapter.PaymentOrderAdapterIF
}

func NewCreatePaymentOrderUseCase(
	repository event_store.OrderEventRepository,
	collectionAccountQuery queries.GetCollectionAccountByRouteUsecaseIF,
	cardQuery queriesCard.GetCardByUserUsecaseIF,
	paymentOrderAdapter adapter.PaymentOrderAdapterIF,
) CreatePaymentOrderUseCaseIF {
	return &CreatePaymentOrderUseCaseImpl{
		repository:             repository,
		collectionAccountQuery: collectionAccountQuery,
		cardQuery:              cardQuery,
		paymentOrderAdapter:    paymentOrderAdapter,
	}
}

func (p *CreatePaymentOrderUseCaseImpl) CreatePaymentOrder(
	ctx context.Context,
	cmd command.CreatePaymentOrderCommand,
) (response.PaymentOrderResponseDTO, error) {
	return decorators.TraceDecorator[response.PaymentOrderResponseDTO](
		ctx,
		"CreatePaymentOrderUseCaseImpl.CreatePaymentOrder",
		func(ctx context.Context, span decorators.Span) (response.PaymentOrderResponseDTO, error) {
			order := new(aggregate.Order)

			err := p.repository.Get(ctx, cmd.ReferenceOrderID, order)

			if err != nil {
				return response.PaymentOrderResponseDTO{}, err
			}

			collectionAccount, err := p.collectionAccountQuery.GetCollectionAccountByRoute(
				ctx,
				order.CountryCode.Iso3(),
				cmd.AssociatedOrigin.Type.String(),
				cmd.CurrencyCode.Code,
				order.EnterpriseID,
			)

			if err != nil {
				return response.PaymentOrderResponseDTO{}, err
			}

			cmd = cmd.WithCollectionAccount(collectionAccount)

			card, err := p.cardQuery.GetCardByIDAndUserID(ctx,
				cmd.User.ID,
				cmd.Payment.Method.CCData.Data.CardID,
				order.EnterpriseID,
			)

			if err != nil {
				return response.PaymentOrderResponseDTO{}, err
			}

			paymentFlow, err := enums.NewPaymentFlowEnum(card.CardType, order.AllowCapture)

			if err != nil {
				return response.PaymentOrderResponseDTO{}, err
			}

			cmd = cmd.WithPaymentFlow(paymentFlow)

			cmd = cmd.WithCardData(card)

			err = order.StartProcessingOrderPayment(cmd)

			if err != nil {
				return response.PaymentOrderResponseDTO{}, err
			}

			adapterErr := p.paymentOrderAdapter.CreatePaymentOrder(ctx, cmd, order, card)

			if adapterErr != nil {
				return response.PaymentOrderResponseDTO{}, adapterErr
			}

			repositoryErr := p.repository.Save(ctx, order)

			if repositoryErr != nil {
				return response.PaymentOrderResponseDTO{}, repositoryErr
			}

			return response.PaymentOrderResponseDTO{
				ReferenceOrderID: order.ID,
			}, nil
		})
}
