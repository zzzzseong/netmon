package formatter

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/shirou/gopsutil/v3/net"
	"netmon/style"
	"netmon/utils"
)

// ConnectionTableFormatter formats active connections in table format.
type ConnectionTableFormatter struct{}

// NewConnectionTableFormatter creates a new ConnectionTableFormatter instance.
func NewConnectionTableFormatter() *ConnectionTableFormatter {
	return &ConnectionTableFormatter{}
}

// Format formats active connection information as a table.
// It takes a slice of connections and returns a formatted table string with process information.
func (f *ConnectionTableFormatter) Format(connections []net.ConnectionStat) string {
	// Sort connections by remote address for better readability
	sort.Slice(connections, func(i, j int) bool {
		if connections[i].Raddr.IP != connections[j].Raddr.IP {
			return connections[i].Raddr.IP < connections[j].Raddr.IP
		}
		return connections[i].Raddr.Port < connections[j].Raddr.Port
	})

	// Pre-allocate rows slice
	rows := make([][]string, 0, len(connections))

	// Process info cache
	processCache := make(map[int32]utils.ProcessInfo)

	for _, conn := range connections {
		// Protocol
		protocol := utils.ConnectionTypeToString(conn.Type)
		protocolStyled := lipgloss.NewStyle().Foreground(style.SecondaryColor).Bold(true).Render(protocol)

		// Local address
		localAddr := fmt.Sprintf("%s:%d", conn.Laddr.IP, conn.Laddr.Port)
		localAddrStyled := lipgloss.NewStyle().Foreground(style.InfoColor).Render(localAddr)

		// Remote address
		remoteAddr := fmt.Sprintf("%s:%d", conn.Raddr.IP, conn.Raddr.Port)
		remoteAddrStyled := lipgloss.NewStyle().Foreground(style.WarningColor).Render(remoteAddr)

		// Process info (with caching)
		pid := int32(conn.Pid)
		var processInfo utils.ProcessInfo

		if info, ok := processCache[pid]; ok {
			processInfo = info
		} else {
			processInfo = utils.GetProcessInfo(pid)
			processCache[pid] = processInfo
		}

		// PID styling
		pidStr := fmt.Sprintf("%d", pid)
		if pid > 0 {
			pidStr = lipgloss.NewStyle().Foreground(style.InfoColor).Render(pidStr)
		}

		rows = append(rows, []string{
			protocolStyled,
			localAddrStyled,
			remoteAddrStyled,
			pidStr,
			processInfo.Name,
		})
	}

	// Create table headers
	headerStyle := style.HeaderStyle
	headerStyle = headerStyle.Align(lipgloss.Center)
	styledHeaders := make([]string, len(ConnectionTableColumns))
	for i, col := range ConnectionTableColumns {
		styledHeaders[i] = headerStyle.Width(col.Width).Render(col.Title)
	}

	// Create and style table
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(style.TableBorderStyle).
		StyleFunc(GetTableRowStyle).
		Headers(styledHeaders...).
		Rows(rows...).
		Width(style.TableWidthConnection)

	return t.String()
}
