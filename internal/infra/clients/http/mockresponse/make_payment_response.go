package mockresponse

const SuccessOrderPaymentResponseMock = `{
	"order": {
		"cash_change": 0,
		"currency": "USD",
		"discounts": [],
		"display_items_total_amount": "$100.00",
		"display_shipping_amount": "$10.00",
		"display_sub_total": "$100.00",
		"display_tax_amount": "$15.00",
		"display_total_amount": "$125.00",
		"display_total_discount": "$0.00",
		"gift_card": [],
		"items": [
			{
				"brand": "Nike",
				"category": "Shoes",
				"color": "Black",
				"description": "Running shoes",
				"details_url": "https://example.com/product/123",
				"discounts": [],
				"id": "item-123",
				"image_url": "https://example.com/images/shoes.jpg",
				"isbn": "123456789",
				"manufacturer": "Nike Inc",
				"name": "Air Max",
				"options": "Size 10",
				"quantity": 1,
				"size": "10",
				"sku": "SKU123",
				"tax_amount": {
					"amount": 15,
					"currency": "USD",
					"currency_symbol": "$",
					"display_amount": "$15.00"
				},
				"taxable": true,
				"total_amount": {
					"amount": 100,
					"currency": "USD",
					"currency_symbol": "$",
					"display_amount": "$100.00"
				},
				"type": "physical",
				"unit_price": {
					"amount": 100,
					"currency": "USD",
					"currency_symbol": "$",
					"display_amount": "$100.00"
				},
				"uom": "each",
				"upc": "123456789",
				"weight": {
					"unit": "kg",
					"weight": 0.5
				}
			}
		],
		"items_total_amount": 100,
		"metadata": {
			"custom_field": "custom_value"
		},
		"order_id": "123456",
		"payment": {
			"data": {
				"amount": {
					"amount": 125,
					"currency": "USD",
					"currency_symbol": "$",
					"display_amount": "$125.00"
				},
				"authorization_3ds": {
					"html_content": "",
					"url_challenge": "",
					"version": "2.0"
				},
				"authorization_code": "AUTH123",
				"created_at": "2023-01-01T12:00:00Z",
				"customer": {
					"email": "test@example.com",
					"id": "cust-123"
				},
				"from_card": {
					"card_brand": "VISA",
					"first_six": "411111",
					"last_four": "1111"
				},
				"id": "payment-123",
				"installments": {
					"installment_amount": 125,
					"installment_rate": 0,
					"installment_type": "regular",
					"installments": 1,
					"plan_id": "plan-123",
					"plan_option_id": "option-123"
				},
				"merchant": {
					"id": "merch-123",
					"store_code": "store-123"
				},
				"metadata": {
					"custom_field": "custom_value"
				},
				"method_type": "card",
				"processor": "PaymentProcessor",
				"reason": "",
				"status": "approved",
				"updated_at": "2023-01-01T12:00:00Z"
			}
		},
		"redirect_url": "",
		"scheduled_at": "",
		"shipping": null,
		"shipping_address": {
			"additional_description": "",
			"address_type": "shipping",
			"address1": "123 Main St",
			"address2": "Apt 4B",
			"city": "New York",
			"country": "US",
			"created_at": "2023-01-01T12:00:00Z",
			"first_name": "John",
			"id": 12345,
			"identity_document": "123456789",
			"identity_document_type": "dni",
			"is_default": true,
			"last_name": "Doe",
			"lat": 40.7128,
			"lng": -74.006,
			"phone": "+1234567890",
			"state_name": "NY",
			"updated_at": "2023-01-01T12:00:00Z",
			"user_id": "user-123",
			"zipcode": "10001"
		},
		"shipping_amount": 10,
		"shipping_method": null,
		"shipping_methods": [],
		"shipping_options": null,
		"status": "completed",
		"store_code": "store-123",
		"sub_total": 100,
		"tax_amount": 15,
		"timezone": "America/New_York",
		"total_amount": 125,
		"total_discount": 0,
		"user_instructions": "",
		"webhook_urls": null
	},
	"order_token": "c5c5abc9-67a7-4a44-aa83-ee7e81d7052c"
}`

const SuccessOrderPaymentV2ResponseMock = `{
	"order": {
		"currency": "USD",
		"display_items_total_amount": "$100.00",
		"display_shipping_amount": "$10.00",
		"display_sub_total": "$100.00",
		"display_tax_amount": "$15.00",
		"display_total_amount": "$125.00",
		"display_total_discount": "$0.00",
		"items": [
			{
				"category": "Shoes",
				"description": "Running shoes",
				"id": "item-123",
				"manufacturer": "Nike Inc",
				"name": "Air Max",
				"options": "Size 10",
				"quantity": 1,
				"tax_amount": {
					"amount": 15,
					"currency": "USD",
					"currency_symbol": "$"
				},
				"taxable": true,
				"total_amount": {
					"amount": 100,
					"currency": "USD",
					"currency_symbol": "$"
				},
				"type": "physical",
				"unit_price": {
					"amount": 100,
					"currency": "USD",
					"currency_symbol": "$"
				}
			}
		],
		"items_total_amount": 100,
		"metadata": {
			"custom_field": "custom_value"
		},
		"order_id": "123456",
		"payment": {
			"data": {
				"amount": {
					"amount": 125,
					"currency": "USD"
				},
				"authorization_code": "AUTH123",
				"created_at": "2023-01-01T12:00:00Z",
				"customer": {
					"email": "test@example.com",
					"id": "cust-123"
				},
				"from_card": {
					"card_brand": "VISA",
					"first_six": "411111",
					"last_four": "1111"
				},
				"id": "payment-123",
				"merchant": {
					"id": "merch-123",
					"store_code": "store-123"
				},
				"metadata": {
					"custom_field": "custom_value"
				},
				"method_type": "card",
				"processor": "PaymentProcessor",
				"reason": "",
				"status": "approved",
				"updated_at": "2023-01-01T12:00:00Z"
			}
		},
		"status": "completed",
		"store_code": "store-123",
		"sub_total": 100,
		"tax_amount": 15,
		"timezone": "America/New_York",
		"total_amount": 125,
		"total_discount": 0,
		"webhook_urls": {
			"notify_order": "https://example.com/webhook"
		}
	},
	"order_token": "c5c5abc9-67a7-4a44-aa83-ee7e81d7052c"
}`
