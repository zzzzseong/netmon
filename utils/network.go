package utils

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/net"
)

// 연결 타입을 문자열로 변환 (1=TCP, 2=UDP)
func ConnectionTypeToString(connType uint32) string {
	switch connType {
	case 1:
		return "TCP"
	case 2:
		return "UDP"
	default:
		return "UNKNOWN"
	}
}

// 연결 타입이 UDP인지 확인
func IsUDP(connType uint32) bool {
	return connType == 2
}

// LISTEN 상태인 연결만 필터링
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

