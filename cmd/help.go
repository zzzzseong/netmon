package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"netmon/style"
)

// printCustomHelp prints a beautifully styled help message
func printCustomHelp(cmd *cobra.Command, args []string, cfg Config) {
	// Version information
	versionText := fmt.Sprintf("Version %s", cfg.Version)

	// Description
	descText := "A powerful CLI tool for monitoring network connections and managing processes."

	// Print ASCII Art
	fmt.Print(style.ASCIIStyle.Render(style.ASCIIArt))
	fmt.Println(style.VersionTextStyle.Render(versionText))
	fmt.Println(style.DescTextStyle.Render(descText))

	// Usage section
	usageHeader := lipgloss.NewStyle().
		Foreground(style.PrimaryColor).
		Bold(true).
		MarginBottom(1).
		Render("Usage:")
	fmt.Println(usageHeader)

	// Usage content box
	usageContent := "netmon <command> [arguments]"
	usageBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.PrimaryColor).
		Padding(1, 2).
		MarginBottom(2).
		Render(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")).
				Render(usageContent),
		)
	fmt.Println(usageBox)

	// Commands section
	commandsHeader := lipgloss.NewStyle().
		Foreground(style.SecondaryColor).
		Bold(true).
		MarginBottom(1).
		Render("Commands:")
	fmt.Println(commandsHeader)

	// Collect command list
	var commands []*cobra.Command
	for _, c := range cmd.Commands() {
		// Exclude completion and help commands
		if c.Name() != "completion" && c.Name() != "help" {
			commands = append(commands, c)
		}
	}

	// Calculate longest command length
	maxUsageLen := 0
	for _, c := range commands {
		usageLen := len(fmt.Sprintf("  netmon %s", c.Use))
		if usageLen > maxUsageLen {
			maxUsageLen = usageLen
		}
	}

	// Create command list
	var commandsList strings.Builder
	for i, c := range commands {
		usage := fmt.Sprintf("  netmon %s", c.Use)
		// Add padding
		padding := maxUsageLen - len(usage) + 4
		usage += strings.Repeat(" ", padding)

		cmdLine := style.CommandStyle.Render(usage)
		desc := style.DescStyle.Render(c.Short)
		commandsList.WriteString(cmdLine + desc)
		if i < len(commands)-1 {
			commandsList.WriteString("\n")
		}
	}

	commandsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.SecondaryColor).
		Padding(1, 2).
		Render(commandsList.String())
	fmt.Println(commandsBox)

	// Watch mode tips
	tipsHeader := lipgloss.NewStyle().
		Foreground(style.InfoColor).
		Bold(true).
		MarginBottom(1).
		Render("Watch Mode:")
	fmt.Println(tipsHeader)

	tips := []struct{ cmd, desc string }{
		{"  netmon ls -w", "refresh every 2s (default)"},
		{"  netmon ls -w -n 5", "refresh every 5s"},
	}
	var tipsList strings.Builder
	for i, t := range tips {
		cmdPart := style.CommandStyle.Render(fmt.Sprintf("%-24s", t.cmd))
		descPart := style.DescStyle.Render(t.desc)
		tipsList.WriteString(cmdPart + descPart)
		if i < len(tips)-1 {
			tipsList.WriteString("\n")
		}
	}
	tipsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(style.InfoColor).
		Padding(1, 2).
		MarginBottom(1).
		Render(tipsList.String())
	fmt.Println(tipsBox)

	// Footer
	footerText := "For more information, visit: https://github.com/zzzzseong/netmon"
	footerStyle := style.FooterTextStyle
	fmt.Println(footerStyle.MarginTop(2).Render(footerText))
}

// printCustomUsage prints a custom usage message
func printCustomUsage(cmd *cobra.Command) error {
	usageText := fmt.Sprintf("Usage: %s [command] [arguments]", cmd.Use)
	usageStyle := lipgloss.NewStyle().
		Foreground(style.PrimaryColor).
		Bold(true)
	fmt.Println(usageStyle.Render(usageText))
	return nil
}
