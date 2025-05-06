package value_objects

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
)

func TestNewAssociatedOrigin(t *testing.T) {
	originType := enums.AssociatedOrigin("test_origin")
	associatedOrigin := NewAssociatedOrigin(originType)

	assert.Equal(t, originType, associatedOrigin.Type)
}

func TestNewFromAssociatedOriginString(t *testing.T) {
	t.Run("valid origin string", func(t *testing.T) {
		originString := "DOWNPAYMENT"
		associatedOrigin, err := NewFromAssociatedOriginString(originString)

		assert.NoError(t, err)
		assert.True(t, associatedOrigin.Type.IsValid())
	})

	t.Run("invalid origin string", func(t *testing.T) {
		originString := "invalid_origin"
		_, err := NewFromAssociatedOriginString(originString)

		assert.Error(t, err)
	})
}

func TestAssociatedOrigin_Validate(t *testing.T) {
	t.Run("valid type", func(t *testing.T) {
		originType := enums.AssociatedOrigin("DOWNPAYMENT")
		associatedOrigin := NewAssociatedOrigin(originType)

		err := associatedOrigin.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid type", func(t *testing.T) {
		originType := enums.AssociatedOrigin("invalid_origin")
		associatedOrigin := NewAssociatedOrigin(originType)

		err := associatedOrigin.Validate()
		assert.Error(t, err)
		assert.Equal(t, "type is not valid", err.Error())
	})
}
