package collection_account_route

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	usecases "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/collection_account_route/use_cases"
	request "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_account_route/dto/request"
)

type CollectionAccountRouteHandlerIF interface {
	Create(context echo.Context) error
	Disable(context echo.Context) error
}

type CollectionAccountRouteRoutes struct {
	usecase usecases.CollectionAccountRouteUsecaseIF
}

func NewCollectionAccountRouteHandler(
	usecase usecases.CollectionAccountRouteUsecaseIF,
) CollectionAccountRouteHandlerIF {
	return &CollectionAccountRouteRoutes{
		usecase: usecase,
	}
}

func (p *CollectionAccountRouteRoutes) Create(context echo.Context) error {
	var request request.CollectionAccountRouteRequest

	auth := auth.GetParamsFromEchoContext(context)
	ctx := context.Request().Context()
	enterpriseId := auth.EnterpriseID

	if err := context.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if validationErr := request.Validate(); validationErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, validationErr.Error())
	}

	response, err := p.usecase.Create(ctx, request, enterpriseId)

	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}

	return context.JSON(http.StatusCreated, response)
}

func (p *CollectionAccountRouteRoutes) Disable(context echo.Context) error {
	id := context.Param("id")

	auth := auth.GetParamsFromEchoContext(context)
	ctx := context.Request().Context()
	enterpriseId := auth.EnterpriseID

	response, err := p.usecase.Disable(ctx, id, enterpriseId)

	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}

	return context.JSON(http.StatusOK, response)
}
