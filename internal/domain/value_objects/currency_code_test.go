package value_objects

import (
	"testing"
)

func TestNewCurrencyCode(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		wantErr    bool
		wantCode   string
		wantSymbol string
	}{
		{
			name:       "valid USD code",
			code:       "USD",
			wantErr:    false,
			wantCode:   "USD",
			wantSymbol: "$",
		},
		{
			name:       "valid EUR code",
			code:       "EUR",
			wantErr:    false,
			wantCode:   "EUR",
			wantSymbol: "€",
		},
		{
			name:    "invalid code",
			code:    "XYZ",
			wantErr: true,
		},
		{
			name:    "empty code",
			code:    "",
			wantErr: true,
		},
		{
			name:    "case sensitive code mismatch",
			code:    "usd",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCurrencyCode(tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCurrencyCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Code != tt.wantCode {
					t.Errorf("NewCurrencyCode() got.Code = %v, want %v", got.Code, tt.wantCode)
				}
				if got.Symbol != tt.wantSymbol {
					t.Errorf("NewCurrencyCode() got.Symbol = %v, want %v", got.Symbol, tt.wantSymbol)
				}
			}
		})
	}
}
