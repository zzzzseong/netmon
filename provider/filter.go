package provider

import (
	"net"
	"strings"
)

// ShouldIncludeRoute는 라우트를 출력에 포함할지 결정합니다
func ShouldIncludeRoute(entry RouteEntry) bool {
	dest := entry.Destination

	// 1. default 게이트웨이는 항상 포함
	if dest == "default" {
		return true
	}

	// 2. /32 단일 호스트 라우트 제외
	if strings.HasSuffix(dest, "/32") {
		return false
	}

	// 3. CIDR이 없는 단일 호스트 IP 제외
	if !strings.Contains(dest, "/") {
		return false
	}

	// 4. CIDR 파싱하여 특정 범위 필터링
	_, ipNet, err := net.ParseCIDR(dest)
	if err != nil {
		// 파싱 실패 시 포함 (예: "default" 등)
		return true
	}

	ip := ipNet.IP
	
	// IP 길이 체크 (IPv4만 처리)
	if len(ip) < 4 {
		return false
	}

	// 5-8. 특수 IP 범위 체크를 한 번에 처리
	firstOctet := ip[0]
	switch {
	case firstOctet == 127: // 루프백 (127.0.0.0/8)
		return false
	case firstOctet == 169 && ip[1] == 254: // Link-local (169.254.0.0/16)
		return false
	case firstOctet >= 224 && firstOctet <= 239: // Multicast (224.0.0.0/4)
		return false
	case firstOctet == 255: // Broadcast
		return false
	}

	// 9-10. 서브넷 크기에 따른 필터링
	ones, _ := ipNet.Mask.Size()
	if ones <= 24 {
		return true
	}

	// /25 ~ /31은 게이트웨이가 없는 경우만 포함 (직접 연결된 네트워크)
	return ones <= 31 && (entry.Gateway == "" || entry.Gateway == "-")
}

// FilterRoutes는 RouteEntry 슬라이스를 필터링합니다
func FilterRoutes(entries []RouteEntry) []RouteEntry {
	var filtered []RouteEntry
	for _, entry := range entries {
		if ShouldIncludeRoute(entry) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

