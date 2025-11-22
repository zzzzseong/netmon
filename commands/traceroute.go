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

	"github.com/charmbracelet/lipgloss"
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

	// 헤더 출력
	header := fmt.Sprintf("Tracing route to %s", target)
	headerStyle := lipgloss.NewStyle().
		Foreground(style.PrimaryColor).
		Bold(true)
	fmt.Println(headerStyle.Render(header))
	fmt.Println()

	// 테이블 헤더 출력
	printTableHeader()

	// 실시간으로 traceroute 실행 및 출력
	err := executeTracerouteStreaming(target)
	if err != nil {
		errorMsg := style.ErrorStyle.Render(fmt.Sprintf("\nError: %v", err))
		fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		return err
	}

	fmt.Println()
	return nil
}

// printTableHeader는 테이블 헤더를 출력합니다
func printTableHeader() {
	headerStyle := style.HeaderStyle.Copy()
	
	hopHeader := headerStyle.Width(6).Render("HOP")
	hostHeader := headerStyle.Width(40).Render("HOST")
	rtt1Header := headerStyle.Width(12).Render("RTT 1")
	rtt2Header := headerStyle.Width(12).Render("RTT 2")
	rtt3Header := headerStyle.Width(12).Render("RTT 3")
	
	border := lipgloss.NewStyle().Foreground(style.PrimaryColor).Render("─")
	
	fmt.Printf("%s  %s  %s  %s  %s\n", hopHeader, hostHeader, rtt1Header, rtt2Header, rtt3Header)
	fmt.Printf("%s  %s  %s  %s  %s\n", 
		strings.Repeat(border, 6),
		strings.Repeat(border, 40),
		strings.Repeat(border, 12),
		strings.Repeat(border, 12),
		strings.Repeat(border, 12))
}

// printHopLine은 한 줄의 hop 정보를 출력합니다
func printHopLine(hop TraceHop) {
	// Hop 번호
	hopStr := fmt.Sprintf("%-6d", hop.Hop)
	hopStr = lipgloss.NewStyle().Foreground(style.SecondaryColor).Bold(true).Render(hopStr)

	// Host/IP
	hostStr := hop.Host
	if hostStr == "" {
		hostStr = "*"
	}
	if hop.Status == "timeout" {
		hostStr = "Request timed out"
		hostStr = lipgloss.NewStyle().Foreground(style.WarningColor).Render(fmt.Sprintf("%-40s", hostStr))
	} else {
		hostStr = lipgloss.NewStyle().Foreground(style.InfoColor).Render(fmt.Sprintf("%-40s", hostStr))
	}

	// RTT 값들
	rtt1Str := formatRTTSimple(hop.RTT1)
	rtt2Str := formatRTTSimple(hop.RTT2)
	rtt3Str := formatRTTSimple(hop.RTT3)

	fmt.Printf("%s  %s  %s  %s  %s\n", hopStr, hostStr, rtt1Str, rtt2Str, rtt3Str)
}

// formatRTTSimple은 RTT 값을 간단한 형식으로 포맷팅합니다
func formatRTTSimple(rtt string) string {
	if rtt == "*" {
		return lipgloss.NewStyle().Foreground(style.SubtleColor).Render(fmt.Sprintf("%-12s", "*"))
	}

	// RTT 값에서 숫자만 추출
	rttValue := strings.TrimSuffix(rtt, " ms")
	val, err := strconv.ParseFloat(rttValue, 64)
	if err != nil {
		return lipgloss.NewStyle().Foreground(style.SubtleColor).Render(fmt.Sprintf("%-12s", rtt))
	}

	// RTT 값에 따라 색상 변경
	var color lipgloss.Color
	if val < 30 {
		color = style.SuccessColor
	} else if val < 100 {
		color = style.WarningColor
	} else {
		color = lipgloss.Color("196")
	}

	return lipgloss.NewStyle().Foreground(color).Render(fmt.Sprintf("%-12s", rtt))
}

// executeTracerouteStreaming은 traceroute를 실시간으로 실행하고 출력합니다
func executeTracerouteStreaming(target string) error {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// Windows: tracert
		cmd = exec.Command("tracert", "-d", "-w", "3000", target)
	} else {
		// Unix-like: traceroute
		cmd = exec.Command("traceroute", "-n", "-w", "3", "-q", "3", target)
	}

	// stdout 파이프 생성
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create pipe: %v", err)
	}

	// 명령 시작
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start traceroute: %v", err)
	}

	// 실시간으로 출력 읽기 및 파싱
	scanner := bufio.NewScanner(stdout)
	
	if runtime.GOOS == "windows" {
		parseWindowsTracertStreaming(scanner)
	} else {
		parseUnixTracerouteStreaming(scanner)
	}

	// 프로세스 종료 대기
	if err := cmd.Wait(); err != nil {
		// 부분적인 결과가 출력되었을 수 있으므로 에러를 무시
		return nil
	}

	return nil
}

// parseUnixTracerouteStreaming은 Unix traceroute를 실시간으로 파싱합니다
func parseUnixTracerouteStreaming(scanner *bufio.Scanner) {
	// 첫 줄 건너뛰기 (헤더)
	if scanner.Scan() {
		// traceroute to ... 줄
	}

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

		// 실시간으로 hop 출력
		printHopLine(hop)
	}
}

// parseWindowsTracertStreaming은 Windows tracert를 실시간으로 파싱합니다
func parseWindowsTracertStreaming(scanner *bufio.Scanner) {
	// 헤더 건너뛰기
	headerSkipped := 0
	for scanner.Scan() && headerSkipped < 4 {
		headerSkipped++
	}

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

		if strings.Contains(rest, "timed out") || strings.Contains(rest, "* * *") {
			hop.Status = "timeout"
		} else {
			if ipMatch := ipRegex.FindString(rest); ipMatch != "" {
				hop.IP = ipMatch
				hop.Host = ipMatch
				hop.Status = "success"
			}

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

		// 실시간으로 hop 출력
		printHopLine(hop)
	}
}


