package response

type PaymentConceptResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

func NewPaymentConceptResponse(id, name, code, description, createdAt string) PaymentConceptResponse {
	return PaymentConceptResponse{
		ID:          id,
		Name:        name,
		Code:        code,
		Description: description,
		CreatedAt:   createdAt,
	}
}
