package value_objects

import (
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
)

type PaymentDevice struct {
}

func NewDevicePaymentMethod() PaymentMethod {
	return PaymentMethod{
		Type: enums.PaymentDevice,
	}
}

func (t PaymentDevice) Validate() error {
	return nil
}
