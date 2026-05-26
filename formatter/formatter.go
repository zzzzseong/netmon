package formatter

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v3/net"
	"netmon/style"
	"netmon/utils"
)

// 미리 정의된 스타일 (재사용)
var (
	protocolStyle = lipgloss.NewStyle().Foreground(style.SecondaryColor).Bold(true)
	pidStyle      = lipgloss.NewStyle().Foreground(style.InfoColor)
	cpuStyle      = lipgloss.NewStyle().Foreground(style.WarningColor)
	memStyle      = lipgloss.NewStyle().Foreground(style.InfoColor)
)

// PortTableFormatter formats port listings in table format.
// It holds a ProcessCache so that *process.Process objects are reused across
// Format calls, enabling gopsutil to compute accurate CPU deltas in watch mode.
type PortTableFormatter struct {
	cache *utils.ProcessCache
}

// NewPortTableFormatter creates a new PortTableFormatter instance.
func NewPortTableFormatter() *PortTableFormatter {
	return &PortTableFormatter{cache: utils.NewProcessCache()}
}

// getStatusStyled 상태에 따른 스타일 적용
func getStatusStyled(status string) string {
	if status == "" {
		status = "LISTEN"
	}

	var color lipgloss.Color
	switch status {
	case "LISTEN":
		color = style.SuccessColor
	case "ESTABLISHED":
		color = style.InfoColor
	default:
		color = style.WarningColor
	}

	return lipgloss.NewStyle().Foreground(color).Bold(true).Render(status)
}

// Format formats connection information as a table.
// It takes a slice of connections and returns a formatted table string with process information.
func (f *PortTableFormatter) Format(connections interface{}) string {
	// 슬라이스 또는 맵을 받을 수 있도록 처리
	var connSlice []net.ConnectionStat
	
	switch v := connections.(type) {
	case []net.ConnectionStat:
		connSlice = v
	case map[string]net.ConnectionStat:
		// 맵을 슬라이스로 변환
		connSlice = make([]net.ConnectionStat, 0, len(v))
		for _, conn := range v {
			connSlice = append(connSlice, conn)
		}
	default:
		return ""
	}

	// 사이즈가 확정된 슬라이스로 메모리 할당 최적화
	rows := make([][]string, 0, len(connSlice))
	
	// per-call dedup cache (avoids double lookup for same PID in one render)
	seen := make(map[int32]utils.ProcessInfo)

	for _, conn := range connSlice {
		// 프로토콜 및 주소
		protocol := utils.ConnectionTypeToString(conn.Type)
		localAddr := fmt.Sprintf("%s:%d", conn.Laddr.IP, conn.Laddr.Port)

		// 상태 스타일링
		statusStr := getStatusStyled(conn.Status)

		pid := int32(conn.Pid)
		var processInfo utils.ProcessInfo
		if info, ok := seen[pid]; ok {
			processInfo = info
		} else {
			processInfo = f.cache.GetProcessInfo(pid)
			seen[pid] = processInfo
		}
		pid32 := int(pid)

		// PID 스타일링
		pidStr := fmt.Sprintf("%d", pid32)
		if pid32 > 0 {
			pidStr = pidStyle.Render(pidStr)
		}

		rows = append(rows, []string{
			protocolStyle.Render(protocol),
			localAddr,
			statusStr,
			pidStr,
			processInfo.Name,
			processInfo.Username,
			cpuStyle.Render(processInfo.CPUPercent),
			memStyle.Render(processInfo.MemPercent),
		})
	}

	return CreateTable(rows, PortTableColumns, style.TableWidthPort)
}

// ProcessInfoFormatter formats process information for display.
type ProcessInfoFormatter struct{}

// NewProcessInfoFormatter creates a new ProcessInfoFormatter instance.
func NewProcessInfoFormatter() *ProcessInfoFormatter {
	return &ProcessInfoFormatter{}
}

// Format formats process information including PID, name, status, and active ports.
// Returns a formatted string with styled process information.
func (f *ProcessInfoFormatter) Format(header string, pid int, processName string, status []string, connections []net.ConnectionStat) string {
	return f.FormatWithCmdline(header, pid, processName, status, connections, "")
}

// FormatWithCmdline formats process information including command line.
func (f *ProcessInfoFormatter) FormatWithCmdline(header string, pid int, processName string, status []string, connections []net.ConnectionStat, cmdline string) string {
	// strings.Builder로 효율적인 문자열 생성
	var builder strings.Builder
	
	// 헤더 추가 (Found by ... 정보)
	if header != "" {
		headerStyle := lipgloss.NewStyle().
			Foreground(style.SecondaryColor).
			Bold(true)
		builder.WriteString(headerStyle.Render("🔍 " + header))
		builder.WriteString("\n\n")
	}
	
	// 프로세스 기본 정보
	builder.WriteString(style.LabelStyle.Render("PID:"))
	builder.WriteString(style.ValueStyle.Render(fmt.Sprintf("%d", pid)))
	builder.WriteString("\n")
	
	builder.WriteString(style.LabelStyle.Render("Name:"))
	builder.WriteString(style.ValueStyle.Render(processName))
	builder.WriteString("\n")
	
	builder.WriteString(style.LabelStyle.Render("Status:"))
	builder.WriteString(style.ValueStyle.Render(fmt.Sprintf("%v", status)))
	builder.WriteString("\n")

	// 명령어 라인 추가 (ps -ef처럼)
	if cmdline != "" && cmdline != "N/A" {
		builder.WriteString("\n")
		cmdlineStyle := lipgloss.NewStyle().Foreground(style.SubtleColor)
		builder.WriteString(cmdlineStyle.Render("Command:"))
		builder.WriteString("\n")
		cmdlineValueStyle := lipgloss.NewStyle().Foreground(style.InfoColor).MarginLeft(2)
		builder.WriteString(cmdlineValueStyle.Render(cmdline))
		builder.WriteString("\n")
	}


	// 연결 정보 추가 (LISTEN과 ESTABLISHED 분리)
	if len(connections) > 0 {
		subtleStyle := lipgloss.NewStyle().Foreground(style.SubtleColor)
		portStyle := lipgloss.NewStyle().Foreground(style.InfoColor).MarginLeft(2)
		connStyle := lipgloss.NewStyle().Foreground(style.WarningColor).MarginLeft(2)
		
		// LISTEN 포트와 ESTABLISHED 연결 분리
		listeningPorts := make([]net.ConnectionStat, 0)
		establishedConns := make([]net.ConnectionStat, 0)
		
		for _, conn := range connections {
			if conn.Status == "LISTEN" || (conn.Status == "" && utils.IsUDP(conn.Type)) {
				listeningPorts = append(listeningPorts, conn)
			} else if conn.Status == "ESTABLISHED" {
				establishedConns = append(establishedConns, conn)
			}
		}
		
		// Listening Ports 표시
		if len(listeningPorts) > 0 {
			builder.WriteString("\n")
			builder.WriteString(subtleStyle.Render("Listening Ports:"))
			builder.WriteString("\n")
			
			for _, conn := range listeningPorts {
				portInfo := fmt.Sprintf("  • %s:%d (%s)",
					conn.Laddr.IP,
					conn.Laddr.Port,
					utils.ConnectionTypeToString(conn.Type))
				builder.WriteString(portStyle.Render(portInfo))
				builder.WriteString("\n")
			}
		}
		
		// Active Connections 표시
		if len(establishedConns) > 0 {
			builder.WriteString("\n")
			builder.WriteString(subtleStyle.Render("Active Connections:"))
			builder.WriteString("\n")
			
			for _, conn := range establishedConns {
				connInfo := fmt.Sprintf("  • %s:%d → %s:%d",
					conn.Laddr.IP,
					conn.Laddr.Port,
					conn.Raddr.IP,
					conn.Raddr.Port)
				builder.WriteString(connStyle.Render(connInfo))
				builder.WriteString("\n")
			}
		}
	}

	infoBox := style.InfoBoxStyle.Render(builder.String())
	return infoBox
}



