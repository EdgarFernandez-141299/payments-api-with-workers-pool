package enums

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProviders(t *testing.T) {
	t.Run("should have correct value for ProvidersDeUna", func(t *testing.T) {
		assert.Equal(t, "DEUNA", string(ProvidersDeUna))
	})

	t.Run("should be able to convert Providers to string", func(t *testing.T) {
		provider := Providers("DEUNA")
		assert.Equal(t, "DEUNA", string(provider))
	})
}
