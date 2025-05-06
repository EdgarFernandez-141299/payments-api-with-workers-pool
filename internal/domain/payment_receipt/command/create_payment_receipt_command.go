package command

import (
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
	"time"
)

type CreatePaymentReceiptCommand struct {
	UserID           string
	EnterpriseID     string
	Email            string
	ReferenceOrderID string
	PaymentID        string
	PaymentStatus    string
	PaymentAmount    value_objects.CurrencyAmount
	PaymentCountry   value_objects.Country
	PaymentMethod    value_objects.PaymentMethod
	PaymentDate      string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
