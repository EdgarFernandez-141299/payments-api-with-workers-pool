package entities

type UserEntity struct {
	ID             string
	UserType       string
	ExternalUserID string
	Email          string
	Address        string
	Zip            string
	City           string
	State          string
	CountryCode    string
	Phone          string
	EnterpriseID   string
}
