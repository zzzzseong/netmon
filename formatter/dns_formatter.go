package formatter

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"netmon/style"
	"netmon/utils"
)

var (
	dnsRecordTypeStyle = lipgloss.NewStyle().Foreground(style.SecondaryColor).Bold(true)
	dnsValueStyle      = lipgloss.NewStyle().Foreground(style.InfoColor)
)

// DNSFormatter formats DNS lookup results.
type DNSFormatter struct{}

// NewDNSFormatter creates a new DNSFormatter instance.
func NewDNSFormatter() *DNSFormatter {
	return &DNSFormatter{}
}

// Format formats DNS lookup results as a table.
// Returns a formatted string with DNS records and response time.
func (f *DNSFormatter) Format(result utils.DNSResult) string {
	var builder strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(style.PrimaryColor).
		Bold(true)
	header := fmt.Sprintf("🔍 DNS Lookup: %s", result.Query)
	builder.WriteString(headerStyle.Render(header))
	builder.WriteString("\n\n")

	// Check for errors
	if result.Error != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(style.DangerColor).
			Bold(true)
		errorMsg := fmt.Sprintf("Error: %s", result.Error.Error())
		builder.WriteString(errorStyle.Render(errorMsg))
		return builder.String()
	}

	// Build table rows, one per record
	rows := [][]string{}
	addRows := func(recordType string, values []string) {
		for _, v := range values {
			rows = append(rows, []string{
				dnsRecordTypeStyle.Render(recordType),
				dnsValueStyle.Render(v),
			})
		}
	}

	addRows("A", result.ARecords)
	addRows("AAAA", result.AAAARecords)
	addRows("PTR", result.PTRRecords)
	if result.CNAMERecord != "" {
		addRows("CNAME", []string{result.CNAMERecord})
	}
	addRows("MX", result.MXRecords)
	addRows("NS", result.NSRecords)
	addRows("TXT", result.TXTRecords)

	// If no records found
	if len(rows) == 0 {
		noResultStyle := lipgloss.NewStyle().
			Foreground(style.WarningColor)
		builder.WriteString(noResultStyle.Render("No DNS records found."))
		return builder.String()
	}

	builder.WriteString(CreateTable(rows, DNSTableColumns))
	builder.WriteString("\n\n")

	// Response time
	responseTimeStyle := lipgloss.NewStyle().
		Foreground(style.SubtleColor)
	responseTime := fmt.Sprintf("Response Time: %dms", result.ResponseTime.Milliseconds())
	builder.WriteString(responseTimeStyle.Render(responseTime))

	return builder.String()
}
