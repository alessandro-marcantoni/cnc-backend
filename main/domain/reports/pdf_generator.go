package reports

import (
	"bytes"
)

// PDFGenerator defines the interface for generating PDF documents
type PDFGenerator interface {
	// GenerateMemberListPDF generates a PDF with the list of all members
	GenerateMemberListPDF(members []MemberSummary, seasonCode string) (*bytes.Buffer, error)

	// GenerateMemberDetailPDF generates a PDF with detailed information about a member
	GenerateMemberDetailPDF(member MemberDetail, facilities []FacilityRental, seasonCode string) (*bytes.Buffer, error)
}

// MemberSummary represents a member in the list report
type MemberSummary struct {
	ID                  int64
	FirstName           string
	LastName            string
	Email               string
	BirthDate           string
	MembershipNumber    int64
	MembershipStatus    string
	MembershipPaid      bool
	HasUnpaidFacilities bool
}

// MemberDetail represents detailed member information
type MemberDetail struct {
	ID           int64
	FirstName    string
	LastName     string
	Email        string
	BirthDate    string
	PhoneNumbers []PhoneNumber
	Addresses    []Address
	Memberships  []Membership
}

// PhoneNumber represents a phone number
type PhoneNumber struct {
	Number string
}

// Address represents a physical address
type Address struct {
	Country      string
	City         string
	Street       string
	StreetNumber string
	ZipCode      string
}

// Membership represents a membership period
type Membership struct {
	ID        int64
	Number    int64
	Status    string
	ValidFrom string
	ExpiresAt string
	Price     float64
	Paid      bool
}

// FacilityRental represents a rented facility
type FacilityRental struct {
	ID                      int64
	FacilityIdentifier      string
	FacilityName            string
	FacilityTypeDescription string
	RentedAt                string
	ExpiresAt               string
	Price                   float64
	Paid                    bool
	BoatName                string
}
