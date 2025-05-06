package request

import "github.com/go-playground/validator/v10"

type CollectionCenterRequest struct {
	Name                string   `json:"name" validate:"required"`
	Description         string   `json:"description" validate:"required"`
	AvailableCurrencies []string `json:"available_currencies" validate:"required"`
}

func (c *CollectionCenterRequest) Validate() error {
	err := validator.New().Struct(c)

	if err != nil {
		return err
	}

	return nil
}
