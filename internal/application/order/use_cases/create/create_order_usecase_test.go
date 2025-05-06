package create

import (
	"context"
	"errors"
	"gitlab.com/clubhub.ai1/go-libraries/eventsourcing"
	errors2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/aggregate"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/event_store"
)

func TestCreateOrder(t *testing.T) {
	userID := "user123"
	user := entities.NewUser(value_objects.NewUserType(value_objects.Member), userID)
	usdCurrencyCode, _ := value_objects.NewCurrencyCode("USD")
	totalAmountAsFloat := decimal.NewFromFloat(100.5)
	enterpriseID := "enterprise123"
	countryCode, _ := value_objects.NewCountryWithCode("MX")
	totalAmount, _ := value_objects.NewCurrencyAmount(usdCurrencyCode, totalAmountAsFloat)

	t.Run("should return an error when order already exist", func(t *testing.T) {
		ctx := context.Background()

		cmd := command.NewCreateOrderCommandBuilder().
			WithReferenceID("1234").
			WithTotalAmount(totalAmount).
			WithPhoneNumber("3222332").
			WithUser(user).
			WithCurrencyCode(usdCurrencyCode).
			WithCountryCode(countryCode).
			WithEnterpriseID(enterpriseID).
			Build()

		mockRepo := event_store.NewOrderEventRepository(t)

		mockRepo.On("Get", mock.Anything, cmd.ReferenceID, mock.AnythingOfType("*aggregate.Order")).
			Return(eventsourcing.ErrAggregateAlreadyExists).Once()

		useCase := NewCreateOrderUseCase(nil, mockRepo)

		_, err := useCase.Create(ctx, cmd)

		assert.Error(t, err, errors2.NewOrderAlreadyExistError(cmd.ReferenceID).Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("error create order event", func(t *testing.T) {
		ctx := context.Background()

		cmd := command.NewCreateOrderCommandBuilder().
			WithReferenceID("1234").
			WithTotalAmount(totalAmount).
			WithPhoneNumber("3222332").
			WithUser(user).
			WithCurrencyCode(usdCurrencyCode).
			WithCountryCode(countryCode).
			WithEnterpriseID(enterpriseID).
			Build()

		orderEventRepository := event_store.NewOrderEventRepository(t)
		orderEventRepository.On("Get", mock.Anything, cmd.ReferenceID, mock.AnythingOfType("*aggregate.Order")).
			Return(nil).Once()

		orderEventRepository.On("Create", mock.Anything, mock.AnythingOfType("*aggregate.Order")).
			Return(errors.New("error saving event")).Once()

		useCase := NewCreateOrderUseCase(nil, orderEventRepository)

		result, err := useCase.Create(ctx, cmd)

		assert.EqualError(t, err, "error saving event")
		assert.Equal(t, "", result.ReferenceOrderID)

		orderEventRepository.AssertExpectations(t)
	})

	t.Run("Create order successful", func(t *testing.T) {
		ctx := context.Background()

		cmd := command.NewCreateOrderCommandBuilder().
			WithReferenceID("1234").
			WithTotalAmount(totalAmount).
			WithPhoneNumber("3222332").
			WithUser(user).
			WithCurrencyCode(usdCurrencyCode).
			WithCountryCode(countryCode).
			WithEnterpriseID(enterpriseID).
			Build()

		mockRepo := event_store.NewOrderEventRepository(t)

		mockRepo.On("Get", mock.Anything, cmd.ReferenceID, mock.AnythingOfType("*aggregate.Order")).
			Return(nil).Once()

		mockRepo.On("Create", mock.Anything, mock.IsType(new(aggregate.Order))).
			Return(nil).
			Once()

		useCase := NewCreateOrderUseCase(nil, mockRepo)

		result, err := useCase.Create(ctx, cmd)

		assert.NoError(t, err)
		assert.Equal(t, cmd.ReferenceID, result.ReferenceOrderID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Create order successful when repository send aggregate not found error", func(t *testing.T) {
		ctx := context.Background()

		cmd := command.NewCreateOrderCommandBuilder().
			WithReferenceID("1234").
			WithTotalAmount(totalAmount).
			WithPhoneNumber("3222332").
			WithUser(user).
			WithCurrencyCode(usdCurrencyCode).
			WithCountryCode(countryCode).
			WithEnterpriseID(enterpriseID).
			Build()

		mockRepo := event_store.NewOrderEventRepository(t)

		mockRepo.On("Get", mock.Anything, cmd.ReferenceID, mock.AnythingOfType("*aggregate.Order")).
			Return(eventsourcing.ErrAggregateNotFound).Once()

		mockRepo.On("Create", mock.Anything, mock.IsType(new(aggregate.Order))).
			Return(nil).
			Once()

		useCase := NewCreateOrderUseCase(nil, mockRepo)

		result, err := useCase.Create(ctx, cmd)

		assert.NoError(t, err)
		assert.Equal(t, cmd.ReferenceID, result.ReferenceOrderID)

		mockRepo.AssertExpectations(t)
	})
}
