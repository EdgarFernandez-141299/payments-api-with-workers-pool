package event_store

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/aggregate"
)

type OrderEventRepository interface {
	Save(oldCtx context.Context, a *aggregate.Order) error
	Get(oldCtx context.Context, id string, a *aggregate.Order) error
}
