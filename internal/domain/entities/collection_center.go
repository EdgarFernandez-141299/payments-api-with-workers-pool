package entities

import (
	"github.com/lib/pq"
	"gitlab.com/clubhub.ai1/gommon/uid"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_center/dto/request"
	"gorm.io/gorm"
)

type CollectionCenterEntity struct {
	ID                  *uid.UniqueID `gorm:"type:varchar(14);primaryKey"`
	Name                string
	AvailableCurrencies pq.StringArray `gorm:"type:text[]"` //
	Description         string
	EnterpriseID        string
	gorm.Model
}

func (CollectionCenterEntity) TableName() string {
	return "collection_center"
}

func NewCollectionCenterEntity(
	collectionCenter request.CollectionCenterRequest,
	enterpriseId string,
) CollectionCenterEntity {
	id := uid.GenerateID()

	return CollectionCenterEntity{
		ID:                  id,
		Name:                collectionCenter.Name,
		AvailableCurrencies: collectionCenter.AvailableCurrencies,
		Description:         collectionCenter.Description,
		EnterpriseID:        enterpriseId,
	}
}

func (c CollectionCenterEntity) IsEmpty() bool {
	return c.ID == nil
}
