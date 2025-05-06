package enums

import (
	"testing"
)

func TestPaymentFlowEnum_String(t *testing.T) {
	tests := []struct {
		name string
		p    PaymentFlowEnum
		want string
	}{
		{
			name: "Autocapture string",
			p:    Autocapture,
			want: "AUTOCAPTURE",
		},
		{
			name: "Capture string",
			p:    Capture,
			want: "CAPTURE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.want {
				t.Errorf("PaymentFlowEnum.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPaymentFlowEnum(t *testing.T) {
	tests := []struct {
		name         string
		cardType     string
		allowCapture bool
		want         PaymentFlowEnum
		wantErr      bool
	}{
		{
			name:         "Debit card should return Autocapture",
			cardType:     DebitCard,
			allowCapture: true,
			want:         Autocapture,
			wantErr:      false,
		},
		{
			name:         "Credit card with allowCapture should return Capture",
			cardType:     CreditCard,
			allowCapture: true,
			want:         Capture,
			wantErr:      false,
		},
		{
			name:         "Credit card without allowCapture should return Autocapture",
			cardType:     CreditCard,
			allowCapture: false,
			want:         Autocapture,
			wantErr:      false,
		},
		{
			name:         "Invalid card type should return error",
			cardType:     "INVALID",
			allowCapture: true,
			want:         "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPaymentFlowEnum(tt.cardType, tt.allowCapture)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPaymentFlowEnum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewPaymentFlowEnum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaymentFlowEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		p    PaymentFlowEnum
		want bool
	}{
		{
			name: "Valid Autocapture",
			p:    Autocapture,
			want: true,
		},
		{
			name: "Valid Capture",
			p:    Capture,
			want: true,
		},
		{
			name: "Invalid payment flow",
			p:    PaymentFlowEnum("INVALID"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.IsValid(); got != tt.want {
				t.Errorf("PaymentFlowEnum.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaymentFlowEnum_DeunaFlowType(t *testing.T) {
	tests := []struct {
		name    string
		p       PaymentFlowEnum
		want    string
		wantErr bool
	}{
		{
			name:    "Autocapture to Deuna flow type",
			p:       Autocapture,
			want:    AutoCapture,
			wantErr: false,
		},
		{
			name:    "Capture to Deuna flow type",
			p:       Capture,
			want:    Authorization,
			wantErr: false,
		},
		{
			name:    "Invalid payment flow",
			p:       PaymentFlowEnum("INVALID"),
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.DeunaFlowType()
			if (err != nil) != tt.wantErr {
				t.Errorf("PaymentFlowEnum.DeunaFlowType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PaymentFlowEnum.DeunaFlowType() = %v, want %v", got, tt.want)
			}
		})
	}
}
