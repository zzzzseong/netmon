package provider

// RouteEntry represents a routing table entry.
type RouteEntry struct {
	Destination string // CIDR notation or "default" for default route
	Gateway     string // Gateway IP address, empty string if no gateway
	Interface   string // Network interface name (e.g., "en0", "eth0")
	Metric      int    // Route metric value (0 means not displayed)
	Source      string // Source IP address (optional)
}

// RouteProvider is an interface for querying routing tables across different operating systems.
// Each OS-specific implementation provides its own way to retrieve routing information.
type RouteProvider interface {
	// GetRoutes retrieves the routing table entries.
	// Returns a slice of RouteEntry and an error if the operation fails.
	GetRoutes() ([]RouteEntry, error)
}

