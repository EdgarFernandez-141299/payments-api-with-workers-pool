package aggregate

import (
	"testing"
	"time"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/events"
	vo "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

type mockEvent struct {
	data interface{}
}

func (m *mockEvent) Data() interface{} {
	return m.data
}

func (m *mockEvent) Type() string {
	return "OrderCreated"
}

func (m *mockEvent) Version() int {
	return 1
}

func (m *mockEvent) Timestamp() time.Time {
	return time.Now()
}

func (m *mockEvent) AggregateID() string {
	return "order-1"
}

func TestOrder_CheckExistence(t *testing.T) {
	tests := []struct {
		name     string
		order    *Order
		expected bool
	}{
		{
			name: "does not exist when ID is empty",
			order: &Order{
				ID: "",
			},
			expected: true,
		},
		{
			name: "exists when ID is not empty",
			order: &Order{
				ID: "order-123",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.order.IsEmpty()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestOrder_CanProcessPayment(t *testing.T) {
	tests := []struct {
		name    string
		order   *Order
		payment entities.PaymentOrder
		want    bool
	}{
		{
			name: "can process valid payment",
			order: &Order{
				Status:      vo.OrderStatusProcessing(),
				TotalAmount: vo.CurrencyAmount{Value: decimal.NewFromFloat(100)},
			},
			payment: entities.PaymentOrder{
				ID:     "payment-1",
				Status: enums.PaymentProcessing,
				Total:  vo.CurrencyAmount{Value: decimal.NewFromFloat(50)},
			},
			want: true,
		},
		{
			name: "cannot process payment when order not in processing",
			order: &Order{
				Status:      vo.OrderStatusFailed(),
				TotalAmount: vo.CurrencyAmount{Value: decimal.NewFromFloat(100)},
			},
			payment: entities.PaymentOrder{
				ID:     "payment-1",
				Status: enums.PaymentProcessing,
				Total:  vo.CurrencyAmount{Value: decimal.NewFromFloat(50)},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.order.CanProcessPayment(tt.payment)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOrder_GetTotalAmount(t *testing.T) {
	order := &Order{
		Currency: vo.CurrencyCode{Code: "USD"},
		OrderPayments: []entities.PaymentOrder{
			{
				ID:     "payment-1",
				Status: enums.PaymentProcessed,
				Total:  vo.CurrencyAmount{Value: decimal.NewFromInt(50)},
			},
			{
				ID:     "payment-2",
				Status: enums.PaymentFailed,
				Total:  vo.CurrencyAmount{Value: decimal.NewFromInt(30)},
			},
		},
	}

	total := order.GetTotalEligiblePayments()
	assert.Equal(t, decimal.NewFromInt(50), total.Value)
}

func TestOrder_StartProcessingOrderPayment(t *testing.T) {
	associatedOrigin, _ := vo.NewFromAssociatedOriginString(enums.Club.String())
	paymentMethod := vo.NewCCPaymentMethod("card-123", "123")
	tests := []struct {
		name    string
		order   *Order
		cmd     command.CreatePaymentOrderCommand
		wantErr bool
	}{
		{
			name: "successful payment processing",
			order: &Order{
				ID:          "order-1",
				Status:      vo.OrderStatusProcessing(),
				TotalAmount: vo.CurrencyAmount{Value: decimal.NewFromFloat(100)},
			},
			cmd: command.CreatePaymentOrderCommand{
				ReferenceOrderID: "order-1",
				ID:               "payment-1",
				Payment: entities.PaymentOrder{
					ID:         "payment-1",
					Status:     enums.PaymentProcessing,
					Total:      vo.CurrencyAmount{Value: decimal.NewFromFloat(50)},
					OriginType: associatedOrigin,
					Method:     paymentMethod,
				},
				User: entities.User{ID: "user-1"},
			},
			wantErr: false,
		},
		{
			name: "invalid payment amount",
			order: &Order{
				ID:          "order-1",
				Status:      vo.OrderStatusProcessing(),
				TotalAmount: vo.CurrencyAmount{Value: decimal.NewFromFloat(100)},
			},
			cmd: command.CreatePaymentOrderCommand{
				ReferenceOrderID: "order-1",
				ID:               "payment-1",
				Payment: entities.PaymentOrder{
					ID:         "payment-1",
					Status:     enums.PaymentProcessing,
					Total:      vo.CurrencyAmount{Value: decimal.NewFromFloat(150)},
					OriginType: associatedOrigin,
					Method:     paymentMethod,
				},
				User: entities.User{ID: "user-1"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.order.StartProcessingOrderPayment(tt.cmd)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestOrder_OrderPaymentProcessed(t *testing.T) {
	order := &Order{
		ID: "order-1",
	}

	cmd := command.CreatePaymentOrderProcessedCommand{
		OrderID:           "order-1",
		PaymentID:         "payment-1",
		AuthorizationCode: "auth-123",
		OrderStatusString: "completed",
	}

	err := order.OrderPaymentProcessed(cmd)
	assert.NoError(t, err)
}

func TestOrder_OrderPaymentFailed(t *testing.T) {
	order := &Order{
		ID: "order-1",
	}

	cmd := command.CreatePaymentOrderFailCommand{
		OrderID:           "order-1",
		PaymentID:         "payment-1",
		PaymentReason:     "insufficient funds",
		OrderStatusString: "failed",
	}

	err := order.OrderPaymentFailed(cmd)
	assert.NoError(t, err)
}

func TestOrder_FindPaymentByID(t *testing.T) {
	usdCode := vo.CurrencyCode{Code: "USD"}
	successfulPayment := entities.PaymentOrder{
		ID: "successful-payment",
		Total: vo.CurrencyAmount{
			Value: decimal.NewFromFloat(50.0),
			Code:  usdCode,
		},
		Status: enums.PaymentProcessed,
	}
	failedPayment := entities.PaymentOrder{
		ID: "voided-payment",
		Total: vo.CurrencyAmount{
			Value: decimal.NewFromFloat(30.0),
			Code:  usdCode,
		},
		Status: enums.PaymentFailed,
	}

	tests := []struct {
		name          string
		order         Order
		paymentID     string
		expectedError bool
		expectedOrder entities.PaymentOrder
	}{
		{
			name: "find existing payment",
			order: Order{
				OrderPayments: []entities.PaymentOrder{successfulPayment, failedPayment},
			},
			paymentID:     "successful-payment",
			expectedError: false,
			expectedOrder: successfulPayment,
		},
		{
			name: "find voided payment",
			order: Order{
				OrderPayments: []entities.PaymentOrder{successfulPayment, failedPayment},
			},
			paymentID:     "voided-payment",
			expectedError: false,
			expectedOrder: failedPayment,
		},
		{
			name: "payment not found",
			order: Order{
				OrderPayments: []entities.PaymentOrder{successfulPayment},
			},
			paymentID:     "non-existent-payment",
			expectedError: true,
			expectedOrder: entities.PaymentOrder{},
		},
		{
			name: "empty payments list",
			order: Order{
				OrderPayments: []entities.PaymentOrder{},
			},
			paymentID:     "any-payment",
			expectedError: true,
			expectedOrder: entities.PaymentOrder{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.order.FindPaymentByID(tt.paymentID)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOrder, result)
			}
		})
	}
}

func TestOrder_WhenOrderCreated(t *testing.T) {
	order := &Order{
		ID: "order-1",
	}

	// Test OrderCreated event
	amount, _ := vo.NewCurrencyAmount(vo.CurrencyCode{Code: "USD"}, decimal.NewFromInt(100))
	countryCode, _ := vo.NewCountryWithCode("MX")

	orderCreatedEvent := events.NewOrderCreatedEventBuilder().
		SetID("order-1").
		SetTotalAmount(amount).
		SetPhoneNumber("1234567890").
		SetUser(entities.User{ID: "user-1"}).
		SetCreatedAt(time.Now()).
		SetCountryCode(countryCode).
		SetEnterpriseID("enterprise-1").
		SetEmail("test@example.com").
		Build()

	WhenOrderCreated(order, *orderCreatedEvent)

	assert.Equal(t, "order-1", order.ID)
	assert.Equal(t, vo.CurrencyCode{Code: "USD"}, order.Currency)
	assert.Equal(t, decimal.NewFromInt(100), order.TotalAmount.Value)
}

func TestOrder_RefundPayment(t *testing.T) {
	tests := []struct {
		name    string
		order   *Order
		cmd     command.RefundTotalCommand
		wantErr bool
	}{
		{
			name: "successful refund for processed payment",
			order: &Order{
				ID: "order-1",
				TotalAmount: vo.CurrencyAmount{
					Value: decimal.NewFromFloat(50.0),
					Code:  vo.CurrencyCode{Code: "USD"},
				},
				OrderPayments: []entities.PaymentOrder{
					{
						ID:     "payment-1",
						Status: enums.PaymentProcessed,
						Total: vo.CurrencyAmount{
							Value: decimal.NewFromFloat(50.0),
							Code:  vo.CurrencyCode{Code: "USD"},
						},
					},
				},
			},
			cmd: command.RefundTotalCommand{
				ReferenceOrderID: "order-1",
				PaymentOrderID:   "payment-1",
				Reason:           "customer request",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.order.RefundPayment(tt.cmd.PaymentOrderID, tt.cmd.Reason)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOrder_RefundPartialPayment(t *testing.T) {
	tests := []struct {
		name    string
		order   *Order
		cmd     command.CreatePartialPaymentRefundCommand
		wantErr bool
	}{
		{
			name: "successful partial refund",
			order: &Order{
				ID: "order-1",
				OrderPayments: []entities.PaymentOrder{
					{
						ID:     "payment-1",
						Status: enums.PaymentProcessed,
						Total:  vo.CurrencyAmount{Value: decimal.NewFromFloat(100)},
					},
				},
			},
			cmd: command.CreatePartialPaymentRefundCommand{
				ReferenceOrderID: "order-1",
				PaymentOrderID:   "payment-1",
				Amount:           decimal.NewFromFloat(50),
				Reason:           "customer request",
			},
			wantErr: false,
		},
		{
			name: "partial refund with invalid payment ID",
			order: &Order{
				ID: "order-1",
				OrderPayments: []entities.PaymentOrder{
					{
						ID:     "payment-1",
						Status: enums.PaymentProcessed,
						Total:  vo.CurrencyAmount{Value: decimal.NewFromFloat(100)},
					},
				},
			},
			cmd: command.CreatePartialPaymentRefundCommand{
				ReferenceOrderID: "order-1",
				PaymentOrderID:   "payment-2",
				Amount:           decimal.NewFromFloat(50),
				Reason:           "customer request",
			},
			wantErr: true,
		},
		{
			name: "partial refund with amount greater than payment",
			order: &Order{
				ID: "order-1",
				OrderPayments: []entities.PaymentOrder{
					{
						ID:     "payment-1",
						Status: enums.PaymentProcessed,
						Total:  vo.CurrencyAmount{Value: decimal.NewFromFloat(100)},
					},
				},
			},
			cmd: command.CreatePartialPaymentRefundCommand{
				ReferenceOrderID: "order-1",
				PaymentOrderID:   "payment-1",
				Amount:           decimal.NewFromFloat(150), // Greater than payment amount
				Reason:           "customer request",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.order.RefundPartialPayment(tt.cmd.PaymentOrderID, tt.cmd.Reason, tt.cmd.Amount)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOrder_OrderPaymentAuthorized(t *testing.T) {
	tests := []struct {
		name    string
		order   *Order
		cmd     command.CreatePaymentOrderAuthorizedCommand
		wantErr bool
	}{
		{
			name: "successful payment authorization",
			order: &Order{
				ID: "order-1",
				OrderPayments: []entities.PaymentOrder{
					{
						ID:     "payment-1",
						Status: enums.PaymentProcessing,
					},
				},
			},
			cmd: command.CreatePaymentOrderAuthorizedCommand{
				OrderID:           "order-1",
				PaymentID:         "payment-1",
				AuthorizationCode: "auth-123",
				PaymentCard: command.CardData{
					CardBrand: "VISA",
					CardLast4: "1234",
					CardType:  "CREDIT",
				},
			},
			wantErr: false,
		},
		{
			name: "authorization for non-existent payment",
			order: &Order{
				ID: "order-1",
				OrderPayments: []entities.PaymentOrder{
					{
						ID:     "payment-1",
						Status: enums.PaymentProcessing,
					},
				},
			},
			cmd: command.CreatePaymentOrderAuthorizedCommand{
				OrderID:           "order-1",
				PaymentID:         "payment-2",
				AuthorizationCode: "auth-123",
				PaymentCard: command.CardData{
					CardBrand: "VISA",
					CardLast4: "1234",
					CardType:  "CREDIT",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.order.OrderPaymentAuthorized(tt.cmd)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOrder_OrderPaymentReleased(t *testing.T) {
	tests := []struct {
		name    string
		order   *Order
		cmd     command.PaymentOrderReleasedCommand
		wantErr bool
	}{
		{
			name: "successful payment release",
			order: &Order{
				ID: "order-1",
				OrderPayments: []entities.PaymentOrder{
					{
						ID:     "payment-1",
						Status: enums.PaymentProcessing,
					},
				},
			},
			cmd: command.PaymentOrderReleasedCommand{
				OrderID:   "order-1",
				PaymentID: "payment-1",
				Reason:    "customer request",
			},
			wantErr: false,
		},
		{
			name: "release for non-existent payment",
			order: &Order{
				ID: "order-1",
				OrderPayments: []entities.PaymentOrder{
					{
						ID:     "payment-1",
						Status: enums.PaymentProcessing,
					},
				},
			},
			cmd: command.PaymentOrderReleasedCommand{
				OrderID:   "order-1",
				PaymentID: "payment-2",
				Reason:    "customer request",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.order.OrderPaymentReleased(tt.cmd)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOrder_OrderPaymentCaptured(t *testing.T) {
	tests := []struct {
		name    string
		order   *Order
		cmd     command.PaymentOrderCapturedCommand
		wantErr bool
	}{
		{
			name: "successful payment capture",
			order: &Order{
				ID: "order-1",
				OrderPayments: []entities.PaymentOrder{
					{
						ID:     "payment-1",
						Status: enums.PaymentProcessing,
					},
				},
			},
			cmd: command.PaymentOrderCapturedCommand{
				OrderID:   "order-1",
				PaymentID: "payment-1",
			},
			wantErr: false,
		},
		{
			name: "capture for non-existent payment",
			order: &Order{
				ID: "order-1",
				OrderPayments: []entities.PaymentOrder{
					{
						ID:     "payment-1",
						Status: enums.PaymentProcessing,
					},
				},
			},
			cmd: command.PaymentOrderCapturedCommand{
				OrderID:   "order-1",
				PaymentID: "payment-2",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.order.OrderPaymentCaptured(tt.cmd)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOrder_GetTotalProcessed(t *testing.T) {
	usdCode := vo.CurrencyCode{Code: "USD"}
	tests := []struct {
		name     string
		order    *Order
		expected vo.CurrencyAmount
	}{
		{
			name: "single processed payment",
			order: &Order{
				Currency: usdCode,
				OrderPayments: []entities.PaymentOrder{
					{
						ID:     "payment-1",
						Status: enums.PaymentProcessed,
						Total:  vo.CurrencyAmount{Value: decimal.NewFromInt(50), Code: usdCode},
					},
				},
			},
			expected: vo.CurrencyAmount{Value: decimal.NewFromInt(50), Code: usdCode},
		},
		{
			name: "multiple payments with different statuses",
			order: &Order{
				Currency: usdCode,
				OrderPayments: []entities.PaymentOrder{
					{
						ID:     "payment-1",
						Status: enums.PaymentProcessed,
						Total:  vo.CurrencyAmount{Value: decimal.NewFromInt(50), Code: usdCode},
					},
					{
						ID:     "payment-2",
						Status: enums.PaymentFailed,
						Total:  vo.CurrencyAmount{Value: decimal.NewFromInt(30), Code: usdCode},
					},
					{
						ID:     "payment-3",
						Status: enums.PaymentProcessing,
						Total:  vo.CurrencyAmount{Value: decimal.NewFromInt(20), Code: usdCode},
					},
				},
			},
			expected: vo.CurrencyAmount{Value: decimal.NewFromInt(50), Code: usdCode},
		},
		{
			name: "no processed payments",
			order: &Order{
				Currency: usdCode,
				OrderPayments: []entities.PaymentOrder{
					{
						ID:     "payment-1",
						Status: enums.PaymentFailed,
						Total:  vo.CurrencyAmount{Value: decimal.NewFromInt(50), Code: usdCode},
					},
				},
			},
			expected: vo.CurrencyAmount{Value: decimal.Zero, Code: usdCode},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.order.GetTotalProcessed()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestOrder_HasPaymentRefundable(t *testing.T) {
	tests := []struct {
		name     string
		order    *Order
		expected bool
	}{
		{
			name: "has refundable payment",
			order: &Order{
				OrderPayments: []entities.PaymentOrder{
					{
						ID:     "payment-1",
						Status: enums.PaymentProcessed,
					},
				},
			},
			expected: true,
		},
		{
			name: "has partially refunded payment",
			order: &Order{
				OrderPayments: []entities.PaymentOrder{
					{
						ID:     "payment-1",
						Status: enums.PartiallyRefunded,
					},
				},
			},
			expected: true,
		},
		{
			name: "has processing payment",
			order: &Order{
				OrderPayments: []entities.PaymentOrder{
					{
						ID:     "payment-1",
						Status: enums.PaymentProcessing,
					},
				},
			},
			expected: true,
		},
		{
			name: "no refundable payments",
			order: &Order{
				OrderPayments: []entities.PaymentOrder{
					{
						ID:     "payment-1",
						Status: enums.PaymentFailed,
					},
				},
			},
			expected: false,
		},
		{
			name:     "empty payments list",
			order:    &Order{OrderPayments: []entities.PaymentOrder{}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.order.HasPaymentRefundable()
			assert.Equal(t, tt.expected, result)
		})
	}
}
