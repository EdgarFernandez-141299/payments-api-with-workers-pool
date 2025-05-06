package fixture

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/aggregate"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/event_store"
	"testing"
)

type OrderEventRepositoryFixtureBuilder struct {
	t                     *testing.T
	m                     *event_store.OrderEventRepository
	referenceOrderID      string
	paymentOrderID        string
	paymentStatus         enums.PaymentStatus
	orderCurrencyAmount   value_objects.CurrencyAmount
	paymentCurrencyAmount value_objects.CurrencyAmount
	webhookUrl            value_objects.WebhookUrl
	err                   error
}

func NewOrderEventRepositoryFixtureBuilder(t *testing.T) *OrderEventRepositoryFixtureBuilder {
	return &OrderEventRepositoryFixtureBuilder{t: t}
}

func NewFromOrderMock(m *event_store.OrderEventRepository) *OrderEventRepositoryFixtureBuilder {
	return &OrderEventRepositoryFixtureBuilder{m: m}
}

func (b *OrderEventRepositoryFixtureBuilder) WithReferenceOrderID(id string) *OrderEventRepositoryFixtureBuilder {
	b.referenceOrderID = id
	return b
}

func (b *OrderEventRepositoryFixtureBuilder) WithPaymentOrderID(id string) *OrderEventRepositoryFixtureBuilder {
	b.paymentOrderID = id
	return b
}

func (b *OrderEventRepositoryFixtureBuilder) WithPaymentStatus(status enums.PaymentStatus) *OrderEventRepositoryFixtureBuilder {
	b.paymentStatus = status
	return b
}

func (b *OrderEventRepositoryFixtureBuilder) WithOrderCurrencyAmount(amount value_objects.CurrencyAmount) *OrderEventRepositoryFixtureBuilder {
	b.orderCurrencyAmount = amount
	return b
}

func (b *OrderEventRepositoryFixtureBuilder) WithWebhookURL(url value_objects.WebhookUrl) *OrderEventRepositoryFixtureBuilder {
	b.webhookUrl = url
	return b
}

func (b *OrderEventRepositoryFixtureBuilder) WithPaymentCurrencyAmount(amount value_objects.CurrencyAmount) *OrderEventRepositoryFixtureBuilder {
	b.paymentCurrencyAmount = amount
	return b
}

func (b *OrderEventRepositoryFixtureBuilder) WithError(err error) *OrderEventRepositoryFixtureBuilder {
	b.err = err
	return b
}

func (b *OrderEventRepositoryFixtureBuilder) Build() *event_store.OrderEventRepository {
	var mockRepo *event_store.OrderEventRepository

	if b.m != nil {
		mockRepo = b.m
	} else {
		mockRepo = event_store.NewOrderEventRepository(b.t)
	}

	mockRepo.On("Get", mock.Anything, b.referenceOrderID, mock.Anything).Once().
		Run(func(args mock.Arguments) {
			o := args.Get(2).(*aggregate.Order)
			o.ID = b.referenceOrderID
			o.TotalAmount = b.orderCurrencyAmount
			o.WebhookUrl = b.webhookUrl
			o.OrderPayments = []entities.PaymentOrder{{
				ID:     b.paymentOrderID,
				Status: b.paymentStatus,
				Total:  b.paymentCurrencyAmount,
			}}
		}).
		Return(b.err)

	return mockRepo
}

func OrderEventRepositoryGetFixture(
	t *testing.T,
	referenceOrderID string,
	paymentOrderID string,
	paymentStatus enums.PaymentStatus,
	orderCurrencyAmount, paymentCurrencyAmount value_objects.CurrencyAmount,
	err error) *event_store.OrderEventRepository {

	return NewOrderEventRepositoryFixtureBuilder(t).
		WithReferenceOrderID(referenceOrderID).
		WithPaymentOrderID(paymentOrderID).
		WithPaymentStatus(paymentStatus).
		WithOrderCurrencyAmount(orderCurrencyAmount).
		WithPaymentCurrencyAmount(paymentCurrencyAmount).
		WithError(err).
		Build()
}
