package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/manifoldco/promptui"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// 스타일 정의
var (
	// 색상 팔레트
	primaryColor   = lipgloss.Color("63")   // 보라색
	secondaryColor = lipgloss.Color("212")  // 핑크색
	successColor   = lipgloss.Color("35")   // 초록색
	warningColor   = lipgloss.Color("220")  // 노란색
	infoColor      = lipgloss.Color("39")   // 파란색
	dangerColor    = lipgloss.Color("196")  // 빨간색
	subtleColor    = lipgloss.Color("241")  // 회색

	// 타이틀 스타일
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			MarginBottom(1)

	// 사용법 스타일
	usageStyle = lipgloss.NewStyle().
			Foreground(subtleColor).
			MarginTop(1)

	// 명령어 스타일
	commandStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)

	// 설명 스타일
	descStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	// 에러 스타일
	errorStyle = lipgloss.NewStyle().
			Foreground(dangerColor).
			Bold(true)

	// 성공 스타일
	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	// 정보 박스 스타일
	infoBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			Margin(1, 0)

	// 라벨 스타일
	labelStyle = lipgloss.NewStyle().
			Foreground(subtleColor).
			Width(12)

	// 값 스타일
	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Bold(true)
)

func main() {
	// 명령어 인자 확인
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "ls":
		listPorts()
	case "kill":
		if len(os.Args) < 3 {
			errorMsg := errorStyle.Render("Error: PID가 필요합니다.")
			usageMsg := usageStyle.Render("사용법: netmon kill <pid>")
			fmt.Fprintf(os.Stderr, "%s\n%s\n", errorMsg, usageMsg)
			os.Exit(1)
		}
		pidStr := os.Args[2]
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			errorMsg := errorStyle.Render(fmt.Sprintf("Error: 유효하지 않은 PID: %s", pidStr))
			fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
			os.Exit(1)
		}
		killProcess(pid)
	default:
		errorMsg := errorStyle.Render(fmt.Sprintf("Error: 알 수 없는 명령어: %s", command))
		fmt.Fprintf(os.Stderr, "%s\n\n", errorMsg)
		printUsage()
		os.Exit(1)
	}
}

// 사용법 출력
func printUsage() {
	title := titleStyle.Render("📡 netmon - 네트워크 모니터링 도구")
	
	usage := usageStyle.Render("사용법:")
	cmd1 := commandStyle.Render("  netmon ls") + "              " + descStyle.Render("- 활성 포트 목록 표시")
	cmd2 := commandStyle.Render("  netmon kill <pid>") + "      " + descStyle.Render("- 프로세스 종료")
	
	fmt.Println(title)
	fmt.Println(usage)
	fmt.Println(cmd1)
	fmt.Println(cmd2)
}

// 연결 타입을 문자열로 변환 (1=TCP, 2=UDP)
func connectionTypeToString(connType uint32) string {
	switch connType {
	case 1:
		return "TCP"
	case 2:
		return "UDP"
	default:
		return "UNKNOWN"
	}
}

// 연결 타입이 UDP인지 확인
func isUDP(connType uint32) bool {
	return connType == 2
}

// 포트 목록 출력 함수
func listPorts() {
	// 모든 네트워크 연결 가져오기
	connections, err := net.Connections("inet")
	if err != nil {
		errorMsg := errorStyle.Render(fmt.Sprintf("Error: 네트워크 연결 정보를 가져올 수 없습니다: %v", err))
		fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		os.Exit(1)
	}

	// LISTEN 상태인 연결만 필터링 (UDP는 상태가 없을 수 있으므로 포트가 열려있으면 포함)
	listeningConns := make(map[string]net.ConnectionStat)
	for _, conn := range connections {
		// TCP는 LISTEN 상태만, UDP는 포트가 열려있으면 모두 포함
		if conn.Status == "LISTEN" || (isUDP(conn.Type) && conn.Laddr.Port > 0) {
			key := fmt.Sprintf("%d:%d", conn.Type, conn.Laddr.Port)
			listeningConns[key] = conn
		}
	}

	// 테이블 행 데이터 준비
	var rows [][]string
	for _, conn := range listeningConns {
		protocol := connectionTypeToString(conn.Type)
		localAddr := fmt.Sprintf("%s:%d", conn.Laddr.IP, conn.Laddr.Port)
		
		// 상태에 따른 색상 적용 (UDP는 상태가 없을 수 있음)
		var statusStr string
		statusStyle := lipgloss.NewStyle()
		if conn.Status == "" {
			statusStr = "LISTEN"
			statusStyle = statusStyle.Foreground(successColor).Bold(true)
		} else {
			switch conn.Status {
			case "LISTEN":
				statusStr = conn.Status
				statusStyle = statusStyle.Foreground(successColor).Bold(true)
			case "ESTABLISHED":
				statusStr = conn.Status
				statusStyle = statusStyle.Foreground(infoColor).Bold(true)
			default:
				statusStr = conn.Status
				statusStyle = statusStyle.Foreground(warningColor).Bold(true)
			}
		}
		statusStr = statusStyle.Render(statusStr)

		// PID와 프로세스 이름 가져오기
		pid := int(conn.Pid)
		processName := "N/A"
		
		if pid > 0 {
			proc, err := process.NewProcess(int32(pid))
			if err == nil {
				name, err := proc.Name()
				if err == nil {
					processName = name
				}
			}
		}

		// 프로토콜 스타일링
		protocolStyle := lipgloss.NewStyle().Foreground(secondaryColor).Bold(true)
		protocolStr := protocolStyle.Render(protocol)

		// PID 스타일링
		pidStr := fmt.Sprintf("%d", pid)
		if pid > 0 {
			pidStr = lipgloss.NewStyle().Foreground(infoColor).Render(pidStr)
		}

		rows = append(rows, []string{
			protocolStr,
			localAddr,
			statusStr,
			pidStr,
			processName,
		})
	}

	// 헤더 스타일 적용
	headerStyle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Align(lipgloss.Center).
		Padding(0, 1)

	styledHeaders := []string{
		headerStyle.Render("PROTOCOL"),
		headerStyle.Render("LOCAL ADDRESS"),
		headerStyle.Render("STATUS"),
		headerStyle.Render("PID"),
		headerStyle.Render("PROCESS NAME"),
	}

	// 테이블 생성 및 스타일링
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(primaryColor)).
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
		Width(120)

	// 타이틀 출력
	// title := titleStyle.Render("🌐 활성 포트 목록")
	// fmt.Println(title)
	fmt.Println(t)
}

// 프로세스 종료 함수
func killProcess(pid int) {
	// 프로세스 정보 가져오기
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		errorMsg := errorStyle.Render(fmt.Sprintf("Error: PID %d의 프로세스를 찾을 수 없습니다: %v", pid, err))
		fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		os.Exit(1)
	}

	// 프로세스 이름 가져오기
	processName, err := proc.Name()
	if err != nil {
		processName = "N/A"
	}

	// 프로세스 상태 정보 가져오기
	status, err := proc.Status()
	if err != nil {
		status = []string{"N/A"}
	}

	// 프로세스가 사용하는 포트 가져오기
	connections, err := net.ConnectionsPid("inet", int32(pid))
	if err != nil {
		connections = []net.ConnectionStat{}
	}

	// 프로세스 정보 박스 생성
	title := titleStyle.Render("⚙️  프로세스 정보")
	
	infoContent := fmt.Sprintf("%s%s\n",
		labelStyle.Render("PID:"),
		valueStyle.Render(fmt.Sprintf("%d", pid)))
	
	infoContent += fmt.Sprintf("%s%s\n",
		labelStyle.Render("이름:"),
		valueStyle.Render(processName))
	
	statusStr := fmt.Sprintf("%v", status)
	infoContent += fmt.Sprintf("%s%s\n",
		labelStyle.Render("상태:"),
		valueStyle.Render(statusStr))
	
	// 포트 정보 추가
	if len(connections) > 0 {
		infoContent += "\n" + lipgloss.NewStyle().Foreground(subtleColor).Render("사용 중인 포트:") + "\n"
		for _, conn := range connections {
			// LISTEN 상태이거나 UDP인 경우 표시
			if conn.Status == "LISTEN" || (conn.Status == "" && isUDP(conn.Type)) {
				portInfo := fmt.Sprintf("  • %s:%d (%s)",
					conn.Laddr.IP,
					conn.Laddr.Port,
					connectionTypeToString(conn.Type))
				portStyle := lipgloss.NewStyle().
					Foreground(infoColor).
					MarginLeft(2)
				infoContent += portStyle.Render(portInfo) + "\n"
			}
		}
	}

	infoBox := infoBoxStyle.Render(infoContent)
	
	fmt.Println()
	fmt.Println(title)
	fmt.Println(infoBox)

	// 인터랙티브 확인 프롬프트
	promptTitle := lipgloss.NewStyle().
		Foreground(warningColor).
		Bold(true).
		MarginTop(1).
		Render("⚠️  이 프로세스를 종료하시겠습니까?")
	fmt.Println(promptTitle)
	
	prompt := promptui.Select{
		Label: "",
		Items: []string{"Kill", "Cancel"},
		Size:  2,
	}

	index, result, err := prompt.Run()
	if err != nil {
		errorMsg := errorStyle.Render(fmt.Sprintf("Error: 프롬프트 오류: %v", err))
		fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		os.Exit(1)
	}

	if index == 0 && result == "Kill" {
		// 프로세스 종료
		err = proc.Kill()
		if err != nil {
			errorMsg := errorStyle.Render(fmt.Sprintf("Error: 프로세스를 종료할 수 없습니다: %v", err))
			fmt.Fprintf(os.Stderr, "\n%s\n", errorMsg)
			os.Exit(1)
		}
		successMsg := successStyle.Render(fmt.Sprintf("✓ 프로세스 %d (%s)가 성공적으로 종료되었습니다.", pid, processName))
		fmt.Printf("\n%s\n", successMsg)
	} else {
		cancelMsg := lipgloss.NewStyle().
			Foreground(subtleColor).
			Render("취소되었습니다.")
		fmt.Printf("\n%s\n", cancelMsg)
	}
}

