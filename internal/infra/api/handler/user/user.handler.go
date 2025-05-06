package user

import (
	"context"
	"net/http"

	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	usecases "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/deuna"

	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	commons "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/common"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/user/dto/request"
)

type UserHandlerIF interface {
	Create(context echo.Context) error
	ValidateUser(echoCtx echo.Context) error
}

type UserHandler struct {
	usecase usecases.CreateUserUseDeunaAdapterIF
}

func NewUserHandler(usecase usecases.CreateUserUseDeunaAdapterIF) UserHandlerIF {
	return &UserHandler{
		usecase: usecase,
	}
}

func (u *UserHandler) Create(echoCtx echo.Context) error {
	return decorators.TraceDecoratorNoReturn(echoCtx.Request().Context(), "UserHandler.Create", func(ctx context.Context, span decorators.Span) error {
		var body request.CreateUserRequest

		authParams := auth.GetParamsFromEchoContext(echoCtx)
		enterpriseId := authParams.EnterpriseID

		newCtx := commons.SetEnterpriseID(ctx, enterpriseId)

		if err := echoCtx.Bind(&body); err != nil {
			span.RecordError(err)

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		response, err := u.usecase.Create(newCtx, body, enterpriseId)

		if err != nil {
			span.RecordError(err)

			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}

		return echoCtx.JSON(http.StatusCreated, response)
	})
}

func (u *UserHandler) ValidateUser(echoCtx echo.Context) error {
	return decorators.TraceDecoratorNoReturn(echoCtx.Request().Context(), "UserHandler.ValidateUser", func(ctx context.Context, span decorators.Span) error {
		authParams := auth.GetParamsFromEchoContext(echoCtx)
		enterpriseId := authParams.EnterpriseID

		request := request.CreateUserRequest{
			UserID:   authParams.UserID,
			UserType: authParams.UserType,
		}

		newCtx := commons.SetEnterpriseID(ctx, enterpriseId)

		response, err := u.usecase.ValidateUser(newCtx, request, enterpriseId)

		if err != nil {
			span.RecordError(err)

			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}

		return echoCtx.JSON(http.StatusCreated, response)
	})
}
