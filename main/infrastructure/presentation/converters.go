package presentation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	facilityrental "github.com/alessandro-marcantoni/cnc-backend/main/domain/facility_rental"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/membership"
	"github.com/alessandro-marcantoni/cnc-backend/main/domain/payment"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]string{"error": msg})
}

func convertPhoneNumbersToPresentation(numbers []membership.PhoneNumber) []PhoneNumber {
	phoneNumbers := make([]PhoneNumber, len(numbers))
	for i, pn := range numbers {
		prefix := ""
		if pn.Prefix != nil {
			prefix = *pn.Prefix
		}
		phoneNumbers[i] = PhoneNumber{
			Prefix: prefix,
			Number: pn.Number,
		}
	}
	return phoneNumbers
}

func convertAddressesToPresentation(a []membership.Address) []Address {
	addresses := make([]Address, len(a))
	for i, addr := range a {
		addresses[i] = Address{
			Country:      addr.Country,
			City:         addr.City,
			Street:       addr.Street,
			StreetNumber: fmt.Sprintf("%s", addr.Number),
			ZipCode:      addr.ZipCode,
		}
	}
	return addresses
}

func convertMembershipToPresentation(m membership.Membership) Membership {
	var p *Payment

	switch casted := m.Payment.(type) {
	case payment.PaymentPaid:
		p = &Payment{
			ID:             casted.ID,
			Amount:         casted.AmountPaid,
			PaidAt:         casted.PaymentDate.Format("2006-01-02T15:04:05Z07:00"),
			Currency:       casted.Currency,
			PaymentMethod:  casted.PaymentMethod,
			TransactionRef: casted.TransactionRef,
		}
	default:
		p = nil
	}

	membership := Membership{
		ID:        m.Id.Value,
		Number:    m.Number,
		Status:    string(m.Status.GetStatus()),
		ValidFrom: m.Status.GetValidFromDate().Format("2006-01-02"),
		ExpiresAt: m.Status.GetValidUntilDate().Format("2006-01-02"),
		PeriodId:  m.Status.GetPeriodId(),
		Price:     m.Price,
		Payment:   p,
	}

	return membership
}

func convertMembershipsToPresentation(memberships []membership.Membership) []Membership {
	presentationMemberships := make([]Membership, len(memberships))
	for i, m := range memberships {
		presentationMemberships[i] = convertMembershipToPresentation(m)
	}
	return presentationMemberships
}

func ConvertMemberToPresentation(domainMember membership.Member) Member {
	birthDate := ""
	if !domainMember.User.BirthDate.IsZero() {
		birthDate = domainMember.User.BirthDate.Format("2006-01-02")
	}

	return Member{
		ID:               domainMember.Id.Value,
		FirstName:        domainMember.User.FirstName,
		LastName:         domainMember.User.LastName,
		BirthDate:        birthDate,
		MembershipNumber: domainMember.Membership.Number,
		MembershipStatus: string(domainMember.Membership.Status.GetStatus()),
		Paid:             domainMember.Membership.Payment.GetStatus() == payment.Paid,
	}
}

func ConvertMembersToPresentation(domainMembers []membership.Member) []Member {
	presentationMembers := make([]Member, len(domainMembers))
	for i, dm := range domainMembers {
		presentationMembers[i] = ConvertMemberToPresentation(dm)
	}
	return presentationMembers
}

func ConvertMemberDetailsToPresentation(domainMember membership.MemberDetails) MemberDetails {
	birthDate := ""
	if !domainMember.User.BirthDate.IsZero() {
		birthDate = domainMember.User.BirthDate.Format("2006-01-02")
	}

	return MemberDetails{
		ID:           domainMember.User.Id.Value,
		FirstName:    domainMember.User.FirstName,
		LastName:     domainMember.User.LastName,
		Email:        domainMember.User.Email.Value,
		BirthDate:    birthDate,
		PhoneNumbers: convertPhoneNumbersToPresentation(domainMember.PhoneNumbers),
		Addresses:    convertAddressesToPresentation(domainMember.Addresses),
		Memberships:  convertMembershipsToPresentation(domainMember.Memberships),
	}
}

func ConvertMemberToSummary(domainMember membership.Member) MemberSummary {
	return MemberSummary{
		ID:    domainMember.User.Id.Value,
		Name:  fmt.Sprintf("%s %s", domainMember.User.FirstName, domainMember.User.LastName),
		Email: domainMember.User.Email.Value,
	}
}

func ConvertRentedFacilityToPresentation(rf facilityrental.RentedFacility) RentedFacility {
	rentedFacility := RentedFacility{
		ID:                      rf.GetId().Value,
		FacilityID:              rf.GetFacility().Id.Value,
		FacilityIdentifier:      rf.GetFacility().Identifier,
		FacilityName:            rf.GetFacility().FacilityType.FacilityName.String(),
		FacilityTypeDescription: rf.GetFacility().FacilityType.Description,
		RentedAt:                rf.GetValidity().FromDate.Format("2006-01-02"),
		Price:                   rf.GetPrice(),
		ExpiresAt:               rf.GetValidity().ToDate.Format("2006-01-02"),
		BoatInfo:                nil,
		Payment:                 nil,
	}

	// Convert payment information
	if rf.GetPayment().GetStatus() == payment.Paid {
		if paymentPaid, ok := rf.GetPayment().(payment.PaymentPaid); ok {
			rentedFacility.Payment = &Payment{
				ID:             paymentPaid.ID,
				Amount:         paymentPaid.AmountPaid,
				Currency:       paymentPaid.Currency,
				PaidAt:         paymentPaid.PaymentDate.Format(time.RFC3339),
				PaymentMethod:  paymentPaid.PaymentMethod,
				TransactionRef: paymentPaid.TransactionRef,
			}
		}
	}

	// Check if this is a boat facility
	if rf.GetType() == facilityrental.BoatFacility {
		if rfWithBoat, ok := rf.(facilityrental.RentedFacilityWithBoat); ok {
			rentedFacility.BoatInfo = &BoatInfo{
				Name:         rfWithBoat.BoatInfo.Name,
				LengthMeters: rfWithBoat.BoatInfo.LengthMeters,
				WidthMeters:  rfWithBoat.BoatInfo.WidthMeters,
			}
		}
	}

	return rentedFacility
}

func ConvertFacilityTypesToPresentation(domainFacilityTypes []facilityrental.FacilityType) []FacilityType {
	presentationFacilityTypes := make([]FacilityType, len(domainFacilityTypes))
	for i, ft := range domainFacilityTypes {
		presentationFacilityTypes[i] = FacilityType{
			ID:             ft.Id.Value,
			Name:           string(ft.FacilityName),
			Description:    ft.Description,
			SuggestedPrice: ft.SuggestedPrice,
		}
	}
	return presentationFacilityTypes
}

func ConvertFacilitiesWithStatusToPresentation(domainFacilities []facilityrental.FacilityWithStatus) []FacilityWithStatus {
	presentationFacilities := make([]FacilityWithStatus, len(domainFacilities))
	for i, f := range domainFacilities {
		var expiresAt *string
		if f.ExpiresAt != nil {
			formatted := f.ExpiresAt.Format("2006-01-02T15:04:05Z07:00")
			expiresAt = &formatted
		}

		presentationFacilities[i] = FacilityWithStatus{
			ID:                      f.Id.Value,
			FacilityTypeID:          f.FacilityTypeId.Value,
			Identifier:              f.Identifier,
			FacilityTypeName:        string(f.FacilityTypeName),
			FacilityTypeDescription: f.FacilityTypeDescription,
			SuggestedPrice:          f.SuggestedPrice,
			IsRented:                f.IsRented,
			ExpiresAt:               expiresAt,
			RentedByMemberId:        f.RentedByMemberId,
			RentedByMemberFirstName: f.RentedByMemberFirstName,
			RentedByMemberLastName:  f.RentedByMemberLastName,
		}
	}
	return presentationFacilities
}

type CreateMemberData struct {
	User             membership.User
	CreateMembership bool
	SeasonId         *int64
	Price            *float64
}

func ConvertCreateMemberRequestToDomain(req CreateMemberRequest) (CreateMemberData, error) {
	// Parse birth date
	birthDate, err := parseDate(req.BirthDate)
	if err != nil {
		return CreateMemberData{}, fmt.Errorf("invalid birth date: %w", err)
	}

	// Create email
	emailResult := membership.NewEmailAddress(req.Email)
	if !emailResult.IsSuccess() {
		return CreateMemberData{}, fmt.Errorf("invalid email: %s", emailResult.Error().Error())
	}

	// Convert phone numbers
	phoneNumbers := make([]membership.PhoneNumber, 0, len(req.PhoneNumbers))
	for _, pn := range req.PhoneNumbers {
		phoneResult := membership.NewPhoneNumber(pn.Prefix, pn.Number)
		if !phoneResult.IsSuccess() {
			return CreateMemberData{}, fmt.Errorf("invalid phone number: %s", phoneResult.Error().Error())
		}
		phoneNumbers = append(phoneNumbers, phoneResult.Value())
	}

	// Convert addresses
	addresses := make([]membership.Address, 0, len(req.Addresses))
	for _, addr := range req.Addresses {
		addresses = append(addresses, membership.Address{
			Country: addr.Country,
			City:    addr.City,
			ZipCode: addr.ZipCode,
			Street:  addr.Street,
			Number:  addr.StreetNumber,
		})
	}

	user := membership.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		BirthDate:    birthDate,
		Email:        emailResult.Value(),
		Addresses:    addresses,
		PhoneNumbers: phoneNumbers,
	}

	return CreateMemberData{
		User:             user,
		CreateMembership: req.CreateMembership,
		SeasonId:         req.SeasonId,
		Price:            req.Price,
	}, nil
}

func parseDate(dateStr string) (t time.Time, err error) {
	return time.Parse("2006-01-02", dateStr)
}
