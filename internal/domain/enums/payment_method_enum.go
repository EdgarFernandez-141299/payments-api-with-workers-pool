package enums

import "errors"

type PaymentMethodEnum string

var ErrInvalidPaymentMethod = errors.New("invalid payment method")

const (
	CCMethod      PaymentMethodEnum = "CCData"
	PaymentDevice PaymentMethodEnum = "DEVICE"
	TokenCard     PaymentMethodEnum = "TOKEN_CARD"
)

func (p PaymentMethodEnum) String() string {
	return string(p)
}

func NewPaymentMethodsFromString(status string) (PaymentMethodEnum, error) {
	switch status {
	case CCMethod.String():
		return CCMethod, nil
	case PaymentDevice.String():
		return PaymentDevice, nil
	case TokenCard.String():
		return TokenCard, nil
	}

	return "", ErrInvalidPaymentMethod
}

func (p PaymentMethodEnum) IsValid() bool {
	switch p {
	case CCMethod, PaymentDevice, TokenCard:
		return true
	}

	return false
}
