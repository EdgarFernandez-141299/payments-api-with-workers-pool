package activities

import (
	"context"
	"errors"
	"gitlab.com/clubhub.ai1/go-libraries/eventsourcing"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository/fixture"
	"testing"
)

func Test_CheckOrderActivity_CheckOrderExistence(t *testing.T) {
	tests := []struct {
		name             string
		mockResponse     bool
		mockError        error
		referenceOrderID string
		wantResponse     bool
		wantError        bool
	}{
		{
			name:             "success - order exists",
			mockResponse:     true,
			mockError:        nil,
			referenceOrderID: "ref123",
			wantResponse:     true,
			wantError:        false,
		},
		{
			name:             "order does not exist - ErrAggregateNotFound",
			mockResponse:     false,
			mockError:        eventsourcing.ErrAggregateNotFound,
			referenceOrderID: "ref456",
			wantResponse:     false,
			wantError:        false,
		},
		{
			name:             "error - other error",
			mockResponse:     false,
			mockError:        errors.New("unexpected error"),
			referenceOrderID: "ref789",
			wantResponse:     false,
			wantError:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := fixture.OrderEventRepositoryGetFixture(
				t, tt.referenceOrderID,
				"",
				enums.PaymentProcessed,
				value_objects.CurrencyAmount{},
				value_objects.CurrencyAmount{},
				tt.mockError)

			activity := NewCheckOrderActivity(mockUseCase)

			result, err := activity.CheckOrderExistence(context.Background(), tt.referenceOrderID)

			if (err != nil) != tt.wantError {
				t.Errorf("unexpected error state: got %v, wantError: %v", err, tt.wantError)
			}

			if result != tt.wantResponse {
				t.Errorf("unexpected response: got %v, wantResponse: %v", result, tt.wantResponse)
			}

			mockUseCase.AssertExpectations(t)
		})
	}
}
