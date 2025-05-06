package entities

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

type PaymentOrder struct {
	ID                string
	OriginType        value_objects.AssociatedOrigin
	Total             value_objects.CurrencyAmount
	TotalRefunded     value_objects.CurrencyAmount
	Status            enums.PaymentStatus
	Method            value_objects.PaymentMethod
	CreatedAt         time.Time
	UpdatedAt         time.Time
	CollectionAccount CollectionAccount
	AuthorizationCode string
	FailureReason     string
	ReceiptUrl        string
	PaymentFlow       string
	PaymentCard       CardData
}

type CardData struct {
	CardBrand string
	CardLast4 string
	CardType  string
}

func NewPaymentOrder(
	id string,
	originType value_objects.AssociatedOrigin,
	total value_objects.CurrencyAmount,
	method value_objects.PaymentMethod) PaymentOrder {
	return PaymentOrder{
		ID:                id,
		OriginType:        originType,
		Total:             total,
		Status:            enums.PaymentProcessing,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		Method:            method,
		AuthorizationCode: "",
		FailureReason:     "",
		ReceiptUrl:        "",
		PaymentFlow:       enums.Autocapture.String(),
	}
}

func (p PaymentOrder) Validate() error {
	if err := p.OriginType.Validate(); err != nil {
		return err
	}

	if err := p.Total.Validate(); err != nil {
		return err
	}

	if !p.Status.IsValid() {
		return errors.New("status is not valid")
	}

	if err := p.Method.Validate(); err != nil {
		return err
	}

	if p.Total.Value.LessThan(decimal.NewFromInt(0)) {
		return errors.New("total amount is not valid")
	}

	return nil
}

func (p PaymentOrder) CanRefund() bool {
	return (p.Status == enums.PartiallyRefunded || p.Status == enums.PaymentProcessed) &&
		p.GetTotalRefundable().GreaterThan(decimal.Zero)
}

func (p *PaymentOrder) GetTotalRefundable() decimal.Decimal {
	return p.Total.Value.Sub(p.TotalRefunded.Value)
}

func (p *PaymentOrder) GetTotalRefunded() decimal.Decimal {
	return p.TotalRefunded.Value
}

func (p *PaymentOrder) CanRefundAmount(amount decimal.Decimal) bool {
	if p.Status == enums.PaymentProcessed {
		return !p.TotalRefunded.Value.GreaterThan(p.Total.Value.Sub(amount))
	}

	return true
}

func (p *PaymentOrder) SetPaymentCard(paymentCard CardData) {
	p.PaymentCard = paymentCard
}

func (p *PaymentOrder) SetPaymentStatus(status enums.PaymentStatus) {
	p.Status = status
}

func (p *PaymentOrder) SetReceiptUrl(invoiceUrl string) {
	p.ReceiptUrl = invoiceUrl
}

func (p *PaymentOrder) SetPaymentFlow(paymentFlow string) {
	p.PaymentFlow = paymentFlow
}

func (p PaymentOrder) Equals(other PaymentOrder) bool {
	return p.ID == other.ID
}
