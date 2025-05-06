package adapter

import (
	"context"

	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	orderAdapter "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters"
	orderResponse "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/resources/dto/response"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/resources/dto/request"
	loginAdapter "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/aggregate"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
	paymentOrderEntity "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_order/entities"
)

const (
	storeCode  = "all"
	creditCard = "credit_card"
)

type PaymentOrderAdapterIF interface {
	CreatePaymentOrder(
		ctx context.Context,
		cmd command.CreatePaymentOrderCommand,
		order *aggregate.Order,
		card entities.Card,
	) error
}

type PaymentOrderAdapter struct {
	resource     resources.DeunaPaymentResourceIF
	orderAdapter orderAdapter.OrderAdapterIF
	loginAdapter loginAdapter.DeunaLoginAdapter
	readModel    repository.PaymentOrderRepositoryIF
}

func NewPaymentOrderAdapter(
	resource resources.DeunaPaymentResourceIF,
	orderAdapter orderAdapter.OrderAdapterIF,
	loginAdapter loginAdapter.DeunaLoginAdapter,
	readModel repository.PaymentOrderRepositoryIF,
) PaymentOrderAdapterIF {
	return &PaymentOrderAdapter{
		resource:     resource,
		orderAdapter: orderAdapter,
		loginAdapter: loginAdapter,
		readModel:    readModel,
	}
}

func (a *PaymentOrderAdapter) CreatePaymentOrder(
	oldCtx context.Context,
	cmd command.CreatePaymentOrderCommand,
	order *aggregate.Order,
	card entities.Card,
) error {
	return decorators.TraceDecoratorNoReturn(
		oldCtx,
		"PaymentOrderAdapter.CreatePaymentOrder",
		func(ctx context.Context, span decorators.Span) error {
			orderCreated, err := a.createOrderDeUNA(ctx, cmd)

			if err != nil {
				return err
			}

			err = a.createPaymentOrderDeUNA(ctx, cmd, orderCreated.Token, card, order)

			if err != nil {
				return err
			}

			purchaseOrderEntity := paymentOrderEntity.NewPaymentOrderEntity().
				SetOrderID(order.ID).
				SetAssociatedOrigin(cmd.AssociatedOrigin.Type.String()).
				SetPaymentMethod(cmd.Payment.Method.Type.String()).
				SetCurrencyCode(cmd.CurrencyCode.Code).
				SetCountryCode(cmd.CountryCode).
				SetCardID(card.ID).
				SetCollectionAccountID(cmd.CollectionAccount.ID).
				SetTransactionDate().
				SetPaymentOrderID(cmd.Payment.ID).
				SetEnterpriseID(order.EnterpriseID).
				SetTotalAmount(cmd.Payment.Total.Value).
				SetPaymentFlow(cmd.PaymentFlow).
				SetStatus(enums.PaymentProcessing.String()).
				Build()

			return a.readModel.CreatePaymentOrder(ctx, purchaseOrderEntity)
		},
	)
}

func (a *PaymentOrderAdapter) createOrderDeUNA(oldCtx context.Context, cmd command.CreatePaymentOrderCommand) (orderResponse.DeunaOrderResponseDTO, error) {
	return decorators.TraceDecorator[orderResponse.DeunaOrderResponseDTO](
		oldCtx,
		"PaymentOrderAdapter.createOrderDeUNA",
		func(ctx context.Context, span decorators.Span) (orderResponse.DeunaOrderResponseDTO, error) {

			orderDeUNA, err := a.orderAdapter.CreateOrder(
				ctx, cmd.ReferenceOrderID, cmd.Payment.ID, cmd.CurrencyCode.Code,
				cmd.Payment.Total.Value, cmd.PaymentFlow,
			)

			if err != nil {
				return orderResponse.DeunaOrderResponseDTO{}, err
			}

			return orderDeUNA, nil
		})
}

func (a *PaymentOrderAdapter) createPaymentOrderDeUNA(
	oldCtx context.Context,
	cmd command.CreatePaymentOrderCommand,
	orderToken string,
	cardEntity entities.Card,
	order *aggregate.Order,
) error {
	return decorators.TraceDecoratorNoReturn(
		oldCtx,
		"PaymentOrderAdapter.createPaymentOrderDeUNA",
		func(ctx context.Context, span decorators.Span) error {
			medata := map[string]interface{}{}
			card := cmd.Payment.Method.CCData.Data

			newPayment := request.DeunaOrderPaymentRequest{
				OrderToken: orderToken,
				CardID:     cardEntity.ExternalCardID,
				Email:      order.Email,
				CreditCard: &request.CreditCardInfo{
					CardCVV: &card.CVV,
				},
				Metadata:   &medata,
				MethodType: creditCard,
				StoreCode:  storeCode,
			}

			token, err := a.loginAdapter.LoginWitUserID(ctx, cmd.User.ID, order.EnterpriseID)
			if err != nil {
				return err
			}

			makeOrderErr := a.resource.MakeOrderPayment(ctx, newPayment, token)

			return makeOrderErr
		})
}
