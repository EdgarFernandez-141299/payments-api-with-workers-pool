package entities

import (
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

type User struct {
	Type          value_objects.UserType
	ID            string
	userTypeAsStr string
}

func NewUser(userType value_objects.UserType, userID string) User {
	return User{
		Type:          userType,
		ID:            userID,
		userTypeAsStr: userType.String(),
	}
}

func (u User) Equals(other User) bool {
	return u.Type.Equals(other.Type) && u.ID == other.ID
}

func (u User) Validate() error {
	if u.ID == "" {
		return errors.NewInvalidUserIDError(u.ID)
	}

	if !u.Type.IsValid() {
		return errors.NewUnsupportedUserTypeError(u.userTypeAsStr)
	}

	return nil
}
