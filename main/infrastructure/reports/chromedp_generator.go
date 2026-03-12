package reports

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
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

	// Log Chrome path for debugging
	log.Printf("[PDF] Using Chrome binary at: %s", chromePath)

	// Check if Chrome binary exists and is executable
	if _, err := os.Stat(chromePath); os.IsNotExist(err) {
		log.Printf("[PDF] ERROR: Chrome binary not found at %s", chromePath)
		return nil, fmt.Errorf("chrome binary not found at %s: %w", chromePath, err)
	}

	// Try to get Chrome version for diagnostics
	if cmd := exec.Command(chromePath, "--version"); cmd != nil {
		if output, err := cmd.Output(); err == nil {
			log.Printf("[PDF] Chrome version: %s", string(output))
		} else {
			log.Printf("[PDF] WARNING: Could not get Chrome version: %v", err)
		}
	}

	log.Printf("[PDF] Starting PDF generation (landscape=%v, html_length=%d)", landscape, len(html))

	// Get Chrome user data directory from environment or create default
	userDataDir := os.Getenv("CHROME_USER_DATA_DIR")
	if userDataDir == "" {
		userDataDir = "/tmp/chrome-data"
	}

	// Create Chrome options for headless operation in containers
	// Start with DefaultExecAllocatorOptions and add container-specific flags
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(chromePath),
		chromedp.Flag("headless", "new"),

		// Critical container flags - without these Chrome crashes
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-crash-reporter", true),

		// Set user data directory explicitly (important for non-root users)
		chromedp.Flag("user-data-dir", userDataDir),

		// Additional stability flags for containers
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-software-rasterizer", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("hide-scrollbars", true),
	)

	// Create allocator context with options
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	// Create context
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set timeout - increase to 60 seconds for production environments
	timeoutDuration := 60 * time.Second
	if timeoutEnv := os.Getenv("PDF_TIMEOUT_SECONDS"); timeoutEnv != "" {
		if parsedTimeout, err := time.ParseDuration(timeoutEnv + "s"); err == nil {
			timeoutDuration = parsedTimeout
		}
	}
	log.Printf("[PDF] Using timeout: %v", timeoutDuration)

	ctx, cancel = context.WithTimeout(ctx, timeoutDuration)
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
	log.Printf("[PDF] Starting chromedp navigation and rendering")

	err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Printf("[PDF] Setting document content")
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				log.Printf("[PDF] ERROR: Failed to get frame tree: %v", err)
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Printf("[PDF] Printing to PDF")
			var err error
			pdfData, _, err = printParams.Do(ctx)
			if err != nil {
				log.Printf("[PDF] ERROR: Failed to print PDF: %v", err)
			} else {
				log.Printf("[PDF] Successfully generated PDF (size: %d bytes)", len(pdfData))
			}
			return err
		}),
	)

	if err != nil {
		log.Printf("[PDF] FINAL ERROR: %v", err)
		// Add more context about the environment
		log.Printf("[PDF] Environment debug info:")
		log.Printf("[PDF]   - Chrome path: %s", chromePath)
		log.Printf("[PDF]   - Working directory: %s", func() string { wd, _ := os.Getwd(); return wd }())
		log.Printf("[PDF]   - User: %s", os.Getenv("USER"))
		log.Printf("[PDF]   - Available memory: checking /proc/meminfo")
		if memInfo, err := os.ReadFile("/proc/meminfo"); err == nil {
			log.Printf("[PDF]   - Memory info (first 500 chars): %s", string(memInfo[:min(500, len(memInfo))]))
		}
		return nil, fmt.Errorf("chromedp error: %w", err)
	}

	log.Printf("[PDF] PDF generation completed successfully")
	return bytes.NewBuffer(pdfData), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
