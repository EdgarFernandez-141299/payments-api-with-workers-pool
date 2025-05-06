package response

type CreatedUserResponse struct {
	ID                    string `json:"id"`
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	Email                 string `json:"email"`
	PaymentsExternalEmail string `json:"payments_external_email"`
	ExternalUserID        string `json:"external_user_id"`
}

type UserValidatedResponse struct {
	ID                string `json:"id"`
	Email             string `json:"email"`
	ExternalUserID    string `json:"external_user_id"`
	ExternalAuthToken string `json:"external_auth_token"`
}
