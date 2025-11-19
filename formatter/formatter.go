package formatter

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"netmon/style"
	"netmon/utils"
)

// PortTableFormatter는 포트 목록을 테이블 형식으로 포맷팅합니다
type PortTableFormatter struct{}

func NewPortTableFormatter() *PortTableFormatter {
	return &PortTableFormatter{}
}

// Format은 연결 정보를 테이블 형식으로 포맷팅합니다
func (f *PortTableFormatter) Format(connections map[string]net.ConnectionStat) string {
	var rows [][]string
	for _, conn := range connections {
		protocol := utils.ConnectionTypeToString(conn.Type)
		localAddr := fmt.Sprintf("%s:%d", conn.Laddr.IP, conn.Laddr.Port)

		// 상태에 따른 색상 적용 (UDP는 상태가 없을 수 있음)
		var statusStr string
		statusStyle := lipgloss.NewStyle()
		if conn.Status == "" {
			statusStr = "LISTEN"
			statusStyle = statusStyle.Foreground(style.SuccessColor).Bold(true)
		} else {
			switch conn.Status {
			case "LISTEN":
				statusStr = conn.Status
				statusStyle = statusStyle.Foreground(style.SuccessColor).Bold(true)
			case "ESTABLISHED":
				statusStr = conn.Status
				statusStyle = statusStyle.Foreground(style.InfoColor).Bold(true)
			default:
				statusStr = conn.Status
				statusStyle = statusStyle.Foreground(style.WarningColor).Bold(true)
			}
		}
		statusStr = statusStyle.Render(statusStr)

	// PID와 프로세스 이름 가져오기
	pid := int(conn.Pid)
	processName := "N/A"
	username := "N/A"
	cpuPercent := "N/A"
	memPercent := "N/A"

	if pid > 0 {
		proc, err := process.NewProcess(int32(pid))
		if err == nil {
			name, err := proc.Name()
			if err == nil {
				processName = name
			}

			// 사용자 이름 가져오기
			user, err := proc.Username()
			if err == nil {
				username = user
			}

			// CPU 사용률 가져오기 (0.1초 대기)
			cpu, err := proc.CPUPercent()
			if err == nil {
				cpuPercent = fmt.Sprintf("%.1f%%", cpu)
			}

			// 메모리 사용률 가져오기
			mem, err := proc.MemoryPercent()
			if err == nil {
				memPercent = fmt.Sprintf("%.1f%%", mem)
			}
		}
	}

	// 프로토콜 스타일링
	protocolStyle := lipgloss.NewStyle().Foreground(style.SecondaryColor).Bold(true)
	protocolStr := protocolStyle.Render(protocol)

	// PID 스타일링
	pidStr := fmt.Sprintf("%d", pid)
	if pid > 0 {
		pidStr = lipgloss.NewStyle().Foreground(style.InfoColor).Render(pidStr)
	}

	// CPU, 메모리 스타일링
	cpuStyle := lipgloss.NewStyle().Foreground(style.WarningColor)
	memStyle := lipgloss.NewStyle().Foreground(style.InfoColor)

	rows = append(rows, []string{
		protocolStr,
		localAddr,
		statusStr,
		pidStr,
		processName,
		username,
		cpuStyle.Render(cpuPercent),
		memStyle.Render(memPercent),
	})
	}

	// 헤더 스타일 적용 (가운데 정렬 및 너비 설정)
	headerStyle := style.HeaderStyle.Copy().Align(lipgloss.Center)
	styledHeaders := []string{
		headerStyle.Width(10).Render("PROTOCOL"),
		headerStyle.Width(19).Render("LOCAL ADDRESS"),
		headerStyle.Width(10).Render("STATUS"),
		headerStyle.Width(8).Render("PID"),
		headerStyle.Width(25).Render("PROCESS NAME"),
		headerStyle.Width(15).Render("USERNAME"),
		headerStyle.Width(8).Render("CPU %"),
		headerStyle.Width(8).Render("MEM %"),
	}

	// 테이블 생성 및 스타일링
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(style.PrimaryColor)).
		StyleFunc(func(row, col int) lipgloss.Style {
			// 모든 데이터 행에 왼쪽 정렬 적용
			if row%2 == 0 {
				return lipgloss.NewStyle().
					Foreground(lipgloss.Color("252")).
					Align(lipgloss.Left).
					Padding(0, 1)
			} else {
				return lipgloss.NewStyle().
					Foreground(lipgloss.Color("248")).
					Align(lipgloss.Left).
					Padding(0, 1)
			}
		}).
		Headers(styledHeaders...).
		Rows(rows...).
		Width(130)

	return t.String()
}

// ProcessInfoFormatter는 프로세스 정보를 포맷팅합니다
type ProcessInfoFormatter struct{}

func NewProcessInfoFormatter() *ProcessInfoFormatter {
	return &ProcessInfoFormatter{}
}

// Format은 프로세스 정보를 포맷팅합니다
func (f *ProcessInfoFormatter) Format(pid int, processName string, status []string, connections []net.ConnectionStat) string {
	title := style.TitleStyle.Render("⚙️  Process Information")

	infoContent := fmt.Sprintf("%s%s\n",
		style.LabelStyle.Render("PID:"),
		style.ValueStyle.Render(fmt.Sprintf("%d", pid)))

	infoContent += fmt.Sprintf("%s%s\n",
		style.LabelStyle.Render("Name:"),
		style.ValueStyle.Render(processName))

	statusStr := fmt.Sprintf("%v", status)
	infoContent += fmt.Sprintf("%s%s\n",
		style.LabelStyle.Render("Status:"),
		style.ValueStyle.Render(statusStr))

	// 포트 정보 추가
	if len(connections) > 0 {
		infoContent += "\n" + lipgloss.NewStyle().Foreground(style.SubtleColor).Render("Active Ports:") + "\n"
		for _, conn := range connections {
			// LISTEN 상태이거나 UDP인 경우 표시
			if conn.Status == "LISTEN" || (conn.Status == "" && utils.IsUDP(conn.Type)) {
				portInfo := fmt.Sprintf("  • %s:%d (%s)",
					conn.Laddr.IP,
					conn.Laddr.Port,
					utils.ConnectionTypeToString(conn.Type))
				portStyle := lipgloss.NewStyle().
					Foreground(style.InfoColor).
					MarginLeft(2)
				infoContent += portStyle.Render(portInfo) + "\n"
			}
		}
	}

	infoBox := style.InfoBoxStyle.Render(infoContent)

	return fmt.Sprintf("\n%s\n%s", title, infoBox)
}

