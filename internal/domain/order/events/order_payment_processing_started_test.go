package events

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

func TestFromProcessOrderCommand(t *testing.T) {
	t.Run("create payment order event from command", func(t *testing.T) {
		associatedOrign, _ := value_objects.NewFromAssociatedOriginString(enums.Downpayment.String())
		cmd := command.CreatePaymentOrderCommand{
			ReferenceOrderID: "123",
			Payment: entities.PaymentOrder{
				ID:         "123",
				OriginType: associatedOrign,
				Total: value_objects.CurrencyAmount{
					Code:  value_objects.CurrencyCode{Code: "USD"},
					Value: decimal.NewFromFloat(100.0),
				},
			},
		}

		event := FromProcessOrderCommand(cmd)

		assert.Equal(t, event.PaymentOrder, event.PaymentOrder)
		assert.Equal(t, event.PaymentOrder.ID, event.PaymentOrder.ID)
	})
}
