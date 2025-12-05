package formatter

import (
	"fmt"
	"strings"

	"netmon/provider"
)

// FormatLinuxStyle formats a RouteEntry in Linux ip route style.
// It returns a string representation similar to the output of "ip route" command.
func FormatLinuxStyle(entry provider.RouteEntry) string {
	var parts []string

	// Destination
	parts = append(parts, entry.Destination)

	// via Gateway
	if entry.Gateway != "" {
		parts = append(parts, "via", entry.Gateway)
	}

	// dev Interface
	if entry.Interface != "" {
		parts = append(parts, "dev", entry.Interface)
	}

	// metric
	if entry.Metric > 0 {
		parts = append(parts, "metric", fmt.Sprintf("%d", entry.Metric))
	}

	// src
	if entry.Source != "" {
		parts = append(parts, "src", entry.Source)
	}

	return strings.Join(parts, " ")
}

// FormatLinuxStyleRoutes formats a slice of RouteEntry in Linux ip route style.
// Each route is formatted on a separate line.
func FormatLinuxStyleRoutes(entries []provider.RouteEntry) string {
	var lines []string
	for _, entry := range entries {
		lines = append(lines, FormatLinuxStyle(entry))
	}
	return strings.Join(lines, "\n")
}
