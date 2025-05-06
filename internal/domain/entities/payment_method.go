package entities

import (
	"gitlab.com/clubhub.ai1/gommon/uid"
	"gorm.io/gorm"
)

type PaymentMethodEntity struct {
	ID           *uid.UniqueID `gorm:"type:varchar(14);primaryKey"`
	Name         string
	Code         string
	Description  string
	EnterpriseID string
	gorm.Model
}

func (p PaymentMethodEntity) TableName() string {
	return "payment_method"
}

func NewPaymentMethodEntity(name, code, description, enterpriseId string) PaymentMethodEntity {
	return PaymentMethodEntity{
		ID:           uid.GenerateID(),
		Name:         name,
		Code:         code,
		Description:  description,
		EnterpriseID: enterpriseId,
	}
}

type PaymentMethod struct {
	Type string `json:"type"` // CCData, TERMINAL, etc
	Card Card   `json:"card"`
}
