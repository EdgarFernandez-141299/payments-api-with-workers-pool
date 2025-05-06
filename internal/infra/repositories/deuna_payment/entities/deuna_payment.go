package entities

import "time"

type DeunaPaymentEntity struct {
	PaymentID  string    `gorm:"column:payment_id"`
	OrderToken string    `gorm:"column:order_token"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
	DeletedAt  time.Time `gorm:"column:deleted_at"`
}

func (d *DeunaPaymentEntity) TableName() string {
	return "deuna_payment"
}

func NewDeunaPaymentEntity(paymentID, orderToken string) DeunaPaymentEntity {
	return DeunaPaymentEntity{
		PaymentID:  paymentID,
		OrderToken: orderToken,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}

func (e *DeunaPaymentEntity) IsEmptyToken() bool {
	return e.OrderToken == ""
}
