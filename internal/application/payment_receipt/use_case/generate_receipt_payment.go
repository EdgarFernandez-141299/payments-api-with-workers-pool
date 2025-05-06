package use_case

import (
	"context"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/payment_receipt/command"
)

type GenerateReceiptPaymentUseCase interface {
	Generate(ctx context.Context, cmd command.CreatePaymentReceiptCommand) (entities.PaymentReceipt, error)
}
