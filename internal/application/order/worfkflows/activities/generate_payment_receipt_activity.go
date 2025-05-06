package activities

import (
	"context"
	"fmt"
	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	"gitlab.com/clubhub.ai1/go-libraries/saga/errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/event_store"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_receipt/use_case"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/aggregate"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/payment_receipt/command"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
	"time"
)

const GenerateReceiptActivityName = "generate_receipt_activity"

type GeneratePaymentReceiptActivity struct {
	useCase    use_case.GenerateReceiptPaymentUseCase
	repository event_store.OrderEventRepository
}

func NewGeneratePaymentReceiptActivity(
	useCase use_case.GenerateReceiptPaymentUseCase, repository event_store.OrderEventRepository,
) *GeneratePaymentReceiptActivity {
	return &GeneratePaymentReceiptActivity{
		useCase:    useCase,
		repository: repository,
	}
}

func (n *GeneratePaymentReceiptActivity) GenerateReceipt(oldCtx context.Context, referenceID string, paymentID string) (entities.PaymentReceipt, error) {
	return decorators.TraceDecorator[entities.PaymentReceipt](oldCtx, "GeneratePaymentReceiptActivity.GenerateReceipt",
		func(ctx context.Context, span decorators.Span) (entities.PaymentReceipt, error) {

			order := new(aggregate.Order)
			err := n.repository.Get(ctx, referenceID, order)

			if err != nil {
				return entities.PaymentReceipt{}, errors.WrapActivityError(err)
			}

			payment, err := order.FindPaymentByID(paymentID)
			if err != nil {
				return entities.PaymentReceipt{}, errors.WrapActivityError(fmt.Errorf("payment not found"))
			}

			paymentAmount, _ := value_objects.NewCurrencyAmount(order.Currency, payment.Total.Value)

			cmd := command.CreatePaymentReceiptCommand{
				UserID:           order.User.ID,
				EnterpriseID:     order.EnterpriseID,
				Email:            order.Email,
				ReferenceOrderID: referenceID,
				PaymentID:        paymentID,
				PaymentStatus:    string(payment.Status),
				PaymentAmount:    paymentAmount,
				PaymentCountry:   order.CountryCode,
				PaymentMethod:    payment.Method,
				PaymentDate:      payment.CreatedAt.Format(time.RFC3339),
			}

			receipt, genErr := n.useCase.Generate(ctx, cmd)

			if genErr != nil {
				return receipt, errors.WrapActivityError(genErr)
			}

			return receipt, nil
		})
}
