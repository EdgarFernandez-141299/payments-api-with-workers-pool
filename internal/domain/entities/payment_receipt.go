package entities

import (
	"fmt"
	"github.com/shopspring/decimal"
	"gitlab.com/clubhub.ai1/gommon/uid"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/payment_receipt/command"
	"time"
)

var paymentReceiptIdPrefix = "RCPT"

type PaymentReceipt struct {
	ID                  string
	PaymentReceiptID    string
	UserID              string
	EnterpriseID        string
	Email               string
	ReferenceOrderID    string
	PaymentID           string
	PaymentStatus       string
	PaymentAmount       decimal.Decimal
	PaymentCountryCode  string
	PaymentCurrencyCode string
	PaymentMethod       enums.PaymentMethodEnum
	PaymentDate         string
	FileURL             string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (p PaymentReceipt) IsEmpty() bool {
	return p.ID == ""
}

func NewPaymentReceiptEntity(cmd command.CreatePaymentReceiptCommand) PaymentReceipt {
	id := uid.GenerateID()

	return PaymentReceipt{
		ID:                  fmt.Sprintf("%s-%s", paymentReceiptIdPrefix, id.String()),
		PaymentReceiptID:    fmt.Sprintf("%s-%s", paymentReceiptIdPrefix, id.String()),
		UserID:              cmd.UserID,
		EnterpriseID:        cmd.EnterpriseID,
		Email:               cmd.Email,
		ReferenceOrderID:    cmd.ReferenceOrderID,
		PaymentID:           cmd.PaymentID,
		PaymentStatus:       cmd.PaymentStatus,
		PaymentAmount:       cmd.PaymentAmount.Value,
		PaymentCountryCode:  cmd.PaymentCountry.Code,
		PaymentCurrencyCode: cmd.PaymentAmount.Code.Code,
		PaymentMethod:       cmd.PaymentMethod.Type,
		PaymentDate:         cmd.PaymentDate,
		CreatedAt:           cmd.CreatedAt.UTC(),
		UpdatedAt:           cmd.UpdatedAt.UTC(),
	}
}

func (p PaymentReceipt) WithReceiptURL(url string) PaymentReceipt {
	p.FileURL = url
	return p
}
