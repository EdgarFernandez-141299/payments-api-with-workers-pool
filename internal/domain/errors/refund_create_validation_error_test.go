package errors

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

func TestNewRefundCreateValidationError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantMsg    string
		wantFields map[string]interface{}
	}{
		{
			name: "valid_error",
			args: args{
				err: errors.New("validation failed"),
			},
			wantErr:    true,
			wantMsg:    "REFUND_CREATE_VALIDATION_ERROR",
			wantFields: map[string]interface{}{},
		},
		{
			name: "nil_error",
			args: args{
				err: nil,
			},
			wantErr:    true,
			wantMsg:    "REFUND_CREATE_VALIDATION_ERROR",
			wantFields: map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewRefundCreateValidationError(tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRefundCreateValidationError() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				return
			}

			if err == nil {
				t.Fatalf("NewRefundCreateValidationError() returned nil error, but wantErr is true")
				return
			}

			businessErr, ok := err.(*domain.BusinessError)
			if !ok {
				t.Fatalf("error is not of type *domain.BusinessError, got %T", err)
			}

			assert.ErrorContains(t, businessErr, tt.wantMsg)
		})
	}
}
