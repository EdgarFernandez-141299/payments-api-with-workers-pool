package entities

import (
	"time"

	"gitlab.com/clubhub.ai1/gommon/uid"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_account_route/dto/request"
	"gorm.io/gorm"
)

type CollectionAccountRouteEntity struct {
	ID                  *uid.UniqueID `gorm:"type:varchar(14);primaryKey"`
	CollectionAccountID string
	CountryCode         string
	CurrencyCode        string
	EnterpriseID        string
	DisabledAt          *time.Time
	AssociatedOrigin    string
	gorm.Model
}

func (c CollectionAccountRouteEntity) IsEmpty() bool {
	return c.ID == nil
}

func (CollectionAccountRouteEntity) TableName() string {
	return "collection_account_route"
}

func NewCollectionAccountRouteEntity(
	request request.CollectionAccountRouteRequest, enterpriseId string,
) CollectionAccountRouteEntity {
	return CollectionAccountRouteEntity{
		ID:                  uid.GenerateID(),
		CollectionAccountID: request.CollectionAccountID,
		CountryCode:         request.CountryCode,
		CurrencyCode:        request.CurrencyCode,
		EnterpriseID:        enterpriseId,
		AssociatedOrigin:    request.AssociatedOrigin.String(),
	}
}
