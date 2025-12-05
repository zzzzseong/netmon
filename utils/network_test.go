package utils

import (
	"testing"

	"github.com/shirou/gopsutil/v3/net"
)

func TestConnectionTypeToString(t *testing.T) {
	tests := []struct {
		name     string
		connType uint32
		want     string
	}{
		{"TCP", uint32(TCP), "TCP"},
		{"UDP", uint32(UDP), "UDP"},
		{"Unknown", 999, "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConnectionTypeToString(tt.connType); got != tt.want {
				t.Errorf("ConnectionTypeToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUDP(t *testing.T) {
	tests := []struct {
		name     string
		connType uint32
		want     bool
	}{
		{"TCP", uint32(TCP), false},
		{"UDP", uint32(UDP), true},
		{"Unknown", 999, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUDP(tt.connType); got != tt.want {
				t.Errorf("IsUDP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterListeningConnections(t *testing.T) {
	connections := []net.ConnectionStat{
		{Type: uint32(TCP), Status: "LISTEN", Laddr: net.Addr{Port: 80, IP: "0.0.0.0"}},
		{Type: uint32(TCP), Status: "ESTABLISHED", Laddr: net.Addr{Port: 443, IP: "0.0.0.0"}},
		{Type: uint32(UDP), Status: "NONE", Laddr: net.Addr{Port: 53, IP: "0.0.0.0"}},
		{Type: uint32(UDP), Status: "", Laddr: net.Addr{Port: 123, IP: "0.0.0.0"}}, // UDP often has empty status
	}

	t.Run("Only TCP Listen", func(t *testing.T) {
		got := FilterListeningConnections(connections, false)
		if len(got) != 1 {
			t.Errorf("FilterListeningConnections() length = %v, want 1", len(got))
		}
		if _, ok := got["1:80"]; !ok {
			t.Errorf("Expected TCP port 80 to be present")
		}
	})

	t.Run("Include UDP", func(t *testing.T) {
		got := FilterListeningConnections(connections, true)
		if len(got) != 3 {
			t.Errorf("FilterListeningConnections() length = %v, want 3", len(got))
		}
		if _, ok := got["1:80"]; !ok {
			t.Errorf("Expected TCP port 80 to be present")
		}
		if _, ok := got["2:53"]; !ok {
			t.Errorf("Expected UDP port 53 to be present")
		}
		if _, ok := got["2:123"]; !ok {
			t.Errorf("Expected UDP port 123 to be present")
		}
	})
}
