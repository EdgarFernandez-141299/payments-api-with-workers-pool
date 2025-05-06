package group

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/payment_concept"
)

const (
	paymentConceptPathV1 = "/v1/payment-concept"
)

type PaymentConceptRoutesIF interface {
	Resource(c *echo.Group)
}
type PaymentConceptRoutes struct{}

func paymentConceptV1Routes(group *echo.Group, handler payment_concept.PaymentConceptHandlerIF) {
	PaymentConceptRoute := group.Group(paymentConceptPathV1, auth.AuthParamsRetriever(
		auth.WithIgnoringHeaders(map[string]bool{
			auth.HeaderUsername: true,
		},
		)))
	PaymentConceptRoute.POST("", handler.Create)
}

func NewPaymentConceptRoutes(group *echo.Group, handler payment_concept.PaymentConceptHandlerIF) *PaymentConceptRoutes {
	paymentConceptV1Routes(group, handler)
	return &PaymentConceptRoutes{}
}
