package create

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/aggregate"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/adapter"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/event_store"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/queries"
)

var (
	ctx            = context.Background()
	currencyMXN, _ = value_objects.NewCurrencyCode("USD")
)

func TestCreatePaymentOrder(t *testing.T) {
	userType := value_objects.NewUserType(value_objects.Member)
	user := entities.NewUser(userType, "userID")
	paymentMethod := value_objects.NewCCPaymentMethod("cardID", "123")
	associatedOrigin := value_objects.NewAssociatedOrigin(enums.Booking)
	total, _ := value_objects.NewCurrencyAmount(
		value_objects.CurrencyCode{
			Code: "USD",
		},
		decimal.NewFromFloat(100),
	)
	totalOrder, _ := value_objects.NewCurrencyAmount(
		value_objects.CurrencyCode{
			Code: "USD",
		},
		decimal.NewFromFloat(200),
	)

	cardEntity := entities.Card{
		ID:             "cardID",
		ExternalCardID: "externalCardID",
		UserID:         "userID",
		CardType:       enums.CreditCard,
	}

	collectionAccount := entities.CollectionAccount{
		ID:                     "collectionAccountID",
		AccountType:            enums.Payers.String(),
		CollectionCenterID:     "collectionCenterID",
		CurrencyCode:           currencyMXN,
		AccountNumber:          "123456",
		BankName:               "bankName",
		InterbankAccountNumber: "interbankAccountNumber",
		EnterpriseID:           "enterpriseID",
	}
	payment := entities.NewPaymentOrder("paymentID", associatedOrigin, total, paymentMethod)

	commandMockWrong := command.NewCreatePaymentOrderCommandBuilder().
		WithUser(user).
		WithReferenceOrderID("referenceOrderID").
		WithPayment(payment).
		Build()

	commandMockRight := command.NewCreatePaymentOrderCommandBuilder().
		WithUser(user).
		WithReferenceOrderID("referenceOrderID").
		WithPayment(payment).
		WithAssociatedOrigin(associatedOrigin).
		WithCurrencyCode(currencyMXN).
		WithPaymentOrderID("paymentOrderID").
		Build()

	t.Run("should return an error getting order", func(t *testing.T) {
		mockRepo := event_store.NewOrderEventRepository(t)
		mockQueryAccount := queries.NewGetCollectionAccountByRouteUsecaseIF(t)
		mockQueryCard := queries.NewGetCardByUserUsecaseIF(t)
		mockPaymentOrderAdapter := adapter.NewPaymentOrderAdapterIF(t)

		order := new(aggregate.Order)

		mockRepo.On("Get", mock.Anything, "referenceOrderID", order).Once().Return(errors.New("order not found"))

		paymentUseCase := NewCreatePaymentOrderUseCase(mockRepo, mockQueryAccount, mockQueryCard, mockPaymentOrderAdapter)

		_, err := paymentUseCase.CreatePaymentOrder(ctx, commandMockWrong)

		assert.EqualError(t, err, "order not found")

		mockRepo.AssertExpectations(t)
		mockQueryCard.AssertExpectations(t)
		mockQueryAccount.AssertExpectations(t)
	})

	t.Run("should return an error processing payment getting collection account", func(t *testing.T) {
		mockRepo := event_store.NewOrderEventRepository(t)
		mockQueryAccount := queries.NewGetCollectionAccountByRouteUsecaseIF(t)
		mockQueryCard := queries.NewGetCardByUserUsecaseIF(t)
		mockPaymentOrderAdapter := adapter.NewPaymentOrderAdapterIF(t)

		order := new(aggregate.Order)

		mockRepo.On("Get", mock.Anything, "referenceOrderID", order).Once().Return(nil)
		mockQueryAccount.On("GetCollectionAccountByRoute",
			mock.Anything,
			order.CountryCode.Iso3(),
			commandMockWrong.AssociatedOrigin.Type.String(),
			commandMockWrong.CurrencyCode.Code,
			order.EnterpriseID,
		).Return(entities.CollectionAccount{}, assert.AnError)

		paymentUseCase := NewCreatePaymentOrderUseCase(mockRepo, mockQueryAccount, mockQueryCard, mockPaymentOrderAdapter)
		_, err := paymentUseCase.CreatePaymentOrder(ctx, commandMockWrong)

		assert.EqualError(t, err, "assert.AnError general error for testing")
		mockRepo.AssertExpectations(t)
		mockQueryCard.AssertExpectations(t)
		mockQueryAccount.AssertExpectations(t)
	})

	t.Run("should return an error getting card of user", func(t *testing.T) {
		mockRepo := event_store.NewOrderEventRepository(t)
		mockQueryAccount := queries.NewGetCollectionAccountByRouteUsecaseIF(t)
		mockQueryCard := queries.NewGetCardByUserUsecaseIF(t)
		mockPaymentOrderAdapter := adapter.NewPaymentOrderAdapterIF(t)

		order := new(aggregate.Order)

		mockRepo.On("Get", mock.Anything, "referenceOrderID", order).Once().Return(nil)

		mockQueryCard.On("GetCardByIDAndUserID", mock.Anything,
			commandMockWrong.User.ID, commandMockWrong.Payment.Method.CCData.Data.CardID, order.EnterpriseID).
			Return(entities.Card{}, assert.AnError)

		mockQueryAccount.On("GetCollectionAccountByRoute",
			mock.Anything,
			order.CountryCode.Iso3(),
			commandMockWrong.AssociatedOrigin.Type.String(),
			commandMockWrong.CurrencyCode.Code,
			order.EnterpriseID,
		).Return(entities.CollectionAccount{
			ID:                     "collectionAccountID",
			AccountType:            enums.Payers.String(),
			CollectionCenterID:     "collectionCenterID",
			CurrencyCode:           currencyMXN,
			AccountNumber:          "123456",
			BankName:               "bankName",
			InterbankAccountNumber: "interbankAccountNumber",
			EnterpriseID:           "enterpriseID",
		}, nil)

		paymentUseCase := NewCreatePaymentOrderUseCase(mockRepo, mockQueryAccount, mockQueryCard, mockPaymentOrderAdapter)
		_, err := paymentUseCase.CreatePaymentOrder(ctx, commandMockWrong)

		assert.EqualError(t, err, assert.AnError.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return an error processing payment", func(t *testing.T) {
		mockRepo := event_store.NewOrderEventRepository(t)
		mockQueryAccount := queries.NewGetCollectionAccountByRouteUsecaseIF(t)
		mockQueryCard := queries.NewGetCardByUserUsecaseIF(t)
		mockPaymentOrderAdapter := adapter.NewPaymentOrderAdapterIF(t)

		order := new(aggregate.Order)

		mockRepo.On("Get", mock.Anything, "referenceOrderID", order).Once().Return(nil)
		mockQueryAccount.On("GetCollectionAccountByRoute",
			mock.Anything,
			order.CountryCode.Iso3(),
			commandMockWrong.AssociatedOrigin.Type.String(),
			commandMockWrong.CurrencyCode.Code,
			order.EnterpriseID,
		).Return(entities.CollectionAccount{
			ID:                     "collectionAccountID",
			AccountType:            enums.Payers.String(),
			CollectionCenterID:     "collectionCenterID",
			CurrencyCode:           currencyMXN,
			AccountNumber:          "123456",
			BankName:               "bankName",
			InterbankAccountNumber: "interbankAccountNumber",
			EnterpriseID:           "enterpriseID",
		}, nil)

		mockQueryCard.On("GetCardByIDAndUserID", mock.Anything,
			commandMockWrong.User.ID, commandMockWrong.Payment.Method.CCData.Data.CardID, order.EnterpriseID).
			Return(cardEntity, nil)

		paymentUseCase := NewCreatePaymentOrderUseCase(mockRepo, mockQueryAccount, mockQueryCard, mockPaymentOrderAdapter)
		_, err := paymentUseCase.CreatePaymentOrder(ctx, commandMockWrong)

		assert.EqualError(t, err, "payment order cannot be added")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return an error calling adapter", func(t *testing.T) {
		mockRepo := event_store.NewOrderEventRepository(t)
		mockQueryAccount := queries.NewGetCollectionAccountByRouteUsecaseIF(t)
		mockQueryCard := queries.NewGetCardByUserUsecaseIF(t)
		mockPaymentOrderAdapter := adapter.NewPaymentOrderAdapterIF(t)

		order := new(aggregate.Order)

		mockRepo.On("Get", mock.Anything, "referenceOrderID", order).Once().Run(func(args mock.Arguments) {
			o := args.Get(2).(*aggregate.Order)
			o.Status = value_objects.OrderStatusProcessing()
			o.TotalAmount = totalOrder
		}).Return(nil)
		mockQueryAccount.On("GetCollectionAccountByRoute",
			mock.Anything,
			order.CountryCode.Iso3(),
			commandMockRight.AssociatedOrigin.Type.String(),
			commandMockRight.CurrencyCode.Code,
			order.EnterpriseID,
		).Return(collectionAccount, nil)

		mockQueryCard.On("GetCardByIDAndUserID", mock.Anything,
			commandMockWrong.User.ID, commandMockWrong.Payment.Method.CCData.Data.CardID, order.EnterpriseID).
			Return(cardEntity, nil)

		mockPaymentOrderAdapter.On("CreatePaymentOrder", mock.Anything, mock.Anything, mock.Anything, cardEntity).Return(assert.AnError)

		paymentUseCase := NewCreatePaymentOrderUseCase(mockRepo, mockQueryAccount, mockQueryCard, mockPaymentOrderAdapter)
		_, err := paymentUseCase.CreatePaymentOrder(ctx, commandMockRight)

		assert.Error(t, err, assert.AnError.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return an error saving a payment order", func(t *testing.T) {
		mockRepo := event_store.NewOrderEventRepository(t)
		mockQueryAccount := queries.NewGetCollectionAccountByRouteUsecaseIF(t)
		mockQueryCard := queries.NewGetCardByUserUsecaseIF(t)
		mockPaymentOrderAdapter := adapter.NewPaymentOrderAdapterIF(t)

		order := new(aggregate.Order)

		mockRepo.On("Get", mock.Anything, "referenceOrderID", order).Once().Run(func(args mock.Arguments) {
			o := args.Get(2).(*aggregate.Order)
			o.Status = value_objects.OrderStatusProcessing()
			o.TotalAmount = totalOrder
		}).Return(nil)
		mockQueryAccount.On("GetCollectionAccountByRoute",
			mock.Anything,
			order.CountryCode.Iso3(),
			commandMockRight.AssociatedOrigin.Type.String(),
			commandMockRight.CurrencyCode.Code,
			order.EnterpriseID,
		).Return(collectionAccount, nil)

		mockQueryCard.On("GetCardByIDAndUserID", mock.Anything,
			commandMockWrong.User.ID, commandMockWrong.Payment.Method.CCData.Data.CardID, order.EnterpriseID).
			Return(cardEntity, nil)

		mockPaymentOrderAdapter.On("CreatePaymentOrder", mock.Anything, mock.Anything, mock.Anything, cardEntity).
			Return(nil)

		mockRepo.On("Save", mock.Anything, mock.Anything).Return(assert.AnError)

		paymentUseCase := NewCreatePaymentOrderUseCase(mockRepo, mockQueryAccount, mockQueryCard, mockPaymentOrderAdapter)
		_, err := paymentUseCase.CreatePaymentOrder(ctx, commandMockRight)

		assert.EqualError(t, err, "assert.AnError general error for testing")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return an error creating payment flow", func(t *testing.T) {
		mockRepo := event_store.NewOrderEventRepository(t)
		mockQueryAccount := queries.NewGetCollectionAccountByRouteUsecaseIF(t)
		mockQueryCard := queries.NewGetCardByUserUsecaseIF(t)
		mockPaymentOrderAdapter := adapter.NewPaymentOrderAdapterIF(t)

		order := new(aggregate.Order)

		mockRepo.On("Get", mock.Anything, "referenceOrderID", order).Once().Return(nil)
		mockQueryAccount.On("GetCollectionAccountByRoute",
			mock.Anything,
			order.CountryCode.Iso3(),
			commandMockRight.AssociatedOrigin.Type.String(),
			commandMockRight.CurrencyCode.Code,
			order.EnterpriseID,
		).Return(collectionAccount, nil)

		invalidCard := cardEntity
		invalidCard.CardType = "invalid_type"

		mockQueryCard.On("GetCardByIDAndUserID", mock.Anything,
			commandMockRight.User.ID, commandMockRight.Payment.Method.CCData.Data.CardID, order.EnterpriseID).
			Return(invalidCard, nil)

		paymentUseCase := NewCreatePaymentOrderUseCase(mockRepo, mockQueryAccount, mockQueryCard, mockPaymentOrderAdapter)
		_, err := paymentUseCase.CreatePaymentOrder(ctx, commandMockRight)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
		mockQueryCard.AssertExpectations(t)
		mockQueryAccount.AssertExpectations(t)
	})

	t.Run("should save an order successful", func(t *testing.T) {
		mockRepo := event_store.NewOrderEventRepository(t)
		mockQueryAccount := queries.NewGetCollectionAccountByRouteUsecaseIF(t)
		mockQueryCard := queries.NewGetCardByUserUsecaseIF(t)
		mockPaymentOrderAdapter := adapter.NewPaymentOrderAdapterIF(t)

		order := new(aggregate.Order)

		mockRepo.On("Get", mock.Anything, "referenceOrderID", order).Once().Run(func(args mock.Arguments) {
			o := args.Get(2).(*aggregate.Order)
			o.Status = value_objects.OrderStatusProcessing()
			o.TotalAmount = totalOrder
			o.AllowCapture = true
		}).Return(nil)
		mockQueryAccount.On("GetCollectionAccountByRoute",
			mock.Anything,
			order.CountryCode.Iso3(),
			commandMockRight.AssociatedOrigin.Type.String(),
			commandMockRight.CurrencyCode.Code,
			order.EnterpriseID,
		).Return(collectionAccount, nil)

		mockQueryCard.On("GetCardByIDAndUserID", mock.Anything,
			commandMockRight.User.ID, commandMockRight.Payment.Method.CCData.Data.CardID, order.EnterpriseID).
			Return(cardEntity, nil)

		mockPaymentOrderAdapter.On("CreatePaymentOrder", mock.Anything, mock.Anything, mock.Anything, cardEntity).
			Return(nil)

		mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

		paymentUseCase := NewCreatePaymentOrderUseCase(mockRepo, mockQueryAccount, mockQueryCard, mockPaymentOrderAdapter)
		_, err := paymentUseCase.CreatePaymentOrder(ctx, commandMockRight)

		assert.Nil(t, err)
		mockRepo.AssertExpectations(t)
	})
}
