package group

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	card "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/card"
)

const (
	deleteCardRoutePathV1 = "/v1/card"
)

type DeleteCardRoutesIF interface {
	Resource(c *echo.Group)
}

type DeleteCardRoutes struct{}

func deletecardV1Routes(group *echo.Group, handler card.DeleteCardHandlerIF) {
	deletecardRoute := group.Group(deleteCardRoutePathV1, auth.AuthParamsRetriever(
		auth.WithIgnoringHeaders(map[string]bool{
			auth.HeaderUsername: true,
		},
		)))
	deletecardRoute.DELETE("", handler.DeleteCard)
}

func NewDeleteCardRoutes(
	group *echo.Group,
	handler card.DeleteCardHandlerIF,
) *DeleteCardRoutes {
	deletecardV1Routes(group, handler)
	return &DeleteCardRoutes{}
}
