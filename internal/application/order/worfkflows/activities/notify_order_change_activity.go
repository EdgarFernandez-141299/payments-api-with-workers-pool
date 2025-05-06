package activities

import (
	"context"
	"gitlab.com/clubhub.ai1/go-libraries/saga/errors"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/services"
)

const NotifyOrderChangeActivityName = "NotifyOrderChangeActivity"

type NotifyOrderChangeActivity struct {
	orderNotificationStrategy services.OrderNotificationStrategyIF
}

func NewNotifyOrderChangeActivity(orderNotificationStrategy services.OrderNotificationStrategyIF) *NotifyOrderChangeActivity {
	return &NotifyOrderChangeActivity{
		orderNotificationStrategy: orderNotificationStrategy,
	}
}

type NotifyOrderChangeParams struct {
	OrderID   string
	PaymentID string
}

func (n *NotifyOrderChangeActivity) NotifyOrderChange(ctx context.Context, params NotifyOrderChangeParams) error {
	return errors.WrapActivityError(n.orderNotificationStrategy.NotifyChange(ctx, params.OrderID, params.PaymentID))
}
