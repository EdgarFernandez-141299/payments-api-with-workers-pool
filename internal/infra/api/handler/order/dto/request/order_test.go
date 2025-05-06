package dto

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestCreateOrderRequestDTO_Validate(t *testing.T) {
	tests := []struct {
		name    string
		dto     *CreateOrderRequestDTO
		wantErr bool
	}{
		{
			name: "valid request",
			dto: &CreateOrderRequestDTO{
				ReferenceOrderID: "123",
				CurrencyCode:     "USD",
				CountryCode:      "US",
				UserType:         "member",
				UserID:           "user-123",
				TotalAmount:      decimal.NewFromFloat(100.0),
				Email:            "test@example.com",
				BillingAddress: Address{
					ZipCode:     "12345",
					Street:      "123 Main St",
					CountryCode: "US",
					City:        "New York",
				},
				WebhookUrl: "https://example.com/webhook",
			},
			wantErr: false,
		},
		{
			name: "missing required fields",
			dto: &CreateOrderRequestDTO{
				ReferenceOrderID: "123",
			},
			wantErr: true,
		},
		{
			name: "total amount is zero",
			dto: &CreateOrderRequestDTO{
				ReferenceOrderID: "123",
				CurrencyCode:     "USD",
				CountryCode:      "US",
				UserType:         "member",
				UserID:           "user-123",
				TotalAmount:      decimal.NewFromFloat(0),
				Email:            "test@example.com",
				BillingAddress: Address{
					ZipCode:     "12345",
					Street:      "123 Main St",
					CountryCode: "US",
					City:        "New York",
				},
				WebhookUrl: "https://example.com/webhook",
			},
			wantErr: true,
		},
		{
			name: "total amount is negative",
			dto: &CreateOrderRequestDTO{
				ReferenceOrderID: "123",
				CurrencyCode:     "USD",
				CountryCode:      "US",
				UserType:         "member",
				UserID:           "user-123",
				TotalAmount:      decimal.NewFromFloat(-100.0),
				Email:            "test@example.com",
				BillingAddress: Address{
					ZipCode:     "12345",
					Street:      "123 Main St",
					CountryCode: "US",
					City:        "New York",
				},
				WebhookUrl: "https://example.com/webhook",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dto.Validate()
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestCreateOrderRequestDTO_ToCommand(t *testing.T) {
	tests := []struct {
		name         string
		dto          *CreateOrderRequestDTO
		enterpriseID string
		wantErr      bool
	}{
		{
			name: "valid request",
			dto: &CreateOrderRequestDTO{
				ReferenceOrderID: "123",
				CurrencyCode:     "USD",
				CountryCode:      "US",
				UserType:         "member",
				UserID:           "user-123",
				TotalAmount:      decimal.NewFromFloat(100.50),
				PhoneNumber:      "1234567890",
				Email:            "test@example.com",
				BillingAddress: Address{
					ZipCode:     "12345",
					Street:      "123 Main St",
					CountryCode: "US",
					City:        "New York",
				},
			},
			enterpriseID: "enterprise-123",
			wantErr:      false,
		},
		{
			name: "invalid country code",
			dto: &CreateOrderRequestDTO{
				ReferenceOrderID: "123",
				CurrencyCode:     "USD",
				CountryCode:      "INVALID",
				UserType:         "member",
				UserID:           "user-123",
				TotalAmount:      decimal.NewFromFloat(100.50),
				PhoneNumber:      "1234567890",
				Email:            "test@example.com",
				BillingAddress: Address{
					ZipCode:     "12345",
					Street:      "123 Main St",
					CountryCode: "INVALID",
					City:        "New York",
				},
			},
			enterpriseID: "enterprise-123",
			wantErr:      true,
		},
		{
			name: "invalid currency code",
			dto: &CreateOrderRequestDTO{
				ReferenceOrderID: "123",
				CurrencyCode:     "INVALID",
				CountryCode:      "US",
				UserType:         "member",
				UserID:           "user-123",
				TotalAmount:      decimal.NewFromFloat(100.50),
				PhoneNumber:      "1234567890",
				Email:            "test@example.com",
				BillingAddress: Address{
					ZipCode:     "12345",
					Street:      "123 Main St",
					CountryCode: "US",
					City:        "New York",
				},
			},
			enterpriseID: "enterprise-123",
			wantErr:      true,
		},
		{
			name: "missing required fields",
			dto: &CreateOrderRequestDTO{
				ReferenceOrderID: "",
				CurrencyCode:     "",
				CountryCode:      "",
				UserType:         "",
				UserID:           "",
				TotalAmount:      decimal.Decimal{},
				PhoneNumber:      "",
				Email:            "",
				BillingAddress: Address{
					ZipCode:     "",
					Street:      "",
					CountryCode: "",
					City:        "",
				},
			},
			enterpriseID: "enterprise-123",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.dto.ToCommand(tt.enterpriseID)
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}

}
