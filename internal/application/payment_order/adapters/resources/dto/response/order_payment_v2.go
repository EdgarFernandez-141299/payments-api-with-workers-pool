package response

type DeunaOrderPaymentResponseV2 struct {
	Order      OrderDetailsResponseV2 `json:"order"`
	OrderToken string                 `json:"order_token"`
}

type OrderDetailsResponseV2 struct {
	Currency                string                  `json:"currency"`
	DisplayItemsTotalAmount string                  `json:"display_items_total_amount"`
	DisplayShippingAmount   string                  `json:"display_shipping_amount"`
	DisplaySubTotal         string                  `json:"display_sub_total"`
	DisplayTaxAmount        string                  `json:"display_tax_amount"`
	DisplayTotalAmount      string                  `json:"display_total_amount"`
	DisplayTotalDiscount    string                  `json:"display_total_discount"`
	Items                   []ItemResponseDetailsV2 `json:"items"`
	ItemsTotalAmount        int                     `json:"items_total_amount"`
	Metadata                map[string]interface{}  `json:"metadata"`
	OrderID                 string                  `json:"order_id"`
	Payment                 PaymentDetailsV2        `json:"payment"`
	Status                  string                  `json:"status"`
	StoreCode               string                  `json:"store_code"`
	SubTotal                int                     `json:"sub_total"`
	TaxAmount               int                     `json:"tax_amount"`
	Timezone                string                  `json:"timezone"`
	TotalAmount             int                     `json:"total_amount"`
	TotalDiscount           int                     `json:"total_discount"`
	WebhookUrls             *WebhookUrlsV2          `json:"webhook_urls"`
}

type ItemResponseDetailsV2 struct {
	Category     string            `json:"category"`
	Description  string            `json:"description"`
	ID           string            `json:"id"`
	Manufacturer string            `json:"manufacturer"`
	Name         string            `json:"name"`
	Options      string            `json:"options"`
	Quantity     int               `json:"quantity"`
	TaxAmount    CurrencyDetailsV2 `json:"tax_amount"`
	Taxable      bool              `json:"taxable"`
	TotalAmount  CurrencyDetailsV2 `json:"total_amount"`
	Type         string            `json:"type"`
	UnitPrice    CurrencyDetailsV2 `json:"unit_price"`
}

type PaymentDetailsV2 struct {
	Data PaymentDataV2 `json:"data"`
}

type WebhookUrlsV2 struct {
	NotifyOrder string `json:"notify_order"`
}

type PaymentDataV2 struct {
	Amount            AmountV2               `json:"amount"`
	AuthorizationCode string                 `json:"authorization_code"`
	CreatedAt         string                 `json:"created_at"`
	Customer          CustomerV2             `json:"customer"`
	FromCard          CardDetailsV2          `json:"from_card"`
	ID                string                 `json:"id"`
	Merchant          MerchantV2             `json:"merchant"`
	Metadata          map[string]interface{} `json:"metadata"`
	MethodType        string                 `json:"method_type"`
	Processor         string                 `json:"processor"`
	Reason            string                 `json:"reason"`
	Status            string                 `json:"status"`
	UpdatedAt         string                 `json:"updated_at"`
}

type AmountV2 struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

type CustomerV2 struct {
	Email string `json:"email"`
	ID    string `json:"id"`
}

type CardDetailsV2 struct {
	CardBrand string `json:"card_brand"`
	FirstSix  string `json:"first_six"`
	LastFour  string `json:"last_four"`
}

type CurrencyDetailsV2 struct {
	Amount         int    `json:"amount"`
	Currency       string `json:"currency"`
	CurrencySymbol string `json:"currency_symbol"`
}

type MerchantV2 struct {
	ID        string `json:"id"`
	StoreCode string `json:"store_code"`
}
