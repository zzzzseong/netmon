package cmd

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/spf13/cobra"
	"netmon/formatter"
	"netmon/utils"
)

// newConnCmd creates and returns the conn command.
// It displays all active ESTABLISHED connections with process information.
func newConnCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "conn",
		Short: "Show active connections",
		Long:  `Display all active ESTABLISHED TCP connections with process information.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get all network connections
			connections, err := net.Connections("inet")
			if err != nil {
				return fmt.Errorf("failed to get network connection information: %w", err)
			}

			// Filter to only ESTABLISHED connections
			establishedConns := utils.FilterEstablishedConnections(connections)

			// Check if there are any established connections
			if len(establishedConns) == 0 {
				fmt.Println("No active connections found.")
				return nil
			}

			// Format and display the table
			fmtter := formatter.NewConnectionTableFormatter()
			table := fmtter.Format(establishedConns)
			fmt.Println(table)

			return nil
		},
	}
}
