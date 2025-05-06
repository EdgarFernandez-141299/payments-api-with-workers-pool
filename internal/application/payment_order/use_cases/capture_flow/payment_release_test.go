package capture_flow

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/aggregate"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/adapter"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/event_store"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository"
)

func TestPaymentReleaseUseCase_ReleasePayment(t *testing.T) {
	tests := []struct {
		name          string
		orderID       string
		paymentID     string
		reason        string
		setupMocks    func(*event_store.OrderEventRepository, *adapter.PaymentCaptureFlowAdapterIF, *repository.TransactionsRepositoryIF)
		expectedError error
	}{
		{
			name:      "Liberación de pago exitosa",
			orderID:   "order123",
			paymentID: "payment123",
			reason:    "Cliente solicitó cancelación",
			setupMocks: func(mockRepo *event_store.OrderEventRepository, mockAdapter *adapter.PaymentCaptureFlowAdapterIF, mockTransactionsRepo *repository.TransactionsRepositoryIF) {
				order := &aggregate.Order{
					EnterpriseID: "enterprise123",
					Status:       value_objects.OrderStatusProcessing(),
				}

				mockRepo.On("Get", mock.Anything, "order123", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					arg := args.Get(2).(*aggregate.Order)
					*arg = *order
				})
				mockAdapter.On("ReleasePayment", mock.Anything, "order123", "payment123", "Cliente solicitó cancelación").Return(nil)
				mockTransactionsRepo.On("UpdatePaymentOrderStatus", mock.Anything, "order123", "payment123", "enterprise123", mock.Anything).Return(nil)
				mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:      "Error al obtener la orden",
			orderID:   "order123",
			paymentID: "payment123",
			reason:    "Cliente solicitó cancelación",
			setupMocks: func(mockRepo *event_store.OrderEventRepository, mockAdapter *adapter.PaymentCaptureFlowAdapterIF, mockTransactionsRepo *repository.TransactionsRepositoryIF) {
				mockRepo.On("Get", mock.Anything, "order123", mock.Anything).Return(errors.New("error al obtener orden"))
			},
			expectedError: errors.New("error al obtener orden"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := event_store.NewOrderEventRepository(t)
			mockAdapter := adapter.NewPaymentCaptureFlowAdapterIF(t)
			mockTransactionsRepo := repository.NewTransactionsRepositoryIF(t)

			tt.setupMocks(mockRepo, mockAdapter, mockTransactionsRepo)

			useCase := NewPaymentReleaseUseCase(
				mockRepo,
				mockAdapter,
				mockTransactionsRepo,
			)

			err := useCase.ReleasePayment(context.Background(), tt.orderID, tt.paymentID, tt.reason)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockAdapter.AssertExpectations(t)
			mockTransactionsRepo.AssertExpectations(t)
		})
	}
}
