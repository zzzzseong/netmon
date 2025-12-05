package cmd

import (
	"fmt"
	"sort"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/spf13/cobra"
	"netmon/formatter"
	"netmon/provider"
	"netmon/utils"
)

// newStatsCmd creates and returns the stats command.
// It displays network statistics summary.
func newStatsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Show network statistics summary",
		Long:  `Display network statistics summary including connection counts, interfaces, and top processes.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get all network connections
			connections, err := net.Connections("inet")
			if err != nil {
				return fmt.Errorf("failed to get network connection information: %w", err)
			}

			// Count TCP and UDP connections
			tcpCount := 0
			udpCount := 0
			for _, conn := range connections {
				if conn.Status == "ESTABLISHED" {
					if utils.IsUDP(conn.Type) {
						udpCount++
					} else {
						tcpCount++
					}
				}
			}

			// Count listening ports
			listeningConns := utils.FilterListeningConnections(connections, true)
			listeningCount := len(listeningConns)

			// Get network interfaces
			interfaces, err := net.Interfaces()
			if err != nil {
				return fmt.Errorf("failed to get network interface information: %w", err)
			}
			interfaceCount := len(interfaces)

			// Get default gateway from routing table
			routeProvider := provider.NewRouteProvider()
			routes, err := routeProvider.GetRoutes()
			defaultGateway := ""
			if err == nil {
				for _, route := range routes {
					if route.Destination == "default" || route.Destination == "0.0.0.0/0" {
						defaultGateway = route.Gateway
						break
					}
				}
			}

			// Get top processes by connection count
			topProcesses := getTopProcessesByConnections(connections, 5)

			// Create stats object
			stats := formatter.NetworkStats{
				TCPConnections:    tcpCount,
				UDPConnections:    udpCount,
				ListeningPorts:    listeningCount,
				NetworkInterfaces: interfaceCount,
				DefaultGateway:    defaultGateway,
				TopProcesses:      topProcesses,
			}

			// Format and display
			fmtter := formatter.NewStatsFormatter()
			output := fmtter.Format(stats)
			fmt.Println(output)

			return nil
		},
	}
}

// getTopProcessesByConnections returns the top N processes by connection count.
func getTopProcessesByConnections(connections []net.ConnectionStat, topN int) []formatter.ProcessConnectionCount {
	// Count connections per process
	processConnCount := make(map[string]int)
	processNames := make(map[int32]string)

	for _, conn := range connections {
		if conn.Status == "ESTABLISHED" && conn.Pid > 0 {
			// Get process name if not cached
			if _, ok := processNames[conn.Pid]; !ok {
				processInfo := utils.GetProcessInfo(conn.Pid)
				processNames[conn.Pid] = processInfo.Name
			}

			name := processNames[conn.Pid]
			if name != "N/A" {
				processConnCount[name]++
			}
		}
	}

	// Convert to slice and sort
	type processCount struct {
		name  string
		count int
	}

	processes := make([]processCount, 0, len(processConnCount))
	for name, count := range processConnCount {
		processes = append(processes, processCount{name: name, count: count})
	}

	// Sort by count descending
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].count > processes[j].count
	})

	// Take top N
	if len(processes) > topN {
		processes = processes[:topN]
	}

	// Convert to formatter type
	result := make([]formatter.ProcessConnectionCount, len(processes))
	for i, p := range processes {
		result[i] = formatter.ProcessConnectionCount{
			Name:  p.name,
			Count: p.count,
		}
	}

	return result
}
