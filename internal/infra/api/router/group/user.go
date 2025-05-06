package group

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/user"
)

const (
	userPathV1 = "/v1/user"
)

type UserRoutesIF interface {
	Resource(c *echo.Group)
}

type UserRoutes struct{}

func userV1Routes(group *echo.Group, handler user.UserHandlerIF) {
	userRoute := group.Group(userPathV1, auth.AuthParamsRetriever(
		auth.WithIgnoringHeaders(map[string]bool{
			auth.HeaderUsername: true,
		},
		)))
	userRoute.POST("", handler.Create)
	userRoute.GET("/validate", handler.ValidateUser)
}

func NewUserRoutes(group *echo.Group, handler user.UserHandlerIF) *UserRoutes {
	userV1Routes(group, handler)
	return &UserRoutes{}
}
