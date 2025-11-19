package formatter

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
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

// PortTableFormatter는 포트 목록을 테이블 형식으로 포맷팅합니다
type PortTableFormatter struct{}

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

// Format은 연결 정보를 테이블 형식으로 포맷팅합니다
func (f *PortTableFormatter) Format(connections map[string]net.ConnectionStat) string {
	var rows [][]string
	
	for _, conn := range connections {
		// 프로토콜 및 주소
		protocol := utils.ConnectionTypeToString(conn.Type)
		localAddr := fmt.Sprintf("%s:%d", conn.Laddr.IP, conn.Laddr.Port)

		// 상태 스타일링
		statusStr := getStatusStyled(conn.Status)

		// 프로세스 정보 가져오기
		pid := int(conn.Pid)
		processInfo := utils.GetProcessInfo(int32(pid))

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

	return createTable(rows, []tableColumn{
		{width: 10, title: "PROTOCOL"},
		{width: 19, title: "LOCAL ADDRESS"},
		{width: 10, title: "STATUS"},
		{width: 8, title: "PID"},
		{width: 25, title: "PROCESS NAME"},
		{width: 15, title: "USERNAME"},
		{width: 8, title: "CPU %"},
		{width: 8, title: "MEM %"},
	}, 130)
}

// tableColumn 테이블 컬럼 정의
type tableColumn struct {
	width int
	title string
}

// createTable 테이블 생성 헬퍼 함수
func createTable(rows [][]string, columns []tableColumn, totalWidth int) string {
	// 헤더 생성
	headerStyle := style.HeaderStyle.Copy().Align(lipgloss.Center)
	headers := make([]string, len(columns))
	for i, col := range columns {
		headers[i] = headerStyle.Width(col.width).Render(col.title)
	}

	// 테이블 생성
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(style.PrimaryColor)).
		StyleFunc(getTableRowStyle).
		Headers(headers...).
		Rows(rows...).
		Width(totalWidth)

	return t.String()
}

// getTableRowStyle 테이블 행 스타일 반환
func getTableRowStyle(row, col int) lipgloss.Style {
	color := lipgloss.Color("252")
	if row%2 != 0 {
		color = lipgloss.Color("248")
	}
	return lipgloss.NewStyle().
		Foreground(color).
		Align(lipgloss.Left).
		Padding(0, 1)
}

// ProcessInfoFormatter는 프로세스 정보를 포맷팅합니다
type ProcessInfoFormatter struct{}

func NewProcessInfoFormatter() *ProcessInfoFormatter {
	return &ProcessInfoFormatter{}
}

// Format은 프로세스 정보를 포맷팅합니다
func (f *ProcessInfoFormatter) Format(pid int, processName string, status []string, connections []net.ConnectionStat) string {
	title := style.TitleStyle.Render("⚙️  Process Information")

	// strings.Builder로 효율적인 문자열 생성
	var builder strings.Builder
	
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
	return fmt.Sprintf("\n%s\n%s", title, infoBox)
}

