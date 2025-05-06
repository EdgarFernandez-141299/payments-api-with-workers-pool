package group

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/capture_flow"
)

type CaptureFlowRoutesIF interface {
	Resource(c *echo.Group)
}

type CaptureFlowRoutes struct {
	handler capture_flow.CaptureFlowHandlerIF
}

func CaptureFlowV1Routes(group *echo.Group, handler capture_flow.CaptureFlowHandlerIF) {
	captureRoute := group.Group("/v1/order/payment", auth.AuthParamsRetriever(
		auth.WithIgnoringHeaders(map[string]bool{
			auth.HeaderUsername: true,
		},
		)))
	captureRoute.POST("/capture", handler.PaymentCapture)
	captureRoute.POST("/release", handler.PaymentRelease)
}

func NewCaptureFlowRoutes(group *echo.Group, handler capture_flow.CaptureFlowHandlerIF) *CaptureFlowRoutes {
	CaptureFlowV1Routes(group, handler)
	return &CaptureFlowRoutes{handler: handler}
}
