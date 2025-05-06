package group

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/webhooks"
)

const (
	webhooksPath         = "/v1/webhooks"
	deunaNotifyOrderPath = "/deuna/notify-order"
)

type WebhooksRoutesIF interface {
	Resource(c *echo.Group)
}

type WebhooksRoutes struct{}

func webhooksV1Routes(group *echo.Group, handler webhooks.WebhooksHandlerIF) {
	webhookRoute := group.Group(webhooksPath)
	webhookRoute.POST(deunaNotifyOrderPath, handler.DeunaWebhookNotifyOrder)
}

func NewWebhooksRoutes(group *echo.Group, handler webhooks.WebhooksHandlerIF) *WebhooksRoutes {
	webhooksV1Routes(group, handler)
	return &WebhooksRoutes{}
}
