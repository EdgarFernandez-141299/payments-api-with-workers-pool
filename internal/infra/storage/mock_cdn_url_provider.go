package storage

// MockCDNURLProvider is a mock implementation of CDNURLProvider for testing
type MockCDNURLProvider struct {
	cdnURL string
}

// NewMockCDNURLProvider creates a new MockCDNURLProvider
func NewMockCDNURLProvider(cdnURL string) CDNURLProvider {
	return &MockCDNURLProvider{
		cdnURL: cdnURL,
	}
}

// GetCDNURL returns the CDN URL
func (m *MockCDNURLProvider) GetCDNURL() string {
	return m.cdnURL
}