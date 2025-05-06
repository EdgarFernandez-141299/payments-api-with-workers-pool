package request

type DeunaOrderPaymentRequestV2 struct {
	Order         OrderDetailsV2  `json:"order"`
	OrderType     string          `json:"order_type"`
	PayerInfo     PayerInfoV2     `json:"payer_info"`
	PaymentSource PaymentSourceV2 `json:"payment_source"`
}

type OrderDetailsV2 struct {
	BillingAddress AddressV2              `json:"billing_address"`
	Currency       string                 `json:"currency"`
	Items          []ItemDetailsV2        `json:"items"`
	Metadata       map[string]interface{} `json:"metadata"`
	WebhookUrls    WebhookUrlsV2          `json:"webhook_urls"`
	ExpiresAt      string                 `json:"expires_at"`
	GroupID        string                 `json:"group_id"`
	OrderID        string                 `json:"order_id"`
	StoreCode      string                 `json:"store_code"`
	SubTotal       int                    `json:"sub_total"`
	TotalAmount    int                    `json:"total_amount"`
	TotalTaxAmount int                    `json:"total_tax_amount"`
}

type AddressV2 struct {
	ID                    int    `json:"id"`
	AddressType           string `json:"address_type"`
	Country               string `json:"country"`
	IsDefault             bool   `json:"is_default"`
	AdditionalDescription string `json:"additional_description"`
	Address1              string `json:"address1"`
	Address2              string `json:"address2"`
	City                  string `json:"city"`
	FirstName             string `json:"first_name"`
	Email                 string `json:"email"`
	LastName              string `json:"last_name"`
	Phone                 string `json:"phone"`
	StateName             string `json:"state_name"`
	UserID                string `json:"user_id"`
	Zipcode               string `json:"zipcode"`
}

type ItemDetailsV2 struct {
	ID          string            `json:"id"`
	TaxAmount   CurrencyDetailsV2 `json:"tax_amount"`
	Taxable     bool              `json:"taxable"`
	TotalAmount TotalAmountV2     `json:"total_amount"`
	UnitPrice   UnitPriceV2       `json:"unit_price"`
	Category    string            `json:"category"`
	Description string            `json:"description"`
	Name        string            `json:"name"`
	Quantity    int               `json:"quantity"`
}

type CurrencyDetailsV2 struct {
	Amount         int    `json:"amount"`
	Currency       string `json:"currency"`
	CurrencySymbol string `json:"currency_symbol"`
}

type TotalAmountV2 struct {
	Currency       string `json:"currency"`
	CurrencySymbol string `json:"currency_symbol"`
	Amount         int    `json:"amount"`
	OriginalAmount int    `json:"original_amount"`
	TotalDiscount  int    `json:"total_discount"`
}

type UnitPriceV2 struct {
	CurrencySymbol string `json:"currency_symbol"`
	Amount         int    `json:"amount"`
	Currency       string `json:"currency"`
}

type WebhookUrlsV2 struct {
	NotifyOrder string `json:"notify_order"`
}

type PayerInfoV2 struct {
	Email          string `json:"email"`
	ExternalUserID string `json:"external_user_id"`
	SaveUserInfo   bool   `json:"save_user_info"`
}

type PaymentSourceV2 struct {
	CardInfo   CardInfoV2 `json:"card_info"`
	MethodType string     `json:"method_type"`
	Processor  string     `json:"processor"`
}

type CardInfoV2 struct {
	Address1                  string `json:"address1"`
	CardCVV                   string `json:"card_cvv"`
	CardHolder                string `json:"card_holder"`
	CardHolderDNI             string `json:"card_holder_dni"`
	CardID                    string `json:"card_id"`
	CardNumber                string `json:"card_number"`
	City                      string `json:"city"`
	Country                   string `json:"country"`
	DeferredInstallmentMonths int    `json:"deferred_installment_months"`
	ExpiryMonth               string `json:"expiry_month"`
	ExpiryYear                string `json:"expiry_year"`
	Installments              int    `json:"installments"`
	Phone                     string `json:"phone"`
	State                     string `json:"state"`
	Zip                       string `json:"zip"`
}
