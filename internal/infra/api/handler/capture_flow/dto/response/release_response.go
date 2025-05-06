package response

type ReleaseResponse struct {
	ReferenceOrderID string `json:"reference_order_id"`
	PaymentID        string `json:"payment_id"`
	PaymentStatus    string `json:"payment_status"`
}

func NewReleaseResponse(referenceOrderID string, paymentID string, paymentStatus string) *ReleaseResponse {
	return &ReleaseResponse{
		ReferenceOrderID: referenceOrderID,
		PaymentID:        paymentID,
		PaymentStatus:    paymentStatus,
	}
}
