package request

type CreateUserRequest struct {
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"`
}
