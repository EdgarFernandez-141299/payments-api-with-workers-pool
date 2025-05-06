package webhooks

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

func TestWebhookNotificationResource_SendNotification(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		payload        interface{}
		wantErr        bool
	}{
		{
			name: "respuesta exitosa",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			payload: map[string]string{"status": "success"},
			wantErr: false,
		},
		{
			name: "respuesta de error",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			payload: map[string]string{"status": "error"},
			wantErr: true,
		},
		{
			name: "respuesta no encontrada",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			payload: map[string]string{"status": "not_found"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			webhookUrl := value_objects.NewWebhookUrl(server.URL)
			resource := NewWebhookNotificationResource()

			err := resource.SendNotification(context.Background(), webhookUrl, tt.payload)

			if (err != nil) != tt.wantErr {
				t.Errorf("WebhookNotificationResource.SendNotification() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWebhookNotificationResource_SendNotificationWithInvalidURL(t *testing.T) {
	webhookUrl := value_objects.NewWebhookUrl("not-a-valid-url")
	resource := NewWebhookNotificationResource()

	err := resource.SendNotification(context.Background(), webhookUrl, map[string]string{"status": "test"})

	if err == nil {
		t.Error("WebhookNotificationResource.SendNotification() expected error for invalid URL")
	}
}
