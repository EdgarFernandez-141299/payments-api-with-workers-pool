package response

type PaymentMethodResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	CreatedAt string `json:"created_at"`
}

func NewPaymentMethodResponse(id, name, code, createdAt string) PaymentMethodResponse {
	return PaymentMethodResponse{
		ID:        id,
		Name:      name,
		Code:      code,
		CreatedAt: createdAt,
	}
}
