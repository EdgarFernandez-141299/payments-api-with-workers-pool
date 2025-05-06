package request

import "fmt"

type CaptureRequest struct {
	ReferenceOrderID string `json:"reference_order_id" validate:"required"`
	PaymentID        string `json:"payment_id" validate:"required"`
}

func (r *CaptureRequest) Validate() error {
	if r.ReferenceOrderID == "" {
		return fmt.Errorf("reference_order_id is required")
	}

	if r.PaymentID == "" {
		return fmt.Errorf("payment_id is required")
	}

	return nil
}
