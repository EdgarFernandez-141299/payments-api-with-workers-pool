package dto

import (
	"github.com/samber/lo"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/worfkflows"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
)

type OrderWithPaymentsDTO struct {
	ReferenceOrderID string                   `json:"reference_order_id" validate:"required"`
	CurrencyCode     string                   `json:"currency_code" validate:"required"`
	CountryCode      string                   `json:"country_code" validate:"required"`
	UserType         string                   `json:"user_type" validate:"required"`
	UserID           string                   `json:"user_id" validate:"required"`
	Payments         []PaymentOrderRequestDTO `json:"payments"`
}

func (dto OrderWithPaymentsDTO) ToWorkflowInput() (worfkflows.PaymentWorkflowInput, error) {
	paymentOrderCommandDTO := lo.Map(dto.Payments, func(payment PaymentOrderRequestDTO, _ int) CreatePaymentOrderRequestDTO {
		return CreatePaymentOrderRequestDTO{
			OrderID:                dto.ReferenceOrderID,
			UserID:                 dto.UserID,
			UserType:               dto.UserType,
			PaymentOrderRequestDTO: payment,
			CurrencyCode:           dto.CurrencyCode,
			CountryCode:            dto.CountryCode,
		}
	})

	paymentCommands := make([]command.CreatePaymentOrderCommand, 0)

	for _, payment := range paymentOrderCommandDTO {
		if err := payment.Validate(); err != nil {
			return worfkflows.PaymentWorkflowInput{}, err
		}

		paymentCommand, err := payment.Command()

		if err != nil {
			return worfkflows.PaymentWorkflowInput{}, err
		}

		if err := paymentCommand.Validate(); err != nil {
			return worfkflows.PaymentWorkflowInput{}, err
		}

		paymentCommands = append(paymentCommands, paymentCommand)

	}

	return worfkflows.NewPaymentWorkflowInput(dto.ReferenceOrderID, paymentCommands), nil
}
