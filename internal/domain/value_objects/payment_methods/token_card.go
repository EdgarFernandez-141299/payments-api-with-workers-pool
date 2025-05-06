package paymentmethods

import (
	"fmt"
	"strconv"
	"time"
)

type TokenCard struct {
	Token    string
	CVV      string
	Brand    string
	Last4    string
	Exp      string
	CardType string
	SaveCard bool
	Alias    string
}

const (
	CardTypeCredit = "credit"
	CardTypeDebit  = "debit"

	AmexCard     = "amex"
	VisaCard     = "visa"
	MasterCard   = "mastercard"
	DinersCard   = "diners"
	DiscoverCard = "discover"

	lengthCVV     = 3
	lengthAmexCVV = 4
	lengthLast4   = 4
)

func (t TokenCard) Validate() error {
	if t.Token == "" {
		return fmt.Errorf("token is required for credit card payment")
	}

	if t.Brand == AmexCard {
		if len(t.CVV) != lengthAmexCVV {
			return fmt.Errorf("CVV must be a 4-digit number for American Express")
		}
	} else if len(t.CVV) != lengthCVV {
		return fmt.Errorf("CVV must be a 3-digit number for credit card payment")
	}

	if t.Exp == "" {
		return fmt.Errorf("expiration date is required for credit card payment")
	}

	// Obtiene el año actual y lo convierte a 2 dígitos (ej: 2024 -> 24)
	currentYear := time.Now().Year() % 100

	// Extrae los últimos 2 dígitos del año de expiración de la tarjeta
	// y los convierte de string a int
	expYear, err := strconv.Atoi(t.Exp[2:])
	if err != nil {
		return fmt.Errorf("invalid expiration year format")
	}

	// Compara el año de expiración con el año actual
	// Si el año de expiración es menor, la tarjeta está vencida
	if expYear < currentYear {
		return fmt.Errorf("card is expired")
	}

	return nil
}
