package eventsourcing

import (
	"context"

	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
	entitiesOrder "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/order/entities"

	"gitlab.com/clubhub.ai1/go-libraries/eventsourcing"
	"gitlab.com/clubhub.ai1/go-libraries/eventsourcing/eventstore/dynamo"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/config"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/event_store"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/aggregate"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/aws"

	"gitlab.com/clubhub.ai1/gommon/logger"
)

type OrderEventsourcingRepository struct {
	*eventsourcing.EventRepository
	writRepository repository.OrderWriteRepositoryIF
}

func (b *OrderEventsourcingRepository) Get(oldCtx context.Context, id string, a *aggregate.Order) error {
	return b.EventRepository.GetWithContext(oldCtx, id, a)
}

func (b *OrderEventsourcingRepository) Save(oldCtx context.Context, a *aggregate.Order) error {
	return b.EventRepository.Save(oldCtx, a)
}

func (b *OrderEventsourcingRepository) Create(oldCtx context.Context, order *aggregate.Order) error {
	return decorators.TraceDecoratorNoReturn(
		oldCtx,
		"OrderEventsourcingRepository.Create",
		func(ctx context.Context, span decorators.Span) error {
			entity := entitiesOrder.NewOrderEntityBuilder().
				SetID(order.ID).
				SetReferenceOrderID(order.ID).
				SetCountryCode(order.CountryCode.Code).
				SetCurrencyCode(order.Currency.Code).
				SetTotalAmount(order.TotalAmount.Value).
				SetEnterpriseID(order.EnterpriseID).
				SetStatus(value_objects.OrderStatusProcessing().Get()).
				SetMetadata(order.Metadata).
				SetUserID(order.User.ID).
				SetAllowCapture(order.AllowCapture).
				Build()

			err := b.writRepository.CreateOrder(ctx, entity)

			if err != nil {
				return err
			}

			return b.EventRepository.Save(ctx, order)
		},
	)
}

func newOrderEventsourcingRepository(eventStore *eventsourcing.EventRepository, repo repository.OrderWriteRepositoryIF) *OrderEventsourcingRepository {
	eventRepository := &OrderEventsourcingRepository{
		EventRepository: eventStore,
		writRepository:  repo,
	}

	eventRepository.Register(new(aggregate.Order))

	return eventRepository
}

func NewOrderEventsourcingRepository(logger logger.LoggerInterface, writeRepository repository.OrderWriteRepositoryIF) event_store.OrderEventRepository {
	cfg := config.Config()

	dynamodbClient := aws.NewDynamoClient(cfg, logger)

	tablePrefix := cfg.DynamoDB.PaymentsEventStoreTable

	eventStore := dynamo.NewWithOpts(dynamodbClient, dynamo.WithTableName(tablePrefix))

	eventRepository := eventsourcing.NewEventRepository(eventStore)

	return newOrderEventsourcingRepository(eventRepository, writeRepository)
}
