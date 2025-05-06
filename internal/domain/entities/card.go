package entities

type Card struct {
	ID                string
	ExternalCardID    string
	UserID            string
	CardHolder        string
	Alias             string
	Bin               string
	LastFour          string
	Brand             string
	ExpirationDate    string
	CardType          string
	Status            string
	IsDefault         bool
	IsRecurrent       bool
	RetryAttempts     int
	EnterpriseID      string
	CardFailureReason string
	CardFailureCode   string
}
