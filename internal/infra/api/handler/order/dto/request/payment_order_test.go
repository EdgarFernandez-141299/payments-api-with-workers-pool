package dto

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/constants"
	paymentmethodsDto "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/request/payment_methods"
)

func TestCreatePaymentOrderRequestDTO_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request CreatePaymentOrderRequestDTO
		wantErr bool
	}{
		{
			name: "valid request with credit card payment",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.PaymentMethodCreditCard,
						CreditCard: paymentmethodsDto.CreditCard{
							ID:  "card-123",
							CVV: "123",
						},
					},
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				UserType:     "member",
				CurrencyCode: "USD",
				CountryCode:  "US",
			},
			wantErr: false,
		},
		{
			name: "valid request with terminal payment",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.TerminalPaymentMethod,
					},
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				UserType:     "member",
				CurrencyCode: "USD",
				CountryCode:  "US",
			},
			wantErr: false,
		},
		{
			name: "valid request with token card payment",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.PaymentMethodTokenCard,
						TokenCard: paymentmethodsDto.TokenCardDTO{
							Token:    "token-123",
							CVV:      "123",
							SaveCard: true,
							Card: paymentmethodsDto.TokenCardDataDTO{
								Brand:    "visa",
								Last4:    "1234",
								Exp:      "1225",
								CardType: "credit",
								Alias:    "My Card",
							},
						},
					},
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				UserType:     "member",
				CurrencyCode: "USD",
				CountryCode:  "US",
			},
			wantErr: false,
		},
		{
			name: "invalid request - missing order ID",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.PaymentMethodCreditCard,
					},
				},
				UserID:       "user-123",
				UserType:     "member",
				CurrencyCode: "USD",
				CountryCode:  "US",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing payment order ID",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.PaymentMethodCreditCard,
					},
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				UserType:     "member",
				CurrencyCode: "USD",
				CountryCode:  "US",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing user ID",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.PaymentMethodCreditCard,
					},
				},
				OrderID:      "order-123",
				UserType:     "member",
				CurrencyCode: "USD",
				CountryCode:  "US",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing user type",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.PaymentMethodCreditCard,
					},
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				CurrencyCode: "USD",
				CountryCode:  "US",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing associated origin",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID: "payment-123",
					Amount:         decimal.NewFromFloat(100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.PaymentMethodCreditCard,
					},
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				UserType:     "member",
				CurrencyCode: "USD",
				CountryCode:  "US",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing currency code",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.PaymentMethodCreditCard,
					},
				},
				OrderID:     "order-123",
				UserID:      "user-123",
				UserType:    "member",
				CountryCode: "US",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing payment method",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(100.0),
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				UserType:     "member",
				CurrencyCode: "USD",
				CountryCode:  "US",
			},
			wantErr: true,
		},
		{
			name: "invalid request - amount is less than 0",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(-100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.PaymentMethodCreditCard,
					},
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				UserType:     "member",
				CurrencyCode: "USD",
				CountryCode:  "US",
			},
			wantErr: true,
		},
		{
			name: "invalid request - missing country code",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.PaymentMethodCreditCard,
					},
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				UserType:     "member",
				CurrencyCode: "USD",
			},
			wantErr: true,
		},
		{
			name: "invalid request - credit card missing ID and CVV",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.PaymentMethodCreditCard,
						CreditCard: paymentmethodsDto.CreditCard{
							ID:  "",
							CVV: "",
						},
					},
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				UserType:     "member",
				CurrencyCode: "USD",
				CountryCode:  "US",
			},
			wantErr: true,
		},
		{
			name: "invalid request - amount is zero",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.PaymentMethodCreditCard,
						CreditCard: paymentmethodsDto.CreditCard{
							ID:  "card-123",
							CVV: "123",
						},
					},
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				UserType:     "member",
				CurrencyCode: "USD",
				CountryCode:  "US",
			},
			wantErr: true,
		},
		{
			name: "invalid request - invalid user type",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.PaymentMethodCreditCard,
						CreditCard: paymentmethodsDto.CreditCard{
							ID:  "card-123",
							CVV: "123",
						},
					},
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				UserType:     "invalid_type",
				CurrencyCode: "USD",
				CountryCode:  "US",
			},
			wantErr: true,
		},
		{
			name: "invalid request - invalid associated origin",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "INVALID_ORIGIN",
					Amount:           decimal.NewFromFloat(100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.PaymentMethodCreditCard,
						CreditCard: paymentmethodsDto.CreditCard{
							ID:  "card-123",
							CVV: "123",
						},
					},
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				UserType:     "member",
				CurrencyCode: "USD",
				CountryCode:  "US",
			},
			wantErr: true,
		},
		{
			name: "invalid request - token card missing token and cvv",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.PaymentMethodTokenCard,
						TokenCard: paymentmethodsDto.TokenCardDTO{
							Token:    "",
							CVV:      "",
							SaveCard: true,
							Card: paymentmethodsDto.TokenCardDataDTO{
								Brand:    "visa",
								Last4:    "1234",
								Exp:      "1225",
								CardType: "credit",
								Alias:    "My Card",
							},
						},
					},
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				UserType:     "member",
				CurrencyCode: "USD",
				CountryCode:  "US",
			},
			wantErr: true,
		},
		{
			name: "invalid request - token card save card without alias",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.PaymentMethodTokenCard,
						TokenCard: paymentmethodsDto.TokenCardDTO{
							Token:    "token-123",
							CVV:      "123",
							SaveCard: true,
							Card: paymentmethodsDto.TokenCardDataDTO{
								Brand:    "visa",
								Last4:    "1234",
								Exp:      "1225",
								CardType: "credit",
								Alias:    "",
							},
						},
					},
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				UserType:     "member",
				CurrencyCode: "USD",
				CountryCode:  "US",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreatePaymentOrderRequestDTO_Command(t *testing.T) {
	tests := []struct {
		name    string
		request CreatePaymentOrderRequestDTO
		wantErr bool
	}{
		{
			name: "valid command with credit card payment",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.PaymentMethodCreditCard,
						CreditCard: paymentmethodsDto.CreditCard{
							ID:  "card-123",
							CVV: "123",
						},
					},
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				UserType:     "member",
				CurrencyCode: "USD",
				CountryCode:  "US",
			},
			wantErr: false,
		},
		{
			name: "valid command with terminal payment",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.TerminalPaymentMethod,
					},
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				UserType:     "member",
				CurrencyCode: "USD",
				CountryCode:  "US",
			},
			wantErr: false,
		},
		{
			name: "valid command with token card payment",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.PaymentMethodTokenCard,
						TokenCard: paymentmethodsDto.TokenCardDTO{
							Token:    "token-123",
							CVV:      "123",
							SaveCard: true,
							Card: paymentmethodsDto.TokenCardDataDTO{
								Brand:    "visa",
								Last4:    "1234",
								Exp:      "1225",
								CardType: "credit",
								Alias:    "My Card",
							},
						},
					},
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				UserType:     "member",
				CurrencyCode: "USD",
				CountryCode:  "US",
			},
			wantErr: false,
		},
		{
			name: "invalid command - unsupported user type",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.PaymentMethodCreditCard,
						CreditCard: paymentmethodsDto.CreditCard{
							ID:  "card-123",
							CVV: "123",
						},
					},
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				UserType:     "invalid_type",
				CurrencyCode: "USD",
				CountryCode:  "US",
			},
			wantErr: true,
		},
		{
			name: "invalid command - unsupported payment method",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: "INVALID_METHOD",
					},
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				UserType:     "member",
				CurrencyCode: "USD",
				CountryCode:  "US",
			},
			wantErr: true,
		},
		{
			name: "invalid command - invalid currency code",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.PaymentMethodCreditCard,
						CreditCard: paymentmethodsDto.CreditCard{
							ID:  "card-123",
							CVV: "123",
						},
					},
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				UserType:     "member",
				CurrencyCode: "INVALID",
				CountryCode:  "US",
			},
			wantErr: true,
		},
		{
			name: "invalid command - invalid country code",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.PaymentMethodCreditCard,
						CreditCard: paymentmethodsDto.CreditCard{
							ID:  "card-123",
							CVV: "123",
						},
					},
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				UserType:     "member",
				CurrencyCode: "USD",
				CountryCode:  "INVALID",
			},
			wantErr: true,
		},
		{
			name: "invalid command - invalid amount for currency",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(0.001), // Mínimo para USD es 0.01
					PaymentMethod: PaymentMethodDTO{
						Type: constants.PaymentMethodCreditCard,
						CreditCard: paymentmethodsDto.CreditCard{
							ID:  "card-123",
							CVV: "123",
						},
					},
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				UserType:     "member",
				CurrencyCode: "USD",
				CountryCode:  "US",
			},
			wantErr: false,
		},
		{
			name: "invalid command - token card missing token and cvv",
			request: CreatePaymentOrderRequestDTO{
				PaymentOrderRequestDTO: PaymentOrderRequestDTO{
					PaymentOrderID:   "payment-123",
					AssociatedOrigin: "DOWNPAYMENT",
					Amount:           decimal.NewFromFloat(100.0),
					PaymentMethod: PaymentMethodDTO{
						Type: constants.PaymentMethodTokenCard,
						TokenCard: paymentmethodsDto.TokenCardDTO{
							Token:    "",
							CVV:      "",
							SaveCard: true,
							Card: paymentmethodsDto.TokenCardDataDTO{
								Brand:    "visa",
								Last4:    "1234",
								Exp:      "1225",
								CardType: "credit",
								Alias:    "My Card",
							},
						},
					},
				},
				OrderID:      "order-123",
				UserID:       "user-123",
				UserType:     "member",
				CurrencyCode: "USD",
				CountryCode:  "US",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.request.Command()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
