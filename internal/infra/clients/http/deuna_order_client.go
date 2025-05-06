package http

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	"gitlab.com/clubhub.ai1/go-libraries/observability/clients/http/instrument"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/resources/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/resources/dto/response"
	deunaConfig "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/http/config"
)

const authorizationHeaderFormat = "Bearer %s"

type DeUnaOrderHTTPClient struct {
	instrument.Client
}

func NewDeUnaOrderHTTPClient(config *deunaConfig.DeUnaApiConfig) resources.DeunaOrderResourceIF {
	apiKeyheader := deunaConfig.DeunaApiKeyHeader
	return &DeUnaOrderHTTPClient{
		instrument.NewInstrumentedClient(
			instrument.WithBaseUrl(config.URL),
			instrument.WithRequestTimeout(deUnaTimeout),
			instrument.WithHeaders(map[string]string{
				apiKeyheader:   config.ApiKey,
				"Content-Type": "application/json",
			}),
		),
	}
}

func (d DeUnaOrderHTTPClient) CreateOrder(
	oldCtx context.Context,
	body request.CreateDeunaOrderRequestDTO,
) (response.DeunaOrderResponseDTO, error) {
	return decorators.TraceDecorator[response.DeunaOrderResponseDTO](
		oldCtx,
		"DeUnaOrderHTTPClient.CreateOrder",
		func(ctx context.Context, span decorators.Span) (response.DeunaOrderResponseDTO, error) {
			var res response.DeunaOrderResponseDTO

			httpResponse, err := d.Client.NewRequestWithOptions(
				instrument.WithHeadersLogConfig(true, deunaConfig.DeunaApiKeyHeader),
			).
				SetResult(&res).
				SetBody(body).
				SetContext(ctx).
				Post("/merchants/orders")

			if err != nil {
				return response.DeunaOrderResponseDTO{}, fmt.Errorf("failed to create order request: %w", err)
			}

			if httpResponse.StatusCode() == http.StatusConflict {
				return response.DeunaOrderResponseDTO{}, fmt.Errorf("invalid payment method with status: %d", httpResponse.StatusCode())
			}

			if httpResponse.IsError() || httpResponse.StatusCode() != http.StatusOK {
				return response.DeunaOrderResponseDTO{}, fmt.Errorf("create order request failed with status: %d", httpResponse.StatusCode())
			}

			return res, nil
		},
	)
}

func (d DeUnaOrderHTTPClient) GetOrder(
	oldCtx context.Context,
	orderToken string,
) (response.DeunaOrderResponseDTO, error) {
	return decorators.TraceDecorator[response.DeunaOrderResponseDTO](
		oldCtx,
		"DeUnaOrderHTTPClient.GetOrder",
		func(ctx context.Context, span decorators.Span) (response.DeunaOrderResponseDTO, error) {
			var res response.DeunaOrderResponseDTO

			httpResponse, err := d.Client.NewRequestWithOptions(
				instrument.WithHeadersLogConfig(true, deunaConfig.DeunaApiKeyHeader),
			).
				SetResult(&res).
				SetContext(ctx).
				Get(fmt.Sprintf("/merchants/orders/%s", orderToken))

			if err != nil {
				return response.DeunaOrderResponseDTO{}, fmt.Errorf("failed to get order: %w", err)
			}

			if httpResponse.IsError() || httpResponse.StatusCode() != http.StatusOK {
				return response.DeunaOrderResponseDTO{}, fmt.Errorf("get order failed with status: %d", httpResponse.StatusCode())
			}

			return res, nil
		},
	)
}

func (d DeUnaOrderHTTPClient) ExpireOrder(
	oldCtx context.Context,
	orderToken string,
) error {

	return decorators.TraceDecoratorNoReturn(
		oldCtx,
		"DeUnaOrderHTTPClient.ExpireOrder",
		func(ctx context.Context, span decorators.Span) error {
			httpResponse, err := d.Client.NewRequestWithOptions(
				instrument.WithHeadersLogConfig(true, deunaConfig.DeunaApiKeyHeader),
			).
				SetContext(ctx).
				Delete(fmt.Sprintf("/merchants/orders/%s/expire", orderToken))

			if err != nil {
				return fmt.Errorf("failed to expire order request: %w", err)
			}

			if httpResponse.IsError() || httpResponse.StatusCode() != http.StatusOK {
				return fmt.Errorf("expire order request failed with status: %d", httpResponse.StatusCode())
			}

			return nil
		},
	)
}
