package http

import (
	"context"
	"fmt"
	"gitlab.com/clubhub.ai1/go-libraries/observability/apm"
	"net/http"
	"time"

	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	"gitlab.com/clubhub.ai1/go-libraries/observability/clients/http/instrument"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/response"
	deunaConfig "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/http/config"
)

const (
	deUnaTimeout = 10 * time.Second
)

type DeUnaHTTPClient struct {
	instrument.Client
}

func NewDeUnaHTTPClient(config *deunaConfig.DeUnaApiConfig, tracer apm.Tracer) resources.DeUnaUserResourceIF {
	apiKeyHeader := deunaConfig.DeunaApiKeyHeader

	return &DeUnaHTTPClient{
		instrument.NewInstrumentedClient(
			instrument.WithBaseUrl(config.URL),
			instrument.WithTraceOptions(tracer.GetTracer(), instrument.TraceRequest, instrument.TraceResponse),
			instrument.WithRequestTimeout(deUnaTimeout),
			instrument.WithHeaders(map[string]string{
				apiKeyHeader:   config.ApiKey,
				"Content-Type": "application/json",
			}),
		),
	}
}

func (d DeUnaHTTPClient) CreateUser(
	oldCtx context.Context,
	request request.CreateUserRequestDTO,
) (response.CreatedUserResponse, error) {
	return decorators.TraceDecorator[response.CreatedUserResponse](
		oldCtx,
		"DeUnaHTTPClient.CreateUser",
		func(ctx context.Context, span decorators.Span) (response.CreatedUserResponse, error) {
			var responseBody response.CreatedUserResponse

			record, err := d.Client.NewRequestWithOptions(
				instrument.WithHeadersLogConfig(true, deunaConfig.DeunaApiKeyHeader),
			).
				SetResult(&responseBody).
				SetBody(request).
				SetContext(ctx).
				Post("/users/register")

			if err != nil {
				return response.CreatedUserResponse{}, err
			}

			if record.StatusCode() != http.StatusCreated {
				return response.CreatedUserResponse{}, fmt.Errorf("error creating user")
			}

			return responseBody, nil
		},
	)
}
