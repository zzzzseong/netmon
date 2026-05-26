package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/spf13/cobra"
	"netmon/formatter"
	"netmon/provider"
	"netmon/utils"
)

// newStatsCmd creates and returns the stats command.
// It displays network statistics summary.
func newStatsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Show network statistics summary",
		Long:  `Display network statistics summary including connection counts, interfaces, and top processes.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			watch, _ := cmd.Flags().GetBool("watch")
			interval, _ := cmd.Flags().GetInt("interval")

			fmtter := formatter.NewStatsFormatter()
			routeProvider := provider.NewRouteProvider()

			run := func() error {
				connections, err := net.Connections("inet")
				if err != nil {
					return fmt.Errorf("failed to get network connection information: %w", err)
				}

				tcpCount := 0
				udpCount := 0
				for _, conn := range connections {
					if utils.IsUDP(conn.Type) {
						if conn.Laddr.Port > 0 {
							udpCount++
						}
						continue
					}
					if conn.Status == "ESTABLISHED" {
						tcpCount++
					}
				}

				listeningCount := len(utils.FilterListeningConnections(connections, true))

				interfaces, err := net.Interfaces()
				if err != nil {
					return fmt.Errorf("failed to get network interface information: %w", err)
				}

				defaultGateway := ""
				if routes, err := routeProvider.GetRoutes(); err == nil {
					for _, route := range routes {
						if route.Destination == "default" || route.Destination == "0.0.0.0/0" {
							defaultGateway = route.Gateway
							break
						}
					}
				}

				stats := formatter.NetworkStats{
					TCPConnections:    tcpCount,
					UDPConnections:    udpCount,
					ListeningPorts:    listeningCount,
					NetworkInterfaces: len(interfaces),
					DefaultGateway:    defaultGateway,
					TopProcesses:      getTopProcessesByConnections(connections, 5),
				}

				fmt.Println(fmtter.Format(stats))
				return nil
			}

			if watch {
				return runWithWatch("stats", time.Duration(interval)*time.Second, run)
			}
			return run()
		},
	}

	cmd.Flags().BoolP("watch", "w", false, "Watch mode: refresh output periodically")
	cmd.Flags().IntP("interval", "n", 1, "Refresh interval in seconds (used with -w)")

	return cmd
}

// getTopProcessesByConnections returns the top N processes by connection count.
func getTopProcessesByConnections(connections []net.ConnectionStat, topN int) []formatter.ProcessConnectionCount {
	// Count connections per process
	processConnCount := make(map[string]int)
	processNames := make(map[int32]string)

	for _, conn := range connections {
		if conn.Pid <= 0 {
			continue
		}
		if !utils.IsUDP(conn.Type) && conn.Status != "ESTABLISHED" {
			continue
		}
		if utils.IsUDP(conn.Type) && conn.Laddr.Port == 0 {
			continue
		}

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
