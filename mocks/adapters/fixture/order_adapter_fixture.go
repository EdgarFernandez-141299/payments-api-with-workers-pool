package fixture

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/gommon/uid"
	orderResponse "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/resources/dto/response"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/adapters"
)

func NewCreateOrderFixture(
	t *testing.T, orderID string, paymentID string, currency string, totalValue decimal.Decimal, paymentFlow enums.PaymentFlowEnum, err error,
) (*adapters.OrderAdapterIF, orderResponse.DeunaOrderResponseDTO) {
	orderAdapterMock := adapters.NewOrderAdapterIF(t)

	orderToken, _ := uid.NewUniqueID()

	response := orderResponse.DeunaOrderResponseDTO{
		Token: orderToken.String(),
	}

	orderAdapterMock.On(
		"CreateOrder",
		mock.Anything,
		orderID,
		paymentID,
		currency,
		totalValue,
		paymentFlow).Return(orderResponse.DeunaOrderResponseDTO{
		Token: orderToken.String(),
	}, err).Once()

	return orderAdapterMock, response
}
