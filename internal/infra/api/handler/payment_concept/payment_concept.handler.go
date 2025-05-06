package payment_concept

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	use_cases "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_concept/use_cases"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/payment_concept/dto/request"
)

type PaymentConceptHandlerIF interface {
	Create(context echo.Context) error
}

type PaymentConceptHandler struct {
	paymentConceptUsecase use_cases.PaymentConceptUsecaseIF
}

func NewPaymentConceptHandler(paymentConceptUsecase use_cases.PaymentConceptUsecaseIF) PaymentConceptHandlerIF {
	return &PaymentConceptHandler{
		paymentConceptUsecase: paymentConceptUsecase,
	}
}

func (p *PaymentConceptHandler) Create(context echo.Context) error {
	var request request.PaymentConceptRequest

	authParams := auth.GetParamsFromEchoContext(context)
	ctx := context.Request().Context()
	enterpriseID := authParams.EnterpriseID

	if err := context.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if validationErr := request.Validate(); validationErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, validationErr.Error())
	}

	response, err := p.paymentConceptUsecase.Create(ctx, request, enterpriseID)

	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}

	return context.JSON(http.StatusCreated, response)
}
