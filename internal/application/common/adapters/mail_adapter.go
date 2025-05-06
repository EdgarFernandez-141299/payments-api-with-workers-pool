package adapters

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/notification/constants"
)

type Metadata map[string]any

func (m *Metadata) Add(key, value string) {
	if *m == nil {
		*m = make(Metadata)
	}
	(*m)[key] = value
}

type MailRequest struct {
	Recipient string                          `json:"recipient"`
	Content   string                          `json:"content"`
	Title     string                          `json:"title"`
	Channels  []constants.NotificationChannel `json:"channels"`
	Metadata  Metadata                        `json:"metadata"`
}

type MailAdapterIF interface {
	Send(ctx context.Context, request MailRequest) error
}
