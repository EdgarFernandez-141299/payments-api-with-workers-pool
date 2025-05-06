package value_objects

import (
	"fmt"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	paymentmethodsVo "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects/payment_methods"
)

type SupportedPaymentMethodIF interface {
	Validate() error
}

type SupportedPaymentMethod interface {
	SupportedPaymentMethodIF
	PaymentDevice | CardInfo
}

type PaymentMethodData[T SupportedPaymentMethodIF] struct {
	Data T
}

type PaymentMethod struct {
	Type       enums.PaymentMethodEnum
	CCData     PaymentMethodData[CardInfo]
	DeviceData PaymentMethodData[PaymentDevice]
	TokenCard  PaymentMethodData[paymentmethodsVo.TokenCard]
}

func (p PaymentMethod) GetCC() (CardInfo, error) {
	if p.Type == enums.CCMethod {
		return p.CCData.Data, nil
	}

	return CardInfo{}, fmt.Errorf("Data is not a CreditCard")
}

func (p PaymentMethod) GetDevice() (PaymentDevice, error) {
	if p.Type == enums.PaymentDevice {
		return PaymentDevice{}, nil
	}

	return PaymentDevice{}, fmt.Errorf("Data is not a PaymentDevice")
}

func (p PaymentMethod) IsDevice() bool {
	return p.Type == enums.PaymentDevice
}

func (p PaymentMethod) IsCC() bool {
	return p.Type.String() == enums.CCMethod.String()
}

func (p PaymentMethod) IsTokenCard() bool {
	return p.Type.String() == enums.TokenCard.String()
}

func (p PaymentMethod) GetTokenCard() (paymentmethodsVo.TokenCard, error) {
	if p.Type == enums.TokenCard {
		return p.TokenCard.Data, nil
	}

	return paymentmethodsVo.TokenCard{}, fmt.Errorf("Data is not a TokenCard")
}

func (p PaymentMethod) Validate() error {
	if p.IsCC() {
		if err := p.CCData.Data.Validate(); err != nil {
			return err
		}
		return nil
	}

	if p.IsDevice() {
		if err := p.DeviceData.Data.Validate(); err != nil {
			return err
		}
		return nil
	}

	if p.IsTokenCard() {
		if err := p.TokenCard.Data.Validate(); err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("not supported payment method")
}
