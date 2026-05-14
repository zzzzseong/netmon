package provider

import (
	"net"
	"testing"
)

func TestShouldIncludeRoute(t *testing.T) {
	tests := []struct {
		name     string
		entry    RouteEntry
		expected bool
	}{
		{
			name:     "default route should be included",
			entry:    RouteEntry{Destination: "default"},
			expected: true,
		},
		{
			name:     "/32 host route should be excluded",
			entry:    RouteEntry{Destination: "192.168.1.1/32"},
			expected: false,
		},
		{
			name:     "single IP without CIDR should be excluded",
			entry:    RouteEntry{Destination: "192.168.1.1"},
			expected: false,
		},
		{
			name:     "loopback route should be excluded",
			entry:    RouteEntry{Destination: "127.0.0.0/8"},
			expected: false,
		},
		{
			name:     "link-local route should be excluded",
			entry:    RouteEntry{Destination: "169.254.0.0/16"},
			expected: false,
		},
		{
			name:     "multicast route should be excluded",
			entry:    RouteEntry{Destination: "224.0.0.0/4"},
			expected: false,
		},
		{
			name:     "broadcast route should be excluded",
			entry:    RouteEntry{Destination: "255.255.255.255/32"},
			expected: false,
		},
		{
			name:     "/24 network should be included",
			entry:    RouteEntry{Destination: "192.168.1.0/24"},
			expected: true,
		},
		{
			name:     "/16 network should be included",
			entry:    RouteEntry{Destination: "172.16.0.0/16"},
			expected: true,
		},
		{
			name:     "/25 with gateway should be excluded",
			entry:    RouteEntry{Destination: "192.168.1.0/25", Gateway: "192.168.1.1"},
			expected: false,
		},
		{
			name:     "/25 without gateway should be included",
			entry:    RouteEntry{Destination: "192.168.1.0/25", Gateway: ""},
			expected: true,
		},
		{
			name:     "/30 without gateway should be included",
			entry:    RouteEntry{Destination: "10.0.0.0/30", Gateway: ""},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShouldIncludeRoute(tt.entry)
			if got != tt.expected {
				t.Errorf("ShouldIncludeRoute(%v) = %v, want %v", tt.entry.Destination, got, tt.expected)
			}
		})
	}
}

func TestFilterRoutes(t *testing.T) {
	entries := []RouteEntry{
		{Destination: "default"},
		{Destination: "192.168.1.0/24"},
		{Destination: "192.168.1.1/32"},
		{Destination: "127.0.0.0/8"},
		{Destination: "10.0.0.0/16"},
		{Destination: "169.254.0.0/16"},
	}

	filtered := FilterRoutes(entries)

	// Expected: default, 192.168.1.0/24, 10.0.0.0/16
	expectedCount := 3
	if len(filtered) != expectedCount {
		t.Errorf("FilterRoutes() returned %d routes, want %d", len(filtered), expectedCount)
	}

	// Verify default route is included
	foundDefault := false
	for _, entry := range filtered {
		if entry.Destination == "default" {
			foundDefault = true
			break
		}
	}
	if !foundDefault {
		t.Error("FilterRoutes() did not include default route")
	}

	// Verify /32 route is excluded
	for _, entry := range filtered {
		if entry.Destination == "192.168.1.1/32" {
			t.Error("FilterRoutes() should not include /32 routes")
		}
	}

	// Verify loopback is excluded
	for _, entry := range filtered {
		if entry.Destination == "127.0.0.0/8" {
			t.Error("FilterRoutes() should not include loopback routes")
		}
	}
}

func TestShouldIncludeRoute_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		entry    RouteEntry
		expected bool
	}{
		{
			name:     "empty destination",
			entry:    RouteEntry{Destination: ""},
			expected: false,
		},
		{
			name:     "invalid CIDR",
			entry:    RouteEntry{Destination: "invalid"},
			expected: false,
		},
		{
			name:     "IPv6 route should be excluded",
			entry:    RouteEntry{Destination: "fe80::/10"},
			expected: false,
		},
		{
			name:     "global unicast IPv6 should be excluded",
			entry:    RouteEntry{Destination: "2001:db8::/32"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShouldIncludeRoute(tt.entry)
			if got != tt.expected {
				t.Errorf("ShouldIncludeRoute(%v) = %v, want %v", tt.entry.Destination, got, tt.expected)
			}
		})
	}
}

// Benchmark for performance testing
func BenchmarkShouldIncludeRoute(b *testing.B) {
	entry := RouteEntry{Destination: "192.168.1.0/24"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ShouldIncludeRoute(entry)
	}
}

func BenchmarkFilterRoutes(b *testing.B) {
	entries := make([]RouteEntry, 100)
	for i := 0; i < 100; i++ {
		ip := net.IPv4(byte(i), byte(i), byte(i), 0)
		entries[i] = RouteEntry{
			Destination: ip.String() + "/24",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FilterRoutes(entries)
	}
}
