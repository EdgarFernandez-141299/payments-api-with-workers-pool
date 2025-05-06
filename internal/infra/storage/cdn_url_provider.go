package storage

import "gitlab.com/clubhub.ai1/organization/backend/payments-api/config"

// CDNURLProvider is an interface for accessing the CDN URL
type CDNURLProvider interface {
	GetCDNURL() string
}

// CDNURLProviderImpl is an implementation of CDNURLProvider
type CDNURLProviderImpl struct {
	cdnURL string
}

// NewCDNURLProvider creates a new CDNURLProvider
func NewCDNURLProvider() CDNURLProvider {
	cdnURL := config.Config().Aws.CDN_URL
	return &CDNURLProviderImpl{
		cdnURL: cdnURL,
	}
}

// GetCDNURL returns the CDN URL
func (c *CDNURLProviderImpl) GetCDNURL() string {
	return c.cdnURL
}