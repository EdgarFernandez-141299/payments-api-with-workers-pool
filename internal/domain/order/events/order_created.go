package events

import (
	"time"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"

	vo "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

type OrderCreated struct {
	ID              string
	TotalAmount     vo.CurrencyAmount
	PhoneNumber     string
	User            entities.User
	CreatedAt       time.Time
	CountryCode     vo.Country
	BillingAddress  vo.Address
	ShippingAddress vo.Address
	EnterpriseID    string
	Metadata        map[string]interface{}
	Email           string
	WebhookUrl      vo.WebhookUrl
	AllowCapture    bool
}

type OrderCreatedEventBuilder struct {
	event *OrderCreated
}

func NewOrderCreatedEventBuilder() *OrderCreatedEventBuilder {
	return &OrderCreatedEventBuilder{event: &OrderCreated{}}
}

func (b *OrderCreatedEventBuilder) SetID(id string) *OrderCreatedEventBuilder {
	b.event.ID = id
	return b
}

func (b *OrderCreatedEventBuilder) SetTotalAmount(totalAmount vo.CurrencyAmount) *OrderCreatedEventBuilder {
	b.event.TotalAmount = totalAmount
	return b
}

func (b *OrderCreatedEventBuilder) SetPhoneNumber(phoneNumber string) *OrderCreatedEventBuilder {
	b.event.PhoneNumber = phoneNumber
	return b
}

func (b *OrderCreatedEventBuilder) SetUser(user entities.User) *OrderCreatedEventBuilder {
	b.event.User = user
	return b
}

func (b *OrderCreatedEventBuilder) SetCreatedAt(createdAt time.Time) *OrderCreatedEventBuilder {
	b.event.CreatedAt = createdAt
	return b
}

func (b *OrderCreatedEventBuilder) SetCountryCode(countryCode vo.Country) *OrderCreatedEventBuilder {
	b.event.CountryCode = countryCode
	return b
}

func (b *OrderCreatedEventBuilder) SetBillingAddress(billingAddress vo.Address) *OrderCreatedEventBuilder {
	b.event.BillingAddress = billingAddress
	return b
}

func (b *OrderCreatedEventBuilder) SetShippingAddress(shippingAddress vo.Address) *OrderCreatedEventBuilder {
	b.event.ShippingAddress = shippingAddress
	return b
}

func (b *OrderCreatedEventBuilder) SetEnterpriseID(enterpriseID string) *OrderCreatedEventBuilder {
	b.event.EnterpriseID = enterpriseID
	return b
}

func (b *OrderCreatedEventBuilder) SetEmail(email string) *OrderCreatedEventBuilder {
	b.event.Email = email
	return b
}

func (b *OrderCreatedEventBuilder) SetWebhookUrl(webhookUrl vo.WebhookUrl) *OrderCreatedEventBuilder {
	b.event.WebhookUrl = webhookUrl
	return b
}

func (b *OrderCreatedEventBuilder) SetMetadata(metadata map[string]interface{}) *OrderCreatedEventBuilder {
	b.event.Metadata = metadata
	return b
}

func (b *OrderCreatedEventBuilder) SetAllowCapture(allowCapture bool) *OrderCreatedEventBuilder {
	b.event.AllowCapture = allowCapture
	return b
}

func (b *OrderCreatedEventBuilder) Build() *OrderCreated {
	return b.event
}
