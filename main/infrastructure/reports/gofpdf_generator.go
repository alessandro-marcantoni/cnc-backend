package reports

import (
	"bytes"
	"fmt"
	"time"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain/reports"
	"github.com/jung-kurt/gofpdf"
)

// GoPDFGenerator implements PDFGenerator using the gofpdf library
type GoPDFGenerator struct{}

// NewGoPDFGenerator creates a new PDF generator
func NewGoPDFGenerator() *GoPDFGenerator {
	return &GoPDFGenerator{}
}

// GenerateMemberListPDF generates a PDF with the list of all members
func (g *GoPDFGenerator) GenerateMemberListPDF(members []reports.MemberSummary, seasonCode string) (*bytes.Buffer, error) {
	pdf := gofpdf.New("L", "mm", "A4", "") // Landscape for better table fit
	pdf.SetFont("Arial", "", 10)

	// Add page
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "Lista Soci - Circolo Nautico", "", 1, "C", false, 0, "")

	// Subtitle with season and date
	pdf.SetFont("Arial", "", 10)
	currentDate := time.Now().Format("02/01/2006")
	subtitle := fmt.Sprintf("Stagione %s - Generato il %s", seasonCode, currentDate)
	pdf.CellFormat(0, 8, subtitle, "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Table header
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(200, 200, 200)

	// Column widths (total should be ~277mm for A4 landscape)
	colWidths := []float64{30, 70, 35, 80, 30, 32}
	headers := []string{"N. Tessera", "Nome", "Data di Nascita", "Email", "Pag. Tessera", "Pag. Servizi"}

	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 8, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Table rows
	pdf.SetFont("Arial", "", 8)
	pdf.SetFillColor(240, 240, 240)
	fill := false

	for _, member := range members {
		// Alternate row colors
		if fill {
			pdf.SetFillColor(245, 245, 245)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		// Check if we need a new page
		if pdf.GetY() > 180 {
			pdf.AddPage()
			// Re-print header
			pdf.SetFont("Arial", "B", 9)
			pdf.SetFillColor(200, 200, 200)
			for i, header := range headers {
				pdf.CellFormat(colWidths[i], 8, header, "1", 0, "C", true, 0, "")
			}
			pdf.Ln(-1)
			pdf.SetFont("Arial", "", 8)
		}

		// Row data
		fullName := member.LastName + " " + member.FirstName

		pdf.CellFormat(colWidths[0], 7, fmt.Sprintf("%d", member.MembershipNumber), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colWidths[1], 7, fullName, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(colWidths[2], 7, member.BirthDate, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colWidths[3], 7, member.Email, "1", 0, "L", fill, 0, "")

		membershipPaidText := "No"
		if member.MembershipPaid {
			membershipPaidText = "Si"
		}
		pdf.CellFormat(colWidths[4], 7, membershipPaidText, "1", 0, "C", fill, 0, "")

		facilitiesPaidText := "Si"
		if member.HasUnpaidFacilities {
			facilitiesPaidText = "No"
		}
		pdf.CellFormat(colWidths[5], 7, facilitiesPaidText, "1", 0, "C", fill, 0, "")

		pdf.Ln(-1)
		fill = !fill
	}

	// Footer with total count
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 8, fmt.Sprintf("Totale Soci: %d", len(members)), "", 1, "L", false, 0, "")

	// Write to buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return &buf, nil
}

// GenerateMemberDetailPDF generates a PDF with detailed information about a member
func (g *GoPDFGenerator) GenerateMemberDetailPDF(member reports.MemberDetail, facilities []reports.FacilityRental, seasonCode string) (*bytes.Buffer, error) {
	pdf := gofpdf.New("P", "mm", "A4", "") // Portrait for detail view
	pdf.SetFont("Arial", "", 10)

	// Add page
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "Dettaglio Socio", "", 1, "C", false, 0, "")

	// Subtitle with date
	pdf.SetFont("Arial", "", 10)
	currentDate := time.Now().Format("02/01/2006")
	subtitle := fmt.Sprintf("Generato il %s", currentDate)
	pdf.CellFormat(0, 8, subtitle, "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Personal Information Section
	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(220, 220, 220)
	pdf.CellFormat(0, 8, "Informazioni Personali", "1", 1, "L", true, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(50, 7, "Nome:", "1", 0, "L", false, 0, "")
	pdf.CellFormat(0, 7, member.FirstName, "1", 1, "L", false, 0, "")
	pdf.CellFormat(50, 7, "Cognome:", "1", 0, "L", false, 0, "")
	pdf.CellFormat(0, 7, member.LastName, "1", 1, "L", false, 0, "")
	pdf.CellFormat(50, 7, "Email:", "1", 0, "L", false, 0, "")
	pdf.CellFormat(0, 7, member.Email, "1", 1, "L", false, 0, "")
	pdf.CellFormat(50, 7, "Data di Nascita:", "1", 0, "L", false, 0, "")
	pdf.CellFormat(0, 7, member.BirthDate, "1", 1, "L", false, 0, "")

	// Phone Numbers
	if len(member.PhoneNumbers) > 0 {
		pdf.CellFormat(50, 7, "Telefoni:", "1", 0, "L", false, 0, "")
		phoneStr := ""
		for i, phone := range member.PhoneNumbers {
			if i > 0 {
				phoneStr += ", "
			}
			phoneStr += phone.Prefix + " " + phone.Number
		}
		pdf.CellFormat(0, 7, phoneStr, "1", 1, "L", false, 0, "")
	}

	// Addresses
	if len(member.Addresses) > 0 {
		pdf.CellFormat(50, 7, "Indirizzo:", "1", 0, "L", false, 0, "")
		addr := member.Addresses[0]
		addrStr := fmt.Sprintf("%s %s, %s %s, %s", addr.Street, addr.StreetNumber, addr.ZipCode, addr.City, addr.Country)
		pdf.MultiCell(0, 7, addrStr, "1", "L", false)
	}

	pdf.Ln(5)

	// Memberships Section
	if len(member.Memberships) > 0 {
		pdf.SetFont("Arial", "B", 12)
		pdf.SetFillColor(220, 220, 220)
		pdf.CellFormat(0, 8, "Tessere", "1", 1, "L", true, 0, "")

		pdf.SetFont("Arial", "B", 9)
		pdf.SetFillColor(200, 200, 200)
		pdf.CellFormat(30, 7, "Numero", "1", 0, "C", true, 0, "")
		pdf.CellFormat(25, 7, "Stato", "1", 0, "C", true, 0, "")
		pdf.CellFormat(35, 7, "Valida Da", "1", 0, "C", true, 0, "")
		pdf.CellFormat(35, 7, "Scade Il", "1", 0, "C", true, 0, "")
		pdf.CellFormat(30, 7, "Prezzo", "1", 0, "C", true, 0, "")
		pdf.CellFormat(35, 7, "Pagamento", "1", 1, "C", true, 0, "")

		pdf.SetFont("Arial", "", 9)
		for _, membership := range member.Memberships {
			pdf.CellFormat(30, 6, fmt.Sprintf("%d", membership.Number), "1", 0, "C", false, 0, "")
			pdf.CellFormat(25, 6, membership.Status, "1", 0, "C", false, 0, "")
			pdf.CellFormat(35, 6, membership.ValidFrom, "1", 0, "C", false, 0, "")
			pdf.CellFormat(35, 6, membership.ExpiresAt, "1", 0, "C", false, 0, "")
			pdf.CellFormat(30, 6, fmt.Sprintf("%.2f EUR", membership.Price), "1", 0, "R", false, 0, "")

			paidText := "Non Pagato"
			if membership.Paid {
				paidText = "Pagato"
			}
			pdf.CellFormat(35, 6, paidText, "1", 1, "C", false, 0, "")
		}
		pdf.Ln(5)
	}

	// Rented Facilities Section
	if len(facilities) > 0 {
		pdf.SetFont("Arial", "B", 12)
		pdf.SetFillColor(220, 220, 220)
		pdf.CellFormat(0, 8, fmt.Sprintf("Servizi Affittati - Stagione %s", seasonCode), "1", 1, "L", true, 0, "")

		pdf.SetFont("Arial", "B", 9)
		pdf.SetFillColor(200, 200, 200)
		pdf.CellFormat(30, 7, "Identificativo", "1", 0, "C", true, 0, "")
		pdf.CellFormat(40, 7, "Tipo", "1", 0, "C", true, 0, "")
		pdf.CellFormat(30, 7, "Prezzo", "1", 0, "C", true, 0, "")
		pdf.CellFormat(30, 7, "Pagamento", "1", 0, "C", true, 0, "")
		pdf.CellFormat(60, 7, "Barca", "1", 1, "C", true, 0, "")

		pdf.SetFont("Arial", "", 9)
		totalPrice := 0.0
		for _, facility := range facilities {
			pdf.CellFormat(30, 6, facility.FacilityIdentifier, "1", 0, "C", false, 0, "")
			pdf.CellFormat(40, 6, facility.FacilityName, "1", 0, "L", false, 0, "")
			pdf.CellFormat(30, 6, fmt.Sprintf("%.2f EUR", facility.Price), "1", 0, "R", false, 0, "")

			paidText := "Non Pagato"
			if facility.Paid {
				paidText = "Pagato"
			}
			pdf.CellFormat(30, 6, paidText, "1", 0, "C", false, 0, "")

			boatText := "-"
			if facility.BoatName != "" {
				boatText = facility.BoatName
			}
			pdf.CellFormat(60, 6, boatText, "1", 1, "L", false, 0, "")

			totalPrice += facility.Price
		}

		// Total
		pdf.SetFont("Arial", "B", 9)
		pdf.CellFormat(70, 7, "Totale:", "1", 0, "R", false, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%.2f EUR", totalPrice), "1", 0, "R", false, 0, "")
		pdf.CellFormat(90, 7, "", "1", 1, "L", false, 0, "")
	} else {
		pdf.SetFont("Arial", "I", 10)
		pdf.CellFormat(0, 8, "Nessun servizio affittato per questa stagione", "", 1, "L", false, 0, "")
	}

	// Signature Section
	pdf.Ln(20)
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(95, 6, "Data: ___________________________", "", 0, "L", false, 0, "")
	pdf.CellFormat(95, 6, "Firma: ___________________________", "", 1, "L", false, 0, "")

	// Write to buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return &buf, nil
}
