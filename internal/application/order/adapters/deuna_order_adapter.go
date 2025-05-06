package adapters

import (
	"context"
	"os"

	"github.com/shopspring/decimal"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/resources/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/resources/dto/response"
	repository2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/utils"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
)

type OrderAdapterIF interface {
	CreateOrder(
		ctx context.Context, orderID string, paymentID string, currencyCode string, paymentTotal decimal.Decimal,
		paymentFlow enums.PaymentFlowEnum,
	) (response.DeunaOrderResponseDTO, error)
	GetOrderByReferenceID(ctx context.Context, referenceID string, enterpriseID string) (entities.Order, error)
}

type DeunaOrderAdapter struct {
	readRepository  repository.OrderReadRepositoryIF
	resource        resources.DeunaOrderResourceIF
	deunaRepository repository2.DeunaOrderRepository
}

func NewOrderAdapter(
	readRepository repository.OrderReadRepositoryIF,
	deunaRepository repository2.DeunaOrderRepository,
	resource resources.DeunaOrderResourceIF,
) OrderAdapterIF {
	return &DeunaOrderAdapter{
		readRepository:  readRepository,
		resource:        resource,
		deunaRepository: deunaRepository,
	}
}

func (a *DeunaOrderAdapter) CreateOrder(
	ctx context.Context, orderID string, paymentID string, currencyCode string, paymentTotal decimal.Decimal,
	paymentFlow enums.PaymentFlowEnum,
) (response.DeunaOrderResponseDTO, error) {
	if token, _ := a.deunaRepository.GetTokenByOrderAndPaymentID(ctx, orderID, paymentID); token != "" {
		return response.DeunaOrderResponseDTO{
			Token: token,
		}, nil
	}

	deunaFlowType, err := paymentFlow.DeunaFlowType()

	if err != nil {
		return response.DeunaOrderResponseDTO{}, err
	}

	deunaOrderID := utils.NewDeunaOrderID(orderID, paymentID)

	orderDTO := request.CreateDeunaOrderRequestDTO{
		Order: request.DeunaOrder{
			OrderID:     deunaOrderID.GetID(),
			Currency:    currencyCode,
			TotalAmount: utils.NewDeunaAmount(paymentTotal),
			SubTotal:    utils.NewDeunaAmount(paymentTotal),
			StoreCode:   "all",
			WebhooksURL: request.WebhooksURL{
				NotifyOrder: os.Getenv("DEUNA_NOTIFY_ORDER"),
			},
			Metadata: &map[string]interface{}{
				"additional_data.flow_type": deunaFlowType,
			},
		},
		OrderType: request.DeunaOrderType(request.DeUnaNow.String()),
	}

	deunaResponse, err := a.resource.CreateOrder(ctx, orderDTO)

	if err != nil {
		return response.DeunaOrderResponseDTO{}, err
	}

	createErr := a.deunaRepository.CreatePaymentOrderDeuna(ctx, paymentID, orderID, deunaResponse.Token)

	if createErr != nil {
		return response.DeunaOrderResponseDTO{}, createErr
	}

	return deunaResponse, nil
}

func (a *DeunaOrderAdapter) GetOrderByReferenceID(ctx context.Context, referenceID string, enterpriseID string) (entities.Order, error) {
	orderFound, err := a.readRepository.GetOrderByReferenceID(ctx, referenceID, enterpriseID)

	if err != nil {
		return entities.Order{}, err
	}

	return entities.Order{
		ID: orderFound.ID,
	}, nil
}
