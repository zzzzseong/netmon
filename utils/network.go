package utils

import (
	"fmt"
	"sync"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// ConnectionType represents a network connection type.
type ConnectionType uint32

const (
	// TCP represents a TCP connection
	TCP ConnectionType = 1
	// UDP represents a UDP connection
	UDP ConnectionType = 2
)

// ProcessInfo contains process information including name, username, CPU and memory usage.
type ProcessInfo struct {
	Name       string // Process name
	Username   string // Username running the process
	CPUPercent string // CPU usage percentage as a formatted string
	MemPercent string // Memory usage percentage as a formatted string
	Cmdline    string // Full command line (like ps -ef)
}

// ConnectionTypeToString converts a connection type to its string representation.
// It returns "TCP", "UDP", or "UNKNOWN" based on the connection type.
func ConnectionTypeToString(connType uint32) string {
	switch ConnectionType(connType) {
	case TCP:
		return "TCP"
	case UDP:
		return "UDP"
	default:
		return "UNKNOWN"
	}
}

// IsUDP checks if the connection type is UDP.
func IsUDP(connType uint32) bool {
	return ConnectionType(connType) == UDP
}

// FilterListeningConnections filters connections to only include listening connections.
// If includeUDP is true, UDP connections are also included.
// Returns a map keyed by connection type and port.
func FilterListeningConnections(connections []net.ConnectionStat, includeUDP bool) map[string]net.ConnectionStat {
	listeningConns := make(map[string]net.ConnectionStat)
	for _, conn := range connections {
		// TCP는 LISTEN 상태만 포함
		if conn.Status == "LISTEN" {
			key := fmt.Sprintf("%d:%d", conn.Type, conn.Laddr.Port)
			listeningConns[key] = conn
		} else if includeUDP && IsUDP(conn.Type) && conn.Laddr.Port > 0 {
			// UDP는 includeUDP가 true일 때만 포함
			key := fmt.Sprintf("%d:%d", conn.Type, conn.Laddr.Port)
			listeningConns[key] = conn
		}
	}
	return listeningConns
}

// FilterEstablishedConnections filters connections to only include ESTABLISHED connections.
// Returns a slice of connections that are in ESTABLISHED state.
func FilterEstablishedConnections(connections []net.ConnectionStat) []net.ConnectionStat {
	established := make([]net.ConnectionStat, 0)
	for _, conn := range connections {
		if conn.Status == "ESTABLISHED" {
			established = append(established, conn)
		}
	}
	return established
}


// ProcessCache holds *process.Process objects across calls so that gopsutil can
// compute CPU deltas correctly. CPUPercent() always returns 0 on first call to
// a new process object; reusing the same object gives accurate values.
type ProcessCache struct {
	mu    sync.Mutex
	procs map[int32]*process.Process
}

// NewProcessCache creates an empty ProcessCache.
func NewProcessCache() *ProcessCache {
	return &ProcessCache{procs: make(map[int32]*process.Process)}
}

// GetProcessInfo returns process info, reusing the cached *process.Process so
// that CPU percent is computed as a delta from the previous call.
func (c *ProcessCache) GetProcessInfo(pid int32) ProcessInfo {
	info := ProcessInfo{Name: "N/A", Username: "N/A", CPUPercent: "N/A", MemPercent: "N/A", Cmdline: "N/A"}
	if pid <= 0 {
		return info
	}

	c.mu.Lock()
	proc, ok := c.procs[pid]
	if !ok {
		var err error
		proc, err = process.NewProcess(pid)
		if err != nil {
			c.mu.Unlock()
			return info
		}
		c.procs[pid] = proc
	}
	c.mu.Unlock()

	return fetchProcessInfo(proc)
}

// GetProcessInfo retrieves process information for the given PID.
// Creates a new process object each call — CPUPercent will be 0 on this first
// call. Use ProcessCache.GetProcessInfo for accurate CPU values across calls.
func GetProcessInfo(pid int32) ProcessInfo {
	info := ProcessInfo{Name: "N/A", Username: "N/A", CPUPercent: "N/A", MemPercent: "N/A", Cmdline: "N/A"}
	if pid <= 0 {
		return info
	}
	proc, err := process.NewProcess(pid)
	if err != nil {
		return info
	}
	return fetchProcessInfo(proc)
}

func fetchProcessInfo(proc *process.Process) ProcessInfo {
	info := ProcessInfo{Name: "N/A", Username: "N/A", CPUPercent: "N/A", MemPercent: "N/A", Cmdline: "N/A"}

	if name, err := proc.Name(); err == nil {
		info.Name = name
	}
	if user, err := proc.Username(); err == nil {
		info.Username = user
	}
	if cpu, err := proc.CPUPercent(); err == nil {
		info.CPUPercent = fmt.Sprintf("%.1f%%", cpu)
	}
	if mem, err := proc.MemoryPercent(); err == nil {
		info.MemPercent = fmt.Sprintf("%.1f%%", mem)
	}
	const maxCmdlineLength = 500
	if cmd, err := proc.Cmdline(); err == nil && cmd != "" {
		if len(cmd) > maxCmdlineLength {
			cmd = cmd[:maxCmdlineLength] + "..."
		}
		info.Cmdline = cmd
	}
	return info
}

