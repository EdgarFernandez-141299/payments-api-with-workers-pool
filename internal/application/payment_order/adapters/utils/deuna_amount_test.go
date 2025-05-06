package utils

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestNewDeunaAmount(t *testing.T) {
	tests := []struct {
		name     string
		input    decimal.Decimal
		expected int64
	}{
		{
			name:     "Zero",
			input:    decimal.NewFromInt(0),
			expected: 0,
		},
		{
			name:     "PositiveInteger",
			input:    decimal.NewFromInt(1),
			expected: 100,
		},
		{
			name:     "NegativeInteger",
			input:    decimal.NewFromInt(-2),
			expected: -200,
		},
		{
			name:     "PositiveDecimal",
			input:    decimal.NewFromFloat(1.234),
			expected: 123,
		},
		{
			name:     "NegativeDecimal",
			input:    decimal.NewFromFloat(-3.456),
			expected: -345,
		},
		{
			name:     "RoundingDown",
			input:    decimal.NewFromFloat(7.999),
			expected: 799,
		},
		{
			name:     "RoundingUpIgnored",
			input:    decimal.NewFromFloat(12.3456),
			expected: 1234,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewDeunaAmount(test.input)
			if got != test.expected {
				t.Errorf("NewDeunaAmount(%v) = %v; want %v", test.input, got, test.expected)
			}
		})
	}
}

func TestDeunaAmountToAmount(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{
			name:     "Zero",
			input:    0,
			expected: "0",
		},
		{
			name:     "SmallPositive",
			input:    100,
			expected: "1",
		},
		{
			name:     "SmallNegative",
			input:    -250,
			expected: "-2.5",
		},
		{
			name:     "TruncatedValue",
			input:    1234,
			expected: "12.34",
		},
		{
			name:     "LargePositive",
			input:    9223372036854775807,
			expected: "92233720368547758.07",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := DeunaAmountToAmount(test.input)
			if got.String() != test.expected {
				t.Errorf("DeunaAmountToAmount(%v) = %v; want %v", test.input, got.String(), test.expected)
			}
		})
	}
}
