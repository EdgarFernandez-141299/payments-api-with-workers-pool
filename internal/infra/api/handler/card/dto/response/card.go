package response

type CardResponse struct {
	ID             string `json:"card_id"`
	CardTokenID    string `json:"card_token_id"`
	Alias          string `json:"alias"`
	LastFour       string `json:"last_four"`
	Brand          string `json:"brand"`
	IsDefault      bool   `json:"is_default"`
	IsRecurrent    bool   `json:"is_recurrent"`
	ExpirationDate string `json:"expiration_date"`
	CardType       string `json:"card_type"`
}

type NotificationCardExpiringSoonResponseDTO struct {
	Message string `json:"message"`
}
