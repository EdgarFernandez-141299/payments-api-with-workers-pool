package fixture

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository"
	"testing"
)

func CreateDeunaPaymentOrderRepositoryFixture(
	t *testing.T, orderID, paymentID string, deunaOrderTokenID string, err error,
) *repository.DeunaOrderRepository {
	m := repository.NewDeunaOrderRepository(t)

	m.On("CreatePaymentOrderDeuna", mock.Anything, paymentID, orderID, deunaOrderTokenID).Return(err)

	return m
}
