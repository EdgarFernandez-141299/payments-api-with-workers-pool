package http

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	"gitlab.com/clubhub.ai1/go-libraries/observability/clients/http/instrument"
	commonAdapters "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/common/adapters"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/http/config"
)

// MailResponse is the DTO for the mail service response
type MailResponse struct {
	Success   bool `json:"success"`
	TotalSent int  `json:"total_sent"`
}

const (
	broadcastNotificationPath = "/api/v1/notifications/broadcast"
)

type MailHTTPClientImpl struct {
	client config.MailClientHTTPIF
}

func NewMailHTTPClient(client config.MailClientHTTPIF) commonAdapters.MailAdapterIF {
	return &MailHTTPClientImpl{
		client: client,
	}
}

func (m *MailHTTPClientImpl) Send(oldCtx context.Context, request commonAdapters.MailRequest) error {
	return decorators.TraceDecoratorNoReturn(
		oldCtx,
		"MailHTTPClient.Send",
		func(ctx context.Context, decorators decorators.Span) error {
			var response MailResponse
			httpResponse, err := m.client.NewRequestWithOptions(
				instrument.WithHeadersLogConfig(true),
			).
				SetContext(ctx).
				SetBody(request).
				SetResult(&response).
				Post(broadcastNotificationPath)

			if err != nil {
				return fmt.Errorf("failed to send mail notification: %w", err)
			}

			if httpResponse.StatusCode() != http.StatusCreated {
				return fmt.Errorf("unexpected status code from mail service: %d", httpResponse.StatusCode())
			}

			if !response.Success {
				return fmt.Errorf("mail service reported failure: success=%v, total_sent=%d", response.Success, response.TotalSent)
			}

			return nil
		},
	)
}
