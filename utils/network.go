package utils

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// ConnectionType 상수
const (
	TCP ConnectionType = 1
	UDP ConnectionType = 2
)

type ConnectionType uint32

// ProcessInfo는 프로세스 정보를 담는 구조체
type ProcessInfo struct {
	Name       string
	Username   string
	CPUPercent string
	MemPercent string
}

// ConnectionTypeToString 연결 타입을 문자열로 변환
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

// IsUDP 연결 타입이 UDP인지 확인
func IsUDP(connType uint32) bool {
	return ConnectionType(connType) == UDP
}

// FilterListeningConnections LISTEN 상태인 연결만 필터링
func FilterListeningConnections(connections []net.ConnectionStat) map[string]net.ConnectionStat {
	listeningConns := make(map[string]net.ConnectionStat)
	for _, conn := range connections {
		// TCP는 LISTEN 상태만, UDP는 포트가 열려있으면 모두 포함
		if conn.Status == "LISTEN" || (IsUDP(conn.Type) && conn.Laddr.Port > 0) {
			key := fmt.Sprintf("%d:%d", conn.Type, conn.Laddr.Port)
			listeningConns[key] = conn
		}
	}
	return listeningConns
}

// GetProcessInfo 프로세스 정보를 가져옴
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

	return info
}

