package group

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order"
)

const (
	OrderPathV1 = "/v1/order"
)

type OrderRoutesIF interface {
	Resource(c *echo.Group)
}
type OrderRoutes struct{}

func orderV1Routes(
	group *echo.Group, handler order.OrderHandlerIF,
	payOrder order.OrderPaymentsHandlerIF,
	refund order.PaymentRefundHandlerIF,
) {
	routes := group.Group(OrderPathV1, auth.AuthParamsRetriever(
		auth.WithIgnoringHeaders(map[string]bool{
			auth.HeaderUsername: true,
		},
		)))
	routes.POST("", handler.Create)
	routes.POST("/payment", handler.CreatePaymentOrder)
	routes.POST("/payments", payOrder.Pay)
	routes.POST("/payment/refund", refund.Refund)
	routes.GET("/:order_reference_id", handler.GetOrder)
	routes.GET("/status/:order_reference_id", handler.GetOrderPayments)

}

func NewOrderRoutes(
	group *echo.Group, handler order.OrderHandlerIF,
	orderPaymentHandler order.OrderPaymentsHandlerIF,
	orderRefundHandler order.PaymentRefundHandlerIF,
) *OrderRoutes {
	orderV1Routes(group, handler, orderPaymentHandler, orderRefundHandler)
	return &OrderRoutes{}
}
