package presentation

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain/membership"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]string{"error": msg})
}

func ConvertMemberToPresentation(domainMember membership.Member) Member {
	phoneNumbers := make([]PhoneNumber, len(domainMember.User.PhoneNumbers))
	for i, pn := range domainMember.User.PhoneNumbers {
		prefix := ""
		if pn.Prefix != nil {
			prefix = *pn.Prefix
		}
		phoneNumbers[i] = PhoneNumber{
			Prefix: prefix,
			Number: pn.Number,
		}
	}

	addresses := make([]Address, len(domainMember.User.Addresses))
	for i, addr := range domainMember.User.Addresses {
		zipCode := ""
		if addr.ZipCode != nil {
			zipCode = *addr.ZipCode
		}
		addresses[i] = Address{
			Country:      addr.Country,
			City:         addr.City,
			Street:       addr.Street,
			StreetNumber: fmt.Sprintf("%s", addr.Number),
			ZipCode:      zipCode,
		}
	}

	// Convert the current membership
	membership := Membership{
		ID:        domainMember.Membership.Number.Value,
		Number:    domainMember.Membership.Number.Value,
		Status:    string(domainMember.Membership.Status.GetStatus()),
		ValidFrom: domainMember.Membership.Status.GetValidFromDate().Format("2006-01-02"),
		ExpiresAt: domainMember.Membership.Status.GetValidUntilDate().Format("2006-01-02"),
		Payment: Payment{
			Amount:   domainMember.Membership.Payment.GetAmount(),
			Currency: "EUR", // Default currency
		},
	}

	// Set payment date if available
	if paidAt := domainMember.Membership.Payment.GetPaidAt(); !paidAt.IsZero() {
		membership.Payment.PaidAt = paidAt.Format("2006-01-02T15:04:05Z07:00")
	}

	birthDate := ""
	if !domainMember.User.BirthDate.IsZero() {
		birthDate = domainMember.User.BirthDate.Format("2006-01-02")
	}

	return Member{
		ID:          domainMember.User.Id.Value,
		FirstName:   domainMember.User.FirstName,
		LastName:    domainMember.User.LastName,
		Email:       domainMember.User.Email.Value,
		BirthDate:   birthDate,
		Addresses:   addresses,
		Memberships: membership,
	}
}

func ConvertMembersToPresentation(domainMembers []membership.Member) []Member {
	presentationMembers := make([]Member, len(domainMembers))
	for i, dm := range domainMembers {
		presentationMembers[i] = ConvertMemberToPresentation(dm)
	}
	return presentationMembers
}

func ConvertMemberToSummary(domainMember membership.Member) MemberSummary {
	return MemberSummary{
		ID:    domainMember.User.Id.Value,
		Name:  fmt.Sprintf("%s %s", domainMember.User.FirstName, domainMember.User.LastName),
		Email: domainMember.User.Email.Value,
	}
}
