package projections

import (
	"github.com/shopspring/decimal"
)

type OrderPaymentsProjection struct {
	ReferenceOrderID  string
	UserID            string
	TotalAmount       decimal.Decimal
	CurrencyCode      string
	CountryCode       string
	OrderStatus       string
	Metadata          string
	CardID            string
	PaymentID         string
	PaymentMethod     string
	PaymentStatus     string
	AuthorizationCode string
	PaymentOrderID    string
}
