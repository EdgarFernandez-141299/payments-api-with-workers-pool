package order

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/worfkflows"
	dto "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/request"
	"net/http"
)

type OrderPaymentsHandlerIF interface {
	Pay(context echo.Context) error
}

type OrderPaymentsHandler struct {
	workflow worfkflows.PaymentOrderWorkflow
}

func NewOrderPaymentsHandler(
	workflow worfkflows.PaymentOrderWorkflow,
) OrderPaymentsHandlerIF {
	return &OrderPaymentsHandler{
		workflow: workflow,
	}
}

func (p *OrderPaymentsHandler) Pay(context echo.Context) error {
	var request dto.OrderWithPaymentsDTO
	ctx := context.Request().Context()

	if err := context.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	cmd, err := request.ToWorkflowInput()

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	wfFuture, err := p.workflow.Call(ctx, cmd.OrderID, cmd)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var response *worfkflows.PaymentOrderWorkflowOut

	getErr := wfFuture.Get(ctx, &response)

	if getErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, getErr.Error())
	}

	return context.JSON(http.StatusCreated, response)
}
