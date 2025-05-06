package card

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/usecases"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/card/dto/request"
)

type DeleteCardHandlerIF interface {
	DeleteCard(context echo.Context) error
}

type DeleteCardHandler struct {
	usecase usecases.DeleteCardUseCaseIF
}

func NewDeleteCardHandler(usecase usecases.DeleteCardUseCaseIF) DeleteCardHandlerIF {
	return &DeleteCardHandler{
		usecase: usecase,
	}
}

func (p *DeleteCardHandler) DeleteCard(context echo.Context) error {
	var body request.DeleteCardRequest

	ctx := context.Request().Context()

	authParams := auth.GetParamsFromEchoContext(context)
	enterpriseID := authParams.EnterpriseID
	userLanguage := context.Request().Header.Get("x-user-language")

	if err := context.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if validationErr := body.Validate(); validationErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, validationErr.Error())
	}

	response, err := p.usecase.DeleteCard(ctx, body, enterpriseID, userLanguage)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, response)
}
