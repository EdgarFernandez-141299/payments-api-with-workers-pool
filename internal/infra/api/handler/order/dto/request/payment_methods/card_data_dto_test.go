package paymentmethods

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreditCard_Validate(t *testing.T) {
	tests := []struct {
		name    string
		card    CreditCard
		wantErr bool
	}{
		{
			name: "valid credit card",
			card: CreditCard{
				ID:  "123456789",
				CVV: "123",
			},
			wantErr: false,
		},
		{
			name: "empty id",
			card: CreditCard{
				ID:  "",
				CVV: "123",
			},
			wantErr: true,
		},
		{
			name: "empty cvv",
			card: CreditCard{
				ID:  "123456789",
				CVV: "",
			},
			wantErr: true,
		},
		{
			name: "empty id and cvv",
			card: CreditCard{
				ID:  "",
				CVV: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.card.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
