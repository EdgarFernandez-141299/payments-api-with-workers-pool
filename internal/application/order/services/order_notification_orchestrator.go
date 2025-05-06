package services

import (
	"context"
	"fmt"
	"sync"

	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/event_store"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/aggregate"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/integration_events"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

type BillsApiWebhookUrl value_objects.WebhookUrl

type OrderNotificationStrategyIF interface {
	NotifyChange(ctx context.Context, orderID string, paymentID string) error
}

type OrderNotificationOrchestrator struct {
	orderRepository          event_store.OrderEventRepository
	orderFailedNotification  OrderFailedNotificationServiceIF
	orderSuccessNotification OrderProcessedNotificationServiceIF
	billsApiWebHookUrl       BillsApiWebhookUrl
}

func NewOrderNotificationOrchestrator(
	orderRepository event_store.OrderEventRepository,
	orderFailedNotification OrderFailedNotificationServiceIF,
	orderSuccessNotification OrderProcessedNotificationServiceIF,
	billsApiWebHookUrl BillsApiWebhookUrl,
) OrderNotificationStrategyIF {
	return &OrderNotificationOrchestrator{
		orderRepository:          orderRepository,
		orderFailedNotification:  orderFailedNotification,
		orderSuccessNotification: orderSuccessNotification,
		billsApiWebHookUrl:       billsApiWebHookUrl,
	}
}

func (o *OrderNotificationOrchestrator) notifyErrorBillsApi(
	oldCtx context.Context, failEvent integration_events.OrderFailedPaidIntegrationEvent,
) error {
	return decorators.TraceDecoratorNoReturn(oldCtx,
		"OrderNotificationOrchestrator.notifyErrorBillsApi",
		func(ctx context.Context, span decorators.Span) error {
			return o.orderFailedNotification.Notify(ctx, value_objects.WebhookUrl(o.billsApiWebHookUrl), failEvent)
		})
}

func (o *OrderNotificationOrchestrator) notifySuccessBillsApi(oldCtx context.Context, ev integration_events.OrderPaymentProcessedIntegrationEvent) error {
	return decorators.TraceDecoratorNoReturn(oldCtx,
		"OrderNotificationOrchestrator.notifySuccessBillsApi",
		func(ctx context.Context, span decorators.Span) error {
			return o.orderSuccessNotification.Notify(ctx, value_objects.WebhookUrl(o.billsApiWebHookUrl), ev)
		})
}

func (o *OrderNotificationOrchestrator) NotifyChange(oldCtx context.Context, orderID string, paymentID string) error {
	return decorators.TraceDecoratorNoReturn(oldCtx,
		"OrderNotificationOrchestrator.NotifyChange",
		func(ctx context.Context, span decorators.Span) error {
			var order aggregate.Order
			if err := o.orderRepository.Get(ctx, orderID, &order); err != nil {
				return errors.NewOrderNotFoundError(orderID)
			}

			payment, err := order.FindPaymentByID(paymentID)
			if err != nil {
				return errors.NewOrderPaymentNotFoundError(orderID, paymentID)
			}

			params := integration_events.IntegrationEventsParams{
				ReferenceOrderID:   order.ID,
				ReferencePaymentID: payment.ID,
				AssociatedPayment:  payment.OriginType.Type.String(),
				TotalOrderAmount:   order.TotalAmount.Value.InexactFloat64(),
				Currency:           order.Currency.Code,
				UserID:             order.User.ID,
				UserType:           order.User.Type.String(),
				EnterpriseID:       order.EnterpriseID,
				TotalOrderPaid:     order.GetTotalProcessed().Value.InexactFloat64(),
				TotalPaymentAmount: order.TotalAmount.Value.InexactFloat64(),
				Metadata:           order.Metadata,
				PaymentFlow:        payment.PaymentFlow,
				ReceiptUrl:         payment.ReceiptUrl,
				CardData: integration_events.CardData{
					CardNumber:    payment.PaymentCard.CardLast4,
					CardType:      payment.PaymentCard.CardBrand,
					MethodPayment: payment.PaymentCard.CardType,
				},
			}

			numOfWorkers := 2
			errChan := make(chan error, numOfWorkers)
			wg := &sync.WaitGroup{}

			if payment.Status.IsFailure() {
				failEvent := integration_events.NewOrderFailedIntegrationEvent(params, payment.FailureReason, payment.Status.String())
				wg.Add(2)

				go func(w *sync.WaitGroup) {
					errChan <- o.orderFailedNotification.Notify(ctx, order.WebhookUrl, failEvent)
					w.Done()
				}(wg)

				go func(w *sync.WaitGroup) {
					billsError := o.notifyErrorBillsApi(ctx, failEvent)
					errChan <- fmt.Errorf("billsError: %v", billsError)
					w.Done()
				}(wg)

			} else {
				successEvent := integration_events.NewOrderPaymentProcessedIntegrationEvent(params, payment.AuthorizationCode, payment.Status.String())
				wg.Add(2)

				go func(w *sync.WaitGroup) {
					errChan <- o.orderSuccessNotification.Notify(ctx, order.WebhookUrl, successEvent)
					w.Done()
				}(wg)

				go func(w *sync.WaitGroup) {
					errChan <- o.notifySuccessBillsApi(ctx, successEvent)
					w.Done()
				}(wg)
			}

			go func() {
				wg.Wait()
				close(errChan)
			}()

			var finalErr error = nil

			for incomingErr := range errChan {
				if incomingErr != nil {
					finalErr = incomingErr
				}
			}

			return finalErr
		})
}
