package webhooks

import (
	"fmt"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/config"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/services"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

func NewBillsApiWebhookUrl() services.BillsApiWebhookUrl {
	cfgUrl := fmt.Sprintf("%s/api/payments/register-payment", config.Config().BillApi.URL)
	webhookUrl := value_objects.NewWebhookUrl(cfgUrl)

	return services.BillsApiWebhookUrl(webhookUrl)
}
