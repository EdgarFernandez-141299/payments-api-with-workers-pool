package response

type AirlineInformation struct {
	BookingItems *[]BookingItem `json:"booking_items"`
}

type BookingItem struct {
	Airline3DigitCode     string       `json:"airline_3_digit_code"`
	AirlineIataDesignator string       `json:"airline_iata_designator"`
	ETicket               bool         `json:"e_ticket"`
	Legs                  []Leg        `json:"legs"`
	Passenger             Passenger    `json:"passenger"`
	PNR                   string       `json:"pnr"`
	TicketNumber          string       `json:"ticket_number"`
	TicketingTravelAgency TravelAgency `json:"ticketing_travel_agency"`
}

type Leg struct {
	CarrierCode   string      `json:"carrier_code"`
	CarrierName   string      `json:"carrier_name"`
	Destination   FlightPoint `json:"destination"`
	FareBasisCode string      `json:"fare_basis_code"`
	FlightNumber  string      `json:"flight_number"`
	Origin        FlightPoint `json:"origin"`
	SeatLocation  string      `json:"seat_location"`
	StopoverCode  string      `json:"stopover_code"`
}

type FlightPoint struct {
	Date     string `json:"date"`
	IataCode string `json:"iata_code"`
	Time     string `json:"time"`
}

type Passenger struct {
	DateOfBirth            string `json:"date_of_birth"`
	DocumentType           string `json:"document_type"`
	Email                  string `json:"email"`
	FirstName              string `json:"first_name"`
	FrequentFlyerCode      string `json:"frequent_flyer_code"`
	IdentityDocumentNumber string `json:"identity_document_number"`
	LastName               string `json:"last_name"`
	MiddleName             string `json:"middle_name"`
	Phone                  string `json:"phone"`
	Title                  string `json:"title"`
	UserID                 string `json:"user_id"`
}

type TravelAgency struct {
	AgencyInvoiceNumber string `json:"agency_invoice_number"`
	IataCode            string `json:"iata_code"`
	Name                string `json:"name"`
}
