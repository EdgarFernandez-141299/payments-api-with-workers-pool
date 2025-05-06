package value_objects

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestCurrencyAmount_Equals(t *testing.T) {
	tests := []struct {
		name     string
		currency CurrencyAmount
		other    CurrencyAmount
		expected bool
	}{
		{
			name: "Equal amounts with same code and value",
			currency: CurrencyAmount{
				Code:  CurrencyCode{Code: "USD", Symbol: "$"},
				Value: decimal.NewFromInt(100),
			},
			other: CurrencyAmount{
				Code:  CurrencyCode{Code: "USD", Symbol: "$"},
				Value: decimal.NewFromInt(100),
			},
			expected: true,
		},
		{
			name: "Different values with same currency code",
			currency: CurrencyAmount{
				Code:  CurrencyCode{Code: "USD", Symbol: "$"},
				Value: decimal.NewFromInt(100),
			},
			other: CurrencyAmount{
				Code:  CurrencyCode{Code: "USD", Symbol: "$"},
				Value: decimal.NewFromInt(200),
			},
			expected: false,
		},
		{
			name: "Different currency codes with same value",
			currency: CurrencyAmount{
				Code:  CurrencyCode{Code: "USD", Symbol: "$"},
				Value: decimal.NewFromInt(100),
			},
			other: CurrencyAmount{
				Code:  CurrencyCode{Code: "EUR", Symbol: "€"},
				Value: decimal.NewFromInt(100),
			},
			expected: false,
		},
		{
			name: "Different currency codes and values",
			currency: CurrencyAmount{
				Code:  CurrencyCode{Code: "USD", Symbol: "$"},
				Value: decimal.NewFromInt(100),
			},
			other: CurrencyAmount{
				Code:  CurrencyCode{Code: "GBP", Symbol: "£"},
				Value: decimal.NewFromInt(50),
			},
			expected: false,
		},
		{
			name: "Equal amounts with empty currency code and value",
			currency: CurrencyAmount{
				Code:  CurrencyCode{Code: "", Symbol: ""},
				Value: decimal.NewFromInt(0),
			},
			other: CurrencyAmount{
				Code:  CurrencyCode{Code: "", Symbol: ""},
				Value: decimal.NewFromInt(0),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.currency.Equals(tt.other)
			if result != tt.expected {
				t.Errorf("CurrencyAmount.Equals() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCurrencyAmount_Add(t *testing.T) {
	tests := []struct {
		name        string
		initial     CurrencyAmount
		addAmount   decimal.Decimal
		expected    decimal.Decimal
		expectValid bool
	}{
		{
			name: "Add positive amount",
			initial: CurrencyAmount{
				Code:  CurrencyCode{Code: "USD", Symbol: "$"},
				Value: decimal.NewFromFloat(100.00),
			},
			addAmount:   decimal.NewFromFloat(50.00),
			expected:    decimal.NewFromFloat(150.00),
			expectValid: true,
		},
		{
			name: "Add zero amount",
			initial: CurrencyAmount{
				Code:  CurrencyCode{Code: "USD", Symbol: "$"},
				Value: decimal.NewFromFloat(100.00),
			},
			addAmount:   decimal.NewFromFloat(0.00),
			expected:    decimal.NewFromFloat(100.00),
			expectValid: true,
		},
		{
			name: "Add negative amount resulting in valid amount",
			initial: CurrencyAmount{
				Code:  CurrencyCode{Code: "USD", Symbol: "$"},
				Value: decimal.NewFromFloat(100.00),
			},
			addAmount:   decimal.NewFromFloat(-50.00),
			expected:    decimal.NewFromFloat(50.00),
			expectValid: true,
		},
		{
			name: "Add negative amount resulting in invalid amount",
			initial: CurrencyAmount{
				Code:  CurrencyCode{Code: "USD", Symbol: "$"},
				Value: decimal.NewFromFloat(50.00),
			},
			addAmount:   decimal.NewFromFloat(-100.00),
			expected:    decimal.NewFromFloat(-50.00),
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial.Add(tt.addAmount)

			if !result.Value.Equal(tt.expected) {
				t.Errorf("CurrencyAmount.Add() = %v, expected %v", result.Value, tt.expected)
			}

			err := result.Validate()
			if (err == nil) != tt.expectValid {
				t.Errorf("CurrencyAmount.Validate() valid = %v, expectValid %v", err == nil, tt.expectValid)
			}
		})
	}
}

func TestNewCurrencyAmount(t *testing.T) {
	tests := []struct {
		name      string
		code      CurrencyCode
		value     decimal.Decimal
		wantErr   bool
		wantValid bool
	}{
		{
			name:      "Valid positive amount",
			code:      CurrencyCode{Code: "USD", Symbol: "$"},
			value:     decimal.NewFromFloat(100.50),
			wantErr:   false,
			wantValid: true,
		},
		{
			name:      "Valid zero amount",
			code:      CurrencyCode{Code: "EUR", Symbol: "€"},
			value:     decimal.NewFromFloat(0.0),
			wantErr:   false,
			wantValid: true,
		},
		{
			name:      "Invalid negative amount",
			code:      CurrencyCode{Code: "GBP", Symbol: "£"},
			value:     decimal.NewFromFloat(-50.00),
			wantErr:   false,
			wantValid: false,
		},
		{
			name:      "Empty currency code",
			code:      CurrencyCode{Code: "", Symbol: ""},
			value:     decimal.NewFromFloat(10.00),
			wantErr:   false,
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currencyAmount, err := NewCurrencyAmount(tt.code, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCurrencyAmount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			err = currencyAmount.Validate()
			if (err == nil) != tt.wantValid {
				t.Errorf("CurrencyAmount.Validate() valid = %v, wantValid %v", err == nil, tt.wantValid)
			}
		})
	}
}
