package cmd

import (
	"os"
	"testing"

	gonet "github.com/shirou/gopsutil/v3/net"
)

func TestGetTopProcessesByConnections_FiltersAndTopN(t *testing.T) {
	pid := int32(os.Getpid())

	connections := []gonet.ConnectionStat{
		// Counted: TCP established
		{Pid: pid, Type: 1, Status: "ESTABLISHED", Laddr: gonet.Addr{IP: "127.0.0.1", Port: 5000}},
		// Ignored: TCP not established
		{Pid: pid, Type: 1, Status: "TIME_WAIT", Laddr: gonet.Addr{IP: "127.0.0.1", Port: 5001}},
		// Counted: UDP with local port
		{Pid: pid, Type: 2, Status: "", Laddr: gonet.Addr{IP: "127.0.0.1", Port: 53}},
		// Ignored: UDP without local port
		{Pid: pid, Type: 2, Status: "", Laddr: gonet.Addr{IP: "127.0.0.1", Port: 0}},
	}

	top := getTopProcessesByConnections(connections, 1)
	if len(top) != 1 {
		t.Fatalf("expected exactly one top process, got %d", len(top))
	}
	if top[0].Count != 2 {
		t.Fatalf("expected 2 qualifying connections for top process, got %d", top[0].Count)
	}
}
