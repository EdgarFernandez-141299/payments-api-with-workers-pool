package utils

import (
	"fmt"
	"strings"
)

func composeDeunaOrderID(orderID string, paymentID string) string {
	return fmt.Sprintf("%s-%s", orderID, paymentID)
}

func ExtractFromDeunaOrderID(deunaOrderID string) (DeunaOrderID, error) {
	parts := strings.Split(deunaOrderID, "-")
	if len(parts) < 2 {
		return DeunaOrderID{}, fmt.Errorf("invalid deunaOrderID: %s", deunaOrderID)
	}
	orderID := strings.Join(parts[:len(parts)-1], "-")

	paymentID := parts[len(parts)-1]
	return DeunaOrderID{
		composedOrderID:  deunaOrderID,
		paymentID:        paymentID,
		referenceOrderID: orderID,
	}, nil

}

type DeunaOrderID struct {
	composedOrderID  string
	paymentID        string
	referenceOrderID string
}

func NewDeunaOrderID(orderID, paymentID string) *DeunaOrderID {
	return &DeunaOrderID{
		composedOrderID:  composeDeunaOrderID(orderID, paymentID),
		paymentID:        paymentID,
		referenceOrderID: orderID,
	}
}

func (d DeunaOrderID) GetID() string {
	return d.composedOrderID
}

func (d DeunaOrderID) GetOrderID() string {
	return d.referenceOrderID
}

func (d DeunaOrderID) GetPaymentID() string {
	return d.paymentID
}
