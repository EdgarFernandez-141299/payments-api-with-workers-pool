package capture_flow

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/worfkflows"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/utils"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/capture_flow/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/capture_flow/dto/response"
)

type CaptureFlowHandlerIF interface {
	PaymentCapture(context echo.Context) error
	PaymentRelease(context echo.Context) error
}

type CaptureFlowHandler struct {
	workflow worfkflows.PaymentOrderWorkflow
}

func NewCaptureFlowHandler(workflow worfkflows.PaymentOrderWorkflow) CaptureFlowHandlerIF {
	return &CaptureFlowHandler{
		workflow: workflow,
	}
}

func (h *CaptureFlowHandler) PaymentCapture(context echo.Context) error {
	var body request.CaptureRequest

	authParams := auth.GetParamsFromEchoContext(context)

	fmt.Println("authParams", authParams)

	if err := context.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if validationErr := body.Validate(); validationErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, validationErr.Error())
	}

	err := h.sendFlowSignal(context, body.ReferenceOrderID, body.PaymentID, enums.CapturePayment, "")

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusCreated, response.NewCaptureResponse(body.ReferenceOrderID, body.PaymentID, "CAPTURED"))
}

func (h *CaptureFlowHandler) PaymentRelease(context echo.Context) error {
	var body request.ReleaseRequest

	if err := context.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if validationErr := body.Validate(); validationErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, validationErr.Error())
	}

	err := h.sendFlowSignal(context, body.ReferenceOrderID, body.PaymentID, enums.ReleasePayment, body.Reason)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusCreated, response.NewReleaseResponse(body.ReferenceOrderID, body.PaymentID, "RELEASED"))
}

func (h *CaptureFlowHandler) sendFlowSignal(context echo.Context,
	ReferenceOrderID string,
	PaymentID string,
	action enums.PaymentFlowActionEnum,
	reason string) error {
	paymentOrderID := utils.NewDeunaOrderID(ReferenceOrderID, PaymentID).GetID()

	err := h.workflow.SendCaptureFlowSignal(context.Request().Context(),
		paymentOrderID, worfkflows.CompleteCaptureFlowSignal{
			OrderReferecenId: ReferenceOrderID,
			PaymentID:        PaymentID,
			Action:           action,
			Reason:           reason,
		})

	return err
}
