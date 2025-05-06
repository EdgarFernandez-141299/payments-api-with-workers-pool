package utils

import "fmt"

func GeneratePaymentOrderID(orderID string, paymentID string) string {
	return fmt.Sprintf("%s_%s", orderID, paymentID)
}
