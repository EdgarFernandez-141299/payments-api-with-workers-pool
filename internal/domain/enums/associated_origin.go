package enums

import (
	"errors"
	"strings"
)

type AssociatedOrigin string

var ErrInvalidAssociatedOrigin = errors.New("invalid associated origin")

const (
	Downpayment AssociatedOrigin = "DOWNPAYMENT"
	Loan        AssociatedOrigin = "LOAN"
	Club        AssociatedOrigin = "CLUB"
	Booking     AssociatedOrigin = "BOOKING"
)

func (a AssociatedOrigin) String() string {
	return strings.ToUpper(string(a))
}

func NewAssociatedOrigin(origin string) (AssociatedOrigin, error) {
	switch strings.ToUpper(origin) {
	case "DOWNPAYMENT":
		return Downpayment, nil
	case "LOAN":
		return Loan, nil
	case "CLUB":
		return Club, nil
	case "BOOKING":
		return Booking, nil
	}

	return "", ErrInvalidAssociatedOrigin
}

func (a AssociatedOrigin) IsValid() bool {
	switch a {
	case Downpayment, Loan, Club, Booking:
		return true
	}

	return false
}
