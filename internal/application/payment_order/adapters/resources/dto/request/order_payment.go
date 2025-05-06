package request

type DeunaOrderPaymentRequest struct {
	BillingAddress  *Address                `json:"billing_address"`
	OrderToken      string                  `json:"token"`
	CreditCard      *CreditCardInfo         `json:"credit_card"`
	Installment     *InstallmentInfo        `json:"installment"`
	Metadata        *map[string]interface{} `json:"metadata"`
	MethodType      string                  `json:"method_type"`
	ShippingAddress *Address                `json:"shipping_address"`
	SpecificFields  *SpecificFields         `json:"specific_fields"`
	CallbackURL     *string                 `json:"callback_url"`
	CardID          string                  `json:"card_id"`
	Email           string                  `json:"email"`
	SaveUserInfo    bool                    `json:"save_user_info"`
	StoreCode       string                  `json:"store_code"`
}

type Address struct {
	AddressType           string  `json:"address_type"`
	Country               string  `json:"country"`
	IsDefault             bool    `json:"is_default"`
	AdditionalDescription string  `json:"additional_description"`
	Address1              string  `json:"address1"`
	Address2              string  `json:"address2"`
	City                  string  `json:"city"`
	FirstName             string  `json:"first_name"`
	Email                 string  `json:"email"`
	ID                    int     `json:"id"`
	IdentityDocument      string  `json:"identity_document"`
	IdentityDocumentType  string  `json:"identity_document_type"`
	LastName              string  `json:"last_name"`
	Lat                   float64 `json:"lat"`
	Lng                   float64 `json:"lng"`
	Phone                 string  `json:"phone"`
	StateName             string  `json:"state_name"`
	UserID                string  `json:"user_id"`
	Zipcode               string  `json:"zipcode"`
}

type CreditCardInfo struct {
	Address1      *string `json:"address1"`
	CardCVV       *string `json:"card_cvv"`
	CardHolder    *string `json:"card_holder"`
	CardHolderDNI *string `json:"card_holder_dni"`
	City          *string `json:"city"`
	Country       *string `json:"country"`
	ExpiryMonth   *string `json:"expiry_month"`
	ExpiryYear    *string `json:"expiry_year"`
	Phone         *string `json:"phone"`
	State         *string `json:"state"`
	Zip           *string `json:"zip"`
}

type InstallmentInfo struct {
	PlanOptionID string `json:"plan_option_id"`
}

type SpecificFields struct {
	Callbacks Callbacks `json:"callbacks"`
}

type Callbacks struct {
	OnCanceled string `json:"on_canceled"`
	OnFailed   string `json:"on_failed"`
	OnReject   string `json:"on_reject"`
	OnSuccess  string `json:"on_success"`
}
