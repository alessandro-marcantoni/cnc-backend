package reports

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"time"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain/reports"
)

//go:embed templates/member_list.html
var memberListTemplate string

//go:embed templates/member_detail.html
var memberDetailTemplate string

// WkhtmltopdfGenerator implements PDFGenerator using wkhtmltopdf and HTML templates
type WkhtmltopdfGenerator struct {
	wkhtmltopdfPath string
}

// NewWkhtmltopdfGenerator creates a new wkhtmltopdf-based PDF generator
func NewWkhtmltopdfGenerator() *WkhtmltopdfGenerator {
	// Get wkhtmltopdf path from environment or use default
	wkhtmlPath := os.Getenv("WKHTMLTOPDF_PATH")
	if wkhtmlPath == "" {
		wkhtmlPath = "/usr/local/bin/wkhtmltopdf"
	}

	return &WkhtmltopdfGenerator{
		wkhtmltopdfPath: wkhtmlPath,
	}
}

// MemberListTemplateData holds data for the member list template
type MemberListTemplateData struct {
	SeasonCode    string
	GeneratedDate string
	Members       []reports.MemberSummary
	TotalMembers  int
}

// MemberDetailTemplateData holds data for the member detail template
type MemberDetailTemplateData struct {
	SeasonCode           string
	GeneratedDate        string
	Member               reports.MemberDetail
	Facilities           []reports.FacilityRental
	TotalFacilitiesPrice float64
}

// GenerateMemberListPDF generates a PDF with the list of all members using wkhtmltopdf
func (g *WkhtmltopdfGenerator) GenerateMemberListPDF(members []reports.MemberSummary, seasonCode string) (*bytes.Buffer, error) {
	// Prepare template data
	data := MemberListTemplateData{
		SeasonCode:    seasonCode,
		GeneratedDate: time.Now().Format("02/01/2006"),
		Members:       members,
		TotalMembers:  len(members),
	}

	// Parse and execute template
	tmpl, err := template.New("member_list").Parse(memberListTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var htmlBuf bytes.Buffer
	if err := tmpl.Execute(&htmlBuf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	// Generate PDF from HTML
	pdfBuf, err := g.generatePDFFromHTML(htmlBuf.String(), "A4", "Portrait")
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfBuf, nil
}

// GenerateMemberDetailPDF generates a PDF with detailed information about a member using wkhtmltopdf
func (g *WkhtmltopdfGenerator) GenerateMemberDetailPDF(member reports.MemberDetail, facilities []reports.FacilityRental, seasonCode string) (*bytes.Buffer, error) {
	// Calculate total facilities price
	totalPrice := 0.0
	for _, facility := range facilities {
		totalPrice += facility.Price
	}

	// Prepare template data
	data := MemberDetailTemplateData{
		SeasonCode:           seasonCode,
		GeneratedDate:        time.Now().Format("02/01/2006"),
		Member:               member,
		Facilities:           facilities,
		TotalFacilitiesPrice: totalPrice,
	}

	// Parse and execute template
	tmpl, err := template.New("member_detail").Parse(memberDetailTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var htmlBuf bytes.Buffer
	if err := tmpl.Execute(&htmlBuf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	// Generate PDF from HTML
	pdfBuf, err := g.generatePDFFromHTML(htmlBuf.String(), "A4", "Portrait")
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfBuf, nil
}

// generatePDFFromHTML converts HTML to PDF using wkhtmltopdf
func (g *WkhtmltopdfGenerator) generatePDFFromHTML(html string, pageSize string, orientation string) (*bytes.Buffer, error) {
	// Create temporary file for HTML input
	tmpHTMLFile, err := os.CreateTemp("", "html-*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp HTML file: %w", err)
	}
	defer os.Remove(tmpHTMLFile.Name())

	// Write HTML content to temp file
	if _, err := tmpHTMLFile.WriteString(html); err != nil {
		tmpHTMLFile.Close()
		return nil, fmt.Errorf("failed to write HTML to temp file: %w", err)
	}
	tmpHTMLFile.Close()

	// Create temporary file for PDF output
	tmpPDFFile, err := os.CreateTemp("", "pdf-*.pdf")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp PDF file: %w", err)
	}
	tmpPDFFile.Close()
	defer os.Remove(tmpPDFFile.Name())

	// Build wkhtmltopdf command
	args := []string{
		"--page-size", pageSize,
		"--orientation", orientation,
		"--margin-top", "10mm",
		"--margin-bottom", "10mm",
		"--margin-left", "10mm",
		"--margin-right", "10mm",
		"--encoding", "UTF-8",
		"--enable-local-file-access",
		"--no-stop-slow-scripts",
		"--javascript-delay", "100",
		tmpHTMLFile.Name(),
		tmpPDFFile.Name(),
	}

	cmd := exec.Command(g.wkhtmltopdfPath, args...)

	// Capture stderr for debugging
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Execute wkhtmltopdf
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("wkhtmltopdf execution failed: %w, stderr: %s", err, stderr.String())
	}

	// Read the generated PDF
	pdfData, err := os.ReadFile(tmpPDFFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read generated PDF: %w", err)
	}

	return bytes.NewBuffer(pdfData), nil
}
