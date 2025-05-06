package value_objects

import "fmt"

const cvvLength = 3

type CardInfo struct {
	CardID string
	CVV    string
}

func (c CardInfo) Validate() error {
	if c.CardID == "" {
		return fmt.Errorf("CardID is required for credit card payment")
	}

	if len(c.CVV) != cvvLength {
		return fmt.Errorf("CVV must be a 3-digit number for credit card payment")
	}

	return nil
}
