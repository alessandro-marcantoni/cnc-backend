package presentation

type PhoneNumber struct {
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
	ID             int64   `json:"id"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	PaidAt         string  `json:"paidAt"`
	PaymentMethod  string  `json:"paymentMethod,omitempty"`
	TransactionRef string  `json:"transactionRef,omitempty"`
}

type Membership struct {
	ID        int64    `json:"id"`
	Number    int64    `json:"number"`
	Status    string   `json:"status"`
	ValidFrom string   `json:"validFrom"`
	ExpiresAt string   `json:"expiresAt"`
	PeriodId  *int64   `json:"periodId"`
	Payment   *Payment `json:"payment"`
	Price     float64  `json:"price"`
}

type BoatInfo struct {
	Name         string      `json:"name"`
	LengthMeters float64     `json:"lengthMeters"`
	WidthMeters  float64     `json:"widthMeters"`
	Insurances   []Insurance `json:"insurances,omitempty"`
}

type Insurance struct {
	Provider  string `json:"provider"`
	Number    string `json:"number"`
	ExpiresAt string `json:"expiresAt"`
}

type RentedFacility struct {
	ID                      int64     `json:"id"`
	FacilityID              int64     `json:"facilityId"`
	FacilityIdentifier      string    `json:"facilityIdentifier"`
	FacilityName            string    `json:"facilityName"`
	FacilityTypeDescription string    `json:"facilityTypeDescription"`
	RentedAt                string    `json:"rentedAt"`
	ExpiresAt               string    `json:"expiresAt"`
	Price                   float64   `json:"price"`
	Payment                 *Payment  `json:"payment"`
	BoatInfo                *BoatInfo `json:"boatInfo"`
}

type MemberDetails struct {
	ID           int64         `json:"id"`
	FirstName    string        `json:"firstName"`
	LastName     string        `json:"lastName"`
	Email        string        `json:"email"`
	BirthDate    string        `json:"birthDate"`
	TaxCode      string        `json:"taxCode,omitempty"`
	PhoneNumbers []PhoneNumber `json:"phoneNumbers"`
	Addresses    []Address     `json:"addresses"`
	Memberships  []Membership  `json:"memberships"`
}

type Member struct {
	ID                  int64  `json:"id"`
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	BirthDate           string `json:"birthDate"`
	MembershipNumber    int64  `json:"membershipNumber"`
	MembershipStatus    string `json:"membershipStatus"`
	MembershipPaid      bool   `json:"membershipPaid"`
	HasUnpaidFacilities bool   `json:"hasUnpaidFacilities"`
}

type MemberSummary struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type FacilityType struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	SuggestedPrice float64 `json:"suggestedPrice"`
	HasBoat        bool    `json:"hasBoat"`
}

type FacilityWithStatus struct {
	ID                      int64   `json:"id"`
	FacilityTypeID          int64   `json:"facilityTypeId"`
	Identifier              string  `json:"identifier"`
	FacilityTypeName        string  `json:"facilityTypeName"`
	FacilityTypeDescription string  `json:"facilityTypeDescription"`
	SuggestedPrice          float64 `json:"suggestedPrice"`
	IsRented                bool    `json:"isRented"`
	ExpiresAt               *string `json:"expiresAt,omitempty"`
	RentedByMemberId        *int64  `json:"rentedByMemberId,omitempty"`
	RentedByMemberFirstName *string `json:"rentedByMemberFirstName,omitempty"`
	RentedByMemberLastName  *string `json:"rentedByMemberLastName,omitempty"`
}

type CreateMemberRequest struct {
	FirstName        string        `json:"firstName"`
	LastName         string        `json:"lastName"`
	BirthDate        string        `json:"birthDate"`
	Email            string        `json:"email"`
	TaxCode          string        `json:"taxCode,omitempty"`
	PhoneNumbers     []PhoneNumber `json:"phoneNumbers"`
	Addresses        []Address     `json:"addresses"`
	CreateMembership bool          `json:"createMembership"`
	SeasonId         *int64        `json:"seasonId"`
	Price            *float64      `json:"price"`
}

type AddMembershipRequest struct {
	SeasonId       int64   `json:"seasonId"`
	SeasonStartsAt string  `json:"seasonStartsAt"`
	SeasonEndsAt   string  `json:"seasonEndsAt"`
	Price          float64 `json:"price"`
	MemberId       int64   `json:"memberId"`
}

type RentFacilityRequest struct {
	FacilityId int64     `json:"facilityId"`
	MemberId   int64     `json:"memberId"`
	RentedAt   string    `json:"rentedAt"`
	ExpiresAt  string    `json:"expiresAt"`
	SeasonId   int64     `json:"seasonId"`
	Price      float64   `json:"price"`
	BoatInfo   *BoatInfo `json:"boatInfo,omitempty"`
}

type CreatePaymentRequest struct {
	MembershipPeriodId *int64  `json:"membershipPeriodId"`
	RentedFacilityId   *int64  `json:"rentedFacilityId"`
	Amount             float64 `json:"amount"`
	Currency           string  `json:"currency"`
	PaymentMethod      string  `json:"paymentMethod"`
	TransactionRef     *string `json:"transactionRef"`
}

type UpdatePaymentRequest struct {
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	PaymentMethod  string  `json:"paymentMethod"`
	TransactionRef *string `json:"transactionRef"`
}

type PaymentResponse struct {
	ID             int64   `json:"id"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	PaidAt         string  `json:"paidAt"`
	PaymentMethod  string  `json:"paymentMethod"`
	TransactionRef *string `json:"transactionRef,omitempty"`
}

type WaitingListEntry struct {
	ID             int64  `json:"id"`
	MemberId       int64  `json:"memberId"`
	FacilityTypeId int64  `json:"facilityTypeId"`
	QueuedAt       string `json:"queuedAt"`
	Notes          string `json:"notes,omitempty"`
}

type WaitingList struct {
	FacilityTypeId int64              `json:"facilityTypeId"`
	Entries        []WaitingListEntry `json:"entries"`
}

type AddToWaitingListRequest struct {
	MemberId       int64  `json:"memberId"`
	FacilityTypeId int64  `json:"facilityTypeId"`
	Notes          string `json:"notes,omitempty"`
}
