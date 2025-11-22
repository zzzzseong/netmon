package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"netmon/style"
)

// TraceCommand는 traceroute 정보를 표시하는 명령어입니다
type TraceCommand struct{}

// NewTraceCommand는 새로운 TraceCommand를 생성합니다
func NewTraceCommand() *TraceCommand {
	return &TraceCommand{}
}

// Name은 명령어 이름을 반환합니다
func (c *TraceCommand) Name() string {
	return "traceroute"
}

// Description은 명령어 설명을 반환합니다
func (c *TraceCommand) Description() string {
	return "Trace route to network host"
}

// Usage는 명령어 사용법을 반환합니다
func (c *TraceCommand) Usage() string {
	return "traceroute <host>"
}

// TraceHop은 traceroute의 각 홉 정보를 담습니다
type TraceHop struct {
	Hop     int
	Host    string
	IP      string
	RTT1    string
	RTT2    string
	RTT3    string
	Status  string // "success", "timeout", "error"
}

// Execute는 명령어를 실행합니다
func (c *TraceCommand) Execute(args []string) error {
	if len(args) < 1 {
		errorMsg := style.ErrorStyle.Render("Error: Host is required")
		fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		fmt.Println(style.DescStyle.Render("Usage: netmon traceroute <host>"))
		return fmt.Errorf("host is required")
	}

	target := args[0]

	// 로딩 메시지와 스피너를 위한 채널
	done := make(chan bool)
	var hops []TraceHop
	var traceErr error

	// 백그라운드에서 traceroute 실행
	go func() {
		hops, traceErr = executeTraceroute(target)
		done <- true
	}()

	// 스피너 애니메이션 표시
	showSpinner(target, done)

	// 에러 체크
	if traceErr != nil {
		errorMsg := style.ErrorStyle.Render(fmt.Sprintf("Error: %v", traceErr))
		fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		return traceErr
	}

	// 결과가 없으면 에러
	if len(hops) == 0 {
		errorMsg := style.ErrorStyle.Render("Error: No route information received")
		fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		return fmt.Errorf("no route information")
	}

	// 테이블 형식으로 출력
	tableStr := formatTraceTable(hops)
	fmt.Println(tableStr)

	return nil
}

// showSpinner는 로딩 스피너를 표시합니다
func showSpinner(target string, done chan bool) {
	spinnerFrames := []string{"◐", "◓", "◑", "◒"}
	i := 0
	
	spinnerStyle := lipgloss.NewStyle().
		Foreground(style.SecondaryColor).
		Bold(true)
	
	textStyle := lipgloss.NewStyle().
		Foreground(style.PrimaryColor)

	ticker := time.NewTicker(150 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			// 스피너 완전히 지우기
			fmt.Print("\r\033[K")
			return
		case <-ticker.C:
			frame := spinnerFrames[i%len(spinnerFrames)]
			spinner := spinnerStyle.Render(frame)
			text := textStyle.Render(fmt.Sprintf(" Tracing route to %s...", target))
			fmt.Print("\r" + spinner + text)
			i++
		}
	}
}

// executeTraceroute는 플랫폼에 따라 적절한 traceroute 명령을 실행합니다
func executeTraceroute(target string) ([]TraceHop, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// Windows: tracert
		cmd = exec.Command("tracert", "-d", "-w", "3000", target)
	} else {
		// Unix-like: traceroute
		cmd = exec.Command("traceroute", "-n", "-w", "3", "-q", "3", target)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute traceroute: %v", err)
	}

	// 출력 파싱
	return parseTracerouteOutput(string(output), runtime.GOOS)
}

// parseTracerouteOutput은 traceroute 출력을 파싱합니다
func parseTracerouteOutput(output string, goos string) ([]TraceHop, error) {
	scanner := bufio.NewScanner(strings.NewReader(output))

	// Windows와 Unix 계열의 출력 형식이 다름
	if goos == "windows" {
		return parseWindowsTracert(scanner)
	}
	return parseUnixTraceroute(scanner)
}

// parseUnixTraceroute는 Unix 계열 traceroute 출력을 파싱합니다
func parseUnixTraceroute(scanner *bufio.Scanner) ([]TraceHop, error) {
	// 첫 줄 건너뛰기 (헤더)
	if scanner.Scan() {
		// traceroute to ... 줄
	}

	var hops []TraceHop

	// 각 홉 파싱
	// 예: " 1  192.168.1.1  1.234 ms  1.456 ms  1.789 ms"
	// 예: " 2  * * *"
	hopRegex := regexp.MustCompile(`^\s*(\d+)\s+(.*)$`)
	ipRegex := regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+|[0-9a-fA-F:]+)`)
	rttRegex := regexp.MustCompile(`(\d+\.\d+)\s*ms`)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		matches := hopRegex.FindStringSubmatch(line)
		if len(matches) < 3 {
			continue
		}

		hopNum, _ := strconv.Atoi(matches[1])
		rest := matches[2]

		hop := TraceHop{
			Hop:    hopNum,
			RTT1:   "*",
			RTT2:   "*",
			RTT3:   "*",
			Status: "timeout",
		}

		// IP 주소 추출
		if ipMatch := ipRegex.FindString(rest); ipMatch != "" {
			hop.IP = ipMatch
			hop.Host = ipMatch
			hop.Status = "success"
		}

		// RTT 추출
		rttMatches := rttRegex.FindAllStringSubmatch(rest, -1)
		if len(rttMatches) > 0 {
			if len(rttMatches) > 0 {
				hop.RTT1 = rttMatches[0][1] + " ms"
			}
			if len(rttMatches) > 1 {
				hop.RTT2 = rttMatches[1][1] + " ms"
			}
			if len(rttMatches) > 2 {
				hop.RTT3 = rttMatches[2][1] + " ms"
			}
		}

		hops = append(hops, hop)
	}

	return hops, nil
}

// parseWindowsTracert는 Windows tracert 출력을 파싱합니다
func parseWindowsTracert(scanner *bufio.Scanner) ([]TraceHop, error) {
	var hops []TraceHop

	// 헤더 건너뛰기 (처음 몇 줄)
	headerSkipped := 0
	for scanner.Scan() && headerSkipped < 4 {
		headerSkipped++
	}

	// 각 홉 파싱
	// 예: "  1    <1 ms    <1 ms    <1 ms  192.168.1.1"
	// 예: "  2     *        *        *     Request timed out."
	hopRegex := regexp.MustCompile(`^\s*(\d+)\s+(.*)$`)
	ipRegex := regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+)`)
	rttRegex := regexp.MustCompile(`(<?\d+)\s*ms`)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		matches := hopRegex.FindStringSubmatch(line)
		if len(matches) < 3 {
			continue
		}

		hopNum, _ := strconv.Atoi(matches[1])
		rest := matches[2]

		hop := TraceHop{
			Hop:    hopNum,
			RTT1:   "*",
			RTT2:   "*",
			RTT3:   "*",
			Status: "timeout",
		}

		// timeout 체크
		if strings.Contains(rest, "timed out") || strings.Contains(rest, "* * *") {
			hop.Status = "timeout"
		} else {
			// IP 주소 추출
			if ipMatch := ipRegex.FindString(rest); ipMatch != "" {
				hop.IP = ipMatch
				hop.Host = ipMatch
				hop.Status = "success"
			}

			// RTT 추출
			rttMatches := rttRegex.FindAllStringSubmatch(rest, -1)
			if len(rttMatches) > 0 {
				if len(rttMatches) > 0 {
					hop.RTT1 = rttMatches[0][1] + " ms"
				}
				if len(rttMatches) > 1 {
					hop.RTT2 = rttMatches[1][1] + " ms"
				}
				if len(rttMatches) > 2 {
					hop.RTT3 = rttMatches[2][1] + " ms"
				}
			}
		}

		hops = append(hops, hop)
	}

	return hops, nil
}

// formatTraceTable은 traceroute 정보를 테이블 형식으로 포맷팅합니다
func formatTraceTable(hops []TraceHop) string {
	var rows [][]string

	for _, hop := range hops {
		// Hop 번호
		hopStr := strconv.Itoa(hop.Hop)
		hopStr = lipgloss.NewStyle().Foreground(style.SecondaryColor).Bold(true).Render(hopStr)

		// Host/IP
		hostStr := hop.Host
		if hostStr == "" {
			hostStr = "*"
		}
		if hop.Status == "timeout" {
			hostStr = lipgloss.NewStyle().Foreground(style.WarningColor).Render("Request timed out")
		} else {
			hostStr = lipgloss.NewStyle().Foreground(style.InfoColor).Render(hostStr)
		}

		// RTT 값들 (색상으로 속도 표시)
		rtt1Str := formatRTT(hop.RTT1)
		rtt2Str := formatRTT(hop.RTT2)
		rtt3Str := formatRTT(hop.RTT3)

		rows = append(rows, []string{
			hopStr,
			hostStr,
			rtt1Str,
			rtt2Str,
			rtt3Str,
		})
	}

	// 헤더 스타일 적용
	headerStyle := style.HeaderStyle.Copy().Align(lipgloss.Center)
	styledHeaders := []string{
		headerStyle.Width(6).Render("HOP"),
		headerStyle.Width(40).Render("HOST"),
		headerStyle.Width(12).Render("RTT 1"),
		headerStyle.Width(12).Render("RTT 2"),
		headerStyle.Width(12).Render("RTT 3"),
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
		Width(100)

	return t.String()
}

// formatRTT는 RTT 값을 색상으로 포맷팅합니다
func formatRTT(rtt string) string {
	if rtt == "*" {
		return lipgloss.NewStyle().Foreground(style.SubtleColor).Render("*")
	}

	// RTT 값에서 숫자만 추출
	rttValue := strings.TrimSuffix(rtt, " ms")
	val, err := strconv.ParseFloat(rttValue, 64)
	if err != nil {
		return lipgloss.NewStyle().Foreground(style.SubtleColor).Render(rtt)
	}

	// RTT 값에 따라 색상 변경
	var color lipgloss.Color
	if val < 30 {
		// 빠름: 초록색
		color = style.SuccessColor
	} else if val < 100 {
		// 보통: 노란색
		color = style.WarningColor
	} else {
		// 느림: 빨간색
		color = lipgloss.Color("196")
	}

	return lipgloss.NewStyle().Foreground(color).Render(rtt)
}

