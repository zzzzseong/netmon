package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"netmon/style"
)

var (
	// Version is the application version, set at build time
	Version = "dev"
	// BuildDate is the build date, set at build time
	BuildDate = "unknown"
)

// newVersionCmd creates and returns the version command.
// It displays version and build information.
func newVersionCmd() *cobra.Command {
	return &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display version and build information.`,
	Run: func(cmd *cobra.Command, args []string) {
		asciiStyle := lipgloss.NewStyle().
			Foreground(style.PrimaryColor).
			Bold(true)

		versionText := fmt.Sprintf("Version: %s", Version)
		versionStyle := lipgloss.NewStyle().
			Foreground(style.SecondaryColor).
			Bold(true).
			MarginTop(1)

		buildText := fmt.Sprintf("Build Date: %s", BuildDate)
		buildStyle := lipgloss.NewStyle().
			Foreground(style.SubtleColor).
			MarginBottom(1)

		fmt.Print(asciiStyle.Render(style.ASCIIArt))
		fmt.Println(versionStyle.Render(versionText))
		fmt.Println(buildStyle.Render(buildText))

		footerText := "For more information, visit: https://github.com/zzzzseong/netmon"
		footerStyle := lipgloss.NewStyle().
			Foreground(style.SubtleColor).
			Italic(true)
		fmt.Println(footerStyle.Render(footerText))
	},
	}
}
