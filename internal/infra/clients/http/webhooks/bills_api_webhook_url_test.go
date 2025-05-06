package webhooks

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/config"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/services"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

func TestNewBillsApiWebhookUrl(t *testing.T) {
	// Configurar variables de entorno para la prueba
	expectedUrl := "https://test-bills-api.example.com/webhook"
	os.Setenv("ENVIRONMENT", "test")
	os.Setenv("BILL_API_URL", expectedUrl)

	// Inicializar la configuración
	err := config.Environments()
	require.NoError(t, err, "Cannot initialize config")

	cfg := config.Config()
	assert.NotNil(t, cfg, "Config should not be nil")
	assert.Equal(t, expectedUrl, cfg.BillApi.URL, "Url should be equal")

	// Crear el objeto BillsApiWebhookUrl
	webhookUrl := NewBillsApiWebhookUrl()

	// Verificar que el tipo es correcto
	assert.IsType(t, services.BillsApiWebhookUrl(value_objects.WebhookUrl{}), webhookUrl)

	// Verificar que la URL se configuró correctamente
	webhookUrlStr := value_objects.WebhookUrl(webhookUrl).String()
	assert.NotEmpty(t, webhookUrlStr, "Webhook url should not be empty")
	assert.Equal(t, expectedUrl+"/api/payments/register-payment", webhookUrlStr, "Webhook url should be equal")

	// Limpiar variables de entorno
	os.Unsetenv("ENVIRONMENT")
	os.Unsetenv("BILL_API_URL")
}
