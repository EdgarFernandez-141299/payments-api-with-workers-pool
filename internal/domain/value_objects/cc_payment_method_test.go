package value_objects

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	paymentmethodsVo "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects/payment_methods"
)

func TestNewCCPaymentMethod(t *testing.T) {
	tests := []struct {
		name   string
		cardID string
		cvv    string
		want   PaymentMethod
	}{
		{
			name:   "crear método de pago con tarjeta de crédito válida",
			cardID: "card_123",
			cvv:    "123",
			want: PaymentMethod{
				Type: enums.CCMethod,
				CCData: PaymentMethodData[CardInfo]{
					Data: CardInfo{
						CardID: "card_123",
						CVV:    "123",
					},
				},
			},
		},
		{
			name:   "crear método de pago con valores vacíos",
			cardID: "",
			cvv:    "",
			want: PaymentMethod{
				Type: enums.CCMethod,
				CCData: PaymentMethodData[CardInfo]{
					Data: CardInfo{
						CardID: "",
						CVV:    "",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCCPaymentMethod(tt.cardID, tt.cvv)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewTokenCardPaymentMethod(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		cvv      string
		brand    string
		last4    string
		exp      string
		cardType string
		alias    string
		saveCard bool
		want     PaymentMethod
	}{
		{
			name:     "crear método de pago con tarjeta tokenizada válida",
			token:    "tok_123",
			cvv:      "123",
			brand:    "VISA",
			last4:    "4242",
			exp:      "12/25",
			cardType: "credit",
			alias:    "Mi tarjeta",
			saveCard: true,
			want: PaymentMethod{
				Type: enums.TokenCard,
				TokenCard: PaymentMethodData[paymentmethodsVo.TokenCard]{
					Data: paymentmethodsVo.TokenCard{
						Token:    "tok_123",
						CVV:      "123",
						Brand:    "VISA",
						Last4:    "4242",
						Exp:      "12/25",
						CardType: "credit",
						Alias:    "Mi tarjeta",
						SaveCard: true,
					},
				},
			},
		},
		{
			name:     "crear método de pago con tarjeta tokenizada sin guardar",
			token:    "tok_456",
			cvv:      "456",
			brand:    "MASTERCARD",
			last4:    "1234",
			exp:      "06/24",
			cardType: "debit",
			alias:    "Tarjeta de débito",
			saveCard: false,
			want: PaymentMethod{
				Type: enums.TokenCard,
				TokenCard: PaymentMethodData[paymentmethodsVo.TokenCard]{
					Data: paymentmethodsVo.TokenCard{
						Token:    "tok_456",
						CVV:      "456",
						Brand:    "MASTERCARD",
						Last4:    "1234",
						Exp:      "06/24",
						CardType: "debit",
						Alias:    "Tarjeta de débito",
						SaveCard: false,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTokenCardPaymentMethod(
				tt.token,
				tt.cvv,
				tt.brand,
				tt.last4,
				tt.exp,
				tt.cardType,
				tt.alias,
				tt.saveCard,
			)
			assert.Equal(t, tt.want, got)
		})
	}
}
