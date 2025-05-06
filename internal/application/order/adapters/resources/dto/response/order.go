package response

type DeunaOrderResponseDTO struct {
	Token string `json:"token"`
}

type DeunaOrderDetails struct {
	BillingAddress           Address                `json:"billing_address"`
	AirlineInformation       *AirlineInformation    `json:"airline_information"`
	CashChange               float64                `json:"cash_change"`
	CreatedAt                string                 `json:"created_at"`
	Currency                 string                 `json:"currency"`
	DiscountAmount           float64                `json:"discount_amount"`
	Discounts                []Discount             `json:"discounts"`
	DisplayItemsTotalAmount  string                 `json:"display_items_total_amount"`
	DisplayShippingAmount    string                 `json:"display_shipping_amount"`
	DisplayShippingTaxAmount string                 `json:"display_shipping_tax_amount"`
	DisplaySubTotalAmount    string                 `json:"display_sub_total_amount"`
	DisplayTaxAmount         string                 `json:"display_tax_amount"`
	DisplayTotalAmount       string                 `json:"display_total_amount"`
	DisplayTotalDiscount     string                 `json:"display_total_discount"`
	DisplayTotalTaxAmount    string                 `json:"display_total_tax_amount"`
	GiftCard                 []interface{}          `json:"gift_card"`
	IncludePaymentOptions    []interface{}          `json:"include_payment_options"`
	Items                    []Item                 `json:"items"`
	ItemsTotalAmount         float64                `json:"items_total_amount"`
	Metadata                 map[string]interface{} `json:"metadata"`
	OrderID                  string                 `json:"order_id"`
	PayerInfo                *interface{}           `json:"payer_info"`
	Payment                  PaymentDetails         `json:"payment"`
	PaymentLink              string                 `json:"payment_link"`
	RedirectURL              string                 `json:"redirect_url"`
	RedirectUrls             RedirectURLs           `json:"redirect_urls"`
	ScheduledAt              string                 `json:"scheduled_at"`
	Shipping                 ShippingInfo           `json:"shipping"`
	ShippingAddress          Address                `json:"shipping_address"`
	ShippingAmount           float64                `json:"shipping_amount"`
	ShippingDiscountAmount   float64                `json:"shipping_discount_amount"`
	ShippingMethod           ShippingMethod         `json:"shipping_method"`
	ShippingMethods          []ShippingMethod       `json:"shipping_methods"`
	ShippingOptions          ShippingOptions        `json:"shipping_options"`
	ShippingTaxAmount        float64                `json:"shipping_tax_amount"`
	Status                   string                 `json:"status"`
	StoreCode                string                 `json:"store_code"`
	SubTotal                 float64                `json:"sub_total"`
	TaxAmount                float64                `json:"tax_amount"`
	Timezone                 string                 `json:"timezone"`
	TotalAmount              float64                `json:"total_amount"`
	TotalDiscount            float64                `json:"total_discount"`
	TotalTaxAmount           float64                `json:"total_tax_amount"`
	UpdatedAt                string                 `json:"updated_at"`
	UserID                   string                 `json:"user_id"`
	UserInstructions         string                 `json:"user_instructions"`
	WebhookURLs              WebhookURLs            `json:"webhook_urls"`
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

type Discount struct {
	Amount           float64      `json:"amount"`
	Code             string       `json:"code"`
	Description      string       `json:"description"`
	DetailsURL       string       `json:"details_url"`
	DiscountCategory string       `json:"discount_category"`
	DisplayAmount    string       `json:"display_amount"`
	FreeShipping     FreeShipping `json:"free_shipping"`
	Reference        string       `json:"reference"`
	TargetType       string       `json:"target_type"`
	Type             string       `json:"type"`
}

type FreeShipping struct {
	IsFreeShipping     bool    `json:"is_free_shipping"`
	MaximumCostAllowed float64 `json:"maximum_cost_allowed"`
}

type Item struct {
	Brand        string       `json:"brand"`
	Category     string       `json:"category"`
	Color        string       `json:"color"`
	Description  string       `json:"description"`
	DetailsURL   string       `json:"details_url"`
	Discounts    []Discount   `json:"discounts"`
	ID           string       `json:"id"`
	ImageURL     string       `json:"image_url"`
	ISBN         string       `json:"isbn"`
	Manufacturer string       `json:"manufacturer"`
	Name         string       `json:"name"`
	Options      string       `json:"options"`
	Quantity     int          `json:"quantity"`
	Size         string       `json:"size"`
	Sku          string       `json:"sku"`
	TaxAmount    Money        `json:"tax_amount"`
	Taxable      bool         `json:"taxable"`
	TotalAmount  Money        `json:"total_amount"`
	Type         string       `json:"type"`
	UnitPrice    Money        `json:"unit_price"`
	Uom          string       `json:"uom"`
	Upc          string       `json:"upc"`
	Weight       WeightDetail `json:"weight"`
}

type Money struct {
	Amount                float64 `json:"amount"`
	Currency              string  `json:"currency"`
	CurrencySymbol        string  `json:"currency_symbol"`
	DisplayAmount         string  `json:"display_amount"`
	DisplayOriginalAmount string  `json:"display_original_amount"`
	DisplayTotalDiscount  string  `json:"display_total_discount"`
	OriginalAmount        float64 `json:"original_amount"`
	TotalDiscount         float64 `json:"total_discount"`
}

type WeightDetail struct {
	Unit   string  `json:"unit"`
	Weight float64 `json:"weight"`
}

type PaymentDetails struct {
	Data PaymentData `json:"data"`
}

type PaymentData struct {
	Amount          Money                  `json:"amount"`
	CreatedAt       string                 `json:"created_at"`
	Customer        CustomerInfo           `json:"customer"`
	FromCard        CardInfo               `json:"from_card"`
	ID              string                 `json:"id"`
	Installments    Installments           `json:"installments"`
	Merchant        MerchantInfo           `json:"merchant"`
	Metadata        map[string]interface{} `json:"metadata"`
	MethodType      string                 `json:"method_type"`
	Processor       string                 `json:"processor"`
	Reason          string                 `json:"reason"`
	RoutingStrategy string                 `json:"routing_strategy"`
	Status          string                 `json:"status"`
	UpdatedAt       string                 `json:"updated_at"`
}

type CustomerInfo struct {
	Email string `json:"email"`
	ID    string `json:"id"`
}

type CardInfo struct {
	CardBrand string `json:"card_brand"`
	FirstSix  string `json:"first_six"`
	LastFour  string `json:"last_four"`
}

type Installments struct {
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

type RedirectURLs struct {
	Close    string `json:"close"`
	Error    string `json:"error"`
	Fallback string `json:"fallback"`
	Pending  string `json:"pending"`
	Success  string `json:"success"`
}

type ShippingInfo struct {
	Discounts      []Discount `json:"discounts"`
	OriginalAmount float64    `json:"original_amount"`
	TotalDiscount  float64    `json:"total_discount"`
}

type ShippingMethod struct {
	Code             string        `json:"code"`
	Cost             float64       `json:"cost"`
	DisplayCost      string        `json:"display_cost"`
	DisplayTaxAmount string        `json:"display_tax_amount"`
	MaxDeliveryDate  string        `json:"max_delivery_date"`
	MinDeliveryDate  string        `json:"min_delivery_date"`
	Name             string        `json:"name"`
	Scheduler        []interface{} `json:"scheduler"`
	TaxAmount        float64       `json:"tax_amount"`
}

type ShippingOptions struct {
	Details ShippingDetails `json:"details"`
	Type    string          `json:"type"`
}

type ShippingDetails struct {
	AdditionalDetails  AdditionalDetails  `json:"additional_details"`
	Address            string             `json:"address"`
	AddressCoordinates AddressCoordinates `json:"address_coordinates"`
	Contact            ContactInfo        `json:"contact"`
	StoreName          string             `json:"store_name"`
}

type AdditionalDetails struct {
	PickupTime    string `json:"pickup_time"`
	StockLocation string `json:"stock_location"`
}

type AddressCoordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type ContactInfo struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type WebhookURLs struct {
	ApplyCoupon          string `json:"apply_coupon"`
	GetShippingMethods   string `json:"get_shipping_methods"`
	NotifyOrder          string `json:"notify_order"`
	RemoveCoupon         string `json:"remove_coupon"`
	UpdateShippingMethod string `json:"update_shipping_method"`
}
