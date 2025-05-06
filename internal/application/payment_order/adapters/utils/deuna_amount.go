package utils

import "github.com/shopspring/decimal"

const precisionFloat = 2

type DeunaTotalRefundRequest struct {
	Reason string `json:"reason"`
}

type DeunaPartialRefundRequest struct {
	Amount int64  `json:"amount"`
	Reason string `json:"reason"`
}

func NewDeunaAmount(amount decimal.Decimal) int64 {
	const CentsPerUnit = 100

	return amount.Mul(decimal.NewFromInt(CentsPerUnit)).Truncate(precisionFloat).IntPart()
}

func DeunaAmountToAmount(deunaAmount int64) decimal.Decimal {
	const CentsPerUnit = 100

	return decimal.NewFromInt(deunaAmount).Div(decimal.NewFromInt(CentsPerUnit)).Truncate(precisionFloat)
}
