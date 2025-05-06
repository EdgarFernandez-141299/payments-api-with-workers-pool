package apm

import (
	"gitlab.com/clubhub.ai1/go-libraries/observability/metrics"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/config"
)

func NewMetricProvider() metrics.MetricProvider {
	url := config.Config().Otel.InternalCollectorGrpcURL
	serviceName := config.Config().Otel.ServiceName

	return metrics.NewMetricProvider(metrics.WithProviderConfig(
		metrics.NewOtelCollectorConfig(url, nil, metrics.SecureOtel),
	),
		metrics.WithServiceName(serviceName),
	)
}
