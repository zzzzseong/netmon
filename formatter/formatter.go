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
type PortTableFormatter struct{}

// NewPortTableFormatter creates a new PortTableFormatter instance.
func NewPortTableFormatter() *PortTableFormatter {
	return &PortTableFormatter{}
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
	
	// 프로세스 정보 캐싱을 위한 맵
	processCache := make(map[int32]utils.ProcessInfo)

	for _, conn := range connSlice {
		// 프로토콜 및 주소
		protocol := utils.ConnectionTypeToString(conn.Type)
		localAddr := fmt.Sprintf("%s:%d", conn.Laddr.IP, conn.Laddr.Port)

		// 상태 스타일링
		statusStr := getStatusStyled(conn.Status)

		// 프로세스 정보 가져오기 (캐싱 적용)
		pid := int(conn.Pid)
		var processInfo utils.ProcessInfo
		
		if info, ok := processCache[int32(pid)]; ok {
			processInfo = info
		} else {
			processInfo = utils.GetProcessInfo(int32(pid))
			processCache[int32(pid)] = processInfo
		}

		// PID 스타일링
		pidStr := fmt.Sprintf("%d", pid)
		if pid > 0 {
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

	// 포트 정보 추가
	if len(connections) > 0 {
		subtleStyle := lipgloss.NewStyle().Foreground(style.SubtleColor)
		portStyle := lipgloss.NewStyle().Foreground(style.InfoColor).MarginLeft(2)
		
		builder.WriteString("\n")
		builder.WriteString(subtleStyle.Render("Active Ports:"))
		builder.WriteString("\n")
		
		for _, conn := range connections {
			// LISTEN 상태이거나 UDP인 경우 표시
			if conn.Status == "LISTEN" || (conn.Status == "" && utils.IsUDP(conn.Type)) {
				portInfo := fmt.Sprintf("  • %s:%d (%s)",
					conn.Laddr.IP,
					conn.Laddr.Port,
					utils.ConnectionTypeToString(conn.Type))
				builder.WriteString(portStyle.Render(portInfo))
				builder.WriteString("\n")
			}
		}
	}

	infoBox := style.InfoBoxStyle.Render(builder.String())
	return infoBox
}

