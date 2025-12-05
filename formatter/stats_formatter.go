package formatter

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"netmon/style"
)

// StatsFormatter formats network statistics in summary format.
type StatsFormatter struct{}

// NewStatsFormatter creates a new StatsFormatter instance.
func NewStatsFormatter() *StatsFormatter {
	return &StatsFormatter{}
}

// NetworkStats contains network statistics information.
type NetworkStats struct {
	TCPConnections     int
	UDPConnections     int
	ListeningPorts     int
	NetworkInterfaces  int
	DefaultGateway     string
	TopProcesses       []ProcessConnectionCount
}

// ProcessConnectionCount represents a process and its connection count.
type ProcessConnectionCount struct {
	Name  string
	Count int
}

// Format formats network statistics as a summary box.
// Returns a formatted string with network statistics.
func (f *StatsFormatter) Format(stats NetworkStats) string {
	var builder strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(style.PrimaryColor).
		Bold(true).
		Underline(true)
	builder.WriteString(titleStyle.Render("Network Summary"))
	builder.WriteString("\n\n")

	// Statistics
	labelStyle := lipgloss.NewStyle().
		Foreground(style.SubtleColor).
		Width(30)
	valueStyle := lipgloss.NewStyle().
		Foreground(style.InfoColor).
		Bold(true)

	// TCP Connections
	builder.WriteString(labelStyle.Render("Active TCP Connections:"))
	builder.WriteString(valueStyle.Render(fmt.Sprintf("%d", stats.TCPConnections)))
	builder.WriteString("\n")

	// UDP Connections
	builder.WriteString(labelStyle.Render("Active UDP Connections:"))
	builder.WriteString(valueStyle.Render(fmt.Sprintf("%d", stats.UDPConnections)))
	builder.WriteString("\n")

	// Listening Ports
	builder.WriteString(labelStyle.Render("Listening Ports:"))
	builder.WriteString(valueStyle.Render(fmt.Sprintf("%d", stats.ListeningPorts)))
	builder.WriteString("\n")

	// Network Interfaces
	builder.WriteString(labelStyle.Render("Network Interfaces:"))
	builder.WriteString(valueStyle.Render(fmt.Sprintf("%d", stats.NetworkInterfaces)))
	builder.WriteString("\n")

	// Default Gateway
	gatewayValue := stats.DefaultGateway
	if gatewayValue == "" {
		gatewayValue = "N/A"
	}
	builder.WriteString(labelStyle.Render("Default Gateway:"))
	builder.WriteString(valueStyle.Render(gatewayValue))
	builder.WriteString("\n")

	// Top Processes
	if len(stats.TopProcesses) > 0 {
		builder.WriteString("\n")
		topProcessStyle := lipgloss.NewStyle().
			Foreground(style.SubtleColor)
		builder.WriteString(topProcessStyle.Render("Top Processes by Connections:"))
		builder.WriteString("\n")

		bulletStyle := lipgloss.NewStyle().
			Foreground(style.SecondaryColor).
			MarginLeft(2)

		for _, proc := range stats.TopProcesses {
			procInfo := fmt.Sprintf("• %s (%d connections)", proc.Name, proc.Count)
			builder.WriteString(bulletStyle.Render(procInfo))
			builder.WriteString("\n")
		}
	}

	// Wrap in box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.PrimaryColor).
		Padding(1, 2).
		Margin(1, 0).
		Width(60)

	return boxStyle.Render(builder.String())
}
