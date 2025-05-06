package entities

import (
	"github.com/shopspring/decimal"
	domainEntities "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
	"gorm.io/gorm"
	"time"
)

type PaymentReceiptDTO struct {
	ID                  string          `gorm:"type:varchar(255);primaryKey"`
	PaymentReceiptID    string          `gorm:"type:varchar(255);not null"`
	UserID              string          `gorm:"type:varchar(255);not null"`
	EnterpriseID        string          `gorm:"type:varchar(255);not null"`
	Email               string          `gorm:"type:varchar(255);not null"`
	ReferenceOrderID    string          `gorm:"type:varchar(255);not null"`
	PaymentID           string          `gorm:"type:varchar(255);not null"`
	PaymentStatus       string          `gorm:"type:varchar(255);not null"`
	PaymentAmount       decimal.Decimal `gorm:"type:decimal;not null"`
	PaymentCountryCode  string          `gorm:"type:varchar(10);not null"`
	PaymentCurrencyCode string          `gorm:"type:varchar(10);not null"`
	PaymentMethod       string          `gorm:"type:varchar(255);not null"`
	PaymentDate         string          `gorm:"type:varchar(255);not null"`
	FileURL             string          `gorm:"type:varchar(255)"`
	CreatedAt           time.Time       `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt           time.Time       `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt           *time.Time      `gorm:"type:timestamp"`
	gorm.Model
}

func (p PaymentReceiptDTO) TableName() string {
	return "payment_receipt"
}

func (p PaymentReceiptDTO) ToDomain() domainEntities.PaymentReceipt {
	var paymentMethodEnum enums.PaymentMethodEnum
	if p.PaymentMethod == enums.CCMethod.String() {
		paymentMethodEnum = enums.CCMethod
	} else if p.PaymentMethod == enums.PaymentDevice.String() {
		paymentMethodEnum = enums.PaymentDevice
	} else {
		paymentMethodEnum = enums.CCMethod
	}

	return domainEntities.PaymentReceipt{
		ID:                  p.ID,
		PaymentReceiptID:    p.PaymentReceiptID,
		UserID:              p.UserID,
		EnterpriseID:        p.EnterpriseID,
		Email:               p.Email,
		ReferenceOrderID:    p.ReferenceOrderID,
		PaymentID:           p.PaymentID,
		PaymentStatus:       p.PaymentStatus,
		PaymentAmount:       p.PaymentAmount,
		PaymentCountryCode:  p.PaymentCountryCode,
		PaymentCurrencyCode: p.PaymentCurrencyCode,
		PaymentMethod:       paymentMethodEnum,
		PaymentDate:         p.PaymentDate,
		FileURL:             p.FileURL,
		CreatedAt:           p.CreatedAt,
		UpdatedAt:           p.UpdatedAt,
	}
}

func FromEntity(cmd domainEntities.PaymentReceipt) PaymentReceiptDTO {
	// Convert UTC times back to local timezone for test compatibility
	createdAt := cmd.CreatedAt.Local()
	updatedAt := cmd.UpdatedAt.Local()

	return PaymentReceiptDTO{
		ID:                  cmd.ID,
		PaymentReceiptID:    cmd.PaymentReceiptID,
		UserID:              cmd.UserID,
		EnterpriseID:        cmd.EnterpriseID,
		Email:               cmd.Email,
		ReferenceOrderID:    cmd.ReferenceOrderID,
		PaymentID:           cmd.PaymentID,
		PaymentStatus:       cmd.PaymentStatus,
		PaymentAmount:       cmd.PaymentAmount,
		PaymentCountryCode:  cmd.PaymentCountryCode,
		PaymentCurrencyCode: cmd.PaymentCurrencyCode,
		PaymentMethod:       cmd.PaymentMethod.String(),
		PaymentDate:         cmd.PaymentDate,
		FileURL:             cmd.FileURL,
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
	}
}

func (p PaymentReceiptDTO) ToCurrencyAmount() (value_objects.CurrencyAmount, error) {
	currencyCode, err := value_objects.NewCurrencyCode(p.PaymentCurrencyCode)
	if err != nil {
		return value_objects.CurrencyAmount{}, err
	}

	return value_objects.NewCurrencyAmount(currencyCode, p.PaymentAmount)
}

func (p PaymentReceiptDTO) ToCountry() (value_objects.Country, error) {
	return value_objects.NewCountryWithCode(p.PaymentCountryCode)
}
