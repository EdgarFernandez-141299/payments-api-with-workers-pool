package use_case

import (
	"context"
	"fmt"
	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/common/adapters"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_receipt/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_receipt/interfaces"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/payment_receipt/command"
)

type generateReceiptPaymentUseCase struct {
	paymentReceiptRepository repository.PaymentReceiptRepository
	receiptGenerator         interfaces.ReceiptPaymentGenerator
	storageAdapter           adapters.StorageAdapter
}

func NewGenerateReceiptPaymentUseCase(
	paymentReceiptRepository repository.PaymentReceiptRepository,
	receiptGenerator interfaces.ReceiptPaymentGenerator,
	storageAdapter adapters.StorageAdapter,
) GenerateReceiptPaymentUseCase {
	return &generateReceiptPaymentUseCase{
		paymentReceiptRepository: paymentReceiptRepository,
		receiptGenerator:         receiptGenerator,
		storageAdapter:           storageAdapter,
	}
}

func (g *generateReceiptPaymentUseCase) Generate(oldCtx context.Context, cmd command.CreatePaymentReceiptCommand) (entities.PaymentReceipt, error) {
	return decorators.TraceDecorator[entities.PaymentReceipt](
		oldCtx,
		"GenerateReceiptUseCase.Generate",
		func(ctx context.Context, span decorators.Span) (entities.PaymentReceipt, error) {
			existingReceipt, err := g.paymentReceiptRepository.GetByPaymentID(ctx, cmd.PaymentID)
			if err == nil && !existingReceipt.IsEmpty() {
				return entities.PaymentReceipt{}, errors.NewPaymentReceiptAlreadyExistError(cmd.PaymentID)
			}

			// Generate receipt PDF
			receiptData := interfaces.ReceiptData{
				ReceiptNumber:    cmd.PaymentID,
				TransactionID:    cmd.PaymentID,
				AmountPaid:       cmd.PaymentAmount.Value.String(),
				PaymentDate:      cmd.PaymentDate,
				PaymentMethod:    cmd.PaymentMethod.Type.String(),
				PaymentReference: cmd.ReferenceOrderID,
			}

			pdfReader, err := g.receiptGenerator.GenerateReceiptPaymentPDF(ctx, receiptData)
			if err != nil {
				return entities.PaymentReceipt{}, err
			}

			// Store the PDF
			filePath := fmt.Sprintf("receipts/%s_%s.pdf", cmd.ReferenceOrderID, cmd.PaymentID)
			fileURL, err := g.storageAdapter.Store(ctx, filePath, pdfReader)
			if err != nil {
				return entities.PaymentReceipt{}, err
			}

			newReceipt := entities.NewPaymentReceiptEntity(cmd).WithReceiptURL(fileURL)

			err = g.paymentReceiptRepository.CreatePaymentReceipt(ctx, newReceipt)
			if err != nil {
				return entities.PaymentReceipt{}, err
			}

			return newReceipt, nil
		})
}
