package group

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_center"
)

type CollectionCenterRoutesIF interface {
	Resource(c *echo.Group)
}

type CollectionCenterRoutes struct{}

func CollectionCenterV1Routes(group *echo.Group, handler collection_center.CollectionCenterHandlerIF) {
	collectionCenterRoute := group.Group("/v1/collection-center", auth.AuthParamsRetriever(
		auth.WithIgnoringHeaders(map[string]bool{
			auth.HeaderUsername: true,
		},
		)))
	collectionCenterRoute.POST("", handler.Create)
}

func NewCollectionCenterRoutes(
	c *echo.Group,
	handler collection_center.CollectionCenterHandlerIF,
) *CollectionCenterRoutes {
	CollectionCenterV1Routes(c, handler)
	return &CollectionCenterRoutes{}
}
