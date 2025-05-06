package entities

import (
	"gitlab.com/clubhub.ai1/gommon/uid"
)

type CollectionCenterPaymentMethodEntity struct {
	CollectionCenterID *uid.UniqueID `gorm:"column:collection_center_id"`
	PaymentMethodID    *uid.UniqueID `gorm:"column:payment_method_id"`
}

func NewCollectionCenterPaymentMethodEntity(
	collectionCenterId,
	paymentMethodId *uid.UniqueID,
) CollectionCenterPaymentMethodEntity {
	return CollectionCenterPaymentMethodEntity{
		CollectionCenterID: collectionCenterId,
		PaymentMethodID:    paymentMethodId,
	}
}

func (c CollectionCenterPaymentMethodEntity) TableName() string {
	return "collection_center_payment_method"
}
