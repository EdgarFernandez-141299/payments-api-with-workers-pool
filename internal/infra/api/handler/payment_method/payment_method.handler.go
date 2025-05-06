package payment_method

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	usecases "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_method/use_cases"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/payment_method/dto/request"
)

type PaymentMethodHandlerIF interface {
	Create(context echo.Context) error
}

type PaymentMethodHandler struct {
	paymentMethodUsecase usecases.PaymentMethodUseCasesIF
}

func NewPaymentMethodHandler(paymentMethodUsecase usecases.PaymentMethodUseCasesIF) PaymentMethodHandlerIF {
	return &PaymentMethodHandler{
		paymentMethodUsecase: paymentMethodUsecase,
	}
}

func (p *PaymentMethodHandler) Create(context echo.Context) error {
	var request request.PaymentMethodRequest

	authParams := auth.GetParamsFromEchoContext(context)
	ctx := context.Request().Context()
	enterpriseID := authParams.EnterpriseID

	if err := context.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if validationErr := request.Validate(); validationErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, validationErr.Error())
	}

	creation, err := p.paymentMethodUsecase.Create(ctx, request, enterpriseID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusCreated, creation)
}
