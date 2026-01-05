package persistence

import (
	"encoding/json"
	"time"
)

type GetMemberByIdQueryResult struct {
	MemberID       int64           `json:"member_id"`
	FirstName      string          `json:"first_name"`
	LastName       string          `json:"last_name"`
	DateOfBirth    time.Time       `json:"date_of_birth"`
	Email          string          `json:"email"`
	PhoneNumbers   json.RawMessage `json:"phone_numbers"`
	Addresses      json.RawMessage `json:"addresses"`
	Memberships    json.RawMessage `json:"memberships"`
	RentedServices json.RawMessage `json:"rented_services"`
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
