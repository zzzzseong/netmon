package cmd

import (
	"fmt"

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
			// Get -a option value
			includeUDP, _ := cmd.Flags().GetBool("all")

			// Create formatter
			fmtter := formatter.NewPortTableFormatter()

			// Get all network connections
			connections, err := net.Connections("inet")
			if err != nil {
				return fmt.Errorf("failed to get network connection information: %w", err)
			}

			// Filter only LISTEN connections (include UDP only if -a option is set)
			listeningConns := utils.FilterListeningConnections(connections, includeUDP)

			// Sort by port number
			sortedConns := utils.SortConnectionsByPort(listeningConns)

			// Create table using formatter
			table := fmtter.Format(sortedConns)
			fmt.Println(table)

			return nil
		},
	}

	// Define -a, --all flag
	cmd.Flags().BoolP("all", "a", false, "Include UDP connections")

	return cmd
}
