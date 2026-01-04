package persistence

import (
	"encoding/json"
	"time"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/membership"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/payment"
	"github.com/alessandro-marcantoni/cnc-backend/main/shared/result"
)

func MapToStructs(queryResult QueryResult) result.Result[membership.Member] {
	var phoneNumbers []membership.PhoneNumber
	if err := json.Unmarshal(queryResult.PhoneNumbers, &phoneNumbers); err != nil {
		return result.Err[membership.Member](err)
	}

	var addresses []membership.Address
	if err := json.Unmarshal(queryResult.Addresses, &addresses); err != nil {
		return result.Err[membership.Member](err)
	}

	var memberships []struct {
		MembershipID           int64      `json:"membership_id"`
		ValidFrom              time.Time  `json:"valid_from"`
		ExpiresAt              time.Time  `json:"expires_at"`
		Status                 string     `json:"status"`
		ExclusionDeliberatedAt *time.Time `json:"exclusion_deliberated_at"`
		ExcludedAt             *time.Time `json:"excluded_at"`
		Payment                struct {
			Amount         float64   `json:"amount"`
			Currency       string    `json:"currency"`
			PaidAt         time.Time `json:"paid_at"`
			PaymentMethod  string    `json:"payment_method"`
			TransactionRef string    `json:"transaction_ref"`
		} `json:"payment"`
	}
	if err := json.Unmarshal(queryResult.Memberships, &memberships); err != nil {
		return result.Err[membership.Member](err)
	}

	var rentedServices []struct {
		RentedFacilityID   int64     `json:"rented_facility_id"`
		FacilityIdentifier string    `json:"facility_identifier"`
		FacilityName       string    `json:"facility_name"`
		RentedAt           time.Time `json:"rented_at"`
		ExpiresAt          time.Time `json:"expires_at"`
		Payment            struct {
			Amount         float64   `json:"amount"`
			Currency       string    `json:"currency"`
			PaidAt         time.Time `json:"paid_at"`
			PaymentMethod  string    `json:"payment_method"`
			TransactionRef string    `json:"transaction_ref"`
		} `json:"payment"`
	}
	if err := json.Unmarshal(queryResult.RentedServices, &rentedServices); err != nil {
		return result.Err[membership.Member](err)
	}

	// Map memberships and rented services to domain structs
	var domainMemberships []membership.Membership
	for _, m := range memberships {
		var paymentInfo payment.Payment
		if m.Payment.Amount > 0 {
			paymentInfo = payment.PaymentPaid{
				AmountPaid:  m.Payment.Amount,
				PaymentDate: m.Payment.PaidAt,
			}
		} else {
			paymentInfo = payment.PaymentUnpaid{
				AmountDue: m.Payment.Amount,
				DueDate:   m.ExpiresAt,
			}
		}

		domainMemberships = append(domainMemberships, membership.Membership{
			Number:  domain.Id[membership.Membership]{Value: m.MembershipID},
			Status:  membership.Active{ValidUntilDate: m.ExpiresAt}, // Adjust based on status
			Payment: paymentInfo,
		})
	}

	return result.Ok(membership.Member{
		User: membership.User{
			Id:           domain.Id[membership.User]{Value: queryResult.MemberID},
			FirstName:    queryResult.FirstName,
			LastName:     queryResult.LastName,
			BirthDate:    queryResult.DateOfBirth,
			Email:        membership.EmailAddress{Value: queryResult.Email},
			Addresses:    addresses,
			PhoneNumbers: phoneNumbers,
		},
		Membership: domainMemberships[0], // Assuming one membership for simplicity
	})
}
