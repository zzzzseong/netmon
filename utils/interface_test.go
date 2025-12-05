package utils

import (
	"testing"
)

func TestGetInterfaceIP(t *testing.T) {
	tests := []struct {
		name      string
		ifname    string
		wantError bool
	}{
		{
			name:      "loopback interface",
			ifname:    "lo0",
			wantError: false,
		},
		{
			name:      "non-existent interface",
			ifname:    "nonexistent999",
			wantError: true,
		},
		{
			name:      "empty interface name",
			ifname:    "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip, err := GetInterfaceIP(tt.ifname)
			
			if tt.wantError {
				if err == nil {
					t.Errorf("GetInterfaceIP(%q) expected error, got nil", tt.ifname)
				}
			} else {
				if err != nil {
					// Some systems might not have lo0, try lo
					if tt.ifname == "lo0" {
						ip, err = GetInterfaceIP("lo")
						if err != nil {
							t.Skipf("Neither lo0 nor lo interface found on this system")
						}
					} else {
						t.Errorf("GetInterfaceIP(%q) unexpected error: %v", tt.ifname, err)
					}
				}
				if ip == "" && err == nil {
					t.Errorf("GetInterfaceIP(%q) returned empty IP without error", tt.ifname)
				}
			}
		})
	}
}

func TestGetInterfaceIP_Loopback(t *testing.T) {
	// Try common loopback interface names
	loopbackNames := []string{"lo", "lo0"}
	
	var foundLoopback bool
	for _, name := range loopbackNames {
		ip, err := GetInterfaceIP(name)
		if err == nil {
			foundLoopback = true
			if ip != "127.0.0.1" {
				t.Logf("Loopback interface %s has IP %s (expected 127.0.0.1)", name, ip)
			}
			break
		}
	}
	
	if !foundLoopback {
		t.Skip("No loopback interface found on this system")
	}
}

// Benchmark
func BenchmarkGetInterfaceIP(b *testing.B) {
	// Use loopback as it should exist on all systems
	ifname := "lo0"
	
	// Check if interface exists, otherwise try "lo"
	_, err := GetInterfaceIP(ifname)
	if err != nil {
		ifname = "lo"
		_, err = GetInterfaceIP(ifname)
		if err != nil {
			b.Skip("No loopback interface found")
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetInterfaceIP(ifname)
	}
}
