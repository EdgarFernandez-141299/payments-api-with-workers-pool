package mockrequest

import (
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/resources/dto/request"
)

const CreateOrderRequestMock = `{
  "order": {
    "billing_address": {
      "address_type": "home",
      "country": "ECU",
      "is_default": false
    },
    "currency": "USD",
    "items": [
      {
        "tax_amount": {
          "currency": "USD",
          "currency_symbol": "$"
        },
        "taxable": true,
        "total_amount": {
          "currency": "USD",
          "currency_symbol": "$"
        },
        "unit_price": {
          "currency_symbol": "$"
        }
      }
    ],
    "metadata": {
      "foo": "bar"
    },
    "shipping_address": {
      "is_default": false
    }
  },
  "order_type": "DEUNA_CHECKOUT"
}`

var CreateOrderRequest = request.CreateDeunaOrderRequestDTO{
	Order:     request.DeunaOrder{},
	OrderType: "DEUNA_CHECKOUT",
}
