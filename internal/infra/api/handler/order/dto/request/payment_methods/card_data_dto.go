package paymentmethods

import "errors"

type CreditCard struct {
	ID  string `json:"id"`
	CVV string `json:"cvv"`
}

func (c *CreditCard) Validate() error {
	if c.ID == "" || c.CVV == "" {
		return errors.New("id and cvv are required")
	}

	return nil
}
