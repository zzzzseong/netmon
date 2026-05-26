package cmd

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/spf13/cobra"
	"netmon/formatter"
	"netmon/utils"
)

// newLsCmd creates and returns the ls command.
// It displays all active listening ports with detailed process information.
func newLsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls [-a]",
		Short: "List active ports (-a: include UDP)",
		Long:  `Display all active listening ports with detailed process information.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			includeUDP, _ := cmd.Flags().GetBool("all")
			watch, _ := cmd.Flags().GetBool("watch")
			interval, _ := cmd.Flags().GetInt("interval")

			fmtter := formatter.NewPortTableFormatter()

			run := func() (string, error) {
				connections, err := net.Connections("inet")
				if err != nil {
					return "", fmt.Errorf("failed to get network connection information: %w", err)
				}
				listeningConns := utils.FilterListeningConnections(connections, includeUDP)
				if len(listeningConns) == 0 {
					return "No listening ports found.", nil
				}
				return fmtter.Format(utils.SortConnectionsByPort(listeningConns)), nil
			}

			if watch {
				return runWithWatch("ls", time.Duration(interval)*time.Second, run)
			}
			out, err := run()
			if err != nil {
				return err
			}
			fmt.Println(out)
			return nil
		},
	}

	cmd.Flags().BoolP("all", "a", false, "Include UDP connections")
	cmd.Flags().BoolP("watch", "w", false, "Watch mode: refresh output periodically")
	cmd.Flags().IntP("interval", "n", 1, "Refresh interval in seconds (used with -w)")

	return cmd
}
