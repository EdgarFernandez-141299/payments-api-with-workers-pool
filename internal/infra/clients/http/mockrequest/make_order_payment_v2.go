package mockrequest

import "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/resources/dto/request"

var MakeOrderPaymentV2RequestMock = request.DeunaOrderPaymentRequestV2{
	OrderType: "purchase",
	Order: request.OrderDetailsV2{
		Currency:       "USD",
		StoreCode:      "store-123",
		OrderID:        "order-123",
		SubTotal:       10000,
		TotalAmount:    12500,
		TotalTaxAmount: 1500,
		Items: []request.ItemDetailsV2{
			{
				ID:          "item-123",
				Name:        "Air Max",
				Description: "Running shoes",
				Category:    "Shoes",
				Quantity:    1,
				UnitPrice: request.UnitPriceV2{
					Amount:         10000,
					Currency:       "USD",
					CurrencySymbol: "$",
				},
				TotalAmount: request.TotalAmountV2{
					Amount:         10000,
					Currency:       "USD",
					CurrencySymbol: "$",
					OriginalAmount: 10000,
					TotalDiscount:  0,
				},
				Taxable: true,
				TaxAmount: request.CurrencyDetailsV2{
					Amount:         1500,
					Currency:       "USD",
					CurrencySymbol: "$",
				},
			},
		},
		BillingAddress: request.AddressV2{
			Country:     "US",
			Address1:    "123 Main St",
			City:        "New York",
			StateName:   "NY",
			FirstName:   "John",
			LastName:    "Doe",
			Phone:       "+1234567890",
			Zipcode:     "10001",
			AddressType: "billing",
		},
	},
	PayerInfo: request.PayerInfoV2{
		Email:        "test@example.com",
		SaveUserInfo: true,
	},
	PaymentSource: request.PaymentSourceV2{
		MethodType: "card",
		CardInfo: request.CardInfoV2{
			CardHolder:    "John Doe",
			CardHolderDNI: "123456789",
			CardNumber:    "4111111111111111",
			ExpiryMonth:   "12",
			ExpiryYear:    "2025",
			CardCVV:       "123",
			Country:       "US",
			State:         "NY",
			City:          "New York",
			Address1:      "123 Main St",
			Zip:           "10001",
			Phone:         "+1234567890",
		},
	},
}
