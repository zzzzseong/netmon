package cmd

import (
	"fmt"
	"sort"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/spf13/cobra"
	"netmon/formatter"
)

// newIPCmd creates and returns the ip command.
// It displays network interfaces with IP addresses.
func newIPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ip [-a]",
		Short: "Show network interfaces (-a: include IPv6)",
		Long:  `Display network interfaces with IP addresses.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check -a flag
			showAll, _ := cmd.Flags().GetBool("all")

			// Get network interface information
			interfaces, err := net.Interfaces()
			if err != nil {
				return fmt.Errorf("failed to get network interface information: %w", err)
			}

			// Sort by interface name
			sort.Slice(interfaces, func(i, j int) bool {
				return interfaces[i].Name < interfaces[j].Name
			})

			// Create table using formatter
			fmtter := formatter.NewInterfaceTableFormatter()
			table := fmtter.Format(interfaces, showAll)
			fmt.Println(table)

			return nil
		},
	}

	// Define -a, --all flag
	cmd.Flags().BoolP("all", "a", false, "Show all addresses including IPv6")

	return cmd
}
