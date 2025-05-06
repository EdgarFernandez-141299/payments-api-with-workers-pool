package entities

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

type OrderEntity struct {
	ID               string `gorm:"type:varchar(14);primaryKey"`
	UserID           string
	ReferenceOrderID string
	TotalAmount      decimal.Decimal
	CountryCode      string
	CurrencyCode     string
	Status           string
	EnterpriseID     string
	Metadata         json.RawMessage `gorm:"type:json"`
	AllowCapture     bool
	DeletedAt        *time.Time
}

func (o OrderEntity) IsEmpty() bool {
	return o.ID == ""
}

func (p OrderEntity) TableName() string {
	return "order"
}

func (p *OrderEntity) SetStatus(status string) {
	p.Status = status
}

type OrderEntityBuilder struct {
	entity OrderEntity
}

func NewOrderEntityBuilder() *OrderEntityBuilder {
	return &OrderEntityBuilder{}
}

func (p *OrderEntityBuilder) SetUserID(userID string) *OrderEntityBuilder {
	p.entity.UserID = userID
	return p
}

func (p *OrderEntityBuilder) SetID(id string) *OrderEntityBuilder {
	p.entity.ID = id
	return p
}

func (p *OrderEntityBuilder) SetReferenceOrderID(referenceOrderID string) *OrderEntityBuilder {
	p.entity.ReferenceOrderID = referenceOrderID
	return p
}

func (p *OrderEntityBuilder) SetTotalAmount(totalAmount decimal.Decimal) *OrderEntityBuilder {
	p.entity.TotalAmount = totalAmount
	return p
}

func (p *OrderEntityBuilder) SetCountryCode(countryCode string) *OrderEntityBuilder {
	p.entity.CountryCode = countryCode
	return p
}

func (p *OrderEntityBuilder) SetCurrencyCode(currencyCode string) *OrderEntityBuilder {
	p.entity.CurrencyCode = currencyCode
	return p
}

func (p *OrderEntityBuilder) SetStatus(status string) *OrderEntityBuilder {
	p.entity.Status = status
	return p
}

func (p *OrderEntityBuilder) SetEnterpriseID(enterpriseID string) *OrderEntityBuilder {
	p.entity.EnterpriseID = enterpriseID
	return p
}

func (p *OrderEntityBuilder) SetMetadata(metadata map[string]interface{}) *OrderEntityBuilder {
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		p.entity.Metadata = json.RawMessage("{}")
		return p
	}
	p.entity.Metadata = jsonData

	return p
}

func (p *OrderEntityBuilder) SetMetadataFromMap(metadata map[string]interface{}) *OrderEntityBuilder {
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		p.entity.Metadata = json.RawMessage("{}")
		return p
	}
	p.entity.Metadata = jsonData
	return p
}

func (p *OrderEntityBuilder) SetAllowCapture(allowCapture bool) *OrderEntityBuilder {
	p.entity.AllowCapture = allowCapture
	return p
}

func (p *OrderEntityBuilder) Build() OrderEntity {
	return p.entity
}
