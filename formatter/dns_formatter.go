package formatter

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"netmon/style"
	"netmon/utils"
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

	// Build table rows
	rows := [][]string{}

	// Add A records (IPv4)
	if len(result.ARecords) > 0 {
		recordTypeStyle := lipgloss.NewStyle().Foreground(style.SecondaryColor).Bold(true)
		valueStyle := lipgloss.NewStyle().Foreground(style.InfoColor)
		
		for _, ip := range result.ARecords {
			rows = append(rows, []string{
				recordTypeStyle.Render("A"),
				valueStyle.Render(ip),
			})
		}
	}

	// Add AAAA records (IPv6)
	if len(result.AAAARecords) > 0 {
		recordTypeStyle := lipgloss.NewStyle().Foreground(style.SecondaryColor).Bold(true)
		valueStyle := lipgloss.NewStyle().Foreground(style.InfoColor)
		
		for _, ip := range result.AAAARecords {
			rows = append(rows, []string{
				recordTypeStyle.Render("AAAA"),
				valueStyle.Render(ip),
			})
		}
	}

	// Add PTR records (reverse DNS)
	if len(result.PTRRecords) > 0 {
		recordTypeStyle := lipgloss.NewStyle().Foreground(style.SecondaryColor).Bold(true)
		valueStyle := lipgloss.NewStyle().Foreground(style.InfoColor)
		
		for _, name := range result.PTRRecords {
			rows = append(rows, []string{
				recordTypeStyle.Render("PTR"),
				valueStyle.Render(name),
			})
		}
	}

	// If no records found
	if len(rows) == 0 {
		noResultStyle := lipgloss.NewStyle().
			Foreground(style.WarningColor)
		builder.WriteString(noResultStyle.Render("No DNS records found."))
		return builder.String()
	}

	// Create table headers
	headerStyle = style.HeaderStyle.Copy().Align(lipgloss.Center)
	styledHeaders := make([]string, len(DNSTableColumns))
	for i, col := range DNSTableColumns {
		styledHeaders[i] = headerStyle.Width(col.Width).Render(col.Title)
	}

	// Create and style table
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(style.TableBorderStyle).
		StyleFunc(GetTableRowStyle).
		Headers(styledHeaders...).
		Rows(rows...).
		Width(style.TableWidthDNS)

	builder.WriteString(t.String())
	builder.WriteString("\n\n")

	// Response time
	responseTimeStyle := lipgloss.NewStyle().
		Foreground(style.SubtleColor)
	responseTime := fmt.Sprintf("Response Time: %dms", result.ResponseTime.Milliseconds())
	builder.WriteString(responseTimeStyle.Render(responseTime))

	return builder.String()
}
