package entities

import (
	"reflect"
	"time"

	"gitlab.com/clubhub.ai1/gommon/uid"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/payment_concept/dto/request"
	"gorm.io/gorm"
)

type PaymentConceptEntity struct {
	gorm.Model
	ID           *uid.UniqueID `gorm:"type:varchar(14);primaryKey"`
	Name         string
	Code         string
	Description  string
	EnterpriseID string
}

func (p PaymentConceptEntity) IsEmpty() bool {
	return reflect.DeepEqual(p, PaymentConceptEntity{})
}

func (PaymentConceptEntity) TableName() string {
	return "payment_concept"
}

func NewPaymentConceptEntity(
	request request.PaymentConceptRequest,
	enterpriseId string,
) PaymentConceptEntity {
	return PaymentConceptEntity{
		ID:           uid.GenerateID(),
		Name:         request.Name,
		Code:         request.Code,
		Description:  request.Description,
		EnterpriseID: enterpriseId,
		Model: gorm.Model{
			CreatedAt: time.Now(),
		},
	}
}
