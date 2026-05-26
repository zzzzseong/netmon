package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"netmon/formatter"
	"netmon/provider"
)

// newRouteCmd creates and returns the route command.
// It displays system routing information with smart filtering.
func newRouteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "route",
		Short: "Show routing table information",
		Long:  `Display system routing information with smart filtering.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			watch, _ := cmd.Flags().GetBool("watch")
			interval, _ := cmd.Flags().GetInt("interval")

			routeProvider := provider.NewRouteProvider()
			fmtter := formatter.NewRouteTableFormatter()

			run := func() (string, error) {
				routes, err := routeProvider.GetRoutes()
				if err != nil {
					return "", fmt.Errorf("failed to get routing table information: %w", err)
				}
				return fmtter.Format(provider.FilterRoutes(routes)), nil
			}

			if watch {
				return runWithWatch("route", time.Duration(interval)*time.Second, run)
			}
			out, err := run()
			if err != nil {
				return err
			}
			fmt.Println(out)
			return nil
		},
	}

	cmd.Flags().BoolP("watch", "w", false, "Watch mode: refresh output periodically")
	cmd.Flags().IntP("interval", "n", 1, "Refresh interval in seconds (used with -w)")

	return cmd
}
