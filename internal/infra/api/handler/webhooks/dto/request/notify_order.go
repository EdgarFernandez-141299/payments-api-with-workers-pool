package dto

type WebhookOrderDTO struct {
	Order Order `json:"order"`
}

// Order contiene toda la información del pedido
type Order struct {
	Token            string                 `json:"token"`
	MerchantID       string                 `json:"merchant_id"`
	PaymentMethod    string                 `json:"payment_method"`
	PaymentStatus    string                 `json:"payment_status"`
	Currency         string                 `json:"currency"`
	TaxAmount        int                    `json:"tax_amount"`
	ItemsTotalAmount int                    `json:"items_total_amount"`
	SubTotal         int                    `json:"sub_total"`
	TotalAmount      int                    `json:"total_amount"`
	OrderID          string                 `json:"order_id"`
	TransactionID    string                 `json:"transaction_id"`
	Metadata         map[string]interface{} `json:"metadata"`
	Payment          Payment                `json:"payment"`
	Status           string                 `json:"status"`
	UserID           string                 `json:"user_id"`
	CashChange       int                    `json:"cash_change"`
}

// Metadata contiene metadatos adicionales del pedido
type Metadata struct {
	AdditionalDataTest string `json:"additional_data.test"`
}

// Payment contiene información del pago
type Payment struct {
	Data PaymentData `json:"data"`
}

// PaymentData contiene los datos detallados del pago
type PaymentData struct {
	Metadata              map[string]interface{} `json:"metadata"`
	FromCard              CardInfo               `json:"from_card"`
	Amount                MoneyAmount            `json:"amount"`
	UpdatedAt             string                 `json:"updated_at"`
	MethodType            string                 `json:"method_type"`
	CreatedAt             string                 `json:"created_at"`
	Merchant              PaymentMerchant        `json:"merchant"`
	ID                    string                 `json:"id"`
	Processor             string                 `json:"processor"`
	Customer              PaymentCustomer        `json:"customer"`
	Status                string                 `json:"status"`
	Reason                string                 `json:"reason"`
	ExternalTransactionID string                 `json:"external_transaction_id"`
	Installments          int                    `json:"installments"`
	AuthenticationMethod  string                 `json:"authentication_method"`
	ManualStatus          string                 `json:"manual_status"`
	AuthorizationCode     string                 `json:"authorization_code"`
}

// CardInfo contiene información de la tarjeta de crédito
type CardInfo struct {
	CardBrand  string `json:"card_brand"`
	FirstSix   string `json:"first_six"`
	LastFour   string `json:"last_four"`
	CardHolder string `json:"card_holder"`
}

// MoneyAmount representa una cantidad monetaria con divisa
type MoneyAmount struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

// PaymentMerchant contiene información del comerciante para el pago
type PaymentMerchant struct {
	StoreCode string `json:"store_code"`
	ID        string `json:"id"`
}

// PaymentCustomer contiene información del cliente para el pago
type PaymentCustomer struct {
	Email string `json:"email"`
	ID    string `json:"id"`
}
