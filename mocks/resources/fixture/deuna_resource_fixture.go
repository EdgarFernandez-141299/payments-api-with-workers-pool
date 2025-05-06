package fixture

import (
	"github.com/stretchr/testify/mock"
	request2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/resources/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/resources/dto/response"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/resources/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/resources"
	"testing"
)

func MakeOrderPaymentFixture(
	t *testing.T,
	userToken string,
	orderToken string,
	err error,
) *resources.DeunaPaymentResourceIF {
	m := resources.NewDeunaPaymentResourceIF(t)

	m.On("MakeOrderPayment", mock.Anything, mock.MatchedBy(func(req request.DeunaOrderPaymentRequest) bool {
		return req.OrderToken == orderToken
	}), userToken).Return(err).Once()

	return m
}

func CreateOrderResourceFixture(
	t *testing.T, order request2.CreateDeunaOrderRequestDTO, err error,
) *resources.DeunaOrderResourceIF {
	m := resources.NewDeunaOrderResourceIF(t)

	m.On("CreateOrder", mock.Anything, order).Return(response.DeunaOrderResponseDTO{}, err)

	return m
}
