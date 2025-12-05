//go:build linux

package provider

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

// LinuxRouteProvider is a RouteProvider implementation using Linux netlink.
type LinuxRouteProvider struct{}

// NewLinuxRouteProvider creates a new LinuxRouteProvider instance.
func NewLinuxRouteProvider() *LinuxRouteProvider {
	return &LinuxRouteProvider{}
}

// NewRouteProvider returns a RouteProvider implementation for the current OS.
// On Linux, it returns a LinuxRouteProvider.
func NewRouteProvider() RouteProvider {
	return NewLinuxRouteProvider()
}

// GetRoutes retrieves the routing table using netlink.
// Returns a slice of RouteEntry and an error if the operation fails.
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

