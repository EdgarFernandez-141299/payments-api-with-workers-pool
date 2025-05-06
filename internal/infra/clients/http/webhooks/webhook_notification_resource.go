package webhooks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

type WebhookNotificationResource struct {
	client *http.Client
}

func NewWebhookNotificationResource() resources.WebhookNotificationResourceIF {
	return &WebhookNotificationResource{
		client: &http.Client{},
	}
}

func (w *WebhookNotificationResource) SendNotification(oldCtx context.Context, webhookUrl value_objects.WebhookUrl, payload interface{}) error {
	return decorators.TraceDecoratorNoReturn(
		oldCtx,
		"WebhookNotificationResource.SendNotification",
		func(ctx context.Context, span decorators.Span) error {
			var err error

			responseData := map[string]interface{}{
				"payload": payload,
				"webhook": webhookUrl.String(),
				"error":   err,
			}

			jsonPayload, err := json.Marshal(responseData)
			if err != nil {
				return fmt.Errorf("error marshaling payload: %w", err)
			}

			req, err := http.NewRequestWithContext(ctx, "POST", webhookUrl.String(), bytes.NewBuffer(jsonPayload))

			if err != nil {
				responseData["error"] = err
				return fmt.Errorf("error creating request: %v", responseData)
			}

			req.Header.Set("Content-Type", "application/json")

			resp, err := w.client.Do(req)

			if err != nil {
				responseData["error"] = err
				return fmt.Errorf("error sending request: %v", responseData)
			}

			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("notification failed with status code: %d, response: %v", resp.StatusCode, responseData)
			}

			return nil
		},
	)
}
