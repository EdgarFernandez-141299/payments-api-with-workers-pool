package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReleaseRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request ReleaseRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: ReleaseRequest{
				ReferenceOrderID: "123",
				PaymentID:        "456",
				Reason:           "test reason",
			},
			wantErr: false,
		},
		{
			name: "missing reference_order_id",
			request: ReleaseRequest{
				PaymentID: "456",
				Reason:    "test reason",
			},
			wantErr: true,
			errMsg:  "reference_order_id is required",
		},
		{
			name: "missing payment_id",
			request: ReleaseRequest{
				ReferenceOrderID: "123",
				Reason:           "test reason",
			},
			wantErr: true,
			errMsg:  "payment_id is required",
		},
		{
			name: "missing both required fields",
			request: ReleaseRequest{
				Reason: "test reason",
			},
			wantErr: true,
			errMsg:  "reference_order_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
