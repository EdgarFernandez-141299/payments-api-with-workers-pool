package value_objects

import (
	"testing"
)

func TestNewAddress(t *testing.T) {
	tests := []struct {
		name     string
		zipCode  string
		street   string
		city     string
		country  Country
		expected Address
	}{
		{
			name:     "valid address",
			zipCode:  "12345",
			street:   "Main Street",
			city:     "Metropolis",
			country:  Country{Code: "US"},
			expected: Address{ZipCode: "12345", Street: "Main Street", City: "Metropolis", Country: Country{Code: "US"}},
		},
		{
			name:     "empty fields",
			zipCode:  "",
			street:   "",
			city:     "",
			country:  Country{Code: ""},
			expected: Address{ZipCode: "", Street: "", City: "", Country: Country{Code: ""}},
		},
		{
			name:     "special characters",
			zipCode:  "987@#",
			street:   "@#$%Street",
			city:     "C!ty%",
			country:  Country{Code: "#@!"},
			expected: Address{ZipCode: "987@#", Street: "@#$%Street", City: "C!ty%", Country: Country{Code: "#@!"}},
		},
		{
			name:    "long values",
			zipCode: "12345678901234567890",
			street:  "A very very very long street name with many characters",
			city:    "A city with an equally long name",
			country: Country{Code: "USA"},
			expected: Address{
				ZipCode: "12345678901234567890",
				Street:  "A very very very long street name with many characters",
				City:    "A city with an equally long name",
				Country: Country{Code: "USA"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewAddress(tt.zipCode, tt.street, tt.city, tt.country)
			if !result.Equals(tt.expected) {
				t.Errorf("NewAddress() = %#v, expected %#v", result, tt.expected)
			}
		})
	}
}
