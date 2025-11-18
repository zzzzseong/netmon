package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/charmbracelet/lipgloss"
	"netmon/formatter"
	"netmon/style"
)

// ShutdownCommand는 프로세스를 종료하는 명령어입니다
type ShutdownCommand struct {
	formatter *formatter.ProcessInfoFormatter
}

// NewShutdownCommand는 새로운 ShutdownCommand를 생성합니다
func NewShutdownCommand() *ShutdownCommand {
	return &ShutdownCommand{
		formatter: formatter.NewProcessInfoFormatter(),
	}
}

// Name은 명령어 이름을 반환합니다
func (c *ShutdownCommand) Name() string {
	return "shutdown"
}

// Description은 명령어 설명을 반환합니다
func (c *ShutdownCommand) Description() string {
	return "Shutdown a process"
}

// Usage는 명령어 사용법을 반환합니다
func (c *ShutdownCommand) Usage() string {
	return "shutdown <pid>"
}

// Execute는 명령어를 실행합니다
func (c *ShutdownCommand) Execute(args []string) error {
	if len(args) < 1 {
		errorMsg := style.ErrorStyle.Render("Error: PID is required.")
		usageMsg := style.UsageStyle.Render("Usage: netmon shutdown <pid>")
		fmt.Fprintf(os.Stderr, "%s\n%s\n", errorMsg, usageMsg)
		os.Exit(1)
	}

	pidStr := args[0]
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		errorMsg := style.ErrorStyle.Render(fmt.Sprintf("Error: Invalid PID: %s", pidStr))
		fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		os.Exit(1)
	}

	// 프로세스 정보 가져오기
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		errorMsg := style.ErrorStyle.Render(fmt.Sprintf("Error: Process with PID %d not found: %v", pid, err))
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

	// 포맷터를 사용하여 프로세스 정보 출력
	info := c.formatter.Format(pid, processName, status, connections)
	fmt.Println(info)

	// 인터랙티브 확인 프롬프트
	promptTitle := lipgloss.NewStyle().
		Foreground(style.WarningColor).
		Bold(true).
		MarginTop(1).
		Render("⚠️  Do you want to shutdown this process?")
	fmt.Println(promptTitle)

	prompt := promptui.Select{
		Label: "",
		Items: []string{"Shutdown", "Cancel"},
		Size:  2,
	}

	index, result, err := prompt.Run()
	if err != nil {
		errorMsg := style.ErrorStyle.Render(fmt.Sprintf("Error: Prompt error: %v", err))
		fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		os.Exit(1)
	}

	if index == 0 && result == "Shutdown" {
		// 프로세스 종료
		err = proc.Kill()
		if err != nil {
			errorMsg := style.ErrorStyle.Render(fmt.Sprintf("Error: Failed to shutdown process: %v", err))
			fmt.Fprintf(os.Stderr, "\n%s\n", errorMsg)
			os.Exit(1)
		}
		successMsg := style.SuccessStyle.Render(fmt.Sprintf("✓ Process %d (%s) has been successfully shut down.", pid, processName))
		fmt.Printf("\n%s\n", successMsg)
	} else {
		cancelMsg := lipgloss.NewStyle().
			Foreground(style.SubtleColor).
			Render("Cancelled.")
		fmt.Printf("\n%s\n", cancelMsg)
	}

	return nil
}

