package entities

import (
	"time"

	"gitlab.com/clubhub.ai1/gommon/uid"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/card/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/card/dto/response"

	"gorm.io/gorm"
)

type CardEntity struct {
	ID                *uid.UniqueID `gorm:"column:id;type:varchar(36);primaryKey"`
	ExternalCardID    string        `gorm:"column:external_card_id;type:varchar(36)"`
	UserID            string        `gorm:"column:user_id;type:varchar(36)"`
	CardHolder        string        `gorm:"column:card_holder;type:varchar(150)"`
	Alias             string        `gorm:"column:alias;type:varchar(50)"`
	Bin               string        `gorm:"column:bin;type:varchar(16)"`
	LastFour          string        `gorm:"column:last_four;type:varchar(4)"`
	Brand             string        `gorm:"column:brand;type:varchar(25)"`
	ExpirationDate    time.Time     `gorm:"column:expiration_date;type:date"`
	CardType          string        `gorm:"column:card_type;type:varchar(30)"`
	Status            string        `gorm:"column:status;type:varchar(25)"`
	IsDefault         bool          `gorm:"column:is_default;type:boolean"`
	IsRecurrent       bool          `gorm:"column:is_recurrent;type:boolean"`
	RetryAttempts     int           `gorm:"column:retry_attempts;type:integer"`
	EnterpriseID      string        `gorm:"column:enterprise_id;type:varchar(36)"`
	CardFailureReason string        `gorm:"column:card_failure_reason;type:varchar(100)"`
	CardFailureCode   string        `gorm:"column:card_failure_code;type:varchar(20)"`
	gorm.Model
}

type CardEntities []CardEntity

func (CardEntity) TableName() string {
	return "card"
}

func NewCard(
	request request.CardRequest,
	enterpriseID string,
) CardEntity {

	layout := "01/06"
	expirationDateParsed, _ := time.Parse(layout, request.ExpirationDate) //cardDeUNA.ExpirationDate)

	year, month, _ := expirationDateParsed.Date()
	expirationDate := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)

	return CardEntity{
		ID:             uid.GenerateID(),
		ExternalCardID: request.CardId, //cardDeUNA.ID,
		UserID:         request.UserID, //userID,
		//CardHolder:        request.CardHolder,
		Alias:             request.Alias,
		Bin:               request.FirstSix,
		LastFour:          request.LastFour,  //cardDeUNA.LastFour,
		Brand:             request.CardBrand, //cardDeUNA.Company,
		ExpirationDate:    expirationDate,
		CardType:          request.CardType,
		Status:            request.Status,
		IsDefault:         request.IsDefault,
		IsRecurrent:       request.IsRecurrent,
		EnterpriseID:      enterpriseID,
		CardFailureReason: "",
		CardFailureCode:   "",
	}
}

func (c CardEntity) ToDTO() response.CardResponse {
	return response.CardResponse{
		ID:             c.ID.String(),
		CardTokenID:    c.ExternalCardID,
		Alias:          c.Alias,
		LastFour:       c.LastFour,
		Brand:          c.Brand,
		IsDefault:      c.IsDefault,
		IsRecurrent:    c.IsRecurrent,
		ExpirationDate: c.ExpirationDate.Format("01/06"),
		CardType:       c.CardType,
	}
}

func (c CardEntities) ToDTO() []response.CardResponse {
	cards := []response.CardResponse{}

	for _, card := range c {
		cards = append(cards, card.ToDTO())
	}

	return cards
}

// type Card struct {
// 	ID  string `json:"id"`
// 	CVV string `json:"cvv"`
// }
