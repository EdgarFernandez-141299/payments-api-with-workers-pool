package collection_account

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/gommon/router/middleware/auth"
	use_cases "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/collection_account/use_cases"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_account/dto/request"
)

type CollectionAccountHandlerIF interface {
	Create(context echo.Context) error
}

type CollectionAccountHandler struct {
	collectionAccountUsecase use_cases.CollectionAccountUsecaseIF
}

func NewCollectionAccountHandler(
	collectionAccountUsecase use_cases.CollectionAccountUsecaseIF,
) CollectionAccountHandlerIF {
	return &CollectionAccountHandler{
		collectionAccountUsecase: collectionAccountUsecase,
	}
}

func (p *CollectionAccountHandler) Create(context echo.Context) error {
	var request request.CollectionAccountRequest

	authParams := auth.GetParamsFromEchoContext(context)
	ctx := context.Request().Context()
	enterpriseID := authParams.EnterpriseID

	if err := context.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if validationErr := request.Validate(); validationErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, validationErr.Error())
	}

	response, err := p.collectionAccountUsecase.Create(ctx, request, enterpriseID)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusCreated, response)
}
