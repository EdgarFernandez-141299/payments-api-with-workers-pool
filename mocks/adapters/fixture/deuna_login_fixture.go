package fixture

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/gommon/uid"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/adapters"
	"testing"
)

func DeunaLoginFixture(
	t *testing.T, userID string, enterpriseID string, err error,
) (*adapters.DeunaLoginAdapter, string) {
	loginAdapterMock := adapters.NewDeunaLoginAdapter(t)
	token, _ := uid.NewUniqueID()

	loginAdapterMock.On("LoginWitUserID", mock.Anything, userID, enterpriseID).
		Return(token.String(), err).Once()

	return loginAdapterMock, token.String()
}
