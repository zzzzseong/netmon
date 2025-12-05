package utils

import (
	"sort"

	"github.com/shirou/gopsutil/v3/net"
)

// SortConnectionsByPort converts a map of connections to a sorted slice by port number.
// Returns connections sorted in ascending order by local port.
func SortConnectionsByPort(conns map[string]net.ConnectionStat) []net.ConnectionStat {
	// 맵을 슬라이스로 변환
	connSlice := make([]net.ConnectionStat, 0, len(conns))
	for _, conn := range conns {
		connSlice = append(connSlice, conn)
	}

	// 포트 번호로 정렬 (오름차순)
	sort.Slice(connSlice, func(i, j int) bool {
		return connSlice[i].Laddr.Port < connSlice[j].Laddr.Port
	})

	return connSlice
}
