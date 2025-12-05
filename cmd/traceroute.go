package cmd

import (
	"bufio"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"netmon/formatter"
	"netmon/parser"
	"netmon/style"
)

const (
	// tracerouteTimeoutSeconds is the timeout in seconds for Unix traceroute
	tracerouteTimeoutSeconds = 3
	// tracerouteProbeCount is the number of probes per hop for Unix traceroute
	tracerouteProbeCount = 3
	// tracertTimeoutMs is the timeout in milliseconds for Windows tracert
	tracertTimeoutMs = 3000
)

// newTracerouteCmd creates and returns the traceroute command.
// It traces the network path to a destination with animated loading.
func newTracerouteCmd() *cobra.Command {
	return &cobra.Command{
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

			// 포맷터와 파서 생성
			fmtter := formatter.NewTracerouteFormatter()
			parser := parser.NewTracerouteParser()

			// 테이블 헤더 출력
			fmtter.PrintTableHeader()

			// 실시간으로 traceroute 실행 및 출력
			err := executeTracerouteStreaming(target, parser)
			if err != nil {
				return err
			}

			fmt.Println()
			return nil
		},
	}
}

// executeTracerouteStreaming executes traceroute in real-time and streams the output.
func executeTracerouteStreaming(target string, parser *parser.TracerouteParser) error {
	var cmd *exec.Cmd
	var cmdName string

	if runtime.GOOS == "windows" {
		// Windows: tracert
		cmdName = "tracert"
		cmd = exec.Command("tracert", "-d", "-w", strconv.Itoa(tracertTimeoutMs), target)
	} else {
		// Unix-like: traceroute
		cmdName = "traceroute"
		timeoutStr := strconv.Itoa(tracerouteTimeoutSeconds)
		probeStr := strconv.Itoa(tracerouteProbeCount)
		cmd = exec.Command("traceroute", "-n", "-w", timeoutStr, "-q", probeStr, target)
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
		return fmt.Errorf("failed to create pipe: %w", err)
	}

	// 명령 시작
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start traceroute: %w", err)
	}

	// 실시간으로 출력 읽기 및 파싱
	scanner := bufio.NewScanner(stdout)

	if runtime.GOOS == "windows" {
		parser.ParseWindowsTracert(scanner)
	} else {
		parser.ParseUnixTraceroute(scanner)
	}

	// 프로세스 종료 대기
	if err := cmd.Wait(); err != nil {
		// 부분적인 결과가 출력되었을 수 있으므로 에러를 무시
		return nil
	}

	return nil
}

// getTracerouteInstallMessage returns OS-specific installation instructions for traceroute.
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
