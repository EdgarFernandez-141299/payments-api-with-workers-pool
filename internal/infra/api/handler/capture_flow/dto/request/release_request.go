package request

import "fmt"

type ReleaseRequest struct {
	ReferenceOrderID string `json:"reference_order_id" validate:"required"`
	PaymentID        string `json:"payment_id" validate:"required"`
	Reason           string `json:"reason"`
}

func (r *ReleaseRequest) Validate() error {
	if r.ReferenceOrderID == "" {
		return fmt.Errorf("reference_order_id is required")
	}

	if r.PaymentID == "" {
		return fmt.Errorf("payment_id is required")
	}

	return nil
}
