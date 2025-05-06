package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/payment_concept/dto/request"

	mockRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repositories"
)

var (
	ctx = context.TODO()
)

func TestCreate(t *testing.T) {
	t.Run("should return an error creating payment concept", func(t *testing.T) {
		// Arrange
		repositoryMock := mockRepository.NewPaymentConceptRepositoryIF(t)

		requestMock := request.PaymentConceptRequest{
			Name:        "TEST",
			Code:        "TEST",
			Description: "TEST",
		}

		repositoryMock.On("Create", ctx, mock.Anything).Return(errors.New("error creating payment concept"))

		paymentUsecaseMock := NewPaymentConceptUsecase(repositoryMock)
		_, err := paymentUsecaseMock.Create(ctx, requestMock, "enterpriseId")

		assert.NotNil(t, err)
		assert.EqualError(t, err, "error creating payment concept")
	})

	t.Run("should create a payment concept", func(t *testing.T) {
		// Arrange
		repositoryMock := mockRepository.NewPaymentConceptRepositoryIF(t)

		requestMock := request.PaymentConceptRequest{
			Name:        "TEST",
			Code:        "TEST",
			Description: "TEST",
		}

		repositoryMock.On("Create", ctx, mock.Anything).Return(nil)

		paymentUsecaseMock := NewPaymentConceptUsecase(repositoryMock)
		response, err := paymentUsecaseMock.Create(ctx, requestMock, "enterpriseId")

		assert.Nil(t, err)
		assert.Equal(t, response.Name, "TEST")
		assert.Equal(t, response.Code, "TEST")
		assert.Equal(t, response.Description, "TEST")

		repositoryMock.AssertExpectations(t)
	})
}
