package auth

import "gitlab.com/clubhub.ai1/gommon/router/middleware/auth"

type AuthParamsFixture struct {
	params auth.AuthParams
}

func NewAuthParamsFixture(userID string, userType string, enterpriseID string) *AuthParamsFixture {
	return &AuthParamsFixture{
		params: auth.AuthParams{
			UserID:            userID,
			UserName:          "TestUser",
			EnterpriseID:      enterpriseID,
			UserType:          userType,
			ContextAuthorized: true,
			ContextValidated:  true,
		},
	}
}
func (a *AuthParamsFixture) GetParams() auth.AuthParams {
	return a.params
}
