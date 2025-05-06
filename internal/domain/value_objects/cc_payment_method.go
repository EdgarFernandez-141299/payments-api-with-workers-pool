package value_objects

import (
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	paymentmethodsVo "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects/payment_methods"
)

func NewCCPaymentMethod(cardID, cvv string) PaymentMethod {
	return PaymentMethod{
		Type: enums.CCMethod,
		CCData: PaymentMethodData[CardInfo]{
			Data: CardInfo{
				CardID: cardID,
				CVV:    cvv,
			},
		},
	}
}

func NewTokenCardPaymentMethod(token, cvv, brand, last4, exp, cardType, alias string,
	saveCard bool) PaymentMethod {
	return PaymentMethod{
		Type: enums.TokenCard,
		TokenCard: PaymentMethodData[paymentmethodsVo.TokenCard]{
			Data: paymentmethodsVo.TokenCard{
				Token:    token,
				CVV:      cvv,
				Brand:    brand,
				Last4:    last4,
				Exp:      exp,
				CardType: cardType,
				SaveCard: saveCard,
				Alias:    alias,
			},
		},
	}
}
