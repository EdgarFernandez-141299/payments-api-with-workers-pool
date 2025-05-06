package aggregate

import (
	"errors"
	"time"

	errors2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"

	"gitlab.com/clubhub.ai1/go-libraries/eventsourcing"
	aggregate "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/common/eventsourcing"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/events"
	vo "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

type Order struct {
	aggregate.AggregateRoot
	ID            string
	Currency      vo.CurrencyCode
	TotalAmount   vo.CurrencyAmount
	PhoneNumber   string
	User          entities.User
	Status        vo.OrderStatus
	CountryCode   vo.Country
	CreatedAt     time.Time
	OrderPayments []entities.PaymentOrder
	EnterpriseID  string
	Email         string
	Metadata      map[string]interface{}
	WebhookUrl    vo.WebhookUrl
	AllowCapture  bool
}

func (o *Order) StartProcessingOrderPayment(cmd command.CreatePaymentOrderCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	if !o.CanProcessPayment(cmd.Payment) {
		return errors.New("payment order cannot be added")
	}

	event := events.FromProcessOrderCommand(cmd)

	o.TrackChange(o, event)

	return nil
}

func (o *Order) OrderPaymentProcessed(cmd command.CreatePaymentOrderProcessedCommand) error {
	event := events.FromProcessedOrderCommand(cmd)

	o.TrackChange(o, event)

	return nil
}

func (o *Order) OrderPaymentAuthorized(cmd command.CreatePaymentOrderAuthorizedCommand) error {
	_, err := o.FindPaymentByID(cmd.PaymentID)
	if err != nil {
		return err
	}

	event := events.FromAuthorizedOrderCommand(cmd)
	o.TrackChange(o, event)

	return nil
}

func (o *Order) OrderPaymentFailed(cmd command.CreatePaymentOrderFailCommand) error {
	event := events.FromFailedOrderCommand(cmd)

	o.TrackChange(o, event)

	return nil
}

func (o *Order) OrderPaymentReleased(cmd command.PaymentOrderReleasedCommand) error {
	_, err := o.FindPaymentByID(cmd.PaymentID)
	if err != nil {
		return err
	}

	event := events.FromReleasedOrderCommand(cmd)
	o.TrackChange(o, event)

	return nil
}

func (o *Order) OrderPaymentCaptured(cmd command.PaymentOrderCapturedCommand) error {
	_, err := o.FindPaymentByID(cmd.PaymentID)
	if err != nil {
		return err
	}

	event := events.FromCapturedOrderCommand(cmd)
	o.TrackChange(o, event)

	return nil
}

func (o *Order) CanProcessPayment(other entities.PaymentOrder) bool {
	if o.Status != vo.OrderStatusProcessing() {
		return false
	}
	_, found := lo.Find(o.OrderPayments, func(item entities.PaymentOrder) bool {
		return item.ID == other.ID && item.Status != enums.PaymentFailed
	})

	if found {
		return false
	}

	sum := o.GetTotalEligiblePayments().Value.Add(other.Total.Value)

	return sum.LessThanOrEqual(o.TotalAmount.Value)
}

func (o *Order) RefundPayment(paymentID string, reason string) (entities.PaymentOrder, error) {
	payment, findErr := o.FindPaymentByID(paymentID)

	if findErr != nil {
		return payment, findErr
	}

	if !payment.CanRefund() {
		return payment, errors2.NewRefundCanNotAppliedDueToPaymentStatusError()
	}

	event := events.FromRefundOrderCommand(paymentID, reason)
	o.TrackChange(o, event)

	return payment, nil
}

func (o *Order) RefundPartialPayment(paymentID string, reason string, amount decimal.Decimal) error {
	payment, findErr := o.FindPaymentByID(paymentID)

	if findErr != nil {
		return findErr
	}

	if !payment.CanRefundAmount(amount) {
		return errors2.NewRefundCanNotAppliedDueToPaymentStatusError()
	}

	event := events.FromRefundPartialPaymentCommand(paymentID, reason, amount)

	o.TrackChange(o, event)

	return nil
}

func (o *Order) FindPaymentByID(paymentID string) (entities.PaymentOrder, error) {
	for _, payment := range o.OrderPayments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}

	return entities.PaymentOrder{}, errors2.NewOrderPaymentNotFoundError(o.ID, paymentID)
}

func (o *Order) GetTotalEligiblePayments() vo.CurrencyAmount {
	var totalPaid decimal.Decimal

	for _, payment := range o.OrderPayments {
		if payment.Status.IsProcessing() || payment.Status.IsProcessed() {
			totalPaid = totalPaid.Add(payment.Total.Value)
		}
	}

	totalPaidAmount, _ := vo.NewCurrencyAmount(o.Currency, totalPaid)

	return totalPaidAmount
}

func (o *Order) GetTotalProcessed() vo.CurrencyAmount {
	totalPaid := decimal.Zero

	for _, payment := range o.OrderPayments {
		if payment.Status == enums.PaymentProcessed {
			totalPaid = totalPaid.Add(payment.Total.Value)
		}
	}

	totalPaidAmount, _ := vo.NewCurrencyAmount(o.Currency, totalPaid)

	return totalPaidAmount
}

func (o *Order) HasPaymentRefundable() bool {
	for _, payment := range o.OrderPayments {
		if payment.Status.IsProcessing() ||
			payment.Status.IsProcessed() ||
			payment.Status.IsPartiallyRefunded() {
			return true
		}
	}

	return false
}

func (o *Order) IsEmpty() bool {
	return o.ID == ""
}

func (o *Order) Register(r eventsourcing.RegisterFunc) {
	r(
		new(events.OrderCreated),
		new(events.PaymentProcessingStarted),
		new(events.OrderPaymentProcessedEvent),
		new(events.OrderPaymentFailedEvent),
		new(events.OrderPaymentTotalRefundedEvent),
		new(events.OrderPaymentPartialRefundedEvent),
		new(events.OrderPaymentAuthorizedEvent),
		new(events.OrderPaymentCapturedEvent),
		new(events.OrderPaymentReleasedEvent),
	)
}

func (o *Order) Transition(event eventsourcing.Event) {
	switch e := event.Data().(type) {
	case *events.OrderCreated:
		WhenOrderCreated(o, *e)
	case *events.PaymentProcessingStarted:
		WhenOrderPaymentProcessingStarted(o, *e)
	case *events.OrderPaymentProcessedEvent:
		WhenOrderPaymentProcessed(o, *e)
	case *events.OrderPaymentFailedEvent:
		WhenOrderPaymentFailed(o, *e)
	case *events.OrderPaymentTotalRefundedEvent:
		WhenPaymentRefunded(o, *e)
	case *events.OrderPaymentPartialRefundedEvent:
		WhenPaymentPartialRefunded(o, *e)
	case *events.OrderPaymentAuthorizedEvent:
		WhenOrderPaymentAuthorized(o, *e)
	case *events.OrderPaymentCapturedEvent:
		WhenPaymentCaptured(o, *e)
	case *events.OrderPaymentReleasedEvent:
		WhenPaymentReleased(o, *e)
	}
}
