package request

import (
	"fmt"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/config"
)

type CreateUserRequestDTO struct {
	Email            string `json:"email"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Phone            string `json:"phone"`
	IdentityDocument string `json:"identity_document"`
}

func NewUserEmailAlias(userID string) string {
	return fmt.Sprintf("%s@%s", userID, config.Config().App.ClubhubMainHost)
}
