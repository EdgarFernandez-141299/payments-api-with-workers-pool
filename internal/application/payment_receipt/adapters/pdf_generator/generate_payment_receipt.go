package pdf_generator

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_receipt/interfaces"
)

const (
	receiptPaymentHtmlPath = "assets/"
	filename               = "receipt-payment.html"
)

type ReceiptPaymentGeneratorImpl struct {
	pdfUtil interfaces.PDFUtilInterface
}

// NewReceiptPaymentGeneratorImpl creates a new ReceiptPaymentGeneratorImpl
func NewReceiptPaymentGeneratorImpl() interfaces.ReceiptPaymentGenerator {
	baseDir, err := os.Getwd()
	if err != nil {
		baseDir = "."
	}

	templatePath := filepath.Join(baseDir, receiptPaymentHtmlPath)

	return &ReceiptPaymentGeneratorImpl{
		pdfUtil: NewPDFUtil(templatePath),
	}
}

func (g *ReceiptPaymentGeneratorImpl) GenerateReceiptPaymentPDF(ctx context.Context, receiptData interfaces.ReceiptData) (io.Reader, error) {
	dataMap := receiptData.ToMap()

	return g.pdfUtil.GeneratePDF(ctx, filename, dataMap)
}
