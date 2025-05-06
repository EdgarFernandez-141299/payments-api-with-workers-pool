package card

import (
	"context"
	"net/http"
	"strings"

	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"

	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/usecases"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/usecases/queries"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/card/dto/request"
)

type CardHandlerIF interface {
	CreateCard(context echo.Context) error
	GetCardsByUserID(context echo.Context) error
	TriggerExpiringSoonNotifications(context echo.Context) error
}

type CardHandler struct {
	usecase      usecases.CardUsecaseIF
	getCredicard queries.GetCardUsecaseIF
}

func NewCardHandler(usecase usecases.CardUsecaseIF, getCredicard queries.GetCardUsecaseIF) CardHandlerIF {
	return &CardHandler{
		usecase:      usecase,
		getCredicard: getCredicard,
	}
}

func (p *CardHandler) CreateCard(context echo.Context) error {
	var body request.CardRequest

	ctx := context.Request().Context()

	authParams := auth.GetParamsFromEchoContext(context)
	enterpriseID := authParams.EnterpriseID
	preferredLanguage := context.Request().Header.Get("x-user-language")

	if err := context.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if validationErr := body.Validate(); validationErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, validationErr.Error())
	}

	response, err := p.usecase.CreateCard(ctx, body, enterpriseID, preferredLanguage)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusCreated, response)
}

func (p *CardHandler) GetCardsByUserID(reqContext echo.Context) error {
	return decorators.TraceDecoratorNoReturn(reqContext.Request().Context(), "CardHandler.GetCardsByUserId",
		func(ctx context.Context, span decorators.Span) error {
			userID := reqContext.Param("user_id")

			userID = strings.TrimSpace(userID)

			authParams := auth.GetParamsFromEchoContext(reqContext)
			enterpriseID := authParams.EnterpriseID

			response, err := p.getCredicard.GetCardsByUserID(ctx, userID, enterpriseID)

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			return reqContext.JSON(http.StatusOK, response)
		},
	)
}

func (p *CardHandler) TriggerExpiringSoonNotifications(context echo.Context) error {
	var body request.NotificationCardExpiringSoonRequestDTO

	ctx := context.Request().Context()

	if err := context.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if validationErr := body.Validate(); validationErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, validationErr.Error())
	}

	response, err := p.usecase.TriggerCardExpiringSoonNotifications(ctx, body)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, response)
}
