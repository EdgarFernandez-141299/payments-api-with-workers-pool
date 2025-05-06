package mockresponse

import (
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/resources/dto/response"
)

const SuccessOrderResponseMock = `{
  "order": {
    "billing_address": {
      "additional_description": "Descripción adicional",
      "address_type": "home",
      "address1": "Av. Eloy Alfaro 14, Quito 170515, Ecuador",
      "address2": "Av. Eloy Alfaro 14, Quito 170515, Ecuador",
      "city": "Quito",
      "country": "ECU",
      "created_at": "2021-11-03T22:09:09.086990957Z",
      "first_name": "John",
      "id": 1868,
      "identity_document": "1150218418",
      "identity_document_type": "DNI",
      "is_default": false,
      "last_name": "Doe",
      "lat": -0.100032,
      "lng": -78.46956,
      "phone": "593999999999",
      "state_name": "",
      "updated_at": "2021-11-03T22:09:09.087014623Z",
      "user_id": "xxxxx4e2-xxxx-xxxx-xxxx-xxxxx5b7b2e",
      "zipcode": "170515"
    },
    "created_at": "2021-11-03T22:09:09.086990957Z",
    "currency": "USD",
    "display_items_total_amount": "USD 2.99",
    "display_sub_total_amount": "USD 2.99",
    "display_tax_amount": "USD 10.00",
    "display_total_amount": "USD 29.90",
    "display_total_tax_amount": "USD 0.00",
    "items": [
      {
        "brand": "DEUNA Swagstore",
        "category": "booking",
        "description": "",
        "id": "ID",
        "isbn": "ISBN",
        "name": "Booking",
        "quantity": 1,
        "tax_amount": {
          "amount": 754,
          "currency": "USD",
          "currency_symbol": "$"
        },
        "taxable": false,
        "total_amount": {
          "amount": 299,
          "currency": "USD",
          "currency_symbol": "$"
        },
        "type": "Booking",
        "unit_price": {
          "amount": 299,
          "currency": "USD",
          "currency_symbol": "$"
        }
      }
    ],
    "items_total_amount": 299,
    "metadata": {
      "key1": "value1",
      "key2": "value2"
    },
    "order_id": "TESTS-12345",
    "payer_info": null,
    "payment": {
      "data": {
        "amount": {
          "amount": 0,
          "currency": "USD"
        },
        "created_at": "2021-11-03T22:09:09.086990957Z",
        "customer": {
          "email": "email@test.com",
          "id": "12345"
        },
        "from_card": {
          "card_brand": "Visa",
          "first_six": "414141",
          "last_four": "4141"
        },
        "id": "12345",
        "merchant": {
          "id": "12345",
          "store_code": "all"
        },
        "method_type": "credit_card",
        "processor": "",
        "reason": "",
        "routing_strategy": "",
        "status": "pending",
        "updated_at": "0001-01-01 00:00:00 +0000 UTC"
      }
    },
    "shipping_address": {
      "additional_description": "Descripción adicional",
      "address_type": "home",
      "address1": "Av. Eloy Alfaro 14, Quito 170515, Ecuador",
      "address2": "Av. Eloy Alfaro 14, Quito 170515, Ecuador",
      "city": "Quito",
      "country": "ECU",
      "created_at": "2021-11-03T22:09:09.086990957Z",
      "first_name": "John",
      "id": 1868,
      "identity_document": "",
      "identity_document_type": "",
      "is_default": false,
      "last_name": "Doe",
      "lat": -0.100032,
      "lng": -78.46956,
      "phone": "593999999999",
      "state_name": "",
      "updated_at": "2021-11-03T22:09:09.087014623Z",
      "user_id": "xxxxx4e2-xxxx-xxxx-xxxx-xxxxx5b7b2e",
      "zipcode": "170515"
    },
    "status": "pending",
    "store_code": "all",
    "sub_total": 299,
    "tax_amount": 10,
    "timezone": "",
    "total_amount": 2990,
    "total_tax_amount": 0,
    "updated_at": "",
    "user_id": "xxxxx4e2-xxxx-xxxx-xxxx-xxxxx5b7b2e",
    "webhook_urls": {
      "notify_order": "https://webhook.site/xxxxx4e2-xxxx-xxxx-xxxx-xxxxx5b7b2e"
    }
  },
  "token": "c5c5abc9-67a7-4a44-aa83-ee7e81d7052c"
}`

const FailOrderResponseMock = `{
  "data": null,
  "error": {
    "code": "EMA-3002",
    "description": "cannot create order"
  }
}`

const OrderNotFoundResponseMock = `{
  "data": null,
  "error": {
    "code": "EMA-2000",
    "description": "cannot parse token"
  }
}`

const FailOrderResponseMockWithInvalidPaymentMethod = `{
  "data": null,
  "error": {
    "code": "EMA-6004",
    "description": "the provided payment method is not configured for merchant and store: <paymet_method_name>"
  }
}`

var ExpectedSuccessResponse = response.DeunaOrderResponseDTO{
	Token: "c5c5abc9-67a7-4a44-aa83-ee7e81d7052c",
}
