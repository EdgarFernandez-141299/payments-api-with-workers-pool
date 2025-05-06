package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	"gitlab.com/clubhub.ai1/go-libraries/observability/clients/http/instrument"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/resources/dto/request"
	deunaConfig "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/http/config"
)

type DeUnaCaptureFlowHTTPClient struct {
	instrument.Client
}

const (
	captureFlowError = "%s request failed with status: %d"
)

func NewDeUnaCaptureFlowHTTPClient(config *deunaConfig.DeUnaApiConfig) resources.DeunaCaptureFlowResourceIF {
	apiKeyheader := deunaConfig.DeunaApiKeyHeader
	return &DeUnaCaptureFlowHTTPClient{
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

// handleDeUnaResponse maneja la lógica común de respuesta HTTP para los métodos de DeUna
func (d DeUnaCaptureFlowHTTPClient) handleDeUnaResponse(
	httpResponse *resty.Response,
	operation string,
) (bool, error) {
	if httpResponse.IsError() {
		return false, fmt.Errorf(captureFlowError, operation, httpResponse.StatusCode())
	}

	if httpResponse.StatusCode() == http.StatusConflict {
		return false, fmt.Errorf(captureFlowError, operation, httpResponse.StatusCode())
	}

	if httpResponse.StatusCode() == http.StatusOK || httpResponse.StatusCode() == http.StatusNoContent {
		return true, nil
	}

	return false, nil
}

// Void
func (d DeUnaCaptureFlowHTTPClient) Release(
	oldCtx context.Context,
	orderToken string,
	reason string,
) (bool, error) {
	return decorators.TraceDecorator[bool](
		oldCtx,
		"DeUnaCaptureFlowHTTPClient.Release",
		func(ctx context.Context, span decorators.Span) (bool, error) {
			body := request.ReleaseRequestDTO{
				Reason: reason,
			}

			httpResponse, err := d.Client.NewRequestWithOptions(
				instrument.WithHeadersLogConfig(true, deunaConfig.DeunaApiKeyHeader),
			).
				SetBody(body).
				SetContext(ctx).
				Post(fmt.Sprintf("/merchants/orders/%s/void", orderToken))

			if err != nil {
				return false, fmt.Errorf(captureFlowError, "Release", httpResponse.StatusCode())
			}

			return d.handleDeUnaResponse(httpResponse, "Release")
		},
	)
}

func (d DeUnaCaptureFlowHTTPClient) Capture(
	oldCtx context.Context,
	orderToken string,
	amount int64,
) (bool, error) {
	return decorators.TraceDecorator[bool](
		oldCtx,
		"DeUnaCaptureFlowHTTPClient.Capture",
		func(ctx context.Context, span decorators.Span) (bool, error) {
			body := request.CaptureRequestDTO{
				Amount: amount,
			}

			httpResponse, err := d.Client.NewRequestWithOptions(
				instrument.WithHeadersLogConfig(true, deunaConfig.DeunaApiKeyHeader),
			).
				SetBody(body).
				SetContext(ctx).
				Post(fmt.Sprintf("/merchants/orders/%s/capture", orderToken))

			if err != nil {
				return false, fmt.Errorf(captureFlowError, "Capture", httpResponse.StatusCode())
			}

			return d.handleDeUnaResponse(httpResponse, "Capture")
		},
	)
}
