package entities

import (
	"reflect"
	"time"

	"gitlab.com/clubhub.ai1/gommon/uid"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_account/dto/request"
	"gorm.io/gorm"
)

type CollectionAccountEntity struct {
	gorm.Model
	ID                     *uid.UniqueID `gorm:"type:varchar(14);primaryKey"`
	AccountType            string
	CollectionCenterID     string
	CurrencyCode           string
	AccountNumber          string
	BankName               string
	InterbankAccountNumber string
	EnterpriseID           string
}

func (c CollectionAccountEntity) IsEmpty() bool {
	return reflect.DeepEqual(c, CollectionAccountEntity{})
}

func (CollectionAccountEntity) TableName() string {
	return "collection_account"
}

func NewCollectionAccountEntity(
	request request.CollectionAccountRequest,
	enterpriseId string,
) CollectionAccountEntity {
	return CollectionAccountEntity{
		ID:                     uid.GenerateID(),
		AccountType:            request.AccountType.String(),
		CollectionCenterID:     request.CollectionCenterID,
		CurrencyCode:           request.CurrencyCode,
		AccountNumber:          request.AccountNumber,
		BankName:               request.BankName,
		EnterpriseID:           enterpriseId,
		InterbankAccountNumber: request.InterbankAccountNumber,
		Model: gorm.Model{
			CreatedAt: time.Now(),
		},
	}
}
