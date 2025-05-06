package paymentmethods

import (
	"strconv"
	"testing"
	"time"
)

func TestTokenCard_Validate(t *testing.T) {
	// Obtener el año actual para las pruebas
	currentYear := time.Now().Year() % 100
	nextYear := currentYear + 1
	prevYear := currentYear - 1

	tests := []struct {
		name    string
		card    TokenCard
		wantErr bool
		errMsg  string
	}{
		{
			name: "Tarjeta válida Visa",
			card: TokenCard{
				Token: "valid_token",
				CVV:   "123",
				Brand: VisaCard,
				Exp:   "12" + strconv.Itoa(nextYear),
			},
			wantErr: false,
		},
		{
			name: "Tarjeta válida Amex",
			card: TokenCard{
				Token: "valid_token",
				CVV:   "1234",
				Brand: AmexCard,
				Exp:   "12" + strconv.Itoa(nextYear),
			},
			wantErr: false,
		},
		{
			name: "Token vacío",
			card: TokenCard{
				Token: "",
				CVV:   "123",
				Brand: VisaCard,
				Exp:   "12" + strconv.Itoa(nextYear),
			},
			wantErr: true,
			errMsg:  "token is required for credit card payment",
		},
		{
			name: "CVV incorrecto para Visa",
			card: TokenCard{
				Token: "valid_token",
				CVV:   "12",
				Brand: VisaCard,
				Exp:   "12" + strconv.Itoa(nextYear),
			},
			wantErr: true,
			errMsg:  "CVV must be a 3-digit number for credit card payment",
		},
		{
			name: "CVV incorrecto para Amex",
			card: TokenCard{
				Token: "valid_token",
				CVV:   "123",
				Brand: AmexCard,
				Exp:   "12" + strconv.Itoa(nextYear),
			},
			wantErr: true,
			errMsg:  "CVV must be a 4-digit number for American Express",
		},
		{
			name: "Fecha de expiración vacía",
			card: TokenCard{
				Token: "valid_token",
				CVV:   "123",
				Brand: VisaCard,
				Exp:   "",
			},
			wantErr: true,
			errMsg:  "expiration date is required for credit card payment",
		},
		{
			name: "Tarjeta vencida",
			card: TokenCard{
				Token: "valid_token",
				CVV:   "123",
				Brand: VisaCard,
				Exp:   "12" + strconv.Itoa(prevYear),
			},
			wantErr: true,
			errMsg:  "card is expired",
		},
		{
			name: "Formato de año inválido",
			card: TokenCard{
				Token: "valid_token",
				CVV:   "123",
				Brand: VisaCard,
				Exp:   "12XX",
			},
			wantErr: true,
			errMsg:  "invalid expiration year format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.card.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("TokenCard.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("TokenCard.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}
