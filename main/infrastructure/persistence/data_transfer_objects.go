package persistence

import (
	"encoding/json"
	"time"
)

type QueryResult struct {
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
