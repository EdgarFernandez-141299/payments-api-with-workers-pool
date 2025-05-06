package response

type CreatedUserResponse struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}
