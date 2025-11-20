//go:build !linux && !darwin && !freebsd && !openbsd && !netbsd && !windows

package provider

// NewRouteProvider는 현재 OS에 맞는 RouteProvider를 반환합니다
func NewRouteProvider() RouteProvider {
	return NewFallbackRouteProvider()
}

