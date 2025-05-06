package request

import (
	"time"
)

type DeunaOrderType string

const (
	DeUnaCheckout DeunaOrderType = "DEUNA_CHECKOUT"
	DeUnaNow      DeunaOrderType = "DEUNA_NOW"
	PaymentLink   DeunaOrderType = "PAYMENT_LINK"
	AirlineOrder  DeunaOrderType = "AIRLINE_ORDER"
)

type CreateDeunaOrderRequestDTO struct {
	Order     DeunaOrder     `json:"order"`
	OrderType DeunaOrderType `json:"order_type"`
}

func (d DeunaOrderType) String() string {
	return string(d)
}

type DeunaOrder struct {
	BillingAddress        *OrderAddress           `json:"billing_address"`
	Currency              string                  `json:"currency"`
	Items                 []Item                  `json:"items"`
	Metadata              *map[string]interface{} `json:"metadata"`
	PayerInfo             PayerInfo               `json:"payer_info"`
	ShippingAddress       *OrderAddress           `json:"shipping_address"`
	Description           *string                 `json:"description"`
	ExpiresAt             *time.Time              `json:"expires_at"`
	IncludePaymentOptions []PaymentOption         `json:"include_payment_options"`
	OrderID               string                  `json:"order_id"`
	StoreCode             string                  `json:"store_code"`
	SubTotal              int64                   `json:"sub_total"`
	TotalAmount           int64                   `json:"total_amount"`
	Timezone              string                  `json:"timezone"`
	TotalTaxAmount        int64                   `json:"total_tax_amount"`
	TotalDiscount         int64                   `json:"total_discount"`
	ItemsTotalAmount      int                     `json:"items_total_amount"`
	WebhooksURL           WebhooksURL             `json:"webhook_urls"`
}

type WebhooksURL struct {
	NotifyOrder string `json:"notify_order"`
}

type OrderAddress struct {
	ID                    int     `json:"id"`
	AddressType           string  `json:"address_type"`
	Country               string  `json:"country"`
	IsDefault             bool    `json:"is_default"`
	AdditionalDescription string  `json:"additional_description"`
	Address1              string  `json:"address1"`
	Address2              string  `json:"address2"`
	City                  string  `json:"city"`
	Email                 string  `json:"email"`
	FirstName             string  `json:"first_name"`
	LastName              string  `json:"last_name"`
	Lat                   float64 `json:"lat"`
	Lng                   float64 `json:"lng"`
	Phone                 string  `json:"phone"`
	StateName             string  `json:"state_name"`
	UserID                string  `json:"user_id"`
	Zipcode               string  `json:"zipcode"`
}

type Item struct {
	ID          string        `json:"id"`
	TaxAmount   MonetaryValue `json:"tax_amount"`
	Taxable     bool          `json:"taxable"`
	TotalAmount MonetaryValue `json:"total_amount"`
	UnitPrice   MonetaryValue `json:"unit_price"`
	Brand       string        `json:"brand"`
	Category    string        `json:"category"`
	Description string        `json:"description"`
	Name        string        `json:"name"`
	Quantity    int           `json:"quantity"`
}

type MonetaryValue struct {
	Currency       string  `json:"currency"`
	CurrencySymbol string  `json:"currency_symbol"`
	Amount         float64 `json:"amount"`
	OriginalAmount float64 `json:"original_amount,omitempty"`
	TotalDiscount  float64 `json:"total_discount,omitempty"`
}

type PayerInfo struct {
	Email string `json:"email"`
}

type PaymentOption struct {
	PaymentMethod string   `json:"payment_method"`
	Processors    []string `json:"processors"`
}
