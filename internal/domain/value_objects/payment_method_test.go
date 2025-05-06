package value_objects

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	paymentmethodsVo "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects/payment_methods"
)

func TestPaymentMethod_GetCC(t *testing.T) {
	tests := []struct {
		name       string
		payment    PaymentMethod
		expectedCC CardInfo
		expectErr  bool
	}{
		{
			name: "valid CCData payment",
			payment: PaymentMethod{
				Type: enums.CCMethod,
				CCData: PaymentMethodData[CardInfo]{
					Data: CardInfo{CardID: "1234", CVV: "567"},
				},
			},
			expectedCC: CardInfo{CardID: "1234", CVV: "567"},
			expectErr:  false,
		},
		{
			name: "invalid CCData payment type",
			payment: PaymentMethod{
				Type: enums.PaymentDevice,
				CCData: PaymentMethodData[CardInfo]{
					Data: CardInfo{CardID: "1234", CVV: "567"},
				},
			},
			expectedCC: CardInfo{},
			expectErr:  true,
		},
		{
			name: "empty CCData data",
			payment: PaymentMethod{
				Type: enums.CCMethod,
				CCData: PaymentMethodData[CardInfo]{
					Data: CardInfo{},
				},
			},
			expectedCC: CardInfo{},
			expectErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := tt.payment.GetCC()
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if err == nil && cc != tt.expectedCC {
				t.Errorf("expected CCData: %v, got: %v", tt.expectedCC, cc)
			} else if err != nil && !tt.expectErr && !errors.Is(err, nil) {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestPaymentMethod_GetDevice(t *testing.T) {
	tests := []struct {
		name          string
		payment       PaymentMethod
		expectedData  PaymentDevice
		expectedError bool
	}{
		{
			name: "valid Device payment",
			payment: PaymentMethod{
				Type: enums.PaymentDevice,
				DeviceData: PaymentMethodData[PaymentDevice]{
					Data: PaymentDevice{},
				},
			},
			expectedData:  PaymentDevice{},
			expectedError: false,
		},
		{
			name: "invalid Device type",
			payment: PaymentMethod{
				Type: enums.CCMethod,
				DeviceData: PaymentMethodData[PaymentDevice]{
					Data: PaymentDevice{},
				},
			},
			expectedData:  PaymentDevice{},
			expectedError: true,
		},
		{
			name: "empty Device data",
			payment: PaymentMethod{
				Type: enums.PaymentDevice,
				DeviceData: PaymentMethodData[PaymentDevice]{
					Data: PaymentDevice{},
				},
			},
			expectedData:  PaymentDevice{},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			terminal, err := tt.payment.GetDevice()
			if (err != nil) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}
			if err == nil && terminal != tt.expectedData {
				t.Errorf("expected DeviceData: %v, got: %v", tt.expectedData, terminal)
			}
		})
	}
}

func TestPaymentMethod_GetTokenCard(t *testing.T) {
	tests := []struct {
		name          string
		payment       PaymentMethod
		expectedToken paymentmethodsVo.TokenCard
		expectErr     bool
	}{
		{
			name: "valid TokenCard payment",
			payment: PaymentMethod{
				Type: enums.TokenCard,
				TokenCard: PaymentMethodData[paymentmethodsVo.TokenCard]{
					Data: paymentmethodsVo.TokenCard{
						Token:    "token-123",
						CVV:      "123",
						Brand:    "visa",
						Last4:    "1234",
						Exp:      "1225",
						CardType: "credit",
						SaveCard: true,
						Alias:    "My Card",
					},
				},
			},
			expectedToken: paymentmethodsVo.TokenCard{
				Token:    "token-123",
				CVV:      "123",
				Brand:    "visa",
				Last4:    "1234",
				Exp:      "1225",
				CardType: "credit",
				SaveCard: true,
				Alias:    "My Card",
			},
			expectErr: false,
		},
		{
			name: "invalid TokenCard payment type",
			payment: PaymentMethod{
				Type: enums.CCMethod,
				TokenCard: PaymentMethodData[paymentmethodsVo.TokenCard]{
					Data: paymentmethodsVo.TokenCard{
						Token: "token-123",
					},
				},
			},
			expectedToken: paymentmethodsVo.TokenCard{},
			expectErr:     true,
		},
		{
			name: "empty TokenCard data",
			payment: PaymentMethod{
				Type: enums.TokenCard,
				TokenCard: PaymentMethodData[paymentmethodsVo.TokenCard]{
					Data: paymentmethodsVo.TokenCard{},
				},
			},
			expectedToken: paymentmethodsVo.TokenCard{},
			expectErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := tt.payment.GetTokenCard()
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if err == nil && token != tt.expectedToken {
				t.Errorf("expected TokenCard: %v, got: %v", tt.expectedToken, token)
			}
		})
	}
}

func Test_MarshalUnmarshalPaymentMethod(t *testing.T) {
	t.Run("should marshal and unmarshal payment method CCData", func(t *testing.T) {
		cc := NewCCPaymentMethod("1234", "567")

		b, err := json.Marshal(cc)

		assert.NoError(t, err)
		assert.NotNil(t, b)

		deserializedCC := PaymentMethod{}

		err = json.Unmarshal(b, &deserializedCC)

		assert.NoError(t, err)
		assert.Equal(t, cc, deserializedCC)
	})

	t.Run("should marshal and unmarshal payment method Device", func(t *testing.T) {
		device := NewDevicePaymentMethod()

		b, err := json.Marshal(device)

		assert.NoError(t, err)
		assert.NotNil(t, b)

		deserializedDevice := PaymentMethod{}

		err = json.Unmarshal(b, &deserializedDevice)

		assert.NoError(t, err)
		assert.Equal(t, device, deserializedDevice)
	})

	t.Run("should marshal and unmarshal payment method TokenCard", func(t *testing.T) {
		tokenCard := PaymentMethod{
			Type: enums.TokenCard,
			TokenCard: PaymentMethodData[paymentmethodsVo.TokenCard]{
				Data: paymentmethodsVo.TokenCard{
					Token:    "token-123",
					CVV:      "123",
					Brand:    "visa",
					Last4:    "1234",
					Exp:      "1225",
					CardType: "credit",
					SaveCard: true,
					Alias:    "My Card",
				},
			},
		}

		b, err := json.Marshal(tokenCard)

		assert.NoError(t, err)
		assert.NotNil(t, b)

		deserializedTokenCard := PaymentMethod{}

		err = json.Unmarshal(b, &deserializedTokenCard)

		assert.NoError(t, err)
		assert.Equal(t, tokenCard, deserializedTokenCard)
	})
}

func Test_IsDevice(t *testing.T) {
	payment := PaymentMethod{
		Type: enums.PaymentDevice,
		DeviceData: PaymentMethodData[PaymentDevice]{
			Data: PaymentDevice{},
		},
	}

	assert.True(t, payment.IsDevice())
}

func Test_IsCC(t *testing.T) {
	payment := PaymentMethod{
		Type: enums.CCMethod,
		CCData: PaymentMethodData[CardInfo]{
			Data: CardInfo{},
		},
	}

	assert.True(t, payment.IsCC())
}

func TestPaymentMethod_Validate(t *testing.T) {
	t.Run("should return error if payment method is not supported", func(t *testing.T) {
		payment := PaymentMethod{
			Type:   enums.CCMethod,
			CCData: PaymentMethodData[CardInfo]{},
		}

		err := payment.Validate()

		assert.EqualError(t, err, "CardID is required for credit card payment")
	})
	t.Run("should return payment method Credit card successful", func(t *testing.T) {
		payment := PaymentMethod{
			Type: enums.CCMethod,
			CCData: PaymentMethodData[CardInfo]{
				Data: CardInfo{
					CVV:    "111",
					CardID: "123",
				},
			},
		}

		err := payment.Validate()

		assert.Nil(t, err, "CardID is required for credit card payment")
	})

	t.Run("should return valid payment method", func(t *testing.T) {
		payment := PaymentMethod{
			Type:       enums.PaymentDevice,
			DeviceData: PaymentMethodData[PaymentDevice]{},
		}

		err := payment.Validate()

		assert.Nil(t, err)
	})

	t.Run("should return valid token card payment method", func(t *testing.T) {
		payment := PaymentMethod{
			Type: enums.TokenCard,
			TokenCard: PaymentMethodData[paymentmethodsVo.TokenCard]{
				Data: paymentmethodsVo.TokenCard{
					Token:    "token-123",
					CVV:      "123",
					Brand:    "visa",
					Last4:    "1234",
					Exp:      "1225",
					CardType: "credit",
					SaveCard: true,
					Alias:    "My Card",
				},
			},
		}

		err := payment.Validate()

		assert.Nil(t, err)
	})

	t.Run("should return error for invalid token card payment method", func(t *testing.T) {
		payment := PaymentMethod{
			Type: enums.TokenCard,
			TokenCard: PaymentMethodData[paymentmethodsVo.TokenCard]{
				Data: paymentmethodsVo.TokenCard{
					Token:    "",
					CVV:      "",
					Brand:    "visa",
					Last4:    "1234",
					Exp:      "1225",
					CardType: "credit",
					SaveCard: true,
					Alias:    "My Card",
				},
			},
		}

		err := payment.Validate()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token is required for credit card payment")
	})

	t.Run("should return error for invalid token card CVV", func(t *testing.T) {
		payment := PaymentMethod{
			Type: enums.TokenCard,
			TokenCard: PaymentMethodData[paymentmethodsVo.TokenCard]{
				Data: paymentmethodsVo.TokenCard{
					Token:    "token-123",
					CVV:      "12",
					Brand:    "visa",
					Last4:    "1234",
					Exp:      "1225",
					CardType: "credit",
					SaveCard: true,
					Alias:    "My Card",
				},
			},
		}

		err := payment.Validate()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CVV must be a 3-digit number for credit card payment")
	})

	t.Run("should return error for expired token card", func(t *testing.T) {
		payment := PaymentMethod{
			Type: enums.TokenCard,
			TokenCard: PaymentMethodData[paymentmethodsVo.TokenCard]{
				Data: paymentmethodsVo.TokenCard{
					Token:    "token-123",
					CVV:      "123",
					Brand:    "visa",
					Last4:    "1234",
					Exp:      "1220", // Expired card
					CardType: "credit",
					SaveCard: true,
					Alias:    "My Card",
				},
			},
		}

		err := payment.Validate()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "card is expired")
	})
}
