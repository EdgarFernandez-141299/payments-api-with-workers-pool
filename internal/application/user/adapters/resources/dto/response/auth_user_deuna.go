package response

type DeunaAuthResponseDTO struct {
	RefreshToken string `json:"refresh_token"`
	AuthToken    string `json:"token"`
}
