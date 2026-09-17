//go:build !linux && !darwin && !freebsd && !openbsd && !netbsd && !windows

package provider

import (
	"fmt"
	"runtime"
)

// unsupportedRouteProvider is used on operating systems without a native implementation.
type unsupportedRouteProvider struct{}

// NewRouteProvider returns a RouteProvider implementation for the current OS.
// Release binaries are built only for Linux, macOS and Windows, which all have
// native providers; anything else reports the routing table as unavailable.
func NewRouteProvider() RouteProvider {
	return unsupportedRouteProvider{}
}

// GetRoutes always fails on unsupported operating systems.
func (unsupportedRouteProvider) GetRoutes() ([]RouteEntry, error) {
	return nil, fmt.Errorf("routing table is not supported on %s", runtime.GOOS)
}
