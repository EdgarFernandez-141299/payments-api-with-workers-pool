package pdf_generator

import (
	"context"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_receipt/interfaces"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/pdf_generator"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewReceiptPaymentGeneratorImpl(t *testing.T) {
	// Save current working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "receipt-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to the temporary directory
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	defer os.Chdir(originalWd) // Restore original working directory

	// Test the function
	generator := NewReceiptPaymentGeneratorImpl()

	// Assert that the generator is not nil
	assert.NotNil(t, generator)

	// Assert that the generator is of the correct type
	_, ok := generator.(*ReceiptPaymentGeneratorImpl)
	assert.True(t, ok)

	// We can't directly check the BasePath of the PDFUtil anymore since it's behind an interface
	// Instead, we'll verify that the generator was created successfully
}

func TestReceiptData_ToMap(t *testing.T) {
	// Create a test ReceiptData
	receiptData := interfaces.ReceiptData{
		ReceiptNumber:    "RN123",
		TransactionID:    "TX456",
		AmountPaid:       "$100.00",
		PaymentDate:      "2023-05-01",
		PaymentMethod:    "Credit Card",
		PaymentReference: "REF789",
	}

	// Call the ToMap method
	result := receiptData.ToMap()

	// Assert that the map contains the expected values
	assert.Equal(t, "RN123", result["ReceiptNumber"])
	assert.Equal(t, "TX456", result["TransactionID"])
	assert.Equal(t, "$100.00", result["AmountPaid"])
	assert.Equal(t, "2023-05-01", result["PaymentDate"])
	assert.Equal(t, "Credit Card", result["PaymentMethod"])
	assert.Equal(t, "REF789", result["PaymentReference"])
}

func TestGenerateReceiptPaymentPDF(t *testing.T) {
	// Create a mock PDFUtil
	mockPDFUtil := pdf_generator.NewPDFUtilInterface(t)

	// Create a test ReceiptData
	receiptData := interfaces.ReceiptData{
		ReceiptNumber:    "RN123",
		TransactionID:    "TX456",
		AmountPaid:       "$100.00",
		PaymentDate:      "2023-05-01",
		PaymentMethod:    "Credit Card",
		PaymentReference: "REF789",
	}

	// Create expected data map
	expectedDataMap := map[string]string{
		"ReceiptNumber":    "RN123",
		"TransactionID":    "TX456",
		"AmountPaid":       "$100.00",
		"PaymentDate":      "2023-05-01",
		"PaymentMethod":    "Credit Card",
		"PaymentReference": "REF789",
	}

	// Set up the mock to expect a call to GeneratePDF with the correct parameters
	expectedReader := strings.NewReader("PDF content")
	mockPDFUtil.On("GeneratePDF", mock.Anything, filename, expectedDataMap).Return(expectedReader, nil)

	// Create a ReceiptPaymentGeneratorImpl with the mock PDFUtil
	generator := &ReceiptPaymentGeneratorImpl{
		pdfUtil: mockPDFUtil,
	}

	// Call the method being tested
	result, err := generator.GenerateReceiptPaymentPDF(context.Background(), receiptData)

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the result is the expected reader
	assert.Equal(t, expectedReader, result)

	// Verify that the mock was called as expected
	mockPDFUtil.AssertExpectations(t)
}
