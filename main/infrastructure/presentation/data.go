package presentation

type PhoneNumber struct {
	Prefix string `json:"prefix"`
	Number string `json:"number"`
}

type Address struct {
	Country      string `json:"country"`
	City         string `json:"city"`
	Street       string `json:"street"`
	StreetNumber string `json:"streetNumber"`
	ZipCode      string `json:"zipCode"`
}

type Payment struct {
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	PaidAt         string  `json:"paidAt"`
	PaymentMethod  string  `json:"paymentMethod,omitempty"`
	TransactionRef string  `json:"transactionRef,omitempty"`
}

type Membership struct {
	ID        int64   `json:"id"`
	Number    int64   `json:"number"`
	Status    string  `json:"status"`
	ValidFrom string  `json:"validFrom"`
	ExpiresAt string  `json:"expiresAt"`
	Payment   Payment `json:"payment"`
}

type BoatInfo struct {
	Name         string  `json:"name"`
	LengthMeters float64 `json:"lengthMeters"`
	WidthMeters  float64 `json:"widthMeters"`
}

type RentedFacility struct {
	ID                      int64     `json:"id"`
	FacilityID              int64     `json:"facilityId"`
	FacilityIdentifier      string    `json:"facilityIdentifier"`
	FacilityName            string    `json:"facilityName"`
	FacilityTypeDescription string    `json:"facilityTypeDescription"`
	RentedAt                string    `json:"rentedAt"`
	ExpiresAt               string    `json:"expiresAt"`
	Payment                 *Payment  `json:"payment"`
	BoatInfo                *BoatInfo `json:"boatInfo"`
}

type MemberDetails struct {
	ID           int64         `json:"id"`
	FirstName    string        `json:"firstName"`
	LastName     string        `json:"lastName"`
	Email        string        `json:"email"`
	BirthDate    string        `json:"birthDate"`
	PhoneNumbers []PhoneNumber `json:"phoneNumbers"`
	Addresses    []Address     `json:"addresses"`
	Memberships  []Membership  `json:"memberships"`
}

type Member struct {
	ID         int64      `json:"id"`
	FirstName  string     `json:"firstName"`
	LastName   string     `json:"lastName"`
	Email      string     `json:"email"`
	BirthDate  string     `json:"birthDate"`
	Addresses  []Address  `json:"addresses"`
	Membership Membership `json:"membership"`
}

type MemberSummary struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
