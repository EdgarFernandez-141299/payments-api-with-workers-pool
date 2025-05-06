package infra

import (
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_receipt/adapters/pdf_generator"
	"log"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_receipt/use_case"
	awsClient "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/aws"

	webhooks2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/http/webhooks"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/member"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/storage"

	notificationServices "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/notification/services"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/services"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/worfkflows/activities"

	creditCardAdapters "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/adapters"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/use_cases/create"
	captureFlow "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/use_cases/capture_flow"
	createPaymentOrder "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/use_cases/create"
	postPaymentOrder "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/use_cases/post_payment"

	userAdapters "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/deuna"
	httpconfig "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/http/config"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/eventsourcing"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/worker/workflows"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/capture_flow"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_center"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/webhooks"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/logger"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/worker"

	"github.com/labstack/echo/v4"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/config"
	cardUsecase "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/usecases"
	queryCardUsecase "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/usecases/queries"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/collection_account/queries"
	collectionAccountUsecase "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/collection_account/use_cases"
	collectionAccountRouteUsecase "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/collection_account_route/use_cases"
	collectionCenterUsecase "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/collection_center/use_cases"
	queriesOrder "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/use_cases/queries"
	paymentConceptUsecase "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_concept/use_cases"
	paymentMethodUsecase "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_method/use_cases"
	refundAdapters "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/refund/adapters"
	refundUsecase "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/refund/use_cases"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_account"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_account_route"

	adapterOrder "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters"
	adapterPaymentOrder "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/card"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/payment_concept"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/payment_method"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/user"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/router/group"
	db "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/db"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/http"
	cardRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/card"
	collectionAccountRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/collection_account"
	collecitionRouteRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/collection_account_route"
	collectionCenterRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/collection_center"
	deunaRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/deuna_payment"
	orderRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/order"
	paymentConceptRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_concept"
	paymentMethodRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_method"
	paymentOrderRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_order"
	paymentReceiptRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_receipt"
	refundRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/refund"
	userRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/user"
	"go.uber.org/fx"
)

func Run() {
	if err := config.Environments(); err != nil {
		log.Panicf("There was an error getting the configurations %#v", err)
	}

	fx.New(
		// Infrastructure components layer
		fx.Provide(NewHTTPServer),
		fx.Provide(worker.NewTemporalClient),
		fx.Provide(workflows.NewPaymentOrderWorkflow),
		fx.Provide(NewEchoServer),
		fx.Provide(http.NewMemberHTTPClient),
		fx.Provide(http.NewDeUnaHTTPClient),
		fx.Provide(http.NewDeUnaAuthHTTPClient),
		fx.Provide(http.NewDeUnaCardHTTPClient),
		fx.Provide(http.NewDeUnaLoginHTTPClient),
		fx.Provide(http.NewDeUnaOrderHTTPClient),
		fx.Provide(http.NewDeUnaPaymentHTTPClient),
		fx.Provide(http.NewDeUnaRefundHTTPClient),
		fx.Provide(httpconfig.NewMailHTTPClient),
		fx.Provide(http.NewMailHTTPClient),
		fx.Provide(http.NewDeUnaCaptureFlowHTTPClient),
		fx.Provide(webhooks2.NewBillsApiWebhookUrl),
		fx.Provide(logger.NewLogger),
		fx.Provide(logger.NewTraceProvider),
		fx.Provide(logger.NewLoggerLegacy),
		fx.Provide(db.NewPostgresConnection),
		fx.Provide(awsClient.NewS3Client),
		fx.Provide(storage.NewCDNURLProvider),
		fx.Provide(storage.NewS3StorageAdapter),
		fx.Provide(storage.NewPaymentReceiptBucket),
		fx.Provide(storage.NewDiskStorageAdapter),
		fx.Provide(storage.NewProvideStorageAdapter),
		fx.Provide(httpconfig.DeunaHttpConfig),
		fx.Provide(collectionAccountRepository.NewCollectionAccountRepository),
		fx.Provide(cardRepository.NewCardReadRepository),
		fx.Provide(cardRepository.NewCardWriteRepository),
		fx.Provide(paymentConceptRepository.NewPaymentConceptRepository),
		fx.Provide(paymentMethodRepository.NewPaymentMethodRepository),
		fx.Provide(userRepository.NewUserReadRepository),
		fx.Provide(collecitionRouteRepository.NewCollectionCenterRepositoryIF),
		fx.Provide(collectionCenterRepository.NewCollectionCenterRepository),
		fx.Provide(collectionAccountRouteUsecase.NewCollectionAccountRepositoryIF),
		fx.Provide(collectionAccountRepository.NewCollectionAccountReadRepository),
		fx.Provide(queries.NewGetCollectionAccountByRouteUsecase),
		fx.Provide(paymentOrderRepository.NewPaymentOrderWriteRespository),
		fx.Provide(paymentOrderRepository.NewPaymentOrderReadRepository),
		fx.Provide(paymentOrderRepository.NewUpdateOrderStatusRepository),
		fx.Provide(paymentReceiptRepository.NewPaymentReceiptRepository),
		fx.Provide(refundRepository.NewRefundWriteRepository),
		fx.Provide(refundRepository.NewRefundReadRepository),
		fx.Provide(queriesOrder.NewQueriesOrderUseCase),

		fx.Provide(deunaRepository.NewDeunaPaymentWriteRepository),
		fx.Provide(orderRepository.NewOrderWriteRepository),
		fx.Provide(orderRepository.NewOrderReadRepository),
		fx.Provide(userRepository.NewUserWriteRepository),

		fx.Provide(adapterOrder.NewOrderAdapter),
		fx.Provide(member.NewMemberAdapterIF),

		fx.Provide(eventsourcing.NewOrderEventsourcingRepository),

		fx.Provide(collectionAccountUsecase.NewCollectionAccountUsecase),
		fx.Provide(creditCardAdapters.NewDeunaCardAdapter),
		fx.Provide(userAdapters.NewDeunaLoginAdapter),
		fx.Provide(paymentConceptUsecase.NewPaymentConceptUsecase),
		fx.Provide(deuna.NewCreateUserUseCases),
		fx.Provide(adapterPaymentOrder.NewPaymentOrderAdapter),
		fx.Provide(adapterPaymentOrder.NewPaymentCaptureFlowAdapter),
		fx.Provide(refundAdapters.NewRefundAdapter),
		fx.Provide(refundAdapters.NewPartialRefundAdapter),
		fx.Provide(cardUsecase.NewCardUseCase),
		fx.Provide(cardUsecase.NewDeleteCardUseCase),
		fx.Provide(use_case.NewGenerateReceiptPaymentUseCase),
		fx.Provide(queryCardUsecase.NewGetCardUsecase),
		fx.Provide(queryCardUsecase.NewGetCardByUserUsecase),
		fx.Provide(paymentMethodUsecase.NewPaymentMethodUseCases),
		fx.Provide(collectionCenterUsecase.NewCollectionCenterUsecase),
		fx.Provide(NewEchoGroup),
		fx.Provide(collection_account.NewCollectionAccountHandler),
		fx.Provide(payment_concept.NewPaymentConceptHandler),
		fx.Provide(payment_method.NewPaymentMethodHandler),
		fx.Provide(user.NewUserHandler),
		fx.Provide(card.NewCardHandler),
		fx.Provide(card.NewDeleteCardHandler),
		fx.Provide(activities.NewCheckOrderActivity),
		fx.Provide(activities.NewCreatePaymentOrderActivity),
		fx.Provide(activities.NewPostProcessingPaymentOrderActivity),
		fx.Provide(activities.NewNotifyOrderChangeActivity),
		fx.Provide(activities.NewGeneratePaymentReceiptActivity),
		fx.Provide(activities.NewCapturePaymentActivity),
		fx.Provide(activities.NewReleasePaymentActivity),
		fx.Provide(postPaymentOrder.NewPostProcessingPaymentOrderUseCase),
		fx.Provide(pdf_generator.NewReceiptPaymentGeneratorImpl),
		fx.Provide(group.NewPaymentConceptRoutes),

		// Routes
		fx.Provide(group.NewCollectionAccountRoutes),
		fx.Provide(group.NewOrderRoutes),
		fx.Provide(group.NewPaymentMethodRoutes),
		fx.Provide(group.NewHealthRoutes),
		fx.Provide(group.NewCollectionCenterRoutes),
		fx.Provide(group.NewCollectionAccountRouteRoutes),
		fx.Provide(group.NewUserRoutes),
		fx.Provide(group.NewCardRoutes),
		fx.Provide(group.NewDeleteCardRoutes),
		fx.Provide(group.NewWebhooksRoutes),
		fx.Provide(group.NewCaptureFlowRoutes),

		// Handlers
		fx.Provide(collection_account_route.NewCollectionAccountRouteHandler),
		fx.Provide(order.NewOrderHandler),
		fx.Provide(refundUsecase.NewPartialRefundUse),
		fx.Provide(capture_flow.NewCaptureFlowHandler),

		fx.Provide(order.NewPaymentRefundHandler),
		fx.Provide(collection_center.NewCollectionCenterHandler),
		fx.Provide(order.NewOrderPaymentsHandler),
		fx.Provide(webhooks.NewWebhookHandler),
		// Use cases
		fx.Provide(create.NewCreateOrderUseCase),
		fx.Provide(createPaymentOrder.NewCreatePaymentOrderUseCase),
		fx.Provide(captureFlow.NewPaymentCaptureUseCase),
		fx.Provide(captureFlow.NewPaymentReleaseUseCase),
		fx.Provide(refundUsecase.NewRefundTotalUseCase),

		// Order Services
		fx.Provide(webhooks2.NewWebhookNotificationResource),
		fx.Provide(services.NewOrderNotificationOrchestrator),
		fx.Provide(services.NewOrderFailedNotificationService),
		fx.Provide(services.NewOrderProcessedNotificationService),

		// Notification Services
		fx.Provide(notificationServices.NewNotificationService),

		// infra Init functions
		fx.Invoke(func(*echo.Echo) {}),
		fx.Invoke(func(*group.CollectionAccountRoutes) {}),
		fx.Invoke(func(*group.OrderRoutes) {}),
		fx.Invoke(func(*group.PaymentMethodRoutes) {}),
		fx.Invoke(func(*group.HealthRoutes) {}),
		fx.Invoke(func(*group.CollectionCenterRoutes) {}),
		fx.Invoke(func(*group.CollectionAccountRouteRoutes) {}),
		fx.Invoke(func(*group.UserRoutes) {}),
		fx.Invoke(func(*group.CardRoutes) {}),
		fx.Invoke(func(*group.DeleteCardRoutes) {}),
		fx.Invoke(func(*group.WebhooksRoutes) {}),
		fx.Invoke(func(*group.CaptureFlowRoutes) {}),
	).Run()
}
