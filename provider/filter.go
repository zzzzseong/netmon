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

	// 5. 루프백 (127.0.0.0/8) 제외
	if ip[0] == 127 {
		return false
	}

	// 6. Link-local (169.254.0.0/16) 제외
	if ip[0] == 169 && ip[1] == 254 {
		return false
	}

	// 7. Multicast (224.0.0.0/4) 제외
	if ip[0] >= 224 && ip[0] <= 239 {
		return false
	}

	// 8. Broadcast (255.255.255.255/32) 제외
	if ip[0] == 255 {
		return false
	}

	// 9. 매우 작은 서브넷은 포함 (직접 연결된 네트워크)
	// /16 이하의 네트워크는 일반적으로 유의미한 라우트
	ones, _ := ipNet.Mask.Size()
	if ones <= 24 {
		return true
	}

	// 10. /25 ~ /31은 선택적으로 포함 (Docker 등의 네트워크일 수 있음)
	// 게이트웨이가 없으면 직접 연결된 네트워크일 가능성이 높음
	if ones >= 25 && ones <= 31 {
		return entry.Gateway == "" || entry.Gateway == "-"
	}

	return false
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

