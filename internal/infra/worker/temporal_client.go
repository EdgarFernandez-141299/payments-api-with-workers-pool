package worker

import (
	"context"
	workflows "gitlab.com/clubhub.ai1/go-libraries/saga/client"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/config"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/observability/adapters"
	"go.temporal.io/sdk/client"
	"go.uber.org/fx"
)

func NewTemporalClient(
	lc fx.Lifecycle, logger adapters.Logger,
) client.Client {
	temporalConfig := config.Config().TemporalConfig
	host := temporalConfig.Host
	namespace := temporalConfig.Namespace

	builder := workflows.NewTemporalClientBuilder(logger).
		WithHost(host).
		WithNamespace(namespace).
		WithCertOption(
			*workflows.NewCertOptionsFromConfig(temporalConfig),
		)

	temporalClient, err := builder.Build()

	if err != nil {
		panic(err)
	}

	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				temporalClient.Close()

				return nil
			},
		})

	return temporalClient
}
