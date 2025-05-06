package group

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/payment_method"
)

const (
	paymentMethodPathV1 = "/v1/payment-method"
)

type PaymentMethodRoutesIF interface {
	Resource(c *echo.Group)
}

type PaymentMethodRoutes struct{}

func paymentMethodV1Routes(group *echo.Group, handler payment_method.PaymentMethodHandlerIF) {
	PaymentMethodRoute := group.Group(paymentMethodPathV1, auth.AuthParamsRetriever(
		auth.WithIgnoringHeaders(map[string]bool{
			auth.HeaderUsername: true,
		},
		)))
	PaymentMethodRoute.POST("", handler.Create)
}

func NewPaymentMethodRoutes(group *echo.Group, handler payment_method.PaymentMethodHandlerIF) *PaymentMethodRoutes {
	paymentMethodV1Routes(group, handler)
	return &PaymentMethodRoutes{}
}
