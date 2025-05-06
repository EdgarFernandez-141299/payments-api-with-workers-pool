package response

type CardResponseDTO struct {
	Data CardResponseDataDTO `json:"data"`
}

type CardResponseDataDTO struct {
	ID                      string `json:"id"`
	UserID                  string `json:"user_id"`
	CardHolder              string `json:"card_holder"`
	CardHolderDni           string `json:"card_holder_dni"`
	Company                 string `json:"company"`
	LastFour                string `json:"last_four"`
	FirstSix                string `json:"first_six"`
	ExpirationDate          string `json:"expiration_date"`
	IsValid                 bool   `json:"is_valid"`
	IsExpired               bool   `json:"is_expired"`
	VerifiedBy              string `json:"verified_by"`
	VerifiedWithTransaction string `json:"verified_with_transaction_id"`
	VerifiedAt              string `json:"verified_at"`
	LastUsed                string `json:"last_used"`
	CreatedAt               string `json:"created_at"`
	UpdatedAt               string `json:"updated_at"`
	DeletedAt               string `json:"deleted_at"`
	BankName                string `json:"bank_name"`
	CountryIso              string `json:"country_iso"`
	CardType                string `json:"card_type"`
	Source                  string `json:"source"`
	ZipCode                 string `json:"zip_code"`
	Vault                   string `json:"vault"`
	InternalUserID          string `json:"internal_user_id"`
}

func (c CardResponseDTO) IsEmpty() bool {
	return c.Data.ID == ""
}

type DeleteCardResponseDTO struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
