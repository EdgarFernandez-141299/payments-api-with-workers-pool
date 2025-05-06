package request

type CreateCardRequestDTO struct {
	ExpiryMonth            string                        `json:"expiry_month"`
	ExpiryYear             string                        `json:"expiry_year"`
	CardHolder             string                        `json:"card_holder"`
	CardNumber             string                        `json:"card_number"`
	CardCvv                string                        `json:"card_cvv"`
	AddressFirst           string                        `json:"address1"`
	Zip                    string                        `json:"zip"`
	City                   string                        `json:"city"`
	State                  string                        `json:"state"`
	Country                string                        `json:"country"`
	Phone                  string                        `json:"phone"`
	CardHolderDni          string                        `json:"card_holder_dni"`
	CardVerificationConfig CardVerificationRequestConfig `json:"card_verification_config"`
}

type DeleteCardRequestDTO struct {
	CardId string `json:"card_id"`
	UserId string `json:"user_id"`
}

type CardVerificationRequestConfig struct {
	InvokeCardVerification bool `json:"invoke_card_verification"`
}
