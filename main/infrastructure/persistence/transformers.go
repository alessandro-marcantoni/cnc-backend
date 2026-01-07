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
				Currency:    m.Payment.Currency,
			}
		} else {
			paymentInfo = payment.PaymentUnpaid{}
		}

		var membershipStatus membership.MembershipInfo
		switch {
		case m.Status == "ACTIVE" && !m.Payment.PaidAt.Time.IsZero():
			membershipStatus = membership.Active{
				ValidFromDate:  m.ValidFrom.Time,
				ValidUntilDate: m.ExpiresAt.Time,
			}
		case m.Status == "EXCLUSION_DELIBERATED":
			deliberatedAt := time.Time{}
			if m.ExclusionDeliberatedAt != nil {
				deliberatedAt = m.ExclusionDeliberatedAt.Time
			}
			membershipStatus = membership.Inactive{
				ValidFromDate:  m.ValidFrom.Time,
				ValidUntilDate: m.ExpiresAt.Time,
				ExcludedAt:     deliberatedAt,
			}
		case m.Status == "EXCLUDED":
			excludedAt := time.Time{}
			if m.ExcludedAt != nil {
				excludedAt = m.ExcludedAt.Time
			}
			membershipStatus = membership.Inactive{
				ValidFromDate:  m.ValidFrom.Time,
				ValidUntilDate: m.ExpiresAt.Time,
				ExcludedAt:     excludedAt,
			}
		default:
			// Default to Active if unknown status
			membershipStatus = membership.Active{
				ValidFromDate:  m.ValidFrom.Time,
				ValidUntilDate: m.ExpiresAt.Time,
			}
		}

		domainMemberships = append(domainMemberships, membership.Membership{
			Id:      domain.Id[membership.Membership]{Value: m.MembershipID},
			Number:  m.MembershipNumber,
			Status:  membershipStatus,
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
	// Check if member has a membership
	if queryResult.MembershipNumber == nil {
		return result.Err[membership.Member](errors.RepositoryError{Description: "member has no membership"})
	}

	// Determine membership status based on status string
	var membershipStatus membership.MembershipInfo
	switch {
	case (queryResult.Season == "PAST" || queryResult.Season == "CURRENT") && queryResult.ExclusionDeliberatedAt != nil:
		membershipStatus = membership.Inactive{
			ValidFromDate:  queryResult.SeasonStartsAt,
			ValidUntilDate: queryResult.SeasonEndsAt,
			ExcludedAt:     *queryResult.ExclusionDeliberatedAt,
		}
	default:
		membershipStatus = membership.Active{
			ValidFromDate:  queryResult.SeasonStartsAt,
			ValidUntilDate: queryResult.SeasonEndsAt,
		}
	}

	var paymentInfo payment.Payment
	switch {
	case queryResult.AmountPaid != nil && queryResult.PaidAt != nil:
		paymentInfo = payment.PaymentPaid{
			AmountPaid:  *queryResult.AmountPaid,
			PaymentDate: *queryResult.PaidAt,
			Currency:    *queryResult.Currency,
		}
	default:
		paymentInfo = payment.PaymentUnpaid{}
	}

	domainMembership := membership.Membership{
		Id:      domain.Id[membership.Membership]{Value: queryResult.MemberID},
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
			Email:        membership.EmailAddress{Value: "sample.email@example.com"},
			Addresses:    []membership.Address{},
			PhoneNumbers: []membership.PhoneNumber{},
		},
		Membership: domainMembership,
	})
}

func MapToMemberFromQueryBySeason(queryResult GetMembersBySeasonQueryResult) result.Result[membership.Member] {
	var membershipStatus membership.MembershipInfo
	switch {
	case queryResult.ExclusionDeliberatedAt != nil:
		membershipStatus = membership.Inactive{
			ValidFromDate:  queryResult.SeasonStartsAt,
			ValidUntilDate: queryResult.SeasonEndsAt,
			ExcludedAt:     *queryResult.ExclusionDeliberatedAt,
		}
	default:
		membershipStatus = membership.Active{
			ValidFromDate:  queryResult.SeasonStartsAt,
			ValidUntilDate: queryResult.SeasonEndsAt,
		}
	}

	var paymentInfo payment.Payment
	switch {
	case queryResult.AmountPaid != nil && queryResult.PaidAt != nil:
		paymentInfo = payment.PaymentPaid{
			AmountPaid:  *queryResult.AmountPaid,
			PaymentDate: *queryResult.PaidAt,
			Currency:    *queryResult.Currency,
		}
	default:
		paymentInfo = payment.PaymentUnpaid{}
	}

	domainMembership := membership.Membership{
		Id:      domain.Id[membership.Membership]{Value: queryResult.MemberID},
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
			Email:        membership.EmailAddress{Value: "sample.email@example.com"},
			Addresses:    []membership.Address{},
			PhoneNumbers: []membership.PhoneNumber{},
		},
		Membership: domainMembership,
	})
}
