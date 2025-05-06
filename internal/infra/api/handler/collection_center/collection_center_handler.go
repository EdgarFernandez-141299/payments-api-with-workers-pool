package collection_center

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	use_cases "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/collection_center/use_cases"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_center/dto/request"
)

type CollectionCenterHandlerIF interface {
	Create(context echo.Context) error
}

type CollectionCenterHandler struct {
	collectionCenterUsecase use_cases.CollectionCenterUsecaseIF
}

func NewCollectionCenterHandler(collectionCenterUsecase use_cases.CollectionCenterUsecaseIF) CollectionCenterHandlerIF {
	return &CollectionCenterHandler{
		collectionCenterUsecase: collectionCenterUsecase,
	}
}

func (c *CollectionCenterHandler) Create(context echo.Context) error {
	var request request.CollectionCenterRequest

	authParams := auth.GetParamsFromEchoContext(context)
	ctx := context.Request().Context()
	enterpriseID := authParams.EnterpriseID

	if err := context.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if validationErr := request.Validate(); validationErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, validationErr.Error())
	}

	response, err := c.collectionCenterUsecase.Create(ctx, request, enterpriseID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusCreated, response)
}
