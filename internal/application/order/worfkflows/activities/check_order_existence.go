package activities

import (
	"context"
	"errors"
	"gitlab.com/clubhub.ai1/go-libraries/eventsourcing"
	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/event_store"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/aggregate"
)

const CheckOrderExistenceActivityName = "CheckOrderExistenceActivity"

type CheckOrderActivity struct {
	repository event_store.OrderEventRepository
}

func NewCheckOrderActivity(repository event_store.OrderEventRepository) *CheckOrderActivity {
	return &CheckOrderActivity{
		repository: repository,
	}
}

func (a *CheckOrderActivity) CheckOrderExistence(ctx context.Context, orderID string) (bool, error) {
	return decorators.TraceDecorator[bool](
		ctx,
		"CreatePaymentOrderActivity.CheckOrderExistence",
		func(ctx context.Context, span decorators.Span) (bool, error) {
			o := new(aggregate.Order)
			err := a.repository.Get(ctx, orderID, o)

			if err != nil {
				if errors.Is(err, eventsourcing.ErrAggregateNotFound) {
					return false, nil
				}
				return false, err
			}

			return !o.IsEmpty(), nil
		})
}
