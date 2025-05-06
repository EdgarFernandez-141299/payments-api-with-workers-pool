package config

import (
	"gitlab.com/clubhub.ai1/go-libraries/observability/apm"
	"gitlab.com/clubhub.ai1/go-libraries/observability/clients/http/instrument"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/config"
	"time"
)

const (
	MailServiceTimeout        = 10 * time.Second
	broadcastNotificationPath = "/api/v1/notifications/broadcast"
)

type MailClientHTTPIF interface {
	instrument.Client
}

type MailHTTPClientImpl struct {
	instrument.Client
}

func NewMailHTTPClient(tracer apm.Tracer) MailClientHTTPIF {
	if config.Config().MailService.URL == "" {
		panic("mail service url is empty")
	}

	return &MailHTTPClientImpl{
		instrument.NewInstrumentedClient(
			instrument.WithTraceOptions(tracer.GetTracer(), instrument.TraceRequest, instrument.TraceResponse),
			instrument.WithBaseUrl(config.Config().MailService.URL),
			instrument.WithRequestTimeout(MailServiceTimeout),
		),
	}
}
