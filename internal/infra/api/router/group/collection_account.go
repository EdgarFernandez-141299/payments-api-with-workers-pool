package group

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_account"
)

const (
	collectionAccountPathV1 = "/v1/collection-account"
)

type CollectionAccountRoutesIF interface {
	Resource(c *echo.Group)
}
type CollectionAccountRoutes struct{}

func collectionAccountV1Routes(group *echo.Group, handler collection_account.CollectionAccountHandlerIF) {
	CollectionAccountRoute := group.Group(collectionAccountPathV1, auth.AuthParamsRetriever(
		auth.WithIgnoringHeaders(map[string]bool{
			auth.HeaderUsername: true,
		},
		)))
	CollectionAccountRoute.POST("", handler.Create)
}

func NewCollectionAccountRoutes(
	group *echo.Group, handler collection_account.CollectionAccountHandlerIF,
) *CollectionAccountRoutes {
	collectionAccountV1Routes(group, handler)
	return &CollectionAccountRoutes{}
}
