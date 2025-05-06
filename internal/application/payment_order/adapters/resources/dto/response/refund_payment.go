package response

import "time"

type DeunaRefundPaymentResponse struct {
	Data RefundData `json:"data"`
}

type RefundData struct {
	RefundAmount RefundMonetaryValue `json:"refund_amount"`
	RefundID     string              `json:"refund_id"`
	Refunds      []RefundDetail      `json:"refunds"`
	Status       string              `json:"status"`
}

type RefundDetail struct {
	ExternalTransactionID string              `json:"external_transaction_id"`
	RefundAmount          RefundMonetaryValue `json:"refund_amount"`
	RefundID              string              `json:"refund_id"`
	RefundedOn            time.Time           `json:"refunded_on"`
	Status                string              `json:"status"`
}

type RefundMonetaryValue struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}
