package queries

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"

	log "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/observability/adapters"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/response"
)

type QueriesOrderUseCaseIF interface {
	GetOrderDetail(ctx context.Context, orderID, enterpriseID string) (response.GetOrderResponseDTO, error)
	GetOrderPayments(ctx context.Context, orderID, enterpriseID string) (*response.GetOrderPaymentResponseDTO, error)
}

type QueriesOrderUseCaseImpl struct {
	ReadRepository repository.OrderReadRepositoryIF
	log            log.Logger
}

func NewQueriesOrderUseCase(
	logger log.Logger,
	readRepository repository.OrderReadRepositoryIF,
) QueriesOrderUseCaseIF {
	return &QueriesOrderUseCaseImpl{
		ReadRepository: readRepository,
		log:            logger,
	}
}

func (c *QueriesOrderUseCaseImpl) GetOrderDetail(oldCtx context.Context, orderID, enterpriseID string) (response.GetOrderResponseDTO, error) {
	return decorators.TraceDecorator(
		oldCtx,
		"QueriesOrderUseCaseImpl.GetOrderDetail",
		func(ctx context.Context, span decorators.Span) (response.GetOrderResponseDTO, error) {
			order, err := c.ReadRepository.GetOrderByReferenceID(ctx, orderID, enterpriseID)

			if err != nil {
				return response.GetOrderResponseDTO{}, err
			}

			metadata := make(map[string]string)

			err = json.Unmarshal([]byte(order.Metadata), &metadata)

			if err != nil {
				return response.GetOrderResponseDTO{}, err
			}

			return response.GetOrderResponseDTO{
				ReferenceOrderID: order.ReferenceOrderID,
				Status:           enums.PaymentStatus(order.Status),
				Total:            order.TotalAmount,
				Currency:         order.CurrencyCode,
				CountryCode:      order.CountryCode,
				Metadata:         metadata,
			}, nil
		})
}

func (c *QueriesOrderUseCaseImpl) GetOrderPayments(ctx context.Context, orderID, enterpriseID string) (*response.GetOrderPaymentResponseDTO, error) {
	return decorators.TraceDecorator(
		ctx,
		"QueriesOrderUseCaseImpl.GetOrderPayments",
		func(ctx context.Context, span decorators.Span) (*response.GetOrderPaymentResponseDTO, error) {
			orderPaymentProjection, getOrderPaymentErr := c.ReadRepository.GetOrderPayments(ctx, orderID, enterpriseID)

			if getOrderPaymentErr != nil {
				return nil, getOrderPaymentErr
			}

			if len(orderPaymentProjection) == 0 {
				return nil, errors.New("order not found with id: " + orderID)
			}

			first := orderPaymentProjection[0]
			payments := make([]response.PaymentDTO, 0, len(orderPaymentProjection))

			metadata := make(map[string]string)
			if first.Metadata != "" && json.Unmarshal([]byte(first.Metadata), &metadata) != nil {
				return nil, fmt.Errorf("error unmarshalling metadata for order: %s", orderID)
			}

			for _, payment := range orderPaymentProjection {
				if payment.PaymentID != "" {
					payments = append(payments, response.PaymentDTO{
						ID:                payment.PaymentID,
						Status:            payment.PaymentStatus,
						PaymentMethod:     payment.PaymentMethod,
						CardID:            payment.CardID,
						AuthorizationCode: payment.AuthorizationCode,
						PaymentOrderID:    payment.PaymentOrderID,
					})
				}
			}

			return &response.GetOrderPaymentResponseDTO{
				ReferenceOrderID: first.ReferenceOrderID,
				Status:           enums.PaymentStatus(first.OrderStatus),
				Total:            first.TotalAmount,
				Currency:         first.CurrencyCode,
				CountryCode:      first.CountryCode,
				Metadata:         metadata,
				Payments:         payments,
			}, nil

		},
	)
}
