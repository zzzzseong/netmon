package cmd

import (
	"os"
	"runtime"

	"github.com/shirou/gopsutil/v4/net"
	"netmon/style"
)

// hiddenPIDHint returns a one-line hint when connections are missing their PID
// because the process belongs to another user. This only happens on Linux when
// netmon is not running as root; other platforms expose the PID regardless.
// Returns an empty string when there is nothing to say.
func hiddenPIDHint(connections []net.ConnectionStat) string {
	if runtime.GOOS != "linux" || os.Geteuid() == 0 {
		return ""
	}
	for _, conn := range connections {
		if conn.Pid == 0 {
			return style.DescStyle.Render("Some processes are hidden (PID 0). Run with sudo to see connections owned by other users.")
		}
	}
	return ""
}
