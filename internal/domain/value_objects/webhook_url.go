package value_objects

import (
	"errors"
	"net/url"
	"strings"
)

var (
	ErrInvalidWebhookUrl = errors.New("invalid webhook URL")
	ErrEmptyWebhookUrl   = errors.New("webhook URL cannot be empty")
)

type WebhookUrl struct {
	Url string
}

func NewWebhookUrl(url string) WebhookUrl {
	return WebhookUrl{Url: url}
}

func (w WebhookUrl) IsEmpty() bool {
	return len(w.Url) == 0
}

func (w WebhookUrl) ValidateUrl() error {
	if w.Url == "" {
		return ErrEmptyWebhookUrl
	}

	// Verificar si la URL contiene espacios o caracteres especiales no permitidos
	if strings.ContainsAny(w.Url, " \t\n\r") {
		return ErrInvalidWebhookUrl
	}

	parsedUrl, err := url.Parse(w.Url)
	if err != nil {
		return ErrInvalidWebhookUrl
	}

	if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
		return ErrInvalidWebhookUrl
	}

	return nil
}

func (w WebhookUrl) String() string {
	return w.Url
}
