package paymentmethods

type TokenCardDTO struct {
	Token    string           `json:"token"`
	CVV      string           `json:"cvv"`
	SaveCard bool             `json:"save_card"`
	Card     TokenCardDataDTO `json:"card"`
}

type TokenCardDataDTO struct {
	Brand    string `json:"brand"`
	Last4    string `json:"last4"`
	Exp      string `json:"exp"`
	CardType string `json:"card_type"`
	Alias    string `json:"alias"`
}
