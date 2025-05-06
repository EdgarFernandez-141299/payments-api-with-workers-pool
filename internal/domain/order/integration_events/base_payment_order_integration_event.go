package integration_events

import "errors"

var integrationEventParamError = errors.New("integration event param error")

type baseOrderIntegrationEvent struct {
	Type               string                 `json:"type"`
	ReferenceOrderID   string                 `json:"reference_order_id"`
	PaymentStatus      string                 `json:"payment_status"`
	OrderStatus        string                 `json:"order_status"`
	ReferencePaymentID string                 `json:"reference_payment_id"`
	AssociatedPayment  string                 `json:"associated_payment"`
	TotalOrderAmount   float64                `json:"total_order_amount"`
	Currency           string                 `json:"currency"`
	UserID             string                 `json:"user_id"`
	UserType           string                 `json:"user_type"`
	EnterpriseID       string                 `json:"enterprise_id"`
	TotalOrderPaid     float64                `json:"total_order_paid"`
	TotalPaymentAmount float64                `json:"total_payment_amount"`
	CardData           CardData               `json:"card_data"`
	Metadata           map[string]interface{} `json:"metadata"`
	PaymentFlow        string                 `json:"payment_flow"`
	ReceiptUrl         string                 `json:"invoice_url"`
}

type CardData struct {
	CardNumber    string `json:"card_number"`
	CardType      string `json:"card_type"`
	MethodPayment string `json:"method_payment"`
}

type IntegrationEventsParams struct {
	ReferenceOrderID   string
	ReferencePaymentID string
	AssociatedPayment  string
	TotalOrderAmount   float64
	Currency           string
	UserID             string
	UserType           string
	EnterpriseID       string
	TotalOrderPaid     float64
	TotalPaymentAmount float64
	CardData           CardData
	Metadata           map[string]interface{}
	ReceiptUrl         string
	PaymentFlow        string
}

func PaymentOrderIntegrationEventWithStatus(
	params IntegrationEventsParams,
	eventType string,
	paymentStatus string,
	orderStatus string) baseOrderIntegrationEvent {
	if params.Metadata == nil {
		params.Metadata = make(map[string]interface{})
	}

	return baseOrderIntegrationEvent{
		ReferenceOrderID:   params.ReferenceOrderID,
		PaymentStatus:      paymentStatus,
		OrderStatus:        orderStatus,
		Type:               eventType,
		ReferencePaymentID: params.ReferencePaymentID,
		AssociatedPayment:  params.AssociatedPayment,
		TotalOrderAmount:   params.TotalOrderAmount,
		Currency:           params.Currency,
		UserID:             params.UserID,
		UserType:           params.UserType,
		EnterpriseID:       params.EnterpriseID,
		TotalOrderPaid:     params.TotalOrderPaid,
		TotalPaymentAmount: params.TotalPaymentAmount,
		Metadata:           params.Metadata,
		CardData:           params.CardData,
		ReceiptUrl:         params.ReceiptUrl,
		PaymentFlow:        params.PaymentFlow,
	}
}

func (p baseOrderIntegrationEvent) Validate() error {
	if p.ReferenceOrderID == "" {
		return integrationEventParamError
	}
	if p.PaymentStatus == "" {
		return integrationEventParamError
	}
	if p.OrderStatus == "" {
		return integrationEventParamError
	}
	if p.ReferencePaymentID == "" {
		return integrationEventParamError
	}
	if p.AssociatedPayment == "" {
		return integrationEventParamError
	}
	if p.Currency == "" {
		return integrationEventParamError
	}
	if p.UserID == "" {
		return integrationEventParamError
	}
	if p.UserType == "" {
		return integrationEventParamError
	}
	if p.EnterpriseID == "" {
		return integrationEventParamError
	}

	return nil
}
