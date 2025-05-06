package response

import "github.com/shopspring/decimal"

type RefundResponseDTO struct {
	ReferenceOrderID string          `json:"reference_order_id"`
	PaymentOrderID   string          `json:"payment_order_id"`
	Amount           decimal.Decimal `json:"amount"`
}
