package cmd

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"netmon/style"
)

// tracerouteCmd represents the traceroute command
var tracerouteCmd = &cobra.Command{
	Use:   "traceroute <host>",
	Short: "Trace route to network host",
	Long:  `Trace the network path to a destination with animated loading.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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
			// 에러 메시지는 executeTracerouteStreaming 내부에서 출력하지 않고 여기서 반환
			// 하지만 기존 로직상 에러 메시지 스타일링이 있으므로 유지하거나 리팩토링
			// 여기서는 에러를 반환하여 Cobra가 출력하게 하거나, 직접 출력하고 nil 반환
			// 기존 스타일을 유지하기 위해 직접 출력하고 에러 반환
			return err
		}

		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tracerouteCmd)
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
	var cmdName string

	if runtime.GOOS == "windows" {
		// Windows: tracert
		cmdName = "tracert"
		cmd = exec.Command("tracert", "-d", "-w", "3000", target)
	} else {
		// Unix-like: traceroute
		cmdName = "traceroute"
		cmd = exec.Command("traceroute", "-n", "-w", "3", "-q", "3", target)
	}

	// 명령어가 PATH에 있는지 확인
	if _, err := exec.LookPath(cmdName); err != nil {
		// OS별 설치 안내 메시지 출력
		installMsg := getTracerouteInstallMessage()
		return fmt.Errorf("%s command not found in PATH\n%s", cmdName, installMsg)
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
		for i, match := range rttMatches {
			if i >= 3 {
				break
			}
			rttValue := match[1] + " ms"
			switch i {
			case 0:
				hop.RTT1 = rttValue
			case 1:
				hop.RTT2 = rttValue
			case 2:
				hop.RTT3 = rttValue
			}
		}

		// 실시간으로 hop 출력
		printHopLine(hop)
	}
}

// getTracerouteInstallMessage는 OS별 traceroute 설치 안내 메시지를 반환합니다
func getTracerouteInstallMessage() string {
	var msg string
	
	switch runtime.GOOS {
	case "linux":
		msg = `Please install traceroute using one of the following commands:
  • Debian/Ubuntu: sudo apt-get install traceroute
  • RHEL/CentOS:   sudo yum install traceroute
  • Fedora:        sudo dnf install traceroute
  • Arch Linux:    sudo pacman -S traceroute`
	case "darwin":
		msg = `Please install traceroute using Homebrew:
  brew install traceroute
  
Note: macOS usually includes traceroute by default. If you see this message,
traceroute may have been removed or is not in your PATH.`
	case "windows":
		msg = `Windows includes tracert by default. If you see this message,
please check your system PATH configuration.`
	default:
		msg = fmt.Sprintf(`Please install traceroute for your operating system (%s).
Check your system's package manager for installation instructions.`, runtime.GOOS)
	}
	
	return style.DescStyle.Render(msg)
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
			for i, match := range rttMatches {
				if i >= 3 {
					break
				}
				rttValue := match[1] + " ms"
				switch i {
				case 0:
					hop.RTT1 = rttValue
				case 1:
					hop.RTT2 = rttValue
				case 2:
					hop.RTT3 = rttValue
				}
			}
		}

		// 실시간으로 hop 출력
		printHopLine(hop)
	}
}
