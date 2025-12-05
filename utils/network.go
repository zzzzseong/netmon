package utils

import (
	"fmt"

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


// GetProcessInfo retrieves process information for the given PID.
// It returns a ProcessInfo struct with name, username, CPU and memory usage.
// If the process cannot be found or accessed, it returns default values ("N/A").
// This function silently handles errors to ensure it always returns a valid ProcessInfo.
func GetProcessInfo(pid int32) ProcessInfo {
	// 기본값 설정
	info := ProcessInfo{
		Name:       "N/A",
		Username:   "N/A",
		CPUPercent: "N/A",
		MemPercent: "N/A",
	}

	if pid <= 0 {
		return info
	}

	proc, err := process.NewProcess(pid)
	if err != nil {
		// Process not found or inaccessible, return default values
		return info
	}

	// 프로세스 이름 (필수 정보이므로 먼저 가져옴)
	name, nameErr := proc.Name()
	if nameErr == nil {
		info.Name = name
	}

	// 사용자 이름
	username, userErr := proc.Username()
	if userErr == nil {
		info.Username = username
	}

	// CPU 사용률 (에러 무시)
	cpu, cpuErr := proc.CPUPercent()
	if cpuErr == nil {
		info.CPUPercent = fmt.Sprintf("%.1f%%", cpu)
	}

	// 메모리 사용률 (에러 무시)
	mem, memErr := proc.MemoryPercent()
	if memErr == nil {
		info.MemPercent = fmt.Sprintf("%.1f%%", mem)
	}

	// 명령어 라인 전체 가져오기 (ps -ef처럼)
	cmdline, cmdlineErr := proc.Cmdline()
	if cmdlineErr == nil && cmdline != "" {
		// 너무 긴 명령어 라인은 잘라내기 (터미널 출력 방지)
		const maxCmdlineLength = 500
		if len(cmdline) > maxCmdlineLength {
			cmdline = cmdline[:maxCmdlineLength] + "..."
		}
		info.Cmdline = cmdline
	} else {
		info.Cmdline = "N/A"
	}

	return info
}

