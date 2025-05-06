package enums

import (
	"errors"
	"strings"
)

type PaymentFlowActionEnum string

var ErrInvalidPaymentFlowAction = errors.New("invalid payment flow action")

const (
	CapturePayment PaymentFlowActionEnum = "CAPTURE"
	ReleasePayment PaymentFlowActionEnum = "RELEASE"
)

func (p PaymentFlowActionEnum) String() string {
	return string(p)
}

func NewPaymentFlowActionEnum(action string) (PaymentFlowActionEnum, error) {
	switch strings.ToUpper(action) {
	case CapturePayment.String():
		return CapturePayment, nil
	case ReleasePayment.String():
		return ReleasePayment, nil
	}

	return "", ErrInvalidPaymentFlowAction
}

func (p PaymentFlowActionEnum) IsValid() bool {
	switch p {
	case CapturePayment, ReleasePayment:
		return true
	}

	return false
}

func (p PaymentFlowActionEnum) IsCapture() bool {
	return p == CapturePayment
}

func (p PaymentFlowActionEnum) IsRelease() bool {
	return p == ReleasePayment
}
