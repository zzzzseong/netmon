package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"netmon/formatter"
	"netmon/utils"
)

// newDNSCmd creates and returns the dns command.
// It performs DNS lookups for domains or IP addresses.
func newDNSCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "dns <domain|ip>",
		Short: "Perform DNS lookup",
		Long:  `Perform DNS lookup for a domain (A/AAAA records) or reverse lookup for an IP address (PTR records).`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := args[0]

			var result utils.DNSResult

			// Check if input is an IP address or domain
			if utils.IsIPAddress(query) {
				// Perform reverse lookup
				result = utils.ReverseLookup(query)
			} else {
				// Validate domain
				if err := utils.ValidateDomain(query); err != nil {
					return fmt.Errorf("invalid domain: %w", err)
				}

				// Perform forward lookup
				result = utils.LookupDomain(query)
			}

			// Format and display results
			fmtter := formatter.NewDNSFormatter()
			output := fmtter.Format(result)
			fmt.Println(output)

			return nil
		},
	}
}
