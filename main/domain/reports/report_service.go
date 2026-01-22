package reports

import (
	"bytes"
)

// ReportService orchestrates report generation
type ReportService struct {
	pdfGenerator PDFGenerator
}

// NewReportService creates a new report service
func NewReportService(pdfGenerator PDFGenerator) *ReportService {
	return &ReportService{
		pdfGenerator: pdfGenerator,
	}
}

// GenerateMemberListReport generates a PDF report with all members
func (s *ReportService) GenerateMemberListReport(members []MemberSummary, seasonCode string) (*bytes.Buffer, error) {
	return s.pdfGenerator.GenerateMemberListPDF(members, seasonCode)
}

// GenerateMemberDetailReport generates a PDF report with member details and facilities
func (s *ReportService) GenerateMemberDetailReport(member MemberDetail, facilities []FacilityRental, seasonCode string) (*bytes.Buffer, error) {
	return s.pdfGenerator.GenerateMemberDetailPDF(member, facilities, seasonCode)
}
