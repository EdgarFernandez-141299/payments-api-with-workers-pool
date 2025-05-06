package activities

import (
	"context"
	"errors"
	"testing"

	create2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/create"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/response"
)

func TestCreatePaymentOrder(t *testing.T) {
	testCases := []struct {
		name         string
		cmd          command.CreatePaymentOrderCommand
		setupMock    func() *create2.CreatePaymentOrderUseCaseIF
		expectedResp response.PaymentOrderResponseDTO
		expectedErr  error
	}{
		{
			name: "success",
			cmd: command.CreatePaymentOrderCommand{
				ReferenceOrderID: "ref-123",
				ID:               "id-123",
			},
			setupMock: func() *create2.CreatePaymentOrderUseCaseIF {
				mockUseCase := create2.NewCreatePaymentOrderUseCaseIF(t)
				mockUseCase.On("CreatePaymentOrder", mock.Anything, mock.Anything).
					Return(response.PaymentOrderResponseDTO{
						ReferenceOrderID: "ref-123",
					}, nil)

				return mockUseCase
			},
			expectedResp: response.PaymentOrderResponseDTO{
				ReferenceOrderID: "ref-123",
			},
			expectedErr: nil,
		},
		{
			name: "use case error",
			cmd: command.CreatePaymentOrderCommand{
				ReferenceOrderID: "ref-456",
				ID:               "id-456",
			},
			setupMock: func() *create2.CreatePaymentOrderUseCaseIF {
				mockUseCase := create2.NewCreatePaymentOrderUseCaseIF(t)
				mockUseCase.On("CreatePaymentOrder", mock.Anything, mock.Anything).
					Return(response.PaymentOrderResponseDTO{}, errors.New("use case error"))

				return mockUseCase
			},
			expectedResp: response.PaymentOrderResponseDTO{},
			expectedErr:  errors.New("use case error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := tc.setupMock()
			activity := NewCreatePaymentOrderActivity(m)

			resp, err := activity.CreatePaymentOrder(context.TODO(), tc.cmd)

			assert.Equal(t, tc.expectedResp, resp)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			m.AssertExpectations(t)
		})
	}
}
