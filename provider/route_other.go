//go:build !linux && !darwin && !freebsd && !openbsd && !netbsd && !windows

package provider

// NewRouteProvider returns a RouteProvider implementation for the current OS.
// For unsupported operating systems, it returns a FallbackRouteProvider.
func NewRouteProvider() RouteProvider {
	return NewFallbackRouteProvider()
}

