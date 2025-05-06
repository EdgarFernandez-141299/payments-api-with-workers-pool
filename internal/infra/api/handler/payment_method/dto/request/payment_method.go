package request

import "github.com/go-playground/validator/v10"

type PaymentMethodRequest struct {
	Name        string `json:"name" validate:"required"`
	Code        string `json:"code" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func (p *PaymentMethodRequest) Validate() error {
	err := validator.New().Struct(p)

	if err != nil {
		return err
	}

	return nil
}
