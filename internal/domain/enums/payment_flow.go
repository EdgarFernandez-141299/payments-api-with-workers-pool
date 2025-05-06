package enums

import (
	"errors"
	"strings"
)

type PaymentFlowEnum string

var ErrInvalidPaymentFlow = errors.New("invalid payment flow")

const (
	Autocapture PaymentFlowEnum = "AUTOCAPTURE"
	Capture     PaymentFlowEnum = "CAPTURE"
)

const (
	DebitCard  = "DEBIT_CARD"
	CreditCard = "CREDIT_CARD"

	AutoCapture   = "auto-capture"
	Authorization = "authorization"
)

func (p PaymentFlowEnum) String() string {
	return string(p)
}

func NewPaymentFlowEnum(cardType string, allowCapture bool) (PaymentFlowEnum, error) {

	if strings.ToUpper(cardType) == DebitCard {
		return Autocapture, nil
	}

	if strings.ToUpper(cardType) == CreditCard && allowCapture {
		return Capture, nil
	}

	if !allowCapture {
		return Autocapture, nil
	}

	return "", ErrInvalidPaymentFlow
}

func (p PaymentFlowEnum) IsValid() bool {
	switch p {
	case Autocapture, Capture:
		return true
	}

	return false
}

func (p PaymentFlowEnum) DeunaFlowType() (string, error) {
	switch p {
	case Autocapture:
		return AutoCapture, nil
	case Capture:
		return Authorization, nil
	}

	return "", ErrInvalidPaymentFlow
}
