package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"netmon/style"
)

// newVersionCmd creates and returns the version command.
// It displays version and build information from the provided Config.
func newVersionCmd(cfg Config) *cobra.Command {
	return &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display version and build information.`,
	Run: func(cmd *cobra.Command, args []string) {
		versionText := fmt.Sprintf("Version: %s", cfg.Version)
		buildText := fmt.Sprintf("Build Date: %s", cfg.BuildDate)

		fmt.Print(style.ASCIIStyle.Render(style.ASCIIArt))
		fmt.Println(style.VersionTextStyle.Render(versionText))
		fmt.Println(style.BuildTextStyle.Render(buildText))

		footerText := "For more information, visit: https://github.com/zzzzseong/netmon"
		fmt.Println(style.FooterTextStyle.Render(footerText))
	},
	}
}
