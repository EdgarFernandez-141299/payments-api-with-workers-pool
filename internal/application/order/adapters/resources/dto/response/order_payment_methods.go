package response

type DeunaOrderPaymentMethodsResponseDTO struct {
	Data []PaymentMethod `json:"data"`
}

type PaymentMethod struct {
	Enabled        bool                   `json:"enabled"`
	ExcludeCVV     bool                   `json:"exclude_cvv"`
	InputSchema    []InputField           `json:"input_schema"`
	Labels         map[string]string      `json:"labels"`
	MethodType     string                 `json:"method_type"`
	ProcessorName  string                 `json:"processor_name"`
	SpecificFields map[string]interface{} `json:"specific_fields"`
	Installments   []InstallmentOption    `json:"installments,omitempty"`
}

type InputField struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
	Type     string `json:"type"`
	Always   bool   `json:"always,omitempty"`
}

type InstallmentOption struct {
	BaseTotalAmount         int     `json:"base_total_amount"`
	Currency                string  `json:"currency"`
	ID                      string  `json:"id"`
	InstallmentAmount       int     `json:"installment_amount"`
	Installments            int     `json:"installments"`
	PlanID                  string  `json:"plan_id"`
	Rate                    float64 `json:"rate"`
	TotalAmountWithInterest int     `json:"total_amount_with_interest"`
	Type                    string  `json:"type"`
}
