package provider

import "testing"

func TestConvertToCIDR(t *testing.T) {
	got := convertToCIDR("192.168.1.0", "255.255.255.0")
	if got != "192.168.1.0/24" {
		t.Fatalf("unexpected cidr conversion: got %q", got)
	}
}

func TestParseWindowsRoute_DefaultRoute(t *testing.T) {
	output := `
===========================================================================
IPv4 Route Table
===========================================================================
Active Routes:
Network Destination        Netmask          Gateway       Interface  Metric
          0.0.0.0          0.0.0.0      10.0.0.1      10.0.0.20     25
        10.0.0.0    255.255.255.0         On-link      10.0.0.20    281
===========================================================================
`

	routes := parseWindowsRoute(output)
	if len(routes) < 2 {
		t.Fatalf("expected at least 2 routes, got %d", len(routes))
	}

	if routes[0].Destination != "default" {
		t.Fatalf("expected default destination, got %q", routes[0].Destination)
	}
	if routes[0].Gateway != "10.0.0.1" {
		t.Fatalf("expected default gateway 10.0.0.1, got %q", routes[0].Gateway)
	}
}

func TestParseLinuxRoute_DefaultAndMetric(t *testing.T) {
	output := `
default via 192.168.0.1 dev eth0 proto dhcp src 192.168.0.10 metric 100
10.0.0.0/24 dev eth1 proto kernel scope link src 10.0.0.2 metric 50
`

	routes := parseLinuxRoute(output)
	if len(routes) != 2 {
		t.Fatalf("expected 2 routes, got %d", len(routes))
	}

	if routes[0].Destination != "default" || routes[0].Gateway != "192.168.0.1" || routes[0].Metric != 100 {
		t.Fatalf("unexpected default route parsing: %+v", routes[0])
	}
	if routes[1].Destination != "10.0.0.0/24" || routes[1].Interface != "eth1" || routes[1].Metric != 50 {
		t.Fatalf("unexpected subnet route parsing: %+v", routes[1])
	}
}
