package http

import (
	"context"
	"fmt"
	"gitlab.com/clubhub.ai1/go-libraries/observability/apm"
	"net/http"

	deunaConfig "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/http/config"

	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	"gitlab.com/clubhub.ai1/go-libraries/observability/clients/http/instrument"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/adapters/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/adapters/resources/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/adapters/resources/dto/response"
)

var DeunaEmptyBody = fmt.Errorf("empty body")

type DeUnaCardHTTPClient struct {
	instrument.Client
}

func NewDeUnaCardHTTPClient(config *deunaConfig.DeUnaApiConfig, tracer apm.Tracer) resources.DeunaCardResourceIF {
	apiKeyHeader := deunaConfig.DeunaApiKeyHeader
	return &DeUnaCardHTTPClient{
		instrument.NewInstrumentedClient(
			instrument.WithTraceOptions(tracer.GetTracer(), instrument.TraceRequest, instrument.TraceResponse),
			instrument.WithBaseUrl(config.URL),
			instrument.WithRequestTimeout(deUnaTimeout),
			instrument.WithHeaders(map[string]string{
				apiKeyHeader:   config.ApiKey,
				"Content-Type": "application/json",
			}),
		),
	}
}

func (d DeUnaCardHTTPClient) CreateCard(
	oldCtx context.Context,
	body request.CreateCardRequestDTO,
	userId string,
	token string,
) (response.CardResponseDataDTO, error) {
	return decorators.TraceDecorator[response.CardResponseDataDTO](
		oldCtx,
		"DeUnaCardHTTPClient.CreateCard",
		func(ctx context.Context, span decorators.Span) (response.CardResponseDataDTO, error) {
			httpResponse, err := d.Client.NewRequestWithOptions(
				instrument.WithHeadersLogConfig(true, deunaConfig.DeunaApiKeyHeader),
			).
				SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
				SetResult(new(response.CardResponseDTO)).
				SetBody(body).
				SetContext(ctx).
				Post(fmt.Sprintf("/users/%s/cards", userId))

			if err != nil {
				return response.CardResponseDataDTO{}, fmt.Errorf("failed to create card request: %w", err)
			}

			if httpResponse.StatusCode() == http.StatusConflict {
				return response.CardResponseDataDTO{}, fmt.Errorf("the card already exist with status: %d", httpResponse.StatusCode())
			}

			if httpResponse.StatusCode() == http.StatusUnprocessableEntity {
				return response.CardResponseDataDTO{}, fmt.Errorf("operation denied by anti fraud rules: %d", httpResponse.StatusCode())
			}

			if httpResponse.IsError() || httpResponse.StatusCode() != http.StatusCreated {
				return response.CardResponseDataDTO{}, fmt.Errorf("create card request failed with status: %d", httpResponse.StatusCode())
			}

			responseDTO, ok := httpResponse.Result().(*response.CardResponseDTO)
			if !ok || responseDTO == nil {
				return response.CardResponseDataDTO{}, fmt.Errorf("failed to cast or parse response")
			}

			if responseDTO.IsEmpty() {
				return response.CardResponseDataDTO{}, DeunaEmptyBody
			}

			return responseDTO.Data, nil
		},
	)
}

func (d DeUnaCardHTTPClient) DeleteCard(
	oldCtx context.Context,
	body request.DeleteCardRequestDTO,
	token string,
) (response.DeleteCardResponseDTO, error) {
	return decorators.TraceDecorator[response.DeleteCardResponseDTO](
		oldCtx,
		"DeUnaCardHTTPClient.DeleteCard",
		func(ctx context.Context, span decorators.Span) (response.DeleteCardResponseDTO, error) {
			httpResponse, err := d.Client.NewRequestWithOptions(
				instrument.WithHeadersLogConfig(true, deunaConfig.DeunaApiKeyHeader),
			).
				SetHeaders(map[string]string{
					"Authorization": fmt.Sprintf("Bearer %s", token),
				}).
				SetResult(new(response.DeleteCardResponseDTO)).
				SetContext(ctx).
				Delete(fmt.Sprintf("/users/%s/cards/%s", body.UserId, body.CardId))

			if err != nil {
				return response.DeleteCardResponseDTO{}, err
			}

			if httpResponse.IsError() {
				return response.DeleteCardResponseDTO{}, fmt.Errorf("error deleting card DEUNA")
			}

			result, ok := httpResponse.Result().(*response.DeleteCardResponseDTO)

			if !ok || result == nil {
				return response.DeleteCardResponseDTO{}, fmt.Errorf("error parsing response")
			}

			if httpResponse.StatusCode() != http.StatusNoContent {
				return response.DeleteCardResponseDTO{}, fmt.Errorf("Error deleting card DEUNA")
			}

			return *result, nil
		},
	)
}
