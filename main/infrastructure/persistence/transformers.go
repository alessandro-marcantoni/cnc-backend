package persistence

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/membership"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/payment"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/errors"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

// PgTimestamp handles PostgreSQL timestamp parsing
type PgTimestamp struct {
	time.Time
}

func (pt *PgTimestamp) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		pt.Time = time.Time{}
		return nil
	}

	// Try parsing with different formats
	formats := []string{
		"2006-01-02T15:04:05.999999",
		"2006-01-02T15:04:05.999999Z07:00",
		"2006-01-02T15:04:05Z07:00",
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05.999999",
		"2006-01-02 15:04:05",
	}

	var err error
	for _, format := range formats {
		pt.Time, err = time.Parse(format, s)
		if err == nil {
			return nil
		}
	}

	return err
}

func MapToMemberFromMemberByIdQuery(queryResult GetMemberByIdQueryResult) result.Result[membership.MemberDetails] {
	var phoneNumbers []membership.PhoneNumber
	if err := json.Unmarshal(queryResult.PhoneNumbers, &phoneNumbers); err != nil {
		return result.Err[membership.MemberDetails](err)
	}

	var addresses []membership.Address
	if err := json.Unmarshal(queryResult.Addresses, &addresses); err != nil {
		return result.Err[membership.MemberDetails](err)
	}

	var memberships []struct {
		MembershipID           int64        `json:"membership_id"`
		MembershipNumber       int64        `json:"membership_number"`
		ValidFrom              PgTimestamp  `json:"valid_from"`
		ExpiresAt              PgTimestamp  `json:"expires_at"`
		Status                 string       `json:"status"`
		ExclusionDeliberatedAt *PgTimestamp `json:"exclusion_deliberated_at"`
		ExcludedAt             *PgTimestamp `json:"excluded_at"`
		Payment                struct {
			Amount         float64     `json:"amount"`
			Currency       string      `json:"currency"`
			PaidAt         PgTimestamp `json:"paid_at"`
			PaymentMethod  string      `json:"payment_method"`
			TransactionRef string      `json:"transaction_ref"`
		} `json:"payment"`
	}
	if err := json.Unmarshal(queryResult.Memberships, &memberships); err != nil {
		return result.Err[membership.MemberDetails](err)
	}

	// Map memberships and rented services to domain structs
	var domainMemberships []membership.Membership
	for _, m := range memberships {
		var paymentInfo payment.Payment
		if m.Payment.Amount > 0 {
			paymentInfo = payment.PaymentPaid{
				AmountPaid:  m.Payment.Amount,
				PaymentDate: m.Payment.PaidAt.Time,
			}
		} else {
			paymentInfo = payment.PaymentUnpaid{
				AmountDue: m.Payment.Amount,
				DueDate:   m.ExpiresAt.Time,
			}
		}

		domainMemberships = append(domainMemberships, membership.Membership{
			Id:     domain.Id[membership.Membership]{Value: m.MembershipID},
			Number: m.MembershipNumber,
			Status: membership.Active{
				ValidFromDate:  m.ValidFrom.Time,
				ValidUntilDate: m.ExpiresAt.Time,
			}, // Adjust based on status
			Payment: paymentInfo,
		})
	}

	return result.Ok(membership.MemberDetails{
		User: membership.User{
			Id:           domain.Id[membership.User]{Value: queryResult.MemberID},
			FirstName:    queryResult.FirstName,
			LastName:     queryResult.LastName,
			BirthDate:    queryResult.DateOfBirth,
			Email:        membership.EmailAddress{Value: queryResult.Email},
			Addresses:    addresses,
			PhoneNumbers: phoneNumbers,
		},
		Memberships: domainMemberships, // Assuming one membership for simplicity
	})
}

func MapToMemberFromAllMembersQuery(queryResult GetAllMembersQueryResult) result.Result[membership.Member] {
	// Unmarshal addresses with intermediate struct
	var addressesRaw []struct {
		Country      string `json:"country"`
		City         string `json:"city"`
		Street       string `json:"street"`
		StreetNumber string `json:"street_number"`
	}
	if err := json.Unmarshal(queryResult.Addresses, &addressesRaw); err != nil {
		return result.Err[membership.Member](err)
	}
	addresses := make([]membership.Address, len(addressesRaw))
	for i, addr := range addressesRaw {
		addresses[i] = membership.Address{
			Country: addr.Country,
			City:    addr.City,
			Street:  addr.Street,
			Number:  addr.StreetNumber,
		}
	}

	// Check if member has a membership
	if queryResult.MembershipID == nil || queryResult.MembershipNumber == nil {
		return result.Err[membership.Member](errors.RepositoryError{Description: "member has no membership"})
	}

	// Determine membership status based on status string
	var membershipStatus membership.MembershipInfo
	if queryResult.ExpiresAt != nil && queryResult.ValidFrom != nil {
		switch {
		case queryResult.MembershipStatus != nil && *queryResult.MembershipStatus == "ACTIVE" && queryResult.PaidAt != nil:
			membershipStatus = membership.Active{
				ValidFromDate:  *queryResult.ValidFrom,
				ValidUntilDate: *queryResult.ExpiresAt,
			}
		case queryResult.MembershipStatus != nil && *queryResult.MembershipStatus == "ACTIVE":
			membershipStatus = membership.Unpaid{
				ValidFromDate:  *queryResult.ValidFrom,
				ValidUntilDate: *queryResult.ExpiresAt,
			}
		case queryResult.MembershipStatus != nil && *queryResult.MembershipStatus == "EXCLUSION_DELIBERATED":
			deliberatedAt := time.Time{}
			if queryResult.ExclusionDeliberatedAt != nil {
				deliberatedAt = *queryResult.ExclusionDeliberatedAt
			}
			membershipStatus = membership.ExclusionDeliberated{
				ValidFromDate:  *queryResult.ValidFrom,
				ValidUntilDate: *queryResult.ExpiresAt,
				DecisionDate:   deliberatedAt,
			}
		case queryResult.MembershipStatus != nil && *queryResult.MembershipStatus == "EXCLUDED":
			excludedAt := time.Time{}
			if queryResult.ExcludedAt != nil {
				excludedAt = *queryResult.ExcludedAt
			}
			membershipStatus = membership.Excluded{
				ValidFromDate:  *queryResult.ValidFrom,
				ValidUntilDate: *queryResult.ExpiresAt,
				DecisionDate:   excludedAt,
			}
		default:
			// Default to Active if unknown status
			membershipStatus = membership.Active{
				ValidFromDate:  *queryResult.ValidFrom,
				ValidUntilDate: *queryResult.ExpiresAt,
			}
		}
	} else {
		return result.Err[membership.Member](errors.RepositoryError{Description: "membership has no expiration or valid from date"})
	}

	// Create unpaid payment as default (payment info not included in this query)
	paymentInfo := payment.PaymentUnpaid{
		AmountDue: membership.SuggestedMembershipPrice,
		DueDate:   *queryResult.ExpiresAt,
	}

	domainMembership := membership.Membership{
		Id:      domain.Id[membership.Membership]{Value: *queryResult.MembershipID},
		Number:  *queryResult.MembershipNumber,
		Status:  membershipStatus,
		Payment: paymentInfo,
	}

	return result.Ok(membership.Member{
		User: membership.User{
			Id:           domain.Id[membership.User]{Value: queryResult.MemberID},
			FirstName:    queryResult.FirstName,
			LastName:     queryResult.LastName,
			BirthDate:    queryResult.DateOfBirth,
			Email:        membership.EmailAddress{Value: queryResult.Email},
			Addresses:    addresses,
			PhoneNumbers: []membership.PhoneNumber{},
		},
		Membership: domainMembership,
	})
}
