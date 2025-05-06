package activities

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	postpayment "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/use_cases/post_payment"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
)

type MockPostProcessingPaymentOrderUseCase struct {
	mock.Mock
}

func (m *MockPostProcessingPaymentOrderUseCase) PostProcessPaymentOrder(ctx context.Context, cmd postpayment.PostProcessingPaymentOrderCommand) (enums.PaymentFlowEnum, error) {
	args := m.Called(ctx, cmd)
	return args.Get(0).(enums.PaymentFlowEnum), args.Error(1)
}

func TestPostProcessingPaymentOrderActivity_PostProcessingPaymentOrder(t *testing.T) {
	tests := []struct {
		name          string
		cmd           postpayment.PostProcessingPaymentOrderCommand
		setupMock     func(*MockPostProcessingPaymentOrderUseCase)
		expectedFlow  enums.PaymentFlowEnum
		expectedError error
	}{
		{
			name: "payment processing success",
			cmd: postpayment.PostProcessingPaymentOrderCommand{
				ReferenceOrderID:  "order123",
				PaymentID:         "payment123",
				Status:            enums.PaymentProcessed,
				AuthorizationCode: "auth123",
				OrderStatusString: "completed",
				PaymentReason:     "success",
			},
			setupMock: func(m *MockPostProcessingPaymentOrderUseCase) {
				m.On("PostProcessPaymentOrder", mock.Anything, mock.MatchedBy(func(cmd postpayment.PostProcessingPaymentOrderCommand) bool {
					return cmd.ReferenceOrderID == "order123" &&
						cmd.PaymentID == "payment123" &&
						cmd.Status == enums.PaymentProcessed &&
						cmd.AuthorizationCode == "auth123" &&
						cmd.OrderStatusString == "completed" &&
						cmd.PaymentReason == "success"
				})).Return(enums.Autocapture, nil)
			},
			expectedFlow:  enums.Autocapture,
			expectedError: nil,
		},
		{
			name: "fail processing payment",
			cmd: postpayment.PostProcessingPaymentOrderCommand{
				ReferenceOrderID:  "order123",
				PaymentID:         "payment123",
				Status:            enums.PaymentFailed,
				AuthorizationCode: "",
				OrderStatusString: "failed",
				PaymentReason:     "insufficient funds",
			},
			setupMock: func(m *MockPostProcessingPaymentOrderUseCase) {
				m.On("PostProcessPaymentOrder", mock.Anything, mock.MatchedBy(func(cmd postpayment.PostProcessingPaymentOrderCommand) bool {
					return cmd.ReferenceOrderID == "order123" &&
						cmd.PaymentID == "payment123" &&
						cmd.Status == enums.PaymentFailed &&
						cmd.OrderStatusString == "failed" &&
						cmd.PaymentReason == "insufficient funds"
				})).Return(enums.Capture, errors.New("processing failed"))
			},
			expectedFlow:  enums.Capture,
			expectedError: errors.New("processing failed"),
		},
		{
			name: "validation error",
			cmd: postpayment.PostProcessingPaymentOrderCommand{
				ReferenceOrderID:  "",
				PaymentID:         "",
				Status:            enums.PaymentProcessing,
				AuthorizationCode: "",
				OrderStatusString: "",
				PaymentReason:     "",
			},
			setupMock: func(m *MockPostProcessingPaymentOrderUseCase) {
				m.On("PostProcessPaymentOrder", mock.Anything, mock.MatchedBy(func(cmd postpayment.PostProcessingPaymentOrderCommand) bool {
					return cmd.ReferenceOrderID == "" &&
						cmd.PaymentID == "" &&
						cmd.Status == enums.PaymentProcessing
				})).Return(enums.Capture, errors.New("invalid command"))
			},
			expectedFlow:  enums.Capture,
			expectedError: errors.New("invalid command"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockPostProcessingPaymentOrderUseCase)
			tt.setupMock(mockUseCase)

			activity := NewPostProcessingPaymentOrderActivity(mockUseCase)

			flow, err := activity.PostProcessingPaymentOrder(context.Background(), tt.cmd)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Equal(t, tt.expectedFlow, flow)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedFlow, flow)
			}

			mockUseCase.AssertExpectations(t)
		})
	}
}
