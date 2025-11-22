//go:build darwin || freebsd || openbsd || netbsd

package provider

import (
	"fmt"
	"net"

	"golang.org/x/net/route"
)

// BSDRouteProvider는 BSD routing socket을 사용하는 RouteProvider 구현입니다
type BSDRouteProvider struct{}

// NewBSDRouteProvider는 새로운 BSDRouteProvider를 생성합니다
func NewBSDRouteProvider() *BSDRouteProvider {
	return &BSDRouteProvider{}
}

// NewRouteProvider는 현재 OS에 맞는 RouteProvider를 반환합니다
func NewRouteProvider() RouteProvider {
	return NewBSDRouteProvider()
}

// GetRoutes는 BSD routing socket을 통해 라우팅 테이블을 조회합니다
func (p *BSDRouteProvider) GetRoutes() ([]RouteEntry, error) {
	// FetchRIB를 사용하여 라우팅 테이블 조회
	rib, err := route.FetchRIB(0, route.RIBTypeRoute, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch routing table: %w", err)
	}

	// RIB 메시지 파싱
	msgs, err := route.ParseRIB(route.RIBTypeRoute, rib)
	if err != nil {
		return nil, fmt.Errorf("failed to parse routing table: %w", err)
	}

	var entries []RouteEntry

	for _, msg := range msgs {
		routeMsg, ok := msg.(*route.RouteMessage)
		if !ok {
			continue
		}

		entry := RouteEntry{}

		// 목적지 주소 추출
		if routeMsg.Addrs[0] != nil {
			switch addr := routeMsg.Addrs[0].(type) {
			case *route.Inet4Addr:
				ip := net.IPv4(addr.IP[0], addr.IP[1], addr.IP[2], addr.IP[3])
				
				// 넷마스크와 함께 CIDR 형식으로 변환
				if routeMsg.Addrs[2] != nil { // RTA_NETMASK
					if mask, ok := routeMsg.Addrs[2].(*route.Inet4Addr); ok {
						maskIP := net.IPv4(mask.IP[0], mask.IP[1], mask.IP[2], mask.IP[3])
						prefixLen, _ := net.IPMask(maskIP.To4()).Size()
						
						if ip.Equal(net.IPv4zero) && prefixLen == 0 {
							entry.Destination = "default"
						} else {
							entry.Destination = fmt.Sprintf("%s/%d", ip.String(), prefixLen)
						}
					} else {
						entry.Destination = ip.String()
					}
				} else {
					if ip.Equal(net.IPv4zero) {
						entry.Destination = "default"
					} else {
						entry.Destination = ip.String() + "/32"
					}
				}
			}
		}

		// 게이트웨이 주소 추출
		if routeMsg.Addrs[1] != nil { // RTA_GATEWAY
			switch gw := routeMsg.Addrs[1].(type) {
			case *route.Inet4Addr:
				ip := net.IPv4(gw.IP[0], gw.IP[1], gw.IP[2], gw.IP[3])
				if !ip.Equal(net.IPv4zero) {
					entry.Gateway = ip.String()
				}
			case *route.LinkAddr:
				// Link-level 주소인 경우 인터페이스 이름 추출
				if gw.Name != "" {
					entry.Interface = gw.Name
				}
			}
		}

		// 인터페이스 추출
		if routeMsg.Index > 0 {
			if iface, err := net.InterfaceByIndex(routeMsg.Index); err == nil {
				if entry.Interface == "" {
					entry.Interface = iface.Name
				}
			}
		}

		// Source IP 추출 (RTA_IFA)
		if routeMsg.Addrs[5] != nil {
			if ifa, ok := routeMsg.Addrs[5].(*route.Inet4Addr); ok {
				entry.Source = net.IPv4(ifa.IP[0], ifa.IP[1], ifa.IP[2], ifa.IP[3]).String()
			}
		}

		// Source IP가 없고 인터페이스가 있으면 해당 인터페이스의 IP 가져오기
		if entry.Source == "" && entry.Interface != "" {
			if srcIP := getInterfaceIP(entry.Interface); srcIP != "" {
				entry.Source = srcIP
			}
		}

		// 루프백 및 시스템 라우트 필터링
		if entry.Interface == "lo0" {
			continue
		}

		// 유효한 항목만 추가
		if entry.Destination != "" {
			entries = append(entries, entry)
		}
	}

	return entries, nil
}

