package group

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	card "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/card"
)

const (
	cardRoutePathV1 = "/v1/card"
)

type CardRoutesIF interface {
	Resource(c *echo.Group)
}

type CardRoutes struct{}

func cardV1Routes(group *echo.Group, handler card.CardHandlerIF) {
	cardRoute := group.Group(cardRoutePathV1, auth.AuthParamsRetriever(
		auth.WithIgnoringHeaders(map[string]bool{
			auth.HeaderUsername: true,
		},
		)))
	cardRoute.POST("", handler.CreateCard)
	cardRoute.GET("/by-user/:user_id", handler.GetCardsByUserID)
	cardRoute.POST("/notifications/expiring-soon", handler.TriggerExpiringSoonNotifications)
}

func NewCardRoutes(group *echo.Group, handler card.CardHandlerIF) *CardRoutes {
	cardV1Routes(group, handler)
	return &CardRoutes{}
}
