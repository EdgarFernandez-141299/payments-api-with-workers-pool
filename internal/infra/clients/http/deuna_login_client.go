package http

import (
	"context"
	"errors"
	"gitlab.com/clubhub.ai1/go-libraries/observability/apm"
	deunaConfig "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/http/config"
	"net/http"

	"gitlab.com/clubhub.ai1/go-libraries/observability/clients/http/instrument"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/config"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/response"
)

type DeUnaLoginHTTPClient struct {
	instrument.Client
}

func NewDeUnaLoginHTTPClient(tracer apm.Tracer) resources.DeunaLoginResourceIF {
	apiKeyHeader := deunaConfig.DeunaApiKeyHeader
	return &DeUnaLoginHTTPClient{
		instrument.NewInstrumentedClient(
			instrument.WithBaseUrl(config.Config().DeUnaApi.URL),
			instrument.WithTraceOptions(tracer.GetTracer(), instrument.TraceRequest, instrument.TraceResponse),
			instrument.WithRequestTimeout(deUnaTimeout),
			instrument.WithHeaders(map[string]string{
				apiKeyHeader:   config.Config().DeUnaApi.ApiKey,
				"Content-Type": "application/json",
			}),
		),
	}
}

func (d DeUnaLoginHTTPClient) Login(
	ctx context.Context,
	request request.LoginUserDeUnaRequestDTO,
) (response.LoginResponseDTO, error) {
	var responseBody response.LoginResponseDTO

	record, err := d.Client.NewRequestWithOptions(
		instrument.WithHeadersLogConfig(true, deunaConfig.DeunaApiKeyHeader),
	).
		SetResult(&responseBody).
		SetBody(request).
		SetContext(ctx).
		Post("/users/login?type=guest")

	if err != nil {
		return response.LoginResponseDTO{}, err
	}

	if record.StatusCode() != http.StatusOK {
		return response.LoginResponseDTO{}, errors.New("error login user")
	}

	return responseBody, nil
}
