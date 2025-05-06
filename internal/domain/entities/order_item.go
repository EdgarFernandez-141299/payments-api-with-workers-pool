package entities

import vo "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"

type Item struct {
	ID          string
	Amount      vo.CurrencyAmount
	Description string
}
