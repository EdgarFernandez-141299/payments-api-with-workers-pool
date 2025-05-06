package entities

import (
	"encoding/json"
	"time"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_order/utils"

	"github.com/shopspring/decimal"
	"gitlab.com/clubhub.ai1/gommon/uid"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gorm.io/gorm"
)

type PaymentOrderEntity struct {
	ID                  *uid.UniqueID `gorm:"type:varchar(14);primaryKey"`
	OrderID             string
	AssociatedOrigin    string
	CurrencyCode        string
	CountryCode         string
	CardID              string
	CardDetail          json.RawMessage
	PaymentMethod       string
	CollectionAccountID string
	Metadata            json.RawMessage
	Status              string
	TotalAmount         decimal.Decimal
	Reference           string
	FailureReason       string
	FailureCode         string
	EnterpriseID        string
	IPAddress           string
	DeviceFingerprint   string
	TransactionDate     time.Time
	PaymentOrderID      string
	PaymentFlow         string
	AuthorizedAt        time.Time
	CapturedAt          time.Time
	ReleasedAt          time.Time
	AuthorizationCode   string
	gorm.Model
}

func (p PaymentOrderEntity) TableName() string {
	return "payment"
}

func (p *PaymentOrderEntity) SetStatus(status enums.PaymentStatus) {
	p.Status = status.String()
}

func (p *PaymentOrderEntity) SetAuthorizationCode(authorizationCode string) {
	p.AuthorizationCode = authorizationCode
}

type PaymentOrderEntityBuilder struct {
	entity PaymentOrderEntity
}

func NewPaymentOrderEntity() *PaymentOrderEntityBuilder {
	return &PaymentOrderEntityBuilder{}
}

func (p *PaymentOrderEntityBuilder) SetOrderID(orderID string) *PaymentOrderEntityBuilder {
	p.entity.OrderID = orderID
	return p
}

func (p *PaymentOrderEntityBuilder) SetAssociatedOrigin(associatedOrigin string) *PaymentOrderEntityBuilder {
	p.entity.AssociatedOrigin = associatedOrigin
	return p
}

func (p *PaymentOrderEntityBuilder) SetCurrencyCode(currencyCode string) *PaymentOrderEntityBuilder {
	p.entity.CurrencyCode = currencyCode
	return p
}

func (p *PaymentOrderEntityBuilder) SetCountryCode(countryCode string) *PaymentOrderEntityBuilder {
	p.entity.CountryCode = countryCode
	return p
}

func (p *PaymentOrderEntityBuilder) SetCardID(cardID string) *PaymentOrderEntityBuilder {
	p.entity.CardID = cardID
	return p
}

func (p *PaymentOrderEntityBuilder) SetPaymentMethod(paymentMethod string) *PaymentOrderEntityBuilder {
	p.entity.PaymentMethod = paymentMethod
	return p
}

func (p *PaymentOrderEntityBuilder) SetCollectionAccountID(collectionAccountID string) *PaymentOrderEntityBuilder {
	p.entity.CollectionAccountID = collectionAccountID
	return p
}

func (p *PaymentOrderEntityBuilder) SetCardDetail(cardDetail string) *PaymentOrderEntityBuilder {
	p.entity.CardDetail = json.RawMessage(cardDetail)
	return p
}

func (p *PaymentOrderEntityBuilder) SetMetadata(metadata string) *PaymentOrderEntityBuilder {
	p.entity.Metadata = json.RawMessage(metadata)
	return p
}

func (p *PaymentOrderEntityBuilder) SetStatus(status string) *PaymentOrderEntityBuilder {
	p.entity.Status = status
	return p
}

func (p *PaymentOrderEntityBuilder) SetTotalAmount(totalAmount decimal.Decimal) *PaymentOrderEntityBuilder {
	p.entity.TotalAmount = totalAmount
	return p
}

func (p *PaymentOrderEntityBuilder) SetReference(reference string) *PaymentOrderEntityBuilder {
	p.entity.Reference = reference
	return p
}

func (p *PaymentOrderEntityBuilder) SetFailureReason(failureReason string) *PaymentOrderEntityBuilder {
	p.entity.FailureReason = failureReason
	return p
}

func (p *PaymentOrderEntityBuilder) SetFailureCode(failureCode string) *PaymentOrderEntityBuilder {
	p.entity.FailureCode = failureCode
	return p
}

func (p *PaymentOrderEntityBuilder) SetEnterpriseID(enterpriseID string) *PaymentOrderEntityBuilder {
	p.entity.EnterpriseID = enterpriseID
	return p
}

func (p *PaymentOrderEntityBuilder) SetIPAddress(ipAddress string) *PaymentOrderEntityBuilder {
	p.entity.IPAddress = ipAddress
	return p
}

func (p *PaymentOrderEntityBuilder) SetDeviceFingerprint(deviceFingerprint string) *PaymentOrderEntityBuilder {
	p.entity.DeviceFingerprint = deviceFingerprint
	return p
}

func (p *PaymentOrderEntityBuilder) SetTransactionDate() *PaymentOrderEntityBuilder {
	p.entity.TransactionDate = time.Now()
	return p
}

func (p *PaymentOrderEntityBuilder) SetPaymentOrderID(paymentOrderID string) *PaymentOrderEntityBuilder {
	p.entity.PaymentOrderID = paymentOrderID
	return p
}

func (p *PaymentOrderEntityBuilder) SetPaymentFlow(paymentFlow enums.PaymentFlowEnum) *PaymentOrderEntityBuilder {
	p.entity.PaymentFlow = paymentFlow.String()
	return p
}

func (p *PaymentOrderEntityBuilder) SetAuthorizedAt(authorizedAt time.Time) *PaymentOrderEntityBuilder {
	p.entity.AuthorizedAt = authorizedAt
	return p
}

func (p *PaymentOrderEntityBuilder) SetCapturedAt(capturedAt time.Time) *PaymentOrderEntityBuilder {
	p.entity.CapturedAt = capturedAt
	return p
}

func (p *PaymentOrderEntityBuilder) SetReleasedAt(releasedAt time.Time) *PaymentOrderEntityBuilder {
	p.entity.ReleasedAt = releasedAt
	return p
}

func (p *PaymentOrderEntityBuilder) SetAuthorizationCode(authorizationCode string) *PaymentOrderEntityBuilder {
	p.entity.AuthorizationCode = authorizationCode
	return p
}

func (p *PaymentOrderEntityBuilder) Build() PaymentOrderEntity {
	id, _ := uid.NewUniqueID(
		uid.WithID(utils.GeneratePaymentOrderID(p.entity.OrderID, p.entity.PaymentOrderID)),
	)
	p.entity.ID = id

	return p.entity
}
