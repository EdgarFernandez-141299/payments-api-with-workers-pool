package use_case

import (
	"context"
	"errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/common"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/payment_receipt/command"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/pdf_generator"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository"
)

func TestNewGenerateReceiptPaymentUseCase(t *testing.T) {
	mockRepo := repository.NewPaymentReceiptRepository(t)
	mockGenerator := pdf_generator.NewReceiptPaymentGenerator(t)
	mockStorage := common.NewStorageAdapter(t)

	useCase := NewGenerateReceiptPaymentUseCase(mockRepo, mockGenerator, mockStorage)

	assert.NotNil(t, useCase)
	assert.IsType(t, &generateReceiptPaymentUseCase{}, useCase)
}

func TestGenerateReceiptPaymentUseCase_Generate_Success(t *testing.T) {
	mockRepo := repository.NewPaymentReceiptRepository(t)
	mockGenerator := pdf_generator.NewReceiptPaymentGenerator(t)
	mockStorage := common.NewStorageAdapter(t)
	useCase := NewGenerateReceiptPaymentUseCase(mockRepo, mockGenerator, mockStorage)

	ctx := context.Background()
	paymentID := "payment-123"
	fileURL := "https://example.com/receipts/payment-123.pdf"

	cmd := createTestCommand(paymentID)

	// Mock the PDF generation
	pdfContent := strings.NewReader("PDF content")
	mockGenerator.EXPECT().GenerateReceiptPaymentPDF(mock.Anything, mock.AnythingOfType("interfaces.ReceiptData")).Return(pdfContent, nil)

	// Mock the storage
	mockStorage.On("Store", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(fileURL, nil)

	// Mock the repository
	mockRepo.EXPECT().GetByPaymentID(mock.MatchedBy(func(ctx context.Context) bool { return true }), paymentID).Return(entities.PaymentReceipt{}, errors.New("not found"))
	mockRepo.EXPECT().CreatePaymentReceipt(mock.MatchedBy(func(ctx context.Context) bool { return true }), mock.AnythingOfType("entities.PaymentReceipt")).Return(nil)

	receipt, err := useCase.Generate(ctx, cmd)

	assert.NoError(t, err)
	assert.NotEmpty(t, receipt.ID)
	assert.Equal(t, paymentID, receipt.PaymentID)
	assert.Equal(t, cmd.UserID, receipt.UserID)
	assert.Equal(t, cmd.EnterpriseID, receipt.EnterpriseID)
	assert.Equal(t, cmd.Email, receipt.Email)
	assert.Equal(t, cmd.ReferenceOrderID, receipt.ReferenceOrderID)
	assert.Equal(t, cmd.PaymentStatus, receipt.PaymentStatus)
	assert.Equal(t, cmd.PaymentAmount.Value, receipt.PaymentAmount)
	assert.Equal(t, cmd.PaymentCountry.Code, receipt.PaymentCountryCode)
	assert.Equal(t, cmd.PaymentAmount.Code.Code, receipt.PaymentCurrencyCode)
	assert.Equal(t, cmd.PaymentMethod.Type, receipt.PaymentMethod)
	assert.Equal(t, cmd.PaymentDate, receipt.PaymentDate)
	assert.Equal(t, fileURL, receipt.FileURL)
}

func TestGenerateReceiptPaymentUseCase_Generate_AlreadyExists(t *testing.T) {
	mockRepo := repository.NewPaymentReceiptRepository(t)
	mockGenerator := pdf_generator.NewReceiptPaymentGenerator(t)
	mockStorage := common.NewStorageAdapter(t)
	useCase := NewGenerateReceiptPaymentUseCase(mockRepo, mockGenerator, mockStorage)

	ctx := context.Background()
	paymentID := "payment-123"

	cmd := createTestCommand(paymentID)

	existingReceipt := entities.PaymentReceipt{
		ID:        "RCPT-existing-123",
		PaymentID: paymentID,
	}

	mockRepo.EXPECT().GetByPaymentID(mock.MatchedBy(func(ctx context.Context) bool { return true }), paymentID).Return(existingReceipt, nil)

	_, err := useCase.Generate(ctx, cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "payment receipt for payment payment-123 already exists")
}

func TestGenerateReceiptPaymentUseCase_Generate_PDFGenerationError(t *testing.T) {
	mockRepo := repository.NewPaymentReceiptRepository(t)
	mockGenerator := pdf_generator.NewReceiptPaymentGenerator(t)
	mockStorage := common.NewStorageAdapter(t)
	useCase := NewGenerateReceiptPaymentUseCase(mockRepo, mockGenerator, mockStorage)

	ctx := context.Background()
	paymentID := "payment-123"

	cmd := createTestCommand(paymentID)

	// Mock the repository
	mockRepo.EXPECT().GetByPaymentID(mock.MatchedBy(func(ctx context.Context) bool { return true }), paymentID).Return(entities.PaymentReceipt{}, errors.New("not found"))

	// Mock the PDF generation to fail
	mockGenerator.EXPECT().GenerateReceiptPaymentPDF(mock.Anything, mock.AnythingOfType("interfaces.ReceiptData")).Return(nil, errors.New("pdf generation error"))

	_, err := useCase.Generate(ctx, cmd)

	assert.Error(t, err)
	assert.Equal(t, "pdf generation error", err.Error())
}

func TestGenerateReceiptPaymentUseCase_Generate_StorageError(t *testing.T) {
	mockRepo := repository.NewPaymentReceiptRepository(t)
	mockGenerator := pdf_generator.NewReceiptPaymentGenerator(t)
	mockStorage := common.NewStorageAdapter(t)
	useCase := NewGenerateReceiptPaymentUseCase(mockRepo, mockGenerator, mockStorage)

	ctx := context.Background()
	paymentID := "payment-123"

	cmd := createTestCommand(paymentID)

	// Mock the repository
	mockRepo.EXPECT().GetByPaymentID(mock.MatchedBy(func(ctx context.Context) bool { return true }), paymentID).Return(entities.PaymentReceipt{}, errors.New("not found"))

	// Mock the PDF generation
	pdfContent := strings.NewReader("PDF content")
	mockGenerator.EXPECT().GenerateReceiptPaymentPDF(mock.Anything, mock.AnythingOfType("interfaces.ReceiptData")).Return(pdfContent, nil)

	// Mock the storage to fail
	mockStorage.EXPECT().Store(mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return("", errors.New("storage error"))

	_, err := useCase.Generate(ctx, cmd)

	assert.Error(t, err)
	assert.Equal(t, "storage error", err.Error())
}

func TestGenerateReceiptPaymentUseCase_Generate_RepositoryError(t *testing.T) {
	mockRepo := repository.NewPaymentReceiptRepository(t)
	mockGenerator := pdf_generator.NewReceiptPaymentGenerator(t)
	mockStorage := common.NewStorageAdapter(t)
	useCase := NewGenerateReceiptPaymentUseCase(mockRepo, mockGenerator, mockStorage)

	ctx := context.Background()
	paymentID := "payment-123"
	fileURL := "https://example.com/receipts/payment-123.pdf"

	cmd := createTestCommand(paymentID)

	// Mock the PDF generation
	pdfContent := strings.NewReader("PDF content")
	mockGenerator.EXPECT().GenerateReceiptPaymentPDF(mock.Anything, mock.AnythingOfType("interfaces.ReceiptData")).Return(pdfContent, nil)

	// Mock the storage
	mockStorage.EXPECT().Store(mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(fileURL, nil)

	// Mock the repository
	mockRepo.EXPECT().GetByPaymentID(mock.MatchedBy(func(ctx context.Context) bool { return true }), paymentID).Return(entities.PaymentReceipt{}, errors.New("not found"))
	mockRepo.EXPECT().CreatePaymentReceipt(mock.MatchedBy(func(ctx context.Context) bool { return true }), mock.AnythingOfType("entities.PaymentReceipt")).Return(errors.New("database error"))

	_, err := useCase.Generate(ctx, cmd)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
}

func createTestCommand(paymentID string) command.CreatePaymentReceiptCommand {
	now := time.Now()
	return command.CreatePaymentReceiptCommand{
		UserID:           "user-123",
		EnterpriseID:     "enterprise-123",
		Email:            "test@example.com",
		ReferenceOrderID: "order-123",
		PaymentID:        paymentID,
		PaymentStatus:    "completed",
		PaymentAmount: value_objects.CurrencyAmount{
			Value: decimal.NewFromFloat(100.50),
			Code:  value_objects.CurrencyCode{Code: "USD"},
		},
		PaymentCountry: value_objects.Country{
			Code: "US",
		},
		PaymentMethod: value_objects.PaymentMethod{
			Type: enums.CCMethod,
		},
		PaymentDate: "2023-01-01",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
