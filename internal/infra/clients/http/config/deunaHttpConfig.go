package config

import (
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/config"
)

const DeunaApiKeyHeader = "x-api-key"

type DeUnaApiConfig struct {
	URL          string
	ApiKey       string
	PublicKeyPEM string
}

func DeunaHttpConfig() *DeUnaApiConfig {

	return &DeUnaApiConfig{
		URL:          config.Config().DeUnaApi.URL,
		ApiKey:       config.Config().DeUnaApi.ApiKey,
		PublicKeyPEM: config.Config().DeUnaApi.PublicKeyPEM,
	}
}
