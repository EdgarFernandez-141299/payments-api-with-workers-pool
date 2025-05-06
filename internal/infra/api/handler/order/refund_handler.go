package order

import (
	"net/http"

	"github.com/labstack/echo/v4"
	errorx "gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/interfaces/web"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/refund/use_cases"
	dto "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/pkg/common"
)

const (
	schemaValidationError = "SCHEMA_VALIDATION_ERROR"
)

type PaymentRefundHandlerIF interface {
	Refund(context echo.Context) error
}

type PaymentRefundHandler struct {
	totalRefundUseCase   use_cases.RefundTotalUseCaseIF
	partialRefundUseCase use_cases.PartialRefundUseCaseIF
}

func NewPaymentRefundHandler(
	totalRefundUseCase use_cases.RefundTotalUseCaseIF,
	partialRefundUseCase use_cases.PartialRefundUseCaseIF,
) PaymentRefundHandlerIF {
	return &PaymentRefundHandler{
		totalRefundUseCase:   totalRefundUseCase,
		partialRefundUseCase: partialRefundUseCase,
	}
}

func (p *PaymentRefundHandler) Refund(context echo.Context) error {
	var body dto.RefundDTO
	ctx := context.Request().Context()

	authParams := auth.GetParamsFromEchoContext(context)
	enterpriseID := authParams.EnterpriseID

	if err := context.Bind(&body); err != nil {
		return context.JSON(http.StatusBadRequest, common.MessageError{
			Code:    schemaValidationError,
			Message: err.Error(),
		})
	}

	if err := body.Validate(); err != nil {
		return context.JSON(http.StatusBadRequest, common.MessageError{
			Code:    schemaValidationError,
			Message: err.Error(),
		})
	}

	if body.IsTotal {
		cmd, _ := body.Command()
		response, err := p.totalRefundUseCase.Refund(ctx, cmd, enterpriseID)

		if err != nil {
			return context.JSON(http.StatusBadRequest, errorx.HandleErrorEcho(err))
		}

		return context.JSON(http.StatusOK, response)
	} else {
		cmd, err := body.CommandPartial()
		if err != nil {
			return context.JSON(http.StatusBadRequest, common.MessageError{Code: schemaValidationError, Message: err.Error()})
		}

		response, err := p.partialRefundUseCase.PartialRefund(ctx, cmd, enterpriseID)
		if err != nil {
			return errorx.HandleErrorEcho(err)
		}

		return context.JSON(http.StatusOK, response)
	}
}
