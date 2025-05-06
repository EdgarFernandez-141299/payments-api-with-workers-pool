package entities

import (
	"time"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_order/utils"

	"github.com/shopspring/decimal"
	"gitlab.com/clubhub.ai1/gommon/uid"
)

type RefundEntity struct {
	ID           *uid.UniqueID `gorm:"type:varchar(14);primaryKey"`
	PaymentID    string
	Amount       decimal.Decimal
	Reason       string
	Status       string
	EnterpriseID string
	CreatedAt    time.Time
}

func (r RefundEntity) TableName() string {
	return "refund"
}

func (r *RefundEntity) SetStatus(status string) {
	r.Status = status
}

type RefundEntityBuilder struct {
	refund    *RefundEntity
	paymentID string
	orderID   string
}

func NewRefundEntityBuilder() *RefundEntityBuilder {
	return &RefundEntityBuilder{
		refund: &RefundEntity{},
	}
}

func (b *RefundEntityBuilder) WithPaymentID(paymentID string) *RefundEntityBuilder {
	b.paymentID = paymentID
	return b
}

func (b *RefundEntityBuilder) WithOrderID(orderID string) *RefundEntityBuilder {
	b.orderID = orderID
	return b
}

func (b *RefundEntityBuilder) WithAmount(amount decimal.Decimal) *RefundEntityBuilder {
	b.refund.Amount = amount
	return b
}

func (b *RefundEntityBuilder) WithReason(reason string) *RefundEntityBuilder {
	b.refund.Reason = reason
	return b
}

func (b *RefundEntityBuilder) WithStatus(status string) *RefundEntityBuilder {
	b.refund.Status = status
	return b
}

func (b *RefundEntityBuilder) WithEnterpriseID(enterpriseID string) *RefundEntityBuilder {
	b.refund.EnterpriseID = enterpriseID
	return b
}

func (b *RefundEntityBuilder) Build() *RefundEntity {
	b.refund.ID = uid.GenerateID()
	b.refund.PaymentID = utils.GeneratePaymentOrderID(b.orderID, b.paymentID)

	return b.refund
}
