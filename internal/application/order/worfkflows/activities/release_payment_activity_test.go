package activities

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/capture_flow"
)

func TestReleasePaymentActivity_ReleasePayment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		mockReleasePaymentUseCase := capture_flow.NewPaymentReleaseUseCaseIF(t)
		activity := NewReleasePaymentActivity(mockReleasePaymentUseCase)

		ctx := context.Background()
		orderID := "test-order-id"
		paymentID := "test-payment-id"
		reason := "test-reason"

		mockReleasePaymentUseCase.EXPECT().
			ReleasePayment(ctx, orderID, paymentID, reason).
			Return(nil)

		// Act
		err := activity.ReleasePayment(ctx, orderID, paymentID, reason)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("error_from_use_case", func(t *testing.T) {
		// Arrange
		mockReleasePaymentUseCase := capture_flow.NewPaymentReleaseUseCaseIF(t)
		activity := NewReleasePaymentActivity(mockReleasePaymentUseCase)

		ctx := context.Background()
		orderID := "test-order-id"
		paymentID := "test-payment-id"
		reason := "test-reason"
		expectedError := errors.New("release payment failed")

		mockReleasePaymentUseCase.EXPECT().
			ReleasePayment(ctx, orderID, paymentID, reason).
			Return(expectedError)

		// Act
		err := activity.ReleasePayment(ctx, orderID, paymentID, reason)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectedError.Error())
	})
}
