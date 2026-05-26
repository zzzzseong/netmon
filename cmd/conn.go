package cmd

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/spf13/cobra"
	"netmon/formatter"
	"netmon/utils"
)

// newConnCmd creates and returns the conn command.
// It displays all active ESTABLISHED connections with process information.
func newConnCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "conn",
		Short: "Show active connections",
		Long:  `Display all active ESTABLISHED TCP connections with process information.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			watch, _ := cmd.Flags().GetBool("watch")
			interval, _ := cmd.Flags().GetInt("interval")

			fmtter := formatter.NewConnectionTableFormatter()

			run := func() error {
				connections, err := net.Connections("inet")
				if err != nil {
					return fmt.Errorf("failed to get network connection information: %w", err)
				}

				establishedConns := utils.FilterEstablishedConnections(connections)
				if len(establishedConns) == 0 {
					fmt.Println("No active connections found.")
					return nil
				}

				table := fmtter.Format(establishedConns)
				fmt.Println(table)
				return nil
			}

			if watch {
				return runWithWatch("conn", time.Duration(interval)*time.Second, run)
			}
			return run()
		},
	}

	cmd.Flags().BoolP("watch", "w", false, "Watch mode: refresh output periodically")
	cmd.Flags().IntP("interval", "n", 2, "Refresh interval in seconds (used with -w)")

	return cmd
}
