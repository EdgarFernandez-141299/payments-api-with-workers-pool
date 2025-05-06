package activities

import (
	"context"
	"errors"
	"testing"
	"time"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/event_store"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository/fixture"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	errors2 "gitlab.com/clubhub.ai1/go-libraries/saga/errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/aggregate"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/payment_receipt/command"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/use_case"
)

func TestGeneratePaymentReceiptActivity_GenerateReceipt(t *testing.T) {
	referenceID := "order123"
	paymentID := "payment123"
	userID := "user123"
	enterpriseID := "enterprise123"
	email := "test@example.com"
	countryCode := "US"
	currencyCode := "USD"
	paymentAmount := decimal.NewFromInt(1000)
	paymentStatus := enums.PaymentProcessed
	paymentMethod := value_objects.PaymentMethod{
		Type: "CARD",
	}
	paymentDate := time.Now()

	country, _ := value_objects.NewCountryWithCode(countryCode)
	currency, _ := value_objects.NewCurrencyCode(currencyCode)
	currencyAmount, _ := value_objects.NewCurrencyAmount(currency, paymentAmount)

	tests := []struct {
		name          string
		setupMocks    func(*testing.T) (*event_store.OrderEventRepository, *use_case.GenerateReceiptPaymentUseCase)
		expectedError error
	}{
		{
			name: "generate receipt success",
			setupMocks: func(t *testing.T) (*event_store.OrderEventRepository, *use_case.GenerateReceiptPaymentUseCase) {
				mockRepo := event_store.NewOrderEventRepository(t)

				mockRepo.On("Get", mock.Anything, referenceID, mock.AnythingOfType("*aggregate.Order")).
					Run(func(args mock.Arguments) {
						order := args.Get(2).(*aggregate.Order)
						order.ID = referenceID
						order.User = entities.User{ID: userID}
						order.EnterpriseID = enterpriseID
						order.Email = email
						order.CountryCode = country
						order.Currency = currency
						order.OrderPayments = []entities.PaymentOrder{{
							ID:        paymentID,
							Status:    paymentStatus,
							Total:     currencyAmount,
							Method:    paymentMethod,
							CreatedAt: paymentDate,
						}}
					}).Return(nil)

				mockUseCase := use_case.NewGenerateReceiptPaymentUseCase(t)
				expectedCmd := command.CreatePaymentReceiptCommand{
					UserID:           userID,
					EnterpriseID:     enterpriseID,
					Email:            email,
					ReferenceOrderID: referenceID,
					PaymentID:        paymentID,
					PaymentStatus:    string(paymentStatus),
					PaymentAmount:    currencyAmount,
					PaymentCountry:   country,
					PaymentMethod:    paymentMethod,
					PaymentDate:      paymentDate.Format(time.RFC3339),
				}
				mockUseCase.On("Generate", mock.Anything, expectedCmd).Return(entities.PaymentReceipt{}, nil)

				return mockRepo, mockUseCase
			},
			expectedError: nil,
		},
		{
			name: "repository error",
			setupMocks: func(t *testing.T) (*event_store.OrderEventRepository, *use_case.GenerateReceiptPaymentUseCase) {
				mockRepo := fixture.NewOrderEventRepositoryFixtureBuilder(t).
					WithReferenceOrderID(referenceID).
					WithError(errors.New("repository error")).
					Build()

				mockUseCase := use_case.NewGenerateReceiptPaymentUseCase(t)

				return mockRepo, mockUseCase
			},
			expectedError: errors2.WrapActivityError(errors.New("repository error")),
		},
		{
			name: "payment not found",
			setupMocks: func(t *testing.T) (*event_store.OrderEventRepository, *use_case.GenerateReceiptPaymentUseCase) {
				mockRepo := fixture.NewOrderEventRepositoryFixtureBuilder(t).
					WithReferenceOrderID(referenceID).
					Build()

				mockRepo.On("Get", mock.Anything, referenceID, mock.AnythingOfType("*aggregate.Order")).
					Run(func(args mock.Arguments) {
						order := args.Get(2).(*aggregate.Order)
						order.ID = referenceID
						order.OrderPayments = []entities.PaymentOrder{} // Empty payments
					}).Return(nil)

				mockUseCase := use_case.NewGenerateReceiptPaymentUseCase(t)

				return mockRepo, mockUseCase
			},
			expectedError: errors2.WrapActivityError(errors.New("payment not found")),
		},
		{
			name: "use case error",
			setupMocks: func(t *testing.T) (*event_store.OrderEventRepository, *use_case.GenerateReceiptPaymentUseCase) {
				mockRepo := event_store.NewOrderEventRepository(t)

				mockRepo.On("Get", mock.Anything, referenceID, mock.AnythingOfType("*aggregate.Order")).
					Run(func(args mock.Arguments) {
						order := args.Get(2).(*aggregate.Order)
						order.ID = referenceID
						order.User = entities.User{ID: userID}
						order.EnterpriseID = enterpriseID
						order.Email = email
						order.CountryCode = country
						order.Currency = currency
						order.OrderPayments = []entities.PaymentOrder{{
							ID:        paymentID,
							Status:    paymentStatus,
							Total:     currencyAmount,
							Method:    paymentMethod,
							CreatedAt: paymentDate,
						}}
					}).Return(nil)

				mockUseCase := use_case.NewGenerateReceiptPaymentUseCase(t)
				mockUseCase.On("Generate", mock.Anything, mock.Anything).Return(entities.PaymentReceipt{}, errors.New("use case error"))

				return mockRepo, mockUseCase
			},
			expectedError: errors2.WrapActivityError(errors.New("use case error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo, mockUseCase := tt.setupMocks(t)

			activity := NewGeneratePaymentReceiptActivity(mockUseCase, mockRepo)
			receipt, err := activity.GenerateReceipt(context.Background(), referenceID, paymentID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Empty(t, receipt)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockUseCase.AssertExpectations(t)
		})
	}
}
