package group

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_account_route"
)

const (
	collectionAccountRoutePathV1 = "/v1/collection-account/route"
)

type CollectionAccountRouteRoutesIF interface {
	Resource(c *echo.Group)
}

type CollectionAccountRouteRoutes struct{}

func CollectionAccountRouteV1Routes(
	group *echo.Group,
	handler collection_account_route.CollectionAccountRouteHandlerIF,
) {
	CollectionAccountRoute := group.Group(collectionAccountRoutePathV1, auth.AuthParamsRetriever(
		auth.WithIgnoringHeaders(map[string]bool{
			auth.HeaderUsername: true,
		},
		)))
	CollectionAccountRoute.POST("", handler.Create)
	CollectionAccountRoute.PUT("/:id/disable", handler.Disable)
}

func NewCollectionAccountRouteRoutes(
	group *echo.Group, handler collection_account_route.CollectionAccountRouteHandlerIF,
) *CollectionAccountRouteRoutes {
	CollectionAccountRouteV1Routes(group, handler)
	return &CollectionAccountRouteRoutes{}
}
