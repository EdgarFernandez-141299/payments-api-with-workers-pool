package value_objects

type Address struct {
	ZipCode string
	Street  string
	Country Country
	City    string
}

func NewAddress(zipCode, street, city string, country Country) Address {
	return Address{
		ZipCode: zipCode,
		Street:  street,
		Country: country,
		City:    city,
	}
}

func (a Address) Equals(other Address) bool {
	return a.ZipCode == other.ZipCode &&
		a.Street == other.Street &&
		a.Country == other.Country &&
		a.City == other.City
}
