package fixture

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
	paymentOrderEntity "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_order/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository"
)

func CreatePaymentOrderReadModelRepository(
	t *testing.T, orderID string, cardID string, cmd command.CreatePaymentOrderCommand, err error,
) *repository.PaymentOrderRepositoryIF {
	m := repository.NewPaymentOrderRepositoryIF(t)

	m.On("CreatePaymentOrder", mock.Anything, mock.MatchedBy(func(entity paymentOrderEntity.PaymentOrderEntity) bool {
		return entity.OrderID == orderID &&
			entity.AssociatedOrigin == cmd.AssociatedOrigin.Type.String() &&
			entity.PaymentMethod == cmd.Payment.Method.Type.String() &&
			entity.CurrencyCode == cmd.CurrencyCode.Code &&
			entity.CountryCode == cmd.CountryCode &&
			entity.CardID == cardID &&
			entity.CollectionAccountID == cmd.CollectionAccount.ID &&
			entity.PaymentOrderID == cmd.Payment.ID &&
			entity.TotalAmount.Equal(cmd.Payment.Total.Value) &&
			entity.PaymentFlow == cmd.PaymentFlow.String() &&
			entity.Status == cmd.Payment.Status.String() &&
			!entity.TransactionDate.IsZero()
	})).Return(err).Once()

	return m
}
