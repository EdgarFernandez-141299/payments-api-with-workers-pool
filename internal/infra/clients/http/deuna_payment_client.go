package http

import (
	"context"
	"fmt"
	"gitlab.com/clubhub.ai1/go-libraries/observability/apm"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/utils"
	"net/http"

	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	"gitlab.com/clubhub.ai1/go-libraries/observability/clients/http/instrument"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/resources/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/resources/dto/response"
	deunaConfig "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/http/config"
)

const (
	failedPaymentRefund string = "failed to refund payment"
)

type DeUnaPaymentHTTPClient struct {
	instrument.Client
}

func NewDeUnaPaymentHTTPClient(config *deunaConfig.DeUnaApiConfig, tracer apm.Tracer) resources.DeunaPaymentResourceIF {
	apiKeyHeader := deunaConfig.DeunaApiKeyHeader

	return &DeUnaPaymentHTTPClient{
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

func handleError(err error) error {
	return fmt.Errorf("%s: %w", failedPaymentRefund, err)
}

func (d DeUnaPaymentHTTPClient) MakeOrderPayment(
	oldCtx context.Context,
	body request.DeunaOrderPaymentRequest,
	token string,
) error {
	return decorators.TraceDecoratorNoReturn(
		oldCtx,
		"DeUnaPaymentHTTPClient.MakeOrderPayment",
		func(ctx context.Context, span decorators.Span) error {
			var res response.DeunaOrderPaymentResponse

			httpResponse, err := d.Client.NewRequestWithOptions(
				instrument.WithHeadersLogConfig(true, deunaConfig.DeunaApiKeyHeader),
			).
				SetResult(&res).
				SetBody(body).
				SetContext(ctx).
				SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
				Post("/merchants/transactions/purchase")

			if err != nil {
				return handleError(err)
			}

			if httpResponse.IsError() || httpResponse.StatusCode() != http.StatusOK {
				return fmt.Errorf("make order payment failed with status: %d", httpResponse.StatusCode())
			}

			return nil
		},
	)
}

func (d DeUnaPaymentHTTPClient) MakeOrderPaymentV2(
	oldCtx context.Context,
	body request.DeunaOrderPaymentRequestV2,
) (response.DeunaOrderPaymentResponseV2, error) {
	return decorators.TraceDecorator[response.DeunaOrderPaymentResponseV2](
		oldCtx,
		"DeUnaPaymentHTTPClient.MakeOrderPaymentV2",
		func(ctx context.Context, span decorators.Span) (response.DeunaOrderPaymentResponseV2, error) {
			var res response.DeunaOrderPaymentResponseV2

			httpResponse, err := d.Client.NewRequestWithOptions(
				instrument.WithHeadersLogConfig(true),
			).
				SetResult(&res).
				SetBody(body).
				SetContext(ctx).
				Post("/merchants/orders/purchase")

			if err != nil {
				return response.DeunaOrderPaymentResponseV2{}, handleError(err)
			}

			if httpResponse.IsError() || httpResponse.StatusCode() != http.StatusOK {
				return response.DeunaOrderPaymentResponseV2{}, fmt.Errorf("make order payment failed with status: %d", httpResponse.StatusCode())
			}

			return res, nil
		},
	)
}

func (d DeUnaPaymentHTTPClient) MakeTotalRefund(
	oldCtx context.Context,
	body utils.DeunaTotalRefundRequest,
	orderToken string,
) (response.DeunaRefundPaymentResponse, error) {
	return decorators.TraceDecorator[response.DeunaRefundPaymentResponse](
		oldCtx,
		"DeUnaPaymentHTTPClient.MakeTotalRefund",
		func(ctx context.Context, span decorators.Span) (response.DeunaRefundPaymentResponse, error) {
			var res response.DeunaRefundPaymentResponse

			err := d.MakeRefund(ctx, body, &res, orderToken)

			if err != nil {
				return response.DeunaRefundPaymentResponse{}, handleError(err)
			}

			return res, nil
		},
	)
}

func (d DeUnaPaymentHTTPClient) MakePartialRefund(
	oldCtx context.Context,
	body utils.DeunaPartialRefundRequest,
	orderToken string,
) (response.DeunaRefundPaymentResponse, error) {
	return decorators.TraceDecorator[response.DeunaRefundPaymentResponse](
		oldCtx,
		"DeUnaPaymentHTTPClient.MakePartialRefund",
		func(ctx context.Context, span decorators.Span) (response.DeunaRefundPaymentResponse, error) {
			var res response.DeunaRefundPaymentResponse

			err := d.MakeRefund(ctx, body, &res, orderToken)

			if err != nil {
				return response.DeunaRefundPaymentResponse{}, handleError(err)
			}

			return res, nil
		},
	)
}

func (d DeUnaPaymentHTTPClient) MakeRefund(
	oldCtx context.Context,
	body interface{},
	result interface{},
	orderToken string,
) error {
	return decorators.TraceDecoratorNoReturn(
		oldCtx,
		"DeUnaPaymentHTTPClient.MakeRefund",
		func(ctx context.Context, span decorators.Span) error {
			httpResponse, err := d.Client.NewRequestWithOptions(
				instrument.WithHeadersLogConfig(true),
			).
				SetResult(result).
				SetBody(body).
				SetContext(ctx).
				Post(fmt.Sprintf("/v2/merchants/orders/%s/refund", orderToken))

			if err != nil {
				return handleError(err)
			}

			if httpResponse.IsError() || httpResponse.StatusCode() != http.StatusOK {
				return fmt.Errorf("refund payment failed with status: %d", httpResponse.StatusCode())
			}

			return nil
		},
	)
}
