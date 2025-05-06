package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidCountryCode(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name:     "Valid country code",
			code:     "Mexico",
			expected: true,
		},
		{
			name:     "Invalid country code",
			code:     "InvalidCountry",
			expected: false,
		},
		{
			name:     "Empty string",
			code:     "",
			expected: false,
		},
		{
			name:     "Unknown country",
			code:     "Unknown",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidCountryCode(tt.code)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetCountryIso3ByCode(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "Valid country code",
			code:     "Mexico",
			expected: "MEX",
		},
		{
			name:     "Invalid country code",
			code:     "InvalidCountry",
			expected: "Unknown",
		},
		{
			name:     "Empty string",
			code:     "",
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCountryIso3ByCode(tt.code)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRemoveDuplicateCurrencies(t *testing.T) {
	tests := []struct {
		name       string
		currencies []string
		expected   []string
	}{
		{
			name:       "No duplicates",
			currencies: []string{"USD", "EUR", "MXN"},
			expected:   []string{"USD", "EUR", "MXN"},
		},
		{
			name:       "With duplicates",
			currencies: []string{"USD", "EUR", "MXN", "USD", "EUR"},
			expected:   []string{"USD", "EUR", "MXN"},
		},
		{
			name:       "Empty slice",
			currencies: []string{},
			expected:   []string{},
		},
		{
			name:       "All duplicates",
			currencies: []string{"USD", "USD", "USD"},
			expected:   []string{"USD"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveDuplicateCurrencies(tt.currencies)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateCurrencies(t *testing.T) {
	tests := []struct {
		name        string
		currencies  []string
		expectError bool
	}{
		{
			name:        "Valid currencies",
			currencies:  []string{"USD", "EUR", "MXN"},
			expectError: false,
		},
		{
			name:        "Invalid currency",
			currencies:  []string{"USD", "INVALID", "MXN"},
			expectError: true,
		},
		{
			name:        "Empty slice",
			currencies:  []string{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCurrencies(tt.currencies)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
