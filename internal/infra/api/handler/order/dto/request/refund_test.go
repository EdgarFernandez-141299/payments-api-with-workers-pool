package dto

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
)

func TestRefundDTO_Validate(t *testing.T) {
	tests := []struct {
		name    string
		dto     RefundDTO
		wantErr bool
	}{
		{
			name: "valid DTO with total refund",
			dto: RefundDTO{
				ReferenceOrderID: "order-123",
				PaymentOrderID:   "payment-456",
				Reason:           "customer request",
				Amount:           decimal.NewFromInt(100),
				IsTotal:          true,
			},
			wantErr: false,
		},
		{
			name: "valid DTO with partial refund",
			dto: RefundDTO{
				ReferenceOrderID: "order-123",
				PaymentOrderID:   "payment-456",
				Reason:           "partial refund",
				Amount:           decimal.NewFromInt(50),
				IsTotal:          false,
			},
			wantErr: false,
		},
		{
			name: "missing reference order ID",
			dto: RefundDTO{
				ReferenceOrderID: "",
				PaymentOrderID:   "payment-456",
				Reason:           "customer request",
				Amount:           decimal.NewFromInt(100),
				IsTotal:          true,
			},
			wantErr: true,
		},
		{
			name: "missing payment order ID",
			dto: RefundDTO{
				ReferenceOrderID: "order-123",
				PaymentOrderID:   "",
				Reason:           "customer request",
				Amount:           decimal.NewFromInt(100),
				IsTotal:          true,
			},
			wantErr: true,
		},
		{
			name: "missing reason",
			dto: RefundDTO{
				ReferenceOrderID: "order-123",
				PaymentOrderID:   "payment-456",
				Reason:           "",
				Amount:           decimal.NewFromInt(100),
				IsTotal:          true,
			},
			wantErr: true,
		},
		{
			name: "amount is negative",
			dto: RefundDTO{
				ReferenceOrderID: "order-123",
				PaymentOrderID:   "payment-456",
				Reason:           "customer request",
				Amount:           decimal.NewFromInt(-100),
				IsTotal:          false,
			},
			wantErr: true,
		},
		{
			name: "amount is zero",
			dto: RefundDTO{
				ReferenceOrderID: "order-123",
				PaymentOrderID:   "payment-456",
				Reason:           "customer request",
				Amount:           decimal.NewFromInt(0),
				IsTotal:          false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dto.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRefundDTO_Command(t *testing.T) {
	tests := []struct {
		name    string
		dto     RefundDTO
		want    command.RefundTotalCommand
		wantErr bool
	}{
		{
			name: "create command from valid DTO",
			dto: RefundDTO{
				ReferenceOrderID: "order-123",
				PaymentOrderID:   "payment-456",
				Reason:           "customer request",
				Amount:           decimal.NewFromInt(100),
				IsTotal:          true,
			},
			want: command.RefundTotalCommand{
				ReferenceOrderID: "order-123",
				PaymentOrderID:   "payment-456",
				Reason:           "customer request",
			},
			wantErr: false,
		},
		{
			name: "create command with different values",
			dto: RefundDTO{
				ReferenceOrderID: "order-789",
				PaymentOrderID:   "payment-012",
				Reason:           "product damaged",
				Amount:           decimal.NewFromInt(50),
				IsTotal:          false,
			},
			want: command.RefundTotalCommand{
				ReferenceOrderID: "order-789",
				PaymentOrderID:   "payment-012",
				Reason:           "product damaged",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.dto.Command()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.ReferenceOrderID, got.ReferenceOrderID)
			assert.Equal(t, tt.want.PaymentOrderID, got.PaymentOrderID)
			assert.Equal(t, tt.want.Reason, got.Reason)
		})
	}
}

func TestRefundDTO_CommandPartial(t *testing.T) {
	tests := []struct {
		name    string
		dto     RefundDTO
		want    command.CreatePartialPaymentRefundCommand
		wantErr bool
	}{
		{
			name: "create command from valid DTO",
			dto: RefundDTO{
				ReferenceOrderID: "order-123",
				PaymentOrderID:   "payment-456",
				Reason:           "customer request",
				Amount:           decimal.NewFromInt(50),
				IsTotal:          false,
			},
			want: command.CreatePartialPaymentRefundCommand{
				ReferenceOrderID: "order-123",
				PaymentOrderID:   "payment-456",
				Amount:           decimal.NewFromInt(50),
				Reason:           "customer request",
			},
			wantErr: false,
		},
		{
			name: "create command with different values",
			dto: RefundDTO{
				ReferenceOrderID: "order-789",
				PaymentOrderID:   "payment-012",
				Reason:           "product damaged",
				Amount:           decimal.NewFromInt(30),
				IsTotal:          false,
			},
			want: command.CreatePartialPaymentRefundCommand{
				ReferenceOrderID: "order-789",
				PaymentOrderID:   "payment-012",
				Amount:           decimal.NewFromInt(30),
				Reason:           "product damaged",
			},
			wantErr: false,
		},
		{
			name: "amount is negative",
			dto: RefundDTO{
				ReferenceOrderID: "order-123",
				PaymentOrderID:   "payment-456",
				Reason:           "customer request",
				Amount:           decimal.NewFromInt(-100),
				IsTotal:          false,
			},
			wantErr: true,
		},
		{
			name: "amount is zero",
			dto: RefundDTO{
				ReferenceOrderID: "order-123",
				PaymentOrderID:   "payment-456",
				Reason:           "customer request",
				Amount:           decimal.NewFromInt(0),
				IsTotal:          false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.dto.CommandPartial()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.ReferenceOrderID, got.ReferenceOrderID)
			assert.Equal(t, tt.want.PaymentOrderID, got.PaymentOrderID)
			assert.Equal(t, tt.want.Amount, got.Amount)
			assert.Equal(t, tt.want.Reason, got.Reason)
		})
	}
}
