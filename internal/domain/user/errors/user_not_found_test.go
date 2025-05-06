package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

func TestNewUserNotFoundError(t *testing.T) {
	t.Run("should create error with correct message and code", func(t *testing.T) {
		memberID := "12345"
		originalErr := errors.New("original error")
		err := NewUserNotFoundError(memberID, originalErr)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found with memberId: 12345")
		assert.True(t, domain.IsBusinessErrorCode(err, userNotFoundError))
	})

	t.Run("should create error with nil original error", func(t *testing.T) {
		memberID := "12345"
		err := NewUserNotFoundError(memberID, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found with memberId: 12345")
		assert.True(t, domain.IsBusinessErrorCode(err, userNotFoundError))
	})
}
