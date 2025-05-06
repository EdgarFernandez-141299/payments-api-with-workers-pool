package aggregate

import (
	"testing"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/common/eventsourcing"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
	vo "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
)

const (
	testEnterpriseID = "enterprise1"
)

func TestCreate(t *testing.T) {
	userID := "user123"
	user := entities.NewUser(vo.NewUserType(vo.Member), userID)
	usdCurrencyCode, _ := vo.NewCurrencyCode("USD")
	totalAmountAsFloat := decimal.NewFromFloat(100.5)
	countryCode, _ := vo.NewCountryWithCode("MX")

	totalAmount, _ := vo.NewCurrencyAmount(usdCurrencyCode, totalAmountAsFloat)

	tests := []struct {
		name          string
		cmd           command.CreateOrderCommand
		expectedErr   error
		expectedOrder func() *Order
	}{
		{
			name: "valid command",
			cmd: command.CreateOrderCommand{
				ReferenceID:  "ref123",
				TotalAmount:  totalAmount,
				PhoneNumber:  "1234567890",
				User:         user,
				EnterpriseID: testEnterpriseID,
				CurrencyCode: usdCurrencyCode,
				CountryCode:  countryCode,
				BillingAddress: value_objects.Address{
					Street:  "123 Main St",
					City:    "Anytown",
					ZipCode: "12345",
					Country: countryCode,
				},
			},
			expectedErr: nil,
			expectedOrder: func() *Order {
				o := new(Order)
				_ = o.WithID("ref123")
				o.TotalAmount = totalAmount
				o.PhoneNumber = "1234567890"
				o.User = user
				o.ID = "ref123"
				o.CountryCode = countryCode
				o.EnterpriseID = testEnterpriseID
				return o
			},
		},
		{
			name: "invalid reference ID",
			cmd: command.CreateOrderCommand{
				ReferenceID:  "",
				TotalAmount:  totalAmount,
				PhoneNumber:  "1234567890",
				User:         user,
				EnterpriseID: testEnterpriseID,
				CurrencyCode: usdCurrencyCode,
				CountryCode:  countryCode,
				BillingAddress: value_objects.Address{
					Street:  "123 Main St",
					City:    "Anytown",
					ZipCode: "12345",
					Country: countryCode,
				},
			},
			expectedErr:   eventsourcing.InvalidAggregateID,
			expectedOrder: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := Create(tt.cmd)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Nil(t, order)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOrder().ID, order.ID)
				assert.Equal(t, tt.expectedOrder().TotalAmount, order.TotalAmount)
				assert.Equal(t, tt.expectedOrder().PhoneNumber, order.PhoneNumber)
				assert.Equal(t, tt.expectedOrder().User, order.User)
				assert.Equal(t, tt.expectedOrder().CountryCode, order.CountryCode)
				assert.Equal(t, tt.expectedOrder().EnterpriseID, order.EnterpriseID)
			}
		})
	}
}
