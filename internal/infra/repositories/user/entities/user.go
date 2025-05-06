package entities

import (
	"reflect"
)

type UserEntity struct {
	ID             string `gorm:"type:varchar(14);primaryKey"`
	UserType       string
	ExternalUserID string
	Email          string
	Address        string
	Zip            string
	City           string
	State          string
	CountryCode    string
	Phone          string
	EnterpriseID   string
	EmailAlias     string
}

func (UserEntity) TableName() string {
	return "user"
}

func (u UserEntity) IsEmpty() bool {
	return reflect.DeepEqual(UserEntity{}, u)
}

func NewUserEntity(user UserEntity) UserEntity {
	return UserEntity{
		ID:             user.ID,
		UserType:       user.UserType,
		ExternalUserID: user.ExternalUserID,
		Email:          user.Email,
		Address:        user.Address,
		Zip:            user.Zip,
		City:           user.City,
		State:          user.State,
		CountryCode:    user.CountryCode,
		Phone:          user.Phone,
		EnterpriseID:   user.EnterpriseID,
		EmailAlias:     user.EmailAlias,
	}
}
