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

// FindResult represents a single find result (either by PID or by port)
type FindResult struct {
	Type        string              // "pid" or "port"
	PID         int32               // Process ID
	Port        uint32              // Port number (if found by port)
	Connections []net.ConnectionStat // Connections for this process/port
	ProcessInfo ProcessInfo         // Process information
}

// FindByPID finds a process by PID and returns its information.
func FindByPID(pid int32) (*FindResult, error) {
	_, err := process.NewProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("process with PID %d not found: %w", pid, err)
	}

	// 프로세스가 사용하는 포트 가져오기
	connections, err := net.ConnectionsPid("inet", pid)
	if err != nil {
		connections = []net.ConnectionStat{}
	}

	// LISTEN 상태와 ESTABLISHED 상태 연결 분리
	listeningConns := make([]net.ConnectionStat, 0)
	establishedConns := make([]net.ConnectionStat, 0)
	
	for _, conn := range connections {
		if conn.Status == "LISTEN" || (IsUDP(conn.Type) && conn.Laddr.Port > 0) {
			listeningConns = append(listeningConns, conn)
		} else if conn.Status == "ESTABLISHED" {
			establishedConns = append(establishedConns, conn)
		}
	}

	// LISTEN 연결을 먼저, ESTABLISHED 연결을 나중에 추가
	allConns := append(listeningConns, establishedConns...)

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

	// 모든 네트워크 연결 가져오기
	connections, err := net.Connections("inet")
	if err != nil {
		return nil, fmt.Errorf("failed to get network connection information: %w", err)
	}

	// 해당 포트를 LISTEN하고 있는 PID 찾기
	pidSet := make(map[int32]bool)

	for _, conn := range connections {
		// LISTEN 상태 또는 UDP이면서 포트가 일치하는 경우
		isListening := conn.Status == "LISTEN" || (IsUDP(conn.Type) && conn.Laddr.Port > 0)
		if isListening && conn.Laddr.Port == uint32(port) {
			pidSet[int32(conn.Pid)] = true
		}
	}

	// 찾은 PID들의 모든 연결 정보 가져오기 (LISTEN + ESTABLISHED)
	results := make([]*FindResult, 0, len(pidSet))
	for pid := range pidSet {
		// 해당 PID의 모든 연결 가져오기
		pidConnections, err := net.ConnectionsPid("inet", pid)
		if err != nil {
			pidConnections = []net.ConnectionStat{}
		}

		// LISTEN과 ESTABLISHED 분리
		listeningConns := make([]net.ConnectionStat, 0)
		establishedConns := make([]net.ConnectionStat, 0)

		for _, conn := range pidConnections {
			if conn.Status == "LISTEN" || (IsUDP(conn.Type) && conn.Laddr.Port > 0) {
				listeningConns = append(listeningConns, conn)
			} else if conn.Status == "ESTABLISHED" {
				establishedConns = append(establishedConns, conn)
			}
		}

		// LISTEN 먼저, ESTABLISHED 나중에
		allConns := append(listeningConns, establishedConns...)

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

	// 모든 프로세스 가져오기
	pids, err := process.Pids()
	if err != nil {
		return nil, fmt.Errorf("failed to get process list: %w", err)
	}

	results := make([]*FindResult, 0)
	nameLower := strings.ToLower(name)

	for _, pid := range pids {
		proc, err := process.NewProcess(pid)
		if err != nil {
			continue // 프로세스에 접근할 수 없으면 건너뛰기
		}

		// 프로세스 이름 가져오기
		procName, err := proc.Name()
		if err != nil {
			continue
		}

		// 명령어 라인 가져오기
		cmdline, _ := proc.Cmdline()
		cmdlineLower := strings.ToLower(cmdline)

		// 프로세스 이름 또는 명령어 라인에서 부분 일치 검색 (대소문자 무시)
		matched := strings.Contains(strings.ToLower(procName), nameLower) ||
			(cmdline != "" && strings.Contains(cmdlineLower, nameLower))

		if matched {
			// 프로세스가 사용하는 포트 가져오기
			connections, err := net.ConnectionsPid("inet", pid)
			if err != nil {
				connections = []net.ConnectionStat{}
			}

			// LISTEN과 ESTABLISHED 분리
			listeningConns := make([]net.ConnectionStat, 0)
			establishedConns := make([]net.ConnectionStat, 0)
			
			for _, conn := range connections {
				if conn.Status == "LISTEN" || (IsUDP(conn.Type) && conn.Laddr.Port > 0) {
					listeningConns = append(listeningConns, conn)
				} else if conn.Status == "ESTABLISHED" {
					establishedConns = append(establishedConns, conn)
				}
			}

			// LISTEN 먼저, ESTABLISHED 나중에
			allConns := append(listeningConns, establishedConns...)

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
	// 숫자인지 확인
	if num, err := strconv.Atoi(input); err == nil {
		// 숫자면 PID와 포트 둘 다 검색
		results := make([]*FindResult, 0)
		pidSet := make(map[int32]bool) // 중복 제거를 위한 PID 집합

		// PID로 검색
		if pidResult, err := FindByPID(int32(num)); err == nil {
			results = append(results, pidResult)
			pidSet[pidResult.PID] = true
		}

		// 포트로 검색 (PID로 이미 찾은 것은 제외)
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

	// 문자열이면 프로세스 이름으로 검색
	return FindByProcessName(input)
}


