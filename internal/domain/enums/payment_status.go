package enums

import "errors"

type PaymentStatus string

var ErrInvalidPaymentStatus = errors.New("invalid payment status")

const (
	PaymentProcessing   PaymentStatus = "PROCESSING"
	PaymentNotProcessed PaymentStatus = "NOT_PROCESSED"
	PaymentProcessed    PaymentStatus = "PROCESSED"
	PaymentFailed       PaymentStatus = "FAILED"
	PaymentAuthorized   PaymentStatus = "AUTHORIZED"
	PaymentCanceled     PaymentStatus = "CANCELED"
	PaymentRefunded     PaymentStatus = "REFUNDED"
	PartiallyRefunded   PaymentStatus = "PARTIALLY_REFUNDED"
	PaymentVoided       PaymentStatus = "VOIDED"
	PaymentDenied       PaymentStatus = "DENIED"
)

func (p PaymentStatus) String() string {
	return string(p)
}

func NewPaymentStatusFromString(status string) (PaymentStatus, error) {
	switch status {
	case PaymentProcessing.String():
		return PaymentProcessing, nil
	case PaymentProcessed.String():
		return PaymentProcessed, nil
	case PaymentVoided.String():
		return PaymentCanceled, nil
	case PaymentDenied.String():
		return PaymentDenied, nil
	case PaymentRefunded.String():
		return PaymentRefunded, nil
	case PaymentFailed.String():
		return PaymentFailed, nil
	case PaymentNotProcessed.String():
		return PaymentNotProcessed, nil
	case PartiallyRefunded.String():
		return PartiallyRefunded, nil
	case PaymentAuthorized.String():
		return PaymentAuthorized, nil
	case PaymentCanceled.String():
		return PaymentCanceled, nil
	}

	return "", ErrInvalidPaymentStatus
}

func (p PaymentStatus) IsFailure() bool {
	switch p {
	case PaymentFailed, PaymentDenied, PaymentVoided:
		return true
	}

	return false
}

func (p PaymentStatus) IsValid() bool {
	switch p {
	case PaymentProcessing, PaymentProcessed, PaymentRefunded, PaymentAuthorized:
		return true
	}

	return false
}

func (p PaymentStatus) IsProcessing() bool {
	return p == PaymentProcessing
}

func (p PaymentStatus) IsPartiallyRefunded() bool {
	return p == PartiallyRefunded
}

func (p PaymentStatus) IsProcessed() bool {
	return p == PaymentProcessed
}

func (p PaymentStatus) IsAuthorized() bool {
	return p == PaymentAuthorized
}

func (p PaymentStatus) IsCanceled() bool {
	switch p {
	case PaymentCanceled, PaymentVoided:
		return true
	}

	return false
}
