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
	MemberID               int64           `json:"member_id"`
	FirstName              string          `json:"first_name"`
	LastName               string          `json:"last_name"`
	DateOfBirth            time.Time       `json:"date_of_birth"`
	Email                  string          `json:"email"`
	Addresses              json.RawMessage `json:"addresses"`
	MembershipID           *int64          `json:"membership_id"`
	MembershipNumber       *int64          `json:"membership_number"`
	MembershipPeriodID     *int64          `json:"membership_period_id"`
	ValidFrom              *time.Time      `json:"valid_from"`
	ExpiresAt              *time.Time      `json:"expires_at"`
	MembershipStatus       *string         `json:"membership_status"`
	ExclusionDeliberatedAt *time.Time      `json:"exclusion_deliberated_at"`
	ExcludedAt             *time.Time      `json:"excluded_at"`
	PaidAt                 *time.Time      `json:"paid_at"`
}

type GetRentedFacilitiesByMemberQueryResult struct {
	RentedFacilityID   int64     `json:"rented_facility_id"`
	RentedAt           time.Time `json:"rented_at"`
	ExpiresAt          time.Time `json:"expires_at"`
	FacilityID         int64     `json:"facility_id"`
	FacilityIdentifier string    `json:"facility_identifier"`
	FacilityTypeID     int64     `json:"facility_type_id"`
	FacilityType       string    `json:"facility_type"`
	FacilityTypeDesc   string    `json:"facility_type_description"`
	SuggestedPrice     float64   `json:"suggested_price"`
	BoatID             *int64    `json:"boat_id"`
	BoatName           *string   `json:"boat_name"`
	LengthMeters       *float64  `json:"length_meters"`
	WidthMeters        *float64  `json:"width_meters"`
}
