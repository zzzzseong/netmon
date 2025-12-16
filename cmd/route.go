package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"netmon/formatter"
	"netmon/provider"
)

// newRouteCmd creates and returns the route command.
// It displays system routing information with smart filtering.
func newRouteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "route",
		Short: "Show routing table information",
		Long:  `Display system routing information with smart filtering.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get routing table through RouteProvider
			routeProvider := provider.NewRouteProvider()
			routes, err := routeProvider.GetRoutes()
			if err != nil {
				return fmt.Errorf("failed to get routing table information: %w", err)
			}

			// Filter unnecessary routes (exclude single hosts, multicast, link-local, etc.)
			routes = provider.FilterRoutes(routes)

			// Output in table format using formatter
			fmtter := formatter.NewRouteTableFormatter()
			table := fmtter.Format(routes)
			fmt.Println(table)

			return nil
		},
	}
}
