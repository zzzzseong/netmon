package formatter

import (
	"strings"
	"testing"

	"github.com/shirou/gopsutil/v3/net"
	"netmon/utils"
)

func TestPortTableFormatter_Format(t *testing.T) {
	f := NewPortTableFormatter()

	connections := map[string]net.ConnectionStat{
		"1:8080": {
			Type:   uint32(utils.TCP),
			Laddr:  net.Addr{IP: "127.0.0.1", Port: 8080},
			Status: "LISTEN",
			Pid:    1234,
		},
	}

	// Note: We can't easily mock utils.GetProcessInfo without refactoring, 
	// so the process name/user might be "N/A" or actual values if PID 1234 exists.
	// For this test, we focus on checking if the table contains the IP and Port.
	
	output := f.Format(connections)

	if !strings.Contains(output, "127.0.0.1:8080") {
		t.Errorf("Output should contain local address '127.0.0.1:8080', got:\n%s", output)
	}
	if !strings.Contains(output, "TCP") {
		t.Errorf("Output should contain protocol 'TCP', got:\n%s", output)
	}
	if !strings.Contains(output, "LISTEN") {
		t.Errorf("Output should contain status 'LISTEN', got:\n%s", output)
	}
	if !strings.Contains(output, "1234") {
		t.Errorf("Output should contain PID '1234', got:\n%s", output)
	}
}

func TestProcessInfoFormatter_Format(t *testing.T) {
	f := NewProcessInfoFormatter()

	pid := 5678
	name := "test-process"
	status := []string{"S"}
	connections := []net.ConnectionStat{
		{
			Type:   uint32(utils.TCP),
			Laddr:  net.Addr{IP: "0.0.0.0", Port: 9090},
			Status: "LISTEN",
		},
	}

	output := f.Format(pid, name, status, connections)

	if !strings.Contains(output, "5678") {
		t.Errorf("Output should contain PID '5678'")
	}
	if !strings.Contains(output, "test-process") {
		t.Errorf("Output should contain process name 'test-process'")
	}
	if !strings.Contains(output, "0.0.0.0:9090") {
		t.Errorf("Output should contain port '0.0.0.0:9090'")
	}
}
