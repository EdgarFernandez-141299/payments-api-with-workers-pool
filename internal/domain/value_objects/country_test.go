package value_objects

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	bizErrors "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/errors"
)

func TestNewCountry(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		wantErr   bool
		errorType error
	}{
		{
			name:      "valid country code",
			code:      "US",
			wantErr:   false,
			errorType: nil,
		},
		{
			name:      "valid country code mexico",
			code:      "mexico",
			wantErr:   false,
			errorType: nil,
		},
		{
			name:      "valid country code mx",
			code:      "mx",
			wantErr:   false,
			errorType: nil,
		},
		{
			name:      "empty country code",
			code:      "",
			wantErr:   true,
			errorType: errors.New("dddd"),
		},
		{
			name:      "invalid country code",
			code:      "XYZ",
			wantErr:   true,
			errorType: bizErrors.NewInvalidCountryCodeError("XYZ"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCountryWithCode(tt.code)
			if tt.wantErr {
				assert.Error(t, err)
				// assert.Equal(t, tt.errorType, err)
			} else {
				assert.NoError(t, err)
				// assert.Equal(t, tt.code, country.Code)
			}
		})
	}
}

func TestCountry_Equals(t *testing.T) {
	country1 := newCountry("US")
	country2 := newCountry("US")
	country3 := newCountry("CA")

	assert.True(t, country1.Equals(country2))
	assert.False(t, country1.Equals(country3))
}
