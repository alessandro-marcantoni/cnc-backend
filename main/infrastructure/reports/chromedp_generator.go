package reports

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/alessandro-marcantoni/cnc-backend/main/domain/reports"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

//go:embed templates/member_list.html
var memberListTemplate string

//go:embed templates/member_detail.html
var memberDetailTemplate string

// ChromeDPPDFGenerator implements PDFGenerator using chromedp and HTML templates
type ChromeDPPDFGenerator struct{}

// NewChromeDPPDFGenerator creates a new HTML-based PDF generator
func NewChromeDPPDFGenerator() *ChromeDPPDFGenerator {
	return &ChromeDPPDFGenerator{}
}

// NewWkhtmltopdfPDFGenerator is an alias for backward compatibility
func NewWkhtmltopdfPDFGenerator() *ChromeDPPDFGenerator {
	return NewChromeDPPDFGenerator()
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

// GenerateMemberListPDF generates a PDF with the list of all members using HTML template
func (g *ChromeDPPDFGenerator) GenerateMemberListPDF(members []reports.MemberSummary, seasonCode string) (*bytes.Buffer, error) {
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
	pdfBuf, err := generatePDFFromHTML(htmlBuf.String(), false) // portrait mode
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfBuf, nil
}

// GenerateMemberDetailPDF generates a PDF with detailed information about a member using HTML template
func (g *ChromeDPPDFGenerator) GenerateMemberDetailPDF(member reports.MemberDetail, facilities []reports.FacilityRental, seasonCode string) (*bytes.Buffer, error) {
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
	pdfBuf, err := generatePDFFromHTML(htmlBuf.String(), false) // portrait mode
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfBuf, nil
}

// generatePDFFromHTML converts HTML to PDF using chromedp
func generatePDFFromHTML(html string, landscape bool) (*bytes.Buffer, error) {
	// Get Chrome path from environment variable or use default
	chromePath := os.Getenv("CHROME_BIN")
	if chromePath == "" {
		chromePath = "/usr/bin/chromium" // Default fallback
	}

	// Create Chrome options for headless operation in restricted environments
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(chromePath),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-software-rasterizer", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-setuid-sandbox", true),
	)

	// Create allocator context with options
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	// Create context
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var pdfData []byte

	// Configure PDF print parameters
	printParams := page.PrintToPDF().
		WithPrintBackground(true).
		WithPreferCSSPageSize(false).
		WithMarginTop(0.4).
		WithMarginBottom(0.4).
		WithMarginLeft(0.4).
		WithMarginRight(0.4)

	if landscape {
		printParams = printParams.WithLandscape(true)
	}

	// Navigate to data URL and print to PDF
	err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfData, _, err = printParams.Do(ctx)
			return err
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("chromedp error: %w", err)
	}

	return bytes.NewBuffer(pdfData), nil
}
