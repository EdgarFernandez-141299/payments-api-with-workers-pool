package response

import (
	"github.com/shopspring/decimal"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
)

type OrderResponseDTO struct {
	ReferenceOrderID string `json:"reference_order_id"`
}

type GetOrderResponseDTO struct {
	ReferenceOrderID string              `json:"reference_order_id"`
	Total            decimal.Decimal     `json:"total"`
	Currency         string              `json:"currency"`
	CountryCode      string              `json:"country_code"`
	Status           enums.PaymentStatus `json:"status"`
	Metadata         map[string]string   `json:"metadata"`
}

type PaymentDTO struct {
	ID                string `json:"id,omitempty"`
	Status            string `json:"status,omitempty"`
	PaymentMethod     string `json:"payment_method,omitempty"`
	CardID            string `json:"card_id,omitempty"`
	AuthorizationCode string `json:"authorization_code,omitempty"`
	PaymentOrderID    string `json:"payment_order_id,omitempty"`
}

type GetOrderPaymentResponseDTO struct {
	ReferenceOrderID string              `json:"reference_order_id"`
	Total            decimal.Decimal     `json:"total"`
	Currency         string              `json:"currency"`
	CountryCode      string              `json:"country_code"`
	Status           enums.PaymentStatus `json:"status"`
	Metadata         map[string]string   `json:"metadata,omitempty"`
	Payments         []PaymentDTO        `json:"payments,omitempty"`
}
