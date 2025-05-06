package use_cases

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/event_store"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/aggregate"
	order_entities "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/order/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_order/entities"
	payment_order_entities "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_order/entities"
)

func TestPostProcessingPaymentOrderUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	orderID := "test-order"
	paymentID := "test-payment"
	enterpriseID := "test-enterprise"

	tests := []struct {
		name        string
		cmd         PostProcessingPaymentOrderCommand
		mockSetup   func(t *testing.T) (*event_store.OrderEventRepository, *repository.OrderReadRepositoryIF, *repository.OrderWriteRepositoryIF, *repository.PaymentOrderRepositoryIF, *repository.GetPaymentOrderByReferenceIF)
		wantErr     bool
		expectedErr error
	}{
		{
			name: "successful payment with VISA card",
			cmd: PostProcessingPaymentOrderCommand{
				ReferenceOrderID:  orderID,
				PaymentID:         paymentID,
				Status:            enums.PaymentProcessed,
				AuthorizationCode: "auth123",
				OrderStatusString: "completed",
				PaymentReason:     "",
				PaymentCard: CardData{
					CardBrand: "VISA",
					CardLast4: "1234",
					CardType:  "credit",
				},
			},
			mockSetup: func(t *testing.T) (*event_store.OrderEventRepository, *repository.OrderReadRepositoryIF, *repository.OrderWriteRepositoryIF, *repository.PaymentOrderRepositoryIF, *repository.GetPaymentOrderByReferenceIF) {
				mockOrderRepo := new(event_store.OrderEventRepository)
				mockOrderRepo.On("Get", mock.Anything, orderID, mock.Anything).Run(func(args mock.Arguments) {
					order := args.Get(2).(*aggregate.Order)
					order.EnterpriseID = enterpriseID
				}).Return(nil)
				mockOrderRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

				mockOrderReadRepo := new(repository.OrderReadRepositoryIF)
				orderEntity := order_entities.NewOrderEntityBuilder().
					SetID("test-id").
					SetReferenceOrderID(orderID).
					SetEnterpriseID(enterpriseID).
					SetStatus("completed").
					Build()
				mockOrderReadRepo.On("GetOrderByReferenceID", mock.Anything, orderID, enterpriseID).Return(orderEntity, nil)

				mockOrderWriteRepo := new(repository.OrderWriteRepositoryIF)
				mockOrderWriteRepo.On("UpdateOrder", mock.Anything, mock.Anything).Return(nil)

				mockPaymentOrderRepo := new(repository.PaymentOrderRepositoryIF)
				mockPaymentOrderRepo.On("UpdatePaymentOrder", mock.Anything, mock.Anything).Return(nil)

				mockPaymentOrderReadRepo := new(repository.GetPaymentOrderByReferenceIF)
				paymentEntity := payment_order_entities.NewPaymentOrderEntity().
					SetOrderID(orderID).
					SetPaymentOrderID(paymentID).
					SetEnterpriseID(enterpriseID).
					SetPaymentFlow(enums.Capture).
					Build()
				mockPaymentOrderReadRepo.On("GetPaymentOrderByReference", mock.Anything, orderID, paymentID, enterpriseID).Return(paymentEntity, nil)

				return mockOrderRepo, mockOrderReadRepo, mockOrderWriteRepo, mockPaymentOrderRepo, mockPaymentOrderReadRepo
			},
			wantErr: false,
		},
		{
			name: "failed payment processing",
			cmd: PostProcessingPaymentOrderCommand{
				ReferenceOrderID:  orderID,
				PaymentID:         paymentID,
				Status:            enums.PaymentFailed,
				AuthorizationCode: "",
				OrderStatusString: "failed",
				PaymentReason:     "insufficient funds",
				PaymentCard: CardData{
					CardBrand: "VISA",
					CardLast4: "1234",
					CardType:  "credit",
				},
			},
			mockSetup: func(t *testing.T) (*event_store.OrderEventRepository, *repository.OrderReadRepositoryIF, *repository.OrderWriteRepositoryIF, *repository.PaymentOrderRepositoryIF, *repository.GetPaymentOrderByReferenceIF) {
				mockOrderRepo := new(event_store.OrderEventRepository)
				mockOrderRepo.On("Get", mock.Anything, orderID, mock.Anything).Run(func(args mock.Arguments) {
					order := args.Get(2).(*aggregate.Order)
					order.EnterpriseID = enterpriseID
				}).Return(nil)
				mockOrderRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

				mockOrderReadRepo := new(repository.OrderReadRepositoryIF)
				orderEntity := order_entities.NewOrderEntityBuilder().
					SetID("test-id").
					SetReferenceOrderID(orderID).
					SetEnterpriseID(enterpriseID).
					SetStatus("failed").
					Build()
				mockOrderReadRepo.On("GetOrderByReferenceID", mock.Anything, orderID, enterpriseID).Return(orderEntity, nil)

				mockOrderWriteRepo := new(repository.OrderWriteRepositoryIF)
				mockOrderWriteRepo.On("UpdateOrder", mock.Anything, mock.Anything).Return(nil)

				mockPaymentOrderRepo := new(repository.PaymentOrderRepositoryIF)
				mockPaymentOrderRepo.On("UpdatePaymentOrder", mock.Anything, mock.Anything).Return(nil)

				mockPaymentOrderReadRepo := new(repository.GetPaymentOrderByReferenceIF)
				paymentEntity := payment_order_entities.NewPaymentOrderEntity().
					SetOrderID(orderID).
					SetPaymentOrderID(paymentID).
					SetEnterpriseID(enterpriseID).
					SetPaymentFlow(enums.Capture).
					Build()
				mockPaymentOrderReadRepo.On("GetPaymentOrderByReference", mock.Anything, orderID, paymentID, enterpriseID).Return(paymentEntity, nil)

				return mockOrderRepo, mockOrderReadRepo, mockOrderWriteRepo, mockPaymentOrderRepo, mockPaymentOrderReadRepo
			},
			wantErr: false,
		},
		{
			name: "authorized payment processing",
			cmd: PostProcessingPaymentOrderCommand{
				ReferenceOrderID:  orderID,
				PaymentID:         paymentID,
				Status:            enums.PaymentAuthorized,
				AuthorizationCode: "auth123",
				OrderStatusString: "authorized",
				PaymentReason:     "",
				PaymentCard: CardData{
					CardBrand: "VISA",
					CardLast4: "1234",
					CardType:  "credit",
				},
			},
			mockSetup: func(t *testing.T) (*event_store.OrderEventRepository, *repository.OrderReadRepositoryIF, *repository.OrderWriteRepositoryIF, *repository.PaymentOrderRepositoryIF, *repository.GetPaymentOrderByReferenceIF) {
				mockOrderRepo := new(event_store.OrderEventRepository)
				mockOrderRepo.On("Get", mock.Anything, orderID, mock.Anything).Run(func(args mock.Arguments) {
					order := args.Get(2).(*aggregate.Order)
					order.EnterpriseID = enterpriseID
				}).Return(nil)
				mockOrderRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

				mockOrderReadRepo := new(repository.OrderReadRepositoryIF)
				orderEntity := order_entities.NewOrderEntityBuilder().
					SetID("test-id").
					SetReferenceOrderID(orderID).
					SetEnterpriseID(enterpriseID).
					SetStatus("authorized").
					Build()
				mockOrderReadRepo.On("GetOrderByReferenceID", mock.Anything, orderID, enterpriseID).Return(orderEntity, nil)

				mockOrderWriteRepo := new(repository.OrderWriteRepositoryIF)
				mockOrderWriteRepo.On("UpdateOrder", mock.Anything, mock.Anything).Return(nil)

				mockPaymentOrderRepo := new(repository.PaymentOrderRepositoryIF)
				mockPaymentOrderRepo.On("UpdatePaymentOrder", mock.Anything, mock.Anything).Return(nil)

				mockPaymentOrderReadRepo := new(repository.GetPaymentOrderByReferenceIF)
				paymentEntity := payment_order_entities.NewPaymentOrderEntity().
					SetOrderID(orderID).
					SetPaymentOrderID(paymentID).
					SetEnterpriseID(enterpriseID).
					SetPaymentFlow(enums.Capture).
					Build()
				mockPaymentOrderReadRepo.On("GetPaymentOrderByReference", mock.Anything, orderID, paymentID, enterpriseID).Return(paymentEntity, nil)

				return mockOrderRepo, mockOrderReadRepo, mockOrderWriteRepo, mockPaymentOrderRepo, mockPaymentOrderReadRepo
			},
			wantErr: false,
		},
		{
			name: "error getting order from repository",
			cmd: PostProcessingPaymentOrderCommand{
				ReferenceOrderID:  orderID,
				PaymentID:         paymentID,
				Status:            enums.PaymentProcessed,
				AuthorizationCode: "auth123",
				OrderStatusString: "completed",
				PaymentReason:     "",
				PaymentCard: CardData{
					CardBrand: "VISA",
					CardLast4: "1234",
					CardType:  "credit",
				},
			},
			mockSetup: func(t *testing.T) (*event_store.OrderEventRepository, *repository.OrderReadRepositoryIF, *repository.OrderWriteRepositoryIF, *repository.PaymentOrderRepositoryIF, *repository.GetPaymentOrderByReferenceIF) {
				mockOrderRepo := new(event_store.OrderEventRepository)
				mockOrderRepo.On("Get", mock.Anything, orderID, mock.Anything).Return(errors.New("order not found"))

				mockOrderReadRepo := new(repository.OrderReadRepositoryIF)
				mockOrderWriteRepo := new(repository.OrderWriteRepositoryIF)
				mockPaymentOrderRepo := new(repository.PaymentOrderRepositoryIF)
				mockPaymentOrderReadRepo := new(repository.GetPaymentOrderByReferenceIF)

				return mockOrderRepo, mockOrderReadRepo, mockOrderWriteRepo, mockPaymentOrderRepo, mockPaymentOrderReadRepo
			},
			wantErr:     true,
			expectedErr: errors.New("order not found"),
		},
		{
			name: "error updating payment order",
			cmd: PostProcessingPaymentOrderCommand{
				ReferenceOrderID:  orderID,
				PaymentID:         paymentID,
				Status:            enums.PaymentProcessed,
				AuthorizationCode: "auth123",
				OrderStatusString: "completed",
				PaymentReason:     "",
				PaymentCard: CardData{
					CardBrand: "VISA",
					CardLast4: "1234",
					CardType:  "credit",
				},
			},
			mockSetup: func(t *testing.T) (*event_store.OrderEventRepository, *repository.OrderReadRepositoryIF, *repository.OrderWriteRepositoryIF, *repository.PaymentOrderRepositoryIF, *repository.GetPaymentOrderByReferenceIF) {
				mockOrderRepo := new(event_store.OrderEventRepository)
				mockOrderRepo.On("Get", mock.Anything, orderID, mock.Anything).Run(func(args mock.Arguments) {
					order := args.Get(2).(*aggregate.Order)
					order.EnterpriseID = enterpriseID
				}).Return(nil)

				mockOrderReadRepo := new(repository.OrderReadRepositoryIF)
				mockOrderWriteRepo := new(repository.OrderWriteRepositoryIF)

				mockPaymentOrderRepo := new(repository.PaymentOrderRepositoryIF)
				mockPaymentOrderRepo.On("UpdatePaymentOrder", mock.Anything, mock.Anything).Return(errors.New("failed to update payment order"))

				mockPaymentOrderReadRepo := new(repository.GetPaymentOrderByReferenceIF)
				paymentEntity := payment_order_entities.NewPaymentOrderEntity().
					SetOrderID(orderID).
					SetPaymentOrderID(paymentID).
					SetEnterpriseID(enterpriseID).
					SetPaymentFlow(enums.Capture).
					Build()
				mockPaymentOrderReadRepo.On("GetPaymentOrderByReference", mock.Anything, orderID, paymentID, enterpriseID).Return(paymentEntity, nil)

				return mockOrderRepo, mockOrderReadRepo, mockOrderWriteRepo, mockPaymentOrderRepo, mockPaymentOrderReadRepo
			},
			wantErr:     true,
			expectedErr: errors.New("failed to update payment order"),
		},
		{
			name: "error updating order entity",
			cmd: PostProcessingPaymentOrderCommand{
				ReferenceOrderID:  orderID,
				PaymentID:         paymentID,
				Status:            enums.PaymentProcessed,
				AuthorizationCode: "auth123",
				OrderStatusString: "completed",
				PaymentReason:     "",
				PaymentCard: CardData{
					CardBrand: "VISA",
					CardLast4: "1234",
					CardType:  "credit",
				},
			},
			mockSetup: func(t *testing.T) (*event_store.OrderEventRepository, *repository.OrderReadRepositoryIF, *repository.OrderWriteRepositoryIF, *repository.PaymentOrderRepositoryIF, *repository.GetPaymentOrderByReferenceIF) {
				mockOrderRepo := new(event_store.OrderEventRepository)
				mockOrderRepo.On("Get", mock.Anything, orderID, mock.Anything).Run(func(args mock.Arguments) {
					order := args.Get(2).(*aggregate.Order)
					order.EnterpriseID = enterpriseID
				}).Return(nil)

				mockOrderReadRepo := new(repository.OrderReadRepositoryIF)
				orderEntity := order_entities.NewOrderEntityBuilder().
					SetID("test-id").
					SetReferenceOrderID(orderID).
					SetEnterpriseID(enterpriseID).
					SetStatus("completed").
					Build()
				mockOrderReadRepo.On("GetOrderByReferenceID", mock.Anything, orderID, enterpriseID).Return(orderEntity, nil)

				mockOrderWriteRepo := new(repository.OrderWriteRepositoryIF)
				mockOrderWriteRepo.On("UpdateOrder", mock.Anything, mock.Anything).Return(errors.New("failed to update order"))

				mockPaymentOrderRepo := new(repository.PaymentOrderRepositoryIF)
				mockPaymentOrderRepo.On("UpdatePaymentOrder", mock.Anything, mock.Anything).Return(nil)

				mockPaymentOrderReadRepo := new(repository.GetPaymentOrderByReferenceIF)
				paymentEntity := payment_order_entities.NewPaymentOrderEntity().
					SetOrderID(orderID).
					SetPaymentOrderID(paymentID).
					SetEnterpriseID(enterpriseID).
					SetPaymentFlow(enums.Capture).
					Build()
				mockPaymentOrderReadRepo.On("GetPaymentOrderByReference", mock.Anything, orderID, paymentID, enterpriseID).Return(paymentEntity, nil)

				return mockOrderRepo, mockOrderReadRepo, mockOrderWriteRepo, mockPaymentOrderRepo, mockPaymentOrderReadRepo
			},
			wantErr:     true,
			expectedErr: errors.New("failed to update order"),
		},
		{
			name: "error getting payment order by reference",
			cmd: PostProcessingPaymentOrderCommand{
				ReferenceOrderID:  orderID,
				PaymentID:         paymentID,
				Status:            enums.PaymentProcessed,
				AuthorizationCode: "auth123",
				OrderStatusString: "completed",
				PaymentReason:     "",
				PaymentCard: CardData{
					CardBrand: "VISA",
					CardLast4: "1234",
					CardType:  "credit",
				},
			},
			mockSetup: func(t *testing.T) (*event_store.OrderEventRepository, *repository.OrderReadRepositoryIF, *repository.OrderWriteRepositoryIF, *repository.PaymentOrderRepositoryIF, *repository.GetPaymentOrderByReferenceIF) {
				mockOrderRepo := new(event_store.OrderEventRepository)
				mockOrderRepo.On("Get", mock.Anything, orderID, mock.Anything).Run(func(args mock.Arguments) {
					order := args.Get(2).(*aggregate.Order)
					order.EnterpriseID = enterpriseID
				}).Return(nil)

				mockOrderReadRepo := new(repository.OrderReadRepositoryIF)
				mockOrderWriteRepo := new(repository.OrderWriteRepositoryIF)
				mockPaymentOrderRepo := new(repository.PaymentOrderRepositoryIF)
				mockPaymentOrderReadRepo := new(repository.GetPaymentOrderByReferenceIF)
				mockPaymentOrderReadRepo.On("GetPaymentOrderByReference", mock.Anything, orderID, paymentID, enterpriseID).Return(entities.PaymentOrderEntity{}, errors.New("payment order not found"))

				return mockOrderRepo, mockOrderReadRepo, mockOrderWriteRepo, mockPaymentOrderRepo, mockPaymentOrderReadRepo
			},
			wantErr:     true,
			expectedErr: errors.New("payment order not found"),
		},
		{
			name: "error saving order to event store",
			cmd: PostProcessingPaymentOrderCommand{
				ReferenceOrderID:  orderID,
				PaymentID:         paymentID,
				Status:            enums.PaymentProcessed,
				AuthorizationCode: "auth123",
				OrderStatusString: "completed",
				PaymentReason:     "",
				PaymentCard: CardData{
					CardBrand: "VISA",
					CardLast4: "1234",
					CardType:  "credit",
				},
			},
			mockSetup: func(t *testing.T) (*event_store.OrderEventRepository, *repository.OrderReadRepositoryIF, *repository.OrderWriteRepositoryIF, *repository.PaymentOrderRepositoryIF, *repository.GetPaymentOrderByReferenceIF) {
				mockOrderRepo := new(event_store.OrderEventRepository)
				mockOrderRepo.On("Get", mock.Anything, orderID, mock.Anything).Run(func(args mock.Arguments) {
					order := args.Get(2).(*aggregate.Order)
					order.EnterpriseID = enterpriseID
				}).Return(nil)
				mockOrderRepo.On("Save", mock.Anything, mock.Anything).Return(errors.New("failed to save order"))

				mockOrderReadRepo := new(repository.OrderReadRepositoryIF)
				orderEntity := order_entities.NewOrderEntityBuilder().
					SetID("test-id").
					SetReferenceOrderID(orderID).
					SetEnterpriseID(enterpriseID).
					SetStatus("completed").
					Build()
				mockOrderReadRepo.On("GetOrderByReferenceID", mock.Anything, orderID, enterpriseID).Return(orderEntity, nil)

				mockOrderWriteRepo := new(repository.OrderWriteRepositoryIF)
				mockOrderWriteRepo.On("UpdateOrder", mock.Anything, mock.Anything).Return(nil)

				mockPaymentOrderRepo := new(repository.PaymentOrderRepositoryIF)
				mockPaymentOrderRepo.On("UpdatePaymentOrder", mock.Anything, mock.Anything).Return(nil)

				mockPaymentOrderReadRepo := new(repository.GetPaymentOrderByReferenceIF)
				paymentEntity := payment_order_entities.NewPaymentOrderEntity().
					SetOrderID(orderID).
					SetPaymentOrderID(paymentID).
					SetEnterpriseID(enterpriseID).
					SetPaymentFlow(enums.Capture).
					Build()
				mockPaymentOrderReadRepo.On("GetPaymentOrderByReference", mock.Anything, orderID, paymentID, enterpriseID).Return(paymentEntity, nil)

				return mockOrderRepo, mockOrderReadRepo, mockOrderWriteRepo, mockPaymentOrderRepo, mockPaymentOrderReadRepo
			},
			wantErr:     true,
			expectedErr: errors.New("failed to save order"),
		},
		{
			name: "payment flow is not capture",
			cmd: PostProcessingPaymentOrderCommand{
				ReferenceOrderID:  orderID,
				PaymentID:         paymentID,
				Status:            enums.PaymentProcessed,
				AuthorizationCode: "auth123",
				OrderStatusString: "completed",
				PaymentReason:     "",
				PaymentCard: CardData{
					CardBrand: "VISA",
					CardLast4: "1234",
					CardType:  "credit",
				},
			},
			mockSetup: func(t *testing.T) (*event_store.OrderEventRepository, *repository.OrderReadRepositoryIF, *repository.OrderWriteRepositoryIF, *repository.PaymentOrderRepositoryIF, *repository.GetPaymentOrderByReferenceIF) {
				mockOrderRepo := new(event_store.OrderEventRepository)
				mockOrderRepo.On("Get", mock.Anything, orderID, mock.Anything).Run(func(args mock.Arguments) {
					order := args.Get(2).(*aggregate.Order)
					order.EnterpriseID = enterpriseID
				}).Return(nil)
				mockOrderRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

				mockOrderReadRepo := new(repository.OrderReadRepositoryIF)
				orderEntity := order_entities.NewOrderEntityBuilder().
					SetID("test-id").
					SetReferenceOrderID(orderID).
					SetEnterpriseID(enterpriseID).
					SetStatus("completed").
					Build()
				mockOrderReadRepo.On("GetOrderByReferenceID", mock.Anything, orderID, enterpriseID).Return(orderEntity, nil)

				mockOrderWriteRepo := new(repository.OrderWriteRepositoryIF)
				mockOrderWriteRepo.On("UpdateOrder", mock.Anything, mock.Anything).Return(nil)

				mockPaymentOrderRepo := new(repository.PaymentOrderRepositoryIF)
				mockPaymentOrderRepo.On("UpdatePaymentOrder", mock.Anything, mock.Anything).Return(nil)

				mockPaymentOrderReadRepo := new(repository.GetPaymentOrderByReferenceIF)
				paymentEntity := payment_order_entities.NewPaymentOrderEntity().
					SetOrderID(orderID).
					SetPaymentOrderID(paymentID).
					SetEnterpriseID(enterpriseID).
					SetPaymentFlow(enums.Autocapture).
					Build()
				mockPaymentOrderReadRepo.On("GetPaymentOrderByReference", mock.Anything, orderID, paymentID, enterpriseID).Return(paymentEntity, nil)

				return mockOrderRepo, mockOrderReadRepo, mockOrderWriteRepo, mockPaymentOrderRepo, mockPaymentOrderReadRepo
			},
			wantErr: false,
		},
		{
			name: "payment with Mastercard",
			cmd: PostProcessingPaymentOrderCommand{
				ReferenceOrderID:  orderID,
				PaymentID:         paymentID,
				Status:            enums.PaymentProcessed,
				AuthorizationCode: "auth123",
				OrderStatusString: "completed",
				PaymentReason:     "",
				PaymentCard: CardData{
					CardBrand: "MASTERCARD",
					CardLast4: "5678",
					CardType:  "credit",
				},
			},
			mockSetup: func(t *testing.T) (*event_store.OrderEventRepository, *repository.OrderReadRepositoryIF, *repository.OrderWriteRepositoryIF, *repository.PaymentOrderRepositoryIF, *repository.GetPaymentOrderByReferenceIF) {
				mockOrderRepo := new(event_store.OrderEventRepository)
				mockOrderRepo.On("Get", mock.Anything, orderID, mock.Anything).Run(func(args mock.Arguments) {
					order := args.Get(2).(*aggregate.Order)
					order.EnterpriseID = enterpriseID
				}).Return(nil)
				mockOrderRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

				mockOrderReadRepo := new(repository.OrderReadRepositoryIF)
				orderEntity := order_entities.NewOrderEntityBuilder().
					SetID("test-id").
					SetReferenceOrderID(orderID).
					SetEnterpriseID(enterpriseID).
					SetStatus("completed").
					Build()
				mockOrderReadRepo.On("GetOrderByReferenceID", mock.Anything, orderID, enterpriseID).Return(orderEntity, nil)

				mockOrderWriteRepo := new(repository.OrderWriteRepositoryIF)
				mockOrderWriteRepo.On("UpdateOrder", mock.Anything, mock.Anything).Return(nil)

				mockPaymentOrderRepo := new(repository.PaymentOrderRepositoryIF)
				mockPaymentOrderRepo.On("UpdatePaymentOrder", mock.Anything, mock.Anything).Return(nil)

				mockPaymentOrderReadRepo := new(repository.GetPaymentOrderByReferenceIF)
				paymentEntity := payment_order_entities.NewPaymentOrderEntity().
					SetOrderID(orderID).
					SetPaymentOrderID(paymentID).
					SetEnterpriseID(enterpriseID).
					SetPaymentFlow(enums.Capture).
					Build()
				mockPaymentOrderReadRepo.On("GetPaymentOrderByReference", mock.Anything, orderID, paymentID, enterpriseID).Return(paymentEntity, nil)

				return mockOrderRepo, mockOrderReadRepo, mockOrderWriteRepo, mockPaymentOrderRepo, mockPaymentOrderReadRepo
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrderRepo, mockOrderReadRepo, mockOrderWriteRepo, mockPaymentOrderRepo, mockPaymentOrderReadRepo := tt.mockSetup(t)
			useCase := NewPostProcessingPaymentOrderUseCase(
				mockOrderRepo,
				mockOrderReadRepo,
				mockOrderWriteRepo,
				mockPaymentOrderRepo,
				mockPaymentOrderReadRepo,
			)
			paymentFlow, err := useCase.PostProcessPaymentOrder(ctx, tt.cmd)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, enums.Autocapture, paymentFlow)
				if tt.expectedErr != nil {
					assert.EqualError(t, err, tt.expectedErr.Error())
				}
			} else {
				assert.NoError(t, err)
				if tt.name == "payment flow is not capture" {
					assert.Equal(t, enums.Autocapture, paymentFlow)
				} else {
					assert.Equal(t, enums.Capture, paymentFlow)
				}
				mockOrderRepo.AssertExpectations(t)
				mockOrderReadRepo.AssertExpectations(t)
				mockOrderWriteRepo.AssertExpectations(t)
				mockPaymentOrderRepo.AssertExpectations(t)
				mockPaymentOrderReadRepo.AssertExpectations(t)
			}
		})
	}
}
