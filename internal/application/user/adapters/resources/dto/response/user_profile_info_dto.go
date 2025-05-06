package response

type UserProfileInfoDTO struct {
	UserID             string         `json:"id"`
	ClientID           string         `json:"client_id"`
	FirstName          string         `json:"first_name"`
	LastName           string         `json:"last_name"`
	Email              string         `json:"email"`
	EnterpriseID       string         `json:"enterprise_id"`
	PrimaryPhone       PhoneInfo      `json:"primary_phone"`
	SecondaryPhone     PhoneInfo      `json:"secondary_phone"`
	PreferenceLanguage string         `json:"preference_language"`
	Contract           ContractInfo   `json:"contract"`
	Address            AddressInfo    `json:"address"`
	Travelers          []TravelerInfo `json:"travelers"`
}

type PhoneInfo struct {
	CountryCode string `json:"country_code"`
	Number      string `json:"number"`
	Type        string `json:"type"`
}

type ContractInfo struct {
	ClubID           string `json:"club_id"`
	Currency         string `json:"currency"`
	CurrencyID       string `json:"currency_id"`
	MembershipID     string `json:"membership_id"`
	ContractStatusID int    `json:"contract_status_id"`
	ContractDate     string `json:"contract_date"`
	ExpirationDate   string `json:"expiration_date"`
}

type AddressInfo struct {
	ZipCode     string `json:"zip_code"`
	State       string `json:"state"`
	AddressLine string `json:"address_line"`
	Country     string `json:"country"`
	City        string `json:"city"`
}

type TravelerInfo struct {
	ID               string `json:"id"`
	MemberID         string `json:"member_id"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	PhoneCountryCode string `json:"phone_country_code"`
	PhoneNumber      string `json:"phone_number"`
}
