package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaptureRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request CaptureRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: CaptureRequest{
				ReferenceOrderID: "123",
				PaymentID:        "456",
			},
			wantErr: false,
		},
		{
			name: "missing reference_order_id",
			request: CaptureRequest{
				PaymentID: "456",
			},
			wantErr: true,
			errMsg:  "reference_order_id is required",
		},
		{
			name: "missing payment_id",
			request: CaptureRequest{
				ReferenceOrderID: "123",
			},
			wantErr: true,
			errMsg:  "payment_id is required",
		},
		{
			name:    "missing both required fields",
			request: CaptureRequest{},
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
