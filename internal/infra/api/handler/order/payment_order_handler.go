package order

import (
	"net/http"

	"github.com/labstack/echo/v4"
	dto "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/request"
)

func (p *OrderHandler) CreatePaymentOrder(context echo.Context) error {
	var requestBody dto.CreatePaymentOrderRequestDTO

	ctx := context.Request().Context()

	if err := context.Bind(&requestBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := requestBody.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	cmd, err := requestBody.Command()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	response, err := p.usecaesPaymentOrder.CreatePaymentOrder(ctx, cmd)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return context.JSON(http.StatusCreated, response)
}
