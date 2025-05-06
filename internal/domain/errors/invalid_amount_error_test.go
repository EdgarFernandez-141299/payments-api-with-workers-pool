package errors

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

func TestNewInvalidAmountError(t *testing.T) {
	t.Run("should create error with correct message and code", func(t *testing.T) {
		amount := decimal.NewFromFloat(100.50)
		err := NewInvalidAmountError(amount)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid amount: 100.5")
		assert.True(t, domain.IsBusinessErrorCode(err, invalidAmountErrorCode))
	})
}
