package mockresponse

const SuccessRefundPaymentResponseMock = `{
	"data": {
		"refund_amount": {
			"amount": "125",
			"currency": "USD",
			"currency_symbol": "$",
			"display_amount": "$125.00"
		},
		"refund_id": "refund-123",
		"refunds": [
			{
				"external_transaction_id": "ext-123",
				"refund_amount": {
					"amount": "125",
					"currency": "USD",
					"currency_symbol": "$",
					"display_amount": "$125.00"
				},
				"refund_id": "refund-123",
				"refunded_on": "2023-01-02T12:00:00Z",
				"status": "approved"
			}
		],
		"status": "approved"
	}
}`
