package interfaces

import (
	"context"
	"io"
)

// PDFUtilInterface defines the interface for PDFUtil
type PDFUtilInterface interface {
	GeneratePDF(ctx context.Context, htmlPath string, data map[string]string) (io.Reader, error)
}

// ReceiptData represents the data needed for a payment receipt
type ReceiptData struct {
	ReceiptNumber    string
	TransactionID    string
	AmountPaid       string
	PaymentDate      string
	PaymentMethod    string
	PaymentReference string
}

// ToMap converts ReceiptData to map[string]string
func (r ReceiptData) ToMap() map[string]string {
	return map[string]string{
		"ReceiptNumber":    r.ReceiptNumber,
		"TransactionID":    r.TransactionID,
		"AmountPaid":       r.AmountPaid,
		"PaymentDate":      r.PaymentDate,
		"PaymentMethod":    r.PaymentMethod,
		"PaymentReference": r.PaymentReference,
	}
}

// ReceiptPaymentGenerator defines the interface for generating payment receipts
type ReceiptPaymentGenerator interface {
	GenerateReceiptPaymentPDF(ctx context.Context, receiptData ReceiptData) (io.Reader, error)
}