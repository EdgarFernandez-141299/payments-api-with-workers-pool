package response

type MemberResponse[T interface{}] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type MemberDTO struct {
	ID                 string                `json:"id,omitempty"`
	FirstName          string                `json:"first_name"`
	LastName           string                `json:"last_name"`
	BirthDate          string                `json:"birth_date"`
	Phones             []PhoneDTO            `json:"phones"`
	Emails             []EmailDTO            `json:"emails"`
	BillingInformation BillingInformationDTO `json:"billing_information"`
}

type PhoneDTO struct {
	PhoneType  string `json:"phone_type"`
	Number     string `json:"number"`
	IsDefault  bool   `json:"is_default"`
	CountryID  int    `json:"country_id"`
	Created_At string `json:"created_at"`
	IsFavorite bool   `json:"is_favorite"`
}

type EmailDTO struct {
	Email      string `json:"email"`
	IsDefault  bool   `json:"is_default"`
	Created_At string `json:"created_at"`
	IsCoowner  bool   `json:"is_coowner"`
	IsFavorite bool   `json:"is_favorite"`
}

type BillingInformationDTO struct {
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	CustomOwnerId          string `json:"custom_owner_id"`
	TaxationScheme         string `json:"taxation_scheme"`
	Address                string `json:"address"`
	AddressNumber          string `json:"address_number"`
	InteriorNumber         string `json:"interior_number"`
	PostalCode             string `json:"postal_code"`
	Neighbourhood          string `json:"neighbourhood"`
	Municipality           string `json:"municipality"`
	City                   string `json:"city"`
	State                  string `json:"state"`
	CountryCode            string `json:"country_code"`
	CustomAutoInvoice      string `json:"custom_auto_invoice"`
	CustomInvoiceRecipient string `json:"custom_invoice_recipient"`
	TaxIDRFC               string `json:"tax_id_rfc"`
	Usage                  string `json:"usage_type"`
	TaxationSchemeCode     string `json:"taxation_scheme_code"`
}
