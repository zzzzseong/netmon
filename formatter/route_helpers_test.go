package formatter

import (
	"testing"

	"netmon/provider"
)

func TestFormatLinuxStyle(t *testing.T) {
	tests := []struct {
		name     string
		entry    provider.RouteEntry
		expected string
	}{
		{
			name: "default route with gateway",
			entry: provider.RouteEntry{
				Destination: "default",
				Gateway:     "192.168.1.1",
				Interface:   "eth0",
			},
			expected: "default via 192.168.1.1 dev eth0",
		},
		{
			name: "network route without gateway",
			entry: provider.RouteEntry{
				Destination: "192.168.1.0/24",
				Interface:   "eth0",
				Source:      "192.168.1.100",
			},
			expected: "192.168.1.0/24 dev eth0 src 192.168.1.100",
		},
		{
			name: "route with metric",
			entry: provider.RouteEntry{
				Destination: "10.0.0.0/8",
				Gateway:     "192.168.1.1",
				Interface:   "eth0",
				Metric:      100,
			},
			expected: "10.0.0.0/8 via 192.168.1.1 dev eth0 metric 100",
		},
		{
			name: "complete route entry",
			entry: provider.RouteEntry{
				Destination: "172.16.0.0/12",
				Gateway:     "192.168.1.1",
				Interface:   "eth0",
				Metric:      50,
				Source:      "192.168.1.100",
			},
			expected: "172.16.0.0/12 via 192.168.1.1 dev eth0 metric 50 src 192.168.1.100",
		},
		{
			name: "minimal route entry",
			entry: provider.RouteEntry{
				Destination: "192.168.2.0/24",
			},
			expected: "192.168.2.0/24",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatLinuxStyle(tt.entry)
			if got != tt.expected {
				t.Errorf("FormatLinuxStyle() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestFormatLinuxStyleRoutes(t *testing.T) {
	entries := []provider.RouteEntry{
		{
			Destination: "default",
			Gateway:     "192.168.1.1",
			Interface:   "eth0",
		},
		{
			Destination: "192.168.1.0/24",
			Interface:   "eth0",
			Source:      "192.168.1.100",
		},
	}

	expected := "default via 192.168.1.1 dev eth0\n192.168.1.0/24 dev eth0 src 192.168.1.100"
	got := FormatLinuxStyleRoutes(entries)

	if got != expected {
		t.Errorf("FormatLinuxStyleRoutes() = %q, want %q", got, expected)
	}
}

func TestFormatLinuxStyleRoutes_Empty(t *testing.T) {
	entries := []provider.RouteEntry{}
	expected := ""
	got := FormatLinuxStyleRoutes(entries)

	if got != expected {
		t.Errorf("FormatLinuxStyleRoutes() with empty slice = %q, want %q", got, expected)
	}
}

// Benchmark tests
func BenchmarkFormatLinuxStyle(b *testing.B) {
	entry := provider.RouteEntry{
		Destination: "172.16.0.0/12",
		Gateway:     "192.168.1.1",
		Interface:   "eth0",
		Metric:      50,
		Source:      "192.168.1.100",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FormatLinuxStyle(entry)
	}
}

func BenchmarkFormatLinuxStyleRoutes(b *testing.B) {
	entries := make([]provider.RouteEntry, 10)
	for i := 0; i < 10; i++ {
		entries[i] = provider.RouteEntry{
			Destination: "192.168.1.0/24",
			Gateway:     "192.168.1.1",
			Interface:   "eth0",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FormatLinuxStyleRoutes(entries)
	}
}
