package persistence

import (
	"encoding/json"
	"time"
)

type GetMemberByIdQueryResult struct {
	MemberID     int64           `json:"member_id"`
	FirstName    string          `json:"first_name"`
	LastName     string          `json:"last_name"`
	DateOfBirth  time.Time       `json:"date_of_birth"`
	Email        string          `json:"email"`
	PhoneNumbers json.RawMessage `json:"phone_numbers"`
	Addresses    json.RawMessage `json:"addresses"`
	Memberships  json.RawMessage `json:"memberships"`
}

type GetAllMembersQueryResult struct {
	MemberID               int64      `json:"member_id"`
	FirstName              string     `json:"first_name"`
	LastName               string     `json:"last_name"`
	DateOfBirth            time.Time  `json:"date_of_birth"`
	MembershipNumber       *int64     `json:"membership_number"`
	MembershipStatus       *string    `json:"membership_status"`
	Season                 *string    `json:"season"`
	SeasonStartsAt         *time.Time `json:"season_starts_at"`
	SeasonEndsAt           *time.Time `json:"season_ends_at"`
	ExclusionDeliberatedAt *time.Time `json:"exclusion_deliberated_at"`
	Price                  *float64   `json:"price"`
	AmountPaid             *float64   `json:"amount_paid"`
	PaidAt                 *time.Time `json:"paid_at"`
	Currency               *string    `json:"currency"`
}

type GetMembersBySeasonQueryResult struct {
	MemberID               int64      `json:"member_id"`
	FirstName              string     `json:"first_name"`
	LastName               string     `json:"last_name"`
	Email                  string     `json:"email"`
	DateOfBirth            time.Time  `json:"date_of_birth"`
	MembershipNumber       *int64     `json:"membership_number"`
	MembershipStatus       string     `json:"membership_status"`
	SeasonStartsAt         time.Time  `json:"season_starts_at"`
	SeasonEndsAt           time.Time  `json:"season_ends_at"`
	ExclusionDeliberatedAt *time.Time `json:"exclusion_deliberated_at"`
	Price                  *float64   `json:"price"`
	AmountPaid             *float64   `json:"amount_paid"`
	PaidAt                 *time.Time `json:"paid_at"`
	Currency               *string    `json:"currency"`
	HasUnpaidFacilities    bool       `json:"has_unpaid_facilities"`
}

type GetRentedFacilitiesByMemberQueryResult struct {
	RentedFacilityID   int64      `json:"rented_facility_id"`
	RentedAt           time.Time  `json:"rented_at"`
	ExpiresAt          time.Time  `json:"expires_at"`
	Price              float64    `json:"price"`
	FacilityID         int64      `json:"facility_id"`
	FacilityIdentifier string     `json:"facility_identifier"`
	FacilityTypeID     int64      `json:"facility_type_id"`
	FacilityType       string     `json:"facility_type"`
	FacilityTypeDesc   string     `json:"facility_type_description"`
	SuggestedPrice     float64    `json:"suggested_price"`
	BoatID             *int64     `json:"boat_id"`
	BoatName           *string    `json:"boat_name"`
	LengthMeters       *float64   `json:"length_meters"`
	WidthMeters        *float64   `json:"width_meters"`
	InsuranceID        *int64     `json:"insurance_id"`
	InsuranceProvider  *string    `json:"insurance_provider"`
	InsuranceNumber    *string    `json:"insurance_number"`
	InsuranceExpiresAt *time.Time `json:"insurance_expires_at"`
	PaymentID          *int64     `json:"payment_id"`
	PaymentAmount      *float64   `json:"payment_amount"`
	PaymentCurrency    *string    `json:"payment_currency"`
	PaymentPaidAt      *time.Time `json:"payment_paid_at"`
	PaymentMethod      *string    `json:"payment_method"`
	PaymentNotes       *string    `json:"payment_notes"`
}
