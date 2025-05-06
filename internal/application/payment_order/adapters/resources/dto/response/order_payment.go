package response

type DeunaOrderPaymentResponse struct {
	Order      OrderResponse `json:"order"`
	OrderToken string        `json:"order_token"`
}

type OrderResponse struct {
	CashChange              float64           `json:"cash_change"`
	Currency                string            `json:"currency"`
	Discounts               []interface{}     `json:"discounts"`
	DisplayItemsTotalAmount string            `json:"display_items_total_amount"`
	DisplayShippingAmount   string            `json:"display_shipping_amount"`
	DisplaySubTotal         string            `json:"display_sub_total"`
	DisplayTaxAmount        string            `json:"display_tax_amount"`
	DisplayTotalAmount      string            `json:"display_total_amount"`
	DisplayTotalDiscount    string            `json:"display_total_discount"`
	GiftCard                []interface{}     `json:"gift_card"`
	Items                   []OrderItem       `json:"items"`
	ItemsTotalAmount        float64           `json:"items_total_amount"`
	Metadata                map[string]string `json:"metadata"`
	OrderID                 string            `json:"order_id"`
	Payment                 PaymentDetails    `json:"payment"`
	RedirectURL             string            `json:"redirect_url"`
	ScheduledAt             string            `json:"scheduled_at"`
	Shipping                interface{}       `json:"shipping"`
	ShippingAddress         Address           `json:"shipping_address"`
	ShippingAmount          float64           `json:"shipping_amount"`
	ShippingMethod          interface{}       `json:"shipping_method"`
	ShippingMethods         []interface{}     `json:"shipping_methods"`
	ShippingOptions         interface{}       `json:"shipping_options"`
	Status                  string            `json:"status"`
	StoreCode               string            `json:"store_code"`
	SubTotal                float64           `json:"sub_total"`
	TaxAmount               float64           `json:"tax_amount"`
	Timezone                string            `json:"timezone"`
	TotalAmount             float64           `json:"total_amount"`
	TotalDiscount           float64           `json:"total_discount"`
	UserInstructions        string            `json:"user_instructions"`
	WebhookUrls             interface{}       `json:"webhook_urls"`
}

type Address struct {
	AdditionalDescription string  `json:"additional_description"`
	AddressType           string  `json:"address_type"`
	Address1              string  `json:"address1"`
	Address2              string  `json:"address2"`
	City                  string  `json:"city"`
	Country               string  `json:"country"`
	CreatedAt             string  `json:"created_at"`
	FirstName             string  `json:"first_name"`
	ID                    int     `json:"id"`
	IdentityDocument      string  `json:"identity_document"`
	IdentityDocumentType  string  `json:"identity_document_type"`
	IsDefault             bool    `json:"is_default"`
	LastName              string  `json:"last_name"`
	Lat                   float64 `json:"lat"`
	Lng                   float64 `json:"lng"`
	Phone                 string  `json:"phone"`
	StateName             string  `json:"state_name"`
	UpdatedAt             string  `json:"updated_at"`
	UserID                string  `json:"user_id"`
	Zipcode               string  `json:"zipcode"`
}

type OrderItem struct {
	Brand        string        `json:"brand"`
	Category     string        `json:"category"`
	Color        string        `json:"color"`
	Description  string        `json:"description"`
	DetailsURL   string        `json:"details_url"`
	Discounts    []interface{} `json:"discounts"`
	ID           string        `json:"id"`
	ImageURL     string        `json:"image_url"`
	ISBN         string        `json:"isbn"`
	Manufacturer string        `json:"manufacturer"`
	Name         string        `json:"name"`
	Options      string        `json:"options"`
	Quantity     int           `json:"quantity"`
	Size         string        `json:"size"`
	SKU          string        `json:"sku"`
	TaxAmount    MonetaryValue `json:"tax_amount"`
	Taxable      bool          `json:"taxable"`
	TotalAmount  MonetaryValue `json:"total_amount"`
	Type         string        `json:"type"`
	UnitPrice    MonetaryValue `json:"unit_price"`
	UOM          string        `json:"uom"`
	UPC          string        `json:"upc"`
	Weight       WeightDetail  `json:"weight"`
}

type MonetaryValue struct {
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	CurrencySymbol string  `json:"currency_symbol"`
	DisplayAmount  string  `json:"display_amount"`
}

type WeightDetail struct {
	Unit   string  `json:"unit"`
	Weight float64 `json:"weight"`
}

type PaymentDetails struct {
	Data PaymentData `json:"data"`
}

type PaymentData struct {
	Amount            MonetaryValue      `json:"amount"`
	Authorization3DS  Auth3DS            `json:"authorization_3ds"`
	AuthorizationCode string             `json:"authorization_code"`
	CreatedAt         string             `json:"created_at"`
	Customer          CustomerInfo       `json:"customer"`
	FromCard          CardDetails        `json:"from_card"`
	ID                string             `json:"id"`
	Installments      InstallmentDetails `json:"installments"`
	Merchant          MerchantInfo       `json:"merchant"`
	Metadata          map[string]string  `json:"metadata"`
	MethodType        string             `json:"method_type"`
	Processor         string             `json:"processor"`
	Reason            string             `json:"reason"`
	Status            string             `json:"status"`
	UpdatedAt         string             `json:"updated_at"`
}

type Auth3DS struct {
	HTMLContent  string `json:"html_content"`
	URLChallenge string `json:"url_challenge"`
	Version      string `json:"version"`
}

type CustomerInfo struct {
	Email string `json:"email"`
	ID    string `json:"id"`
}

type CardDetails struct {
	CardBrand string `json:"card_brand"`
	FirstSix  string `json:"first_six"`
	LastFour  string `json:"last_four"`
}

type InstallmentDetails struct {
	InstallmentAmount float64 `json:"installment_amount"`
	InstallmentRate   float64 `json:"installment_rate"`
	InstallmentType   string  `json:"installment_type"`
	Installments      int     `json:"installments"`
	PlanID            string  `json:"plan_id"`
	PlanOptionID      string  `json:"plan_option_id"`
}

type MerchantInfo struct {
	ID        string `json:"id"`
	StoreCode string `json:"store_code"`
}
