package order

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	usecases "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/use_cases/create"
	queries "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/use_cases/queries"
	usecasesPaymentOrder "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/use_cases/create"
	dto "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/request"
)

type OrderHandlerIF interface {
	Create(context echo.Context) error
	CreatePaymentOrder(context echo.Context) error
	GetOrder(context echo.Context) error
	GetOrderPayments(context echo.Context) error
}

type OrderHandler struct {
	usecase             usecases.CreateOrderUseCaseIF
	usecaesPaymentOrder usecasesPaymentOrder.CreatePaymentOrderUseCaseIF
	queries             queries.QueriesOrderUseCaseIF
}

func NewOrderHandler(
	usecase usecases.CreateOrderUseCaseIF,
	usecaesPaymentOrder usecasesPaymentOrder.CreatePaymentOrderUseCaseIF,
	queries queries.QueriesOrderUseCaseIF,
) OrderHandlerIF {
	return &OrderHandler{
		usecase:             usecase,
		usecaesPaymentOrder: usecaesPaymentOrder,
		queries:             queries,
	}
}

func (p *OrderHandler) Create(context echo.Context) error {
	var request dto.CreateOrderRequestDTO
	ctx := context.Request().Context()

	authParams := auth.GetParamsFromEchoContext(context)
	enterpriseID := authParams.EnterpriseID

	if err := context.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := request.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	cmd, err := request.ToCommand(enterpriseID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	response, err := p.usecase.Create(ctx, cmd)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return context.JSON(http.StatusCreated, response)
}

func (p *OrderHandler) GetOrder(context echo.Context) error {
	ctx := context.Request().Context()

	orderID := context.Param("order_reference_id")

	authParams := auth.GetParamsFromEchoContext(context)
	enterpriseID := authParams.EnterpriseID

	response, err := p.queries.GetOrderDetail(ctx, orderID, enterpriseID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return context.JSON(http.StatusOK, response)
}

func (p *OrderHandler) GetOrderPayments(context echo.Context) error {
	ctx := context.Request().Context()

	authParams := auth.GetParamsFromEchoContext(context)
	enterpriseID := authParams.EnterpriseID

	orderID := context.Param("order_reference_id")

	response, err := p.queries.GetOrderPayments(ctx, orderID, enterpriseID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return context.JSON(http.StatusOK, response)
}
