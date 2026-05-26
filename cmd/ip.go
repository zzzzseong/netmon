package cmd

import (
	"fmt"
	"sort"
	"time"

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
			showAll, _ := cmd.Flags().GetBool("all")
			watch, _ := cmd.Flags().GetBool("watch")
			interval, _ := cmd.Flags().GetInt("interval")

			fmtter := formatter.NewInterfaceTableFormatter()

			run := func() error {
				interfaces, err := net.Interfaces()
				if err != nil {
					return fmt.Errorf("failed to get network interface information: %w", err)
				}

				sort.Slice(interfaces, func(i, j int) bool {
					return interfaces[i].Name < interfaces[j].Name
				})

				fmt.Println(fmtter.Format(interfaces, showAll))
				return nil
			}

			if watch {
				return runWithWatch("ip", time.Duration(interval)*time.Second, run)
			}
			return run()
		},
	}

	cmd.Flags().BoolP("all", "a", false, "Show all addresses including IPv6")
	cmd.Flags().BoolP("watch", "w", false, "Watch mode: refresh output periodically")
	cmd.Flags().IntP("interval", "n", 1, "Refresh interval in seconds (used with -w)")

	return cmd
}
