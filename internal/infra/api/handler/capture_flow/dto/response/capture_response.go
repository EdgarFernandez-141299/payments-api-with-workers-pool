package response

type CaptureResponse struct {
	ReferenceOrderID string `json:"reference_order_id"`
	PaymentID        string `json:"payment_id"`
	PaymentStatus    string `json:"payment_status"`
}

func NewCaptureResponse(referenceOrderID string, paymentID string, paymentStatus string) *CaptureResponse {
	return &CaptureResponse{
		ReferenceOrderID: referenceOrderID,
		PaymentID:        paymentID,
		PaymentStatus:    paymentStatus,
	}
}
