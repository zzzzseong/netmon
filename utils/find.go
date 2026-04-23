package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

const (
	// minPort is the minimum valid port number
	minPort = 1
	// maxPort is the maximum valid port number
	maxPort = 65535
)

// SeparateConnections separates connections into LISTEN and ESTABLISHED,
// then returns them combined with LISTEN first.
func SeparateConnections(connections []net.ConnectionStat) []net.ConnectionStat {
	listeningConns := make([]net.ConnectionStat, 0)
	establishedConns := make([]net.ConnectionStat, 0)

	for _, conn := range connections {
		if conn.Status == "LISTEN" || (IsUDP(conn.Type) && conn.Laddr.Port > 0) {
			listeningConns = append(listeningConns, conn)
		} else if conn.Status == "ESTABLISHED" {
			establishedConns = append(establishedConns, conn)
		}
	}

	return append(listeningConns, establishedConns...)
}

// FindResult represents a single find result (either by PID or by port)
type FindResult struct {
	Type        string               // "pid" or "port"
	PID         int32                // Process ID
	Port        uint32               // Port number (if found by port)
	Connections []net.ConnectionStat // Connections for this process/port
	ProcessInfo ProcessInfo          // Process information
}

// FindByPID finds a process by PID and returns its information.
func FindByPID(pid int32) (*FindResult, error) {
	_, err := process.NewProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("process with PID %d not found: %w", pid, err)
	}

	// Get ports used by the process
	connections, err := net.ConnectionsPid("inet", pid)
	if err != nil {
		connections = []net.ConnectionStat{}
	}

	// Separate and order connections: LISTEN first, then ESTABLISHED
	allConns := SeparateConnections(connections)

	processInfo := GetProcessInfo(pid)

	return &FindResult{
		Type:        "pid",
		PID:         pid,
		Connections: allConns,
		ProcessInfo: processInfo,
	}, nil

}

// FindByPort finds processes using a specific port.
func FindByPort(port int) ([]*FindResult, error) {
	if port < minPort || port > maxPort {
		return nil, fmt.Errorf("port number must be between %d and %d", minPort, maxPort)
	}

	// Get all network connections
	connections, err := net.Connections("inet")
	if err != nil {
		return nil, fmt.Errorf("failed to get network connection information: %w", err)
	}

	// Find PIDs listening on the specified port
	pidSet := make(map[int32]bool)

	for _, conn := range connections {
		// Match LISTEN status or UDP with matching port
		isListening := conn.Status == "LISTEN" || (IsUDP(conn.Type) && conn.Laddr.Port > 0)
		if isListening && conn.Laddr.Port == uint32(port) {
			pidSet[int32(conn.Pid)] = true
		}
	}

	// Get all connection information for found PIDs (LISTEN + ESTABLISHED)
	results := make([]*FindResult, 0, len(pidSet))
	for pid := range pidSet {
		// Get all connections for this PID
		pidConnections, err := net.ConnectionsPid("inet", pid)
		if err != nil {
			pidConnections = []net.ConnectionStat{}
		}

		// Separate and order connections: LISTEN first, then ESTABLISHED
		allConns := SeparateConnections(pidConnections)

		processInfo := GetProcessInfo(pid)
		results = append(results, &FindResult{
			Type:        "port",
			PID:         pid,
			Port:        uint32(port),
			Connections: allConns,
			ProcessInfo: processInfo,
		})
	}

	return results, nil
}

// FindByProcessName finds processes by name or command line (partial match).
// It searches both process name and full command line (like ps -ef).
func FindByProcessName(name string) ([]*FindResult, error) {
	if name == "" {
		return nil, fmt.Errorf("process name cannot be empty")
	}

	// Get all processes
	pids, err := process.Pids()
	if err != nil {
		return nil, fmt.Errorf("failed to get process list: %w", err)
	}

	results := make([]*FindResult, 0)
	nameLower := strings.ToLower(name)

	for _, pid := range pids {
		proc, err := process.NewProcess(pid)
		if err != nil {
			continue // Skip if process is not accessible
		}

		// Get process name
		procName, err := proc.Name()
		if err != nil {
			continue
		}

		// Get command line
		cmdline, _ := proc.Cmdline()
		cmdlineLower := strings.ToLower(cmdline)

		// Search for partial match in process name or command line (case-insensitive)
		matched := strings.Contains(strings.ToLower(procName), nameLower) ||
			(cmdline != "" && strings.Contains(cmdlineLower, nameLower))

		if matched {
			// Get ports used by the process
			connections, err := net.ConnectionsPid("inet", pid)
			if err != nil {
				connections = []net.ConnectionStat{}
			}

			// Separate and order connections: LISTEN first, then ESTABLISHED
			allConns := SeparateConnections(connections)

			processInfo := GetProcessInfo(pid)
			results = append(results, &FindResult{
				Type:        "name",
				PID:         pid,
				Connections: allConns,
				ProcessInfo: processInfo,
			})
		}

	}

	return results, nil
}

// FindByInput automatically determines search type and finds processes.
// If input is a number, searches by both PID and port.
// If input is a string, searches by process name.
func FindByInput(input string) ([]*FindResult, error) {
	// Check if input is a number
	if num, err := strconv.Atoi(input); err == nil {
		// If number, search by both PID and port
		results := make([]*FindResult, 0)
		pidSet := make(map[int32]bool) // PID set for deduplication

		// Search by PID
		if pidResult, err := FindByPID(int32(num)); err == nil {
			results = append(results, pidResult)
			pidSet[pidResult.PID] = true
		}

		// Search by port (exclude already found by PID)
		if portResults, err := FindByPort(num); err == nil {
			for _, portResult := range portResults {
				if !pidSet[portResult.PID] {
					results = append(results, portResult)
					pidSet[portResult.PID] = true
				}
			}
		}

		if len(results) == 0 {
			return nil, fmt.Errorf("no process found with PID or port: %s", input)
		}

		return results, nil
	}

	// If string, search by process name
	return FindByProcessName(input)
}
