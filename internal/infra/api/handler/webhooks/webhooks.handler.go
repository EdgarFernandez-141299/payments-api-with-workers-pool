package webhooks

import (
	"fmt"
	"net/http"
	"strings"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/utils"

	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/worfkflows"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	dto "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/webhooks/dto/request"
)

type WebhooksHandlerIF interface {
	DeunaWebhookNotifyOrder(context echo.Context) error
}

type WebhookHandler struct {
	workflow worfkflows.PaymentOrderWorkflow
}

func NewWebhookHandler(workflow worfkflows.PaymentOrderWorkflow) WebhooksHandlerIF {
	return &WebhookHandler{
		workflow: workflow,
	}
}

func (p *WebhookHandler) DeunaWebhookNotifyOrder(context echo.Context) error {
	var orderWebhookRequest dto.WebhookOrderDTO

	if err := context.Bind(&orderWebhookRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	paymentStatus, err := enums.NewPaymentStatusFromString(strings.ToUpper(orderWebhookRequest.Order.Payment.Data.Status))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	deunaOrderID, err := utils.ExtractFromDeunaOrderID(orderWebhookRequest.Order.OrderID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	fmt.Println("orderWebhookRequest: ", orderWebhookRequest)

	signalRequest := worfkflows.PaymentProcessedSignal{
		AuthorizationCode:   orderWebhookRequest.Order.Payment.Data.AuthorizationCode,
		Status:              paymentStatus,
		OrderStatusString:   orderWebhookRequest.Order.Status,
		OrderID:             deunaOrderID.GetOrderID(),
		PaymentID:           deunaOrderID.GetPaymentID(),
		PaymentStatusString: orderWebhookRequest.Order.Payment.Data.Status,
		PaymentReason:       orderWebhookRequest.Order.Payment.Data.Reason,
		PaymentCard: worfkflows.CardData{
			CardBrand: orderWebhookRequest.Order.Payment.Data.FromCard.CardBrand,
			CardLast4: orderWebhookRequest.Order.Payment.Data.FromCard.LastFour,
			CardType:  orderWebhookRequest.Order.PaymentMethod,
		},
	}

	return p.workflow.SendProcessedSignal(context.Request().Context(), orderWebhookRequest.Order.OrderID, signalRequest)
}
