package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/payment_method/dto/request"
	mockRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repositories"
)

var (
	ctx = context.TODO()
)

func TestCreate(t *testing.T) {
	t.Run("should return an error creating payment method", func(t *testing.T) {
		// Arrange
		requestMock := request.PaymentMethodRequest{}
		repositoryMock := mockRepository.NewPaymentMethodRepositoryIF(t)
		repositoryMock.On("Create", ctx, mock.Anything).
			Return(errors.New("error creating payment method"))

		// Act
		paymentUsecaseMethodMock := NewPaymentMethodUseCases(repositoryMock)
		_, err := paymentUsecaseMethodMock.Create(ctx, requestMock, "enterpriseId")

		// Assert
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error creating payment method")
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should create a payment method", func(t *testing.T) {
		// Arrange
		requestMock := request.PaymentMethodRequest{}
		repositoryMock := mockRepository.NewPaymentMethodRepositoryIF(t)
		repositoryMock.On("Create", ctx, mock.Anything).
			Return(nil)

		// Act
		paymentUsecaseMethodMock := NewPaymentMethodUseCases(repositoryMock)
		response, err := paymentUsecaseMethodMock.Create(ctx, requestMock, "enterpriseId")

		// Assert
		assert.Nil(t, err)
		assert.NotNil(t, response)
		repositoryMock.AssertExpectations(t)
	})
}
