package value_objects

import (
	"testing"
)

func TestNewWebhookUrl(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "URL válida con HTTPS",
			url:     "https://example.com/webhook",
			wantErr: false,
		},
		{
			name:    "URL válida con HTTP",
			url:     "http://example.com/webhook",
			wantErr: false,
		},
		{
			name:    "URL válida con puerto",
			url:     "https://example.com:8080/webhook",
			wantErr: false,
		},
		{
			name:    "URL válida con query parameters",
			url:     "https://example.com/webhook?key=value",
			wantErr: false,
		},
		{
			name:    "URL inválida sin protocolo",
			url:     "example.com/webhook",
			wantErr: true,
		},
		{
			name:    "URL inválida con caracteres especiales",
			url:     "https://example.com/webhook with spaces",
			wantErr: true,
		},
		{
			name:    "URL vacía",
			url:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			webhookUrl := NewWebhookUrl(tt.url)
			err := webhookUrl.ValidateUrl()

			if (err != nil) != tt.wantErr {
				t.Errorf("WebhookUrl.ValidateUrl() error = %v, wantErr %v", err, tt.wantErr)
			}

			if webhookUrl.String() != tt.url {
				t.Errorf("WebhookUrl.String() = %v, want %v", webhookUrl.String(), tt.url)
			}
		})
	}
}
