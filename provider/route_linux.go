//go:build linux

package provider

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

// LinuxRouteProvider는 Linux netlink를 사용하는 RouteProvider 구현입니다
type LinuxRouteProvider struct{}

// NewLinuxRouteProvider는 새로운 LinuxRouteProvider를 생성합니다
func NewLinuxRouteProvider() *LinuxRouteProvider {
	return &LinuxRouteProvider{}
}

// NewRouteProvider는 현재 OS에 맞는 RouteProvider를 반환합니다
func NewRouteProvider() RouteProvider {
	return NewLinuxRouteProvider()
}

// GetRoutes는 netlink를 통해 라우팅 테이블을 조회합니다
func (p *LinuxRouteProvider) GetRoutes() ([]RouteEntry, error) {
	// IPv4 라우트만 조회
	routes, err := netlink.RouteList(nil, netlink.FAMILY_V4)
	if err != nil {
		return nil, fmt.Errorf("failed to get routes via netlink: %w", err)
	}

	var entries []RouteEntry
	for _, route := range routes {
		entry := RouteEntry{}

		// Destination 처리
		if route.Dst == nil {
			entry.Destination = "default"
		} else {
			entry.Destination = route.Dst.String()
		}

		// Gateway 처리
		if route.Gw != nil {
			entry.Gateway = route.Gw.String()
		}

		// Interface 처리 (LinkIndex로부터)
		if route.LinkIndex > 0 {
			link, err := netlink.LinkByIndex(route.LinkIndex)
			if err == nil {
				entry.Interface = link.Attrs().Name
			}
		}

		// Metric 처리
		entry.Metric = route.Priority

		// Source 처리
		if route.Src != nil {
			entry.Source = route.Src.String()
		} else if entry.Interface != "" {
			// Source가 없으면 인터페이스에서 IP 주소를 가져옴
			if link, err := netlink.LinkByName(entry.Interface); err == nil {
				addrs, err := netlink.AddrList(link, netlink.FAMILY_V4)
				if err == nil && len(addrs) > 0 && addrs[0].IPNet != nil {
					entry.Source = addrs[0].IPNet.IP.String()
				}
			}
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

